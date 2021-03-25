package iam

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	awsIAM "github.com/aws/aws-sdk-go/service/iam"
	as "github.com/clok/awssession"
	"github.com/clok/kemba"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/kyokomi/emoji/v2"
	"github.com/logrusorgru/aurora/v3"
	"github.com/remeh/sizedwaitgroup"
)

var (
	kiam = kemba.New("gw-aws-audit:iam")
)

func renderUsersReport(users []*User, showOnly string) error {
	kl := kiam.Extend("renderUsersReport")
	kl.Printf("processing %d users", len(users))

	// sort user list
	sort.Slice(users, func(i, j int) bool {
		return strings.ToLower(users[i].UserName()) < strings.ToLower(users[j].UserName())
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{
		"",
		"",
		"",
		"",
		"",
		"",
		"Access Key Details",
	})
	t.AppendHeader(table.Row{
		"User",
		"Status",
		"Age",
		"Console",
		"Last Login",
		"Permissions",
		aurora.Gray(8, "KEY ID | STATUS | AGE | LAST USED | SERVICE"),
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AlignFooter: text.AlignCenter},
		{Number: 2, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 3, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 4, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 5, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 6, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 7, Align: text.AlignCenter, AlignHeader: text.AlignCenter, AlignFooter: text.AlignCenter},
	})

	summaryStats := map[string]int{
		"pass":          0,
		"warn":          0,
		"fail":          0,
		"totalKeys":     0,
		"activeKeys":    0,
		"inactiveKeys":  0,
		"consoleAccess": 0,
	}
	for _, user := range users {
		if user.HasConsoleAccess() {
			summaryStats["consoleAccess"]++
		}
		summaryStats[user.CheckStatus()]++

		if showOnly == "" || showOnly == user.CheckStatus() {
			t.AppendRow([]interface{}{
				user.UserName(),
				user.FormattedCheckStatus(),
				user.CreatedDateDuration(),
				formattedYesNo(user.HasConsoleAccess()),
				user.FormattedLastLoginDateDuration(),
				fmt.Sprintf("G: %d P: %d", len(user.Groups()), len(user.Policies())),
				formattedKeyCount(user.AccessKeysCount()),
			})
		}
		if len(user.accessKeys) > 0 {
			st := table.NewWriter()
			st.SetStyle(table.StyleLight)
			st.Style().Options = table.OptionsNoBorders
			st.SetColumnConfigs([]table.ColumnConfig{
				{Number: 2, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
				{Number: 3, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
				{Number: 4, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
				{Number: 6, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
			})
			summaryStats["totalKeys"] += len(user.accessKeys)
			for _, key := range user.accessKeys {
				switch aws.StringValue(key.status) {
				case "Active":
					summaryStats["activeKeys"]++
				case "Inactive":
					summaryStats["inactiveKeys"]++
				}
				st.AppendRow([]interface{}{
					aws.StringValue(key.id),
					formattedStatus(aws.StringValue(key.status)),
					dateDuration(aws.TimeValue(key.createdDate), 1),
					formattedDateDuration(dateDuration(aws.TimeValue(key.lastUsed.LastUsedDate), 2)),
					aws.StringValue(key.lastUsed.ServiceName),
				})
			}
			if showOnly == "" || showOnly == user.CheckStatus() {
				t.AppendRow(table.Row{
					"",
					"",
					"",
					"",
					"",
					"",
					st.Render(),
				})
			}
		}
		if showOnly == "" || showOnly == user.CheckStatus() {
			t.AppendSeparator()
		}
	}

	f1 := fmt.Sprintf("PASS: %d WARN: %d FAIL: %d", summaryStats["pass"], summaryStats["warn"], summaryStats["fail"])
	f2 := fmt.Sprintf("%d / %d", summaryStats["consoleAccess"], len(users))
	f3 := fmt.Sprintf("ACTIVE: %d INACTIVE: %d TOTAL: %d", summaryStats["activeKeys"], summaryStats["inactiveKeys"], summaryStats["totalKeys"])
	t.AppendFooter(table.Row{f1, "", "", f2, f2, "", f3}, table.RowConfig{AutoMerge: true})
	t.Render()

	return nil
}

func getUser(user string) (*User, error) {
	kl := kiam.Extend("getUser")
	sess, err := as.New()
	if err != nil {
		return nil, err
	}
	client := awsIAM.New(sess)

	var result *awsIAM.GetUserOutput
	result, err = client.GetUser(&awsIAM.GetUserInput{
		UserName: &user,
	})
	if err != nil {
		return nil, err
	}
	kl.Log(result)

	var in *User
	in, err = buildUserData(result.User, true)
	if err != nil {
		return nil, err
	}

	return in, nil
}

func getAllUsers(fullDetails bool) ([]*User, error) {
	kl := kiam.Extend("getAllUsers")
	kmeta := kl.Extend("meta")
	sess, err := as.New()
	if err != nil {
		return nil, err
	}
	client := awsIAM.New(sess)

	var results *awsIAM.ListUsersOutput
	results, err = client.ListUsers(&awsIAM.ListUsersInput{
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return nil, err
	}

	swg := sizedwaitgroup.New(10)

	startTime := time.Now()
	var cur int64
	total := int64(len(results.Users))
	users := make([]*User, total)
	kl.Printf("found %d users", total)
	// TODO: print out status
	for i, user := range results.Users {
		go func(i int, user *awsIAM.User) {
			defer swg.Done()
			kmeta.Printf("[%d] Executing goroutine for user %s", i, aws.StringValue(user.UserName))
			var iu *User
			iu, err = buildUserData(user, fullDetails)
			if err != nil {
				// TODO: Handle panic
				panic(err)
			}

			users[i] = iu
			atomic.AddInt64(&cur, 1)
			dps := float64(cur) / time.Since(startTime).Seconds()
			_, _ = fmt.Fprintf(os.Stderr, "\rUsers pulled: %d / %d DPS: %.2f", cur, total, dps)
		}(i, user)
		swg.Add()
	}

	swg.Wait()
	dps := float64(cur) / time.Since(startTime).Seconds()
	_, _ = fmt.Fprintf(os.Stderr, "\rUsers pulled: %d / %d DPS: %.2f\n", cur, total, dps)

	kl.Log(users)

	return users, nil
}

func buildUserData(user *awsIAM.User, fullDetails bool) (*User, error) {
	kl := kiam.Extend("buildUserData")
	var hasPassword bool
	if user.PasswordLastUsed != nil {
		hasPassword = true
	}

	iu := &User{
		arn:              user.Arn,
		passwordLastUsed: user.PasswordLastUsed,
		createDate:       user.CreateDate,
		userName:         user.UserName,
		userID:           user.UserId,
		hasPassword:      hasPassword,
	}

	if fullDetails {
		err := iu.detectConsoleAccess()
		if err != nil {
			return nil, err
		}

		err = iu.GetPermissions()
		if err != nil {
			return nil, err
		}

		err = iu.GetAccessKeys()
		if err != nil {
			return nil, err
		}
	}

	kl.Log(iu)

	return iu, nil
}

func renderPermissions(permissions []*Permission) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{
		"Type",
		"Name",
		"ARN",
	})

	for _, perm := range permissions {
		t.AppendRow([]interface{}{
			perm.Type,
			perm.Name,
			perm.ARN,
		})
	}

	t.Render()
}

func promptForPermissionsAction(user *User) error {
	detach := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Detach permissions from %s?", user.UserName()),
	}
	err := survey.AskOne(prompt, &detach)
	if err != nil {
		return err
	}
	if detach {
		err = detachPermissions(user.Permissions(), user.UserName())
		if err != nil {
			return err
		}
	}
	return nil
}

func detachPermissions(permissions []*Permission, user string) error {
	kl := kiam.Extend("detachPermissions")
	var options []string
	for _, perm := range permissions {
		options = append(options, perm.ARN)
	}

	sort.Strings(options)

	var toDelete []string
	prompt := &survey.MultiSelect{
		Message: "Select Permissions to Detach from User:",
		Options: options,
	}
	_ = survey.AskOne(prompt, &toDelete)

	kl.Log(toDelete)

	if len(toDelete) == 0 {
		fmt.Println("Nothing to do.")
		return nil
	}

	for _, arn := range toDelete {
		perm := findPermission(permissions, arn)
		switch perm.Type {
		case "GROUP":
			err := removeUserFromGroup(user, perm.Name)
			if err != nil {
				return err
			}
		case "POLICY":
			err := detachPolicyFromUser(user, perm.ARN)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func detachPolicyFromUser(user string, arn string) error {
	kl := kiam.Extend("detachPolicyFromUser")
	sess, err := as.New()
	if err != nil {
		return err
	}
	client := awsIAM.New(sess)

	var results *awsIAM.DetachUserPolicyOutput
	results, err = client.DetachUserPolicy(&awsIAM.DetachUserPolicyInput{
		UserName:  &user,
		PolicyArn: &arn,
	})

	kl.Log(results)

	if err != nil {
		return err
	}

	fmt.Println(emoji.Sprintf(":check_mark_button: Detached %s from user %s", arn, user))

	return nil
}

func removeUserFromGroup(user string, groupName string) error {
	kl := kiam.Extend("removeUserFromGroup")
	sess, err := as.New()
	if err != nil {
		return err
	}
	client := awsIAM.New(sess)

	var results *awsIAM.RemoveUserFromGroupOutput
	results, err = client.RemoveUserFromGroup(&awsIAM.RemoveUserFromGroupInput{
		UserName:  &user,
		GroupName: &groupName,
	})

	kl.Log(results)

	if err != nil {
		return err
	}

	fmt.Println(emoji.Sprintf(":check_mark_button: Removed user %s from group %s", user, groupName))

	return nil
}

func viewUserDetails(user *User) {
	fmt.Println("User Details")
	renderUserDetails(user)

	fmt.Printf("\nPersmissions\n")
	renderPermissions(user.Permissions())

	fmt.Printf("\nAccess Keys\n")
	renderUserAccessKeys(user.AccessKeys())
}

func modifyUser(user *User) error {
	action := ""
	prompt := &survey.Select{
		Message: "What would you like to modify:",
		Options: []string{"ACCESS KEYS", "PERMISSIONS"},
	}
	err := survey.AskOne(prompt, &action)
	if err != nil {
		return err
	}

	switch action {
	case "ACCESS KEYS":
		err = promptForAccessKeyAction(user)
	case "PERMISSIONS":
		err = promptForPermissionsAction(user)
	}
	if err != nil {
		return err
	}

	return nil
}

func renderUserDetails(user *User) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{
		"User",
		"Status",
		"Age",
		"Console",
		"Last Login",
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AlignFooter: text.AlignCenter},
		{Number: 2, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 3, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 4, Align: text.AlignRight, AlignHeader: text.AlignRight},
		{Number: 5, Align: text.AlignRight, AlignHeader: text.AlignRight},
	})
	t.AppendRow([]interface{}{
		user.UserName(),
		user.FormattedCheckStatus(),
		user.CreatedDateDuration(),
		formattedYesNo(user.HasConsoleAccess()),
		user.FormattedLastLoginDateDuration(),
	})
	t.Render()
}

func promptForAccessKeyAction(user *User) error {
	takeAction := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("Adjust status of Access Keys for user %s?", user.UserName()),
	}
	err := survey.AskOne(prompt, &takeAction)
	if err != nil {
		return err
	}

	if takeAction {
		action := ""
		prompt := &survey.Select{
			Message: "Choose action to take:",
			Options: []string{"ACTIVATE", "DEACTIVATE", "DELETE"},
		}
		err = survey.AskOne(prompt, &action)
		if err != nil {
			return err
		}

		err = actionOnUserAccessKey(user.AccessKeys(), action)
		if err != nil {
			return err
		}
	}

	return nil
}

func actionOnUserAccessKey(keys []*AccessKey, action string) error {
	kl := kiam.Extend("actionOnUserAccessKey")
	kl.Printf("ACTION: %s", action)
	var options []string
	for _, key := range keys {
		options = append(options, *key.id)
	}

	sort.Strings(options)

	var toAction []string
	prompt := &survey.MultiSelect{
		Message: fmt.Sprintf("Select Access Keys to %s for User:", aurora.Red(action)),
		Options: options,
	}
	_ = survey.AskOne(prompt, &toAction)

	kl.Log(toAction)

	if len(toAction) == 0 {
		fmt.Println("Nothing to do.")
		return nil
	}

	fmt.Printf("Access Keys to %s:\n", aurora.Red(action))
	for _, id := range toAction {
		fmt.Printf("\t%s\n", id)
	}

	doIt := false
	confirm := &survey.Confirm{
		Message: fmt.Sprintf("Are you sure you want to %s these Access Keys?", aurora.Red(action)),
		Default: false,
	}
	_ = survey.AskOne(confirm, &doIt)

	if doIt {
		for _, id := range toAction {
			key := findAccessKey(keys, id)
			var err error
			switch action {
			case "DEACTIVATE":
				err = key.Deactivate()
			case "ACTIVATE":
				err = key.Activate()
			case "DELETE":
				err = key.Delete()
			default:
				return fmt.Errorf("action type not defined: %s", action)
			}
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Println("Exiting.")
	}
	fmt.Printf("%s complete\n", aurora.Red(action))

	return nil
}

func renderUserAccessKeys(keys []*AccessKey) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{
		"Key ID",
		"Status",
		"Age",
		"Last Used",
		"Service",
	})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 3, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 4, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 6, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
	})
	for _, key := range keys {
		t.AppendRow([]interface{}{
			aws.StringValue(key.id),
			formattedStatus(aws.StringValue(key.status)),
			dateDuration(aws.TimeValue(key.createdDate), 1),
			formattedDateDuration(dateDuration(aws.TimeValue(key.lastUsed.LastUsedDate), 2)),
			aws.StringValue(key.lastUsed.ServiceName),
		})
	}
	t.Render()
}
