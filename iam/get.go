package iam

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	as "github.com/clok/awssession"
	"github.com/clok/kemba"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/kyokomi/emoji/v2"
	"github.com/logrusorgru/aurora/v3"
	"github.com/remeh/sizedwaitgroup"
	"github.com/urfave/cli/v2"
	"os"
	"sort"
	"strings"
)

var (
	kiam = kemba.New("gw-aws-audit:iam")
)

func ListUsers(users []*User, showOnly string) error {
	kl := kiam.Extend("list-users")
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

func GetAllUsersWithAccessKeyData(includePermissions bool) ([]*User, error) {
	kl := kiam.Extend("GetAllUsersWithAccessKeyData")
	kmeta := kl.Extend("meta")
	sess, err := as.New()
	if err != nil {
		return nil, err
	}
	client := iam.New(sess)

	var results *iam.ListUsersOutput
	results, err = client.ListUsers(&iam.ListUsersInput{
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return nil, err
	}

	swg := sizedwaitgroup.New(10)

	users := make([]*User, len(results.Users))
	kl.Printf("found %d users", len(results.Users))
	for i, user := range results.Users {
		go func(i int, user *iam.User) {
			defer swg.Done()
			kmeta.Printf("[%d] Executing goroutine for user %s", i, aws.StringValue(user.UserName))
			var iu *User
			iu, err = buildUserData(user, includePermissions)
			if err != nil {
				// TODO: Handle panic
				panic(err)
			}

			users[i] = iu
		}(i, user)
		swg.Add()
	}

	swg.Wait()

	kl.Log(users)

	return users, nil
}

func buildUserData(user *iam.User, includePermissions bool) (*User, error) {
	hasConsole, err := hasConsoleAccess(user.UserName)
	if err != nil {
		return nil, err
	}

	var hasPassword bool
	if user.PasswordLastUsed != nil {
		hasPassword = true
	}

	keys, err := getAccessKeys(user.UserName)
	if err != nil {
		return nil, err
	}

	var permissions []*permission
	if includePermissions {
		permissions, err = GetUserPermissions(aws.StringValue(user.UserName))
		if err != nil {
			return nil, err
		}
	}

	// TODO: Get attached Policies and Groups
	iu := &User{
		arn:              user.Arn,
		hasPassword:      hasPassword,
		hasConsoleAccess: hasConsole,
		passwordLastUsed: user.PasswordLastUsed,
		createDate:       user.CreateDate,
		userName:         user.UserName,
		userID:           user.UserId,
		accessKeys:       keys,
		permissions:      permissions,
	}
	return iu, nil
}

func getAccessKeys(user *string) ([]*accessKey, error) {
	kl := kiam.Extend("get-access-keys")
	sess, err := as.New()
	if err != nil {
		return nil, err
	}
	client := iam.New(sess)

	var results *iam.ListAccessKeysOutput
	results, err = client.ListAccessKeys(&iam.ListAccessKeysInput{
		UserName: user,
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return nil, err
	}

	kl.Printf("found %d access keys", len(results.AccessKeyMetadata))
	var keys []*accessKey
	for _, key := range results.AccessKeyMetadata {

		// Look up the last used
		var data *iam.GetAccessKeyLastUsedOutput
		data, err = client.GetAccessKeyLastUsed(&iam.GetAccessKeyLastUsedInput{
			AccessKeyId: key.AccessKeyId,
		})
		if err != nil {
			return nil, err
		}

		keys = append(keys, &accessKey{
			id:          key.AccessKeyId,
			createdDate: key.CreateDate,
			status:      key.Status,
			lastUsed:    data.AccessKeyLastUsed,
		})
	}

	kl.Log(keys)

	return keys, nil
}

func ListPermissions(permissions []*permission) error {
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

	return nil
}

func DetachPermissions(permissions []*permission, user string) error {
	kl := kiam.Extend("detach-permissions")
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
	kl := kiam.Extend("detach-policy-from-user")
	sess, err := as.New()
	if err != nil {
		return err
	}
	client := iam.New(sess)

	var results *iam.DetachUserPolicyOutput
	results, err = client.DetachUserPolicy(&iam.DetachUserPolicyInput{
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
	kl := kiam.Extend("remove-user-from-group")
	sess, err := as.New()
	if err != nil {
		return err
	}
	client := iam.New(sess)

	var results *iam.RemoveUserFromGroupOutput
	results, err = client.RemoveUserFromGroup(&iam.RemoveUserFromGroupInput{
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

func GetUserPermissions(user string) ([]*permission, error) {
	kl := kiam.Extend("get-user-permissions")
	groups, err := getGroupsForUser(&user)
	if err != nil {
		return nil, err
	}

	var permissions []*permission
	for _, group := range groups.Groups {
		permissions = append(permissions, &permission{
			Type: "GROUP",
			ARN:  aws.StringValue(group.Arn),
			Name: aws.StringValue(group.GroupName),
		})
	}

	policies, err := getUserPolicies(&user)
	if err != nil {
		return nil, err
	}

	for _, policy := range policies.AttachedPolicies {
		permissions = append(permissions, &permission{
			Type: "POLICY",
			ARN:  aws.StringValue(policy.PolicyArn),
			Name: aws.StringValue(policy.PolicyName),
		})
	}

	kl.Log(permissions)

	return permissions, nil
}

func getGroupsForUser(user *string) (*iam.ListGroupsForUserOutput, error) {
	kl := kiam.Extend("get-groups-for-user")
	sess, err := as.New()
	if err != nil {
		return nil, err
	}
	client := iam.New(sess)

	var results *iam.ListGroupsForUserOutput
	results, err = client.ListGroupsForUser(&iam.ListGroupsForUserInput{
		UserName: user,
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return nil, err
	}
	kl.Log(results)
	return results, nil
}

func getUserPolicies(user *string) (*iam.ListAttachedUserPoliciesOutput, error) {
	kl := kiam.Extend("get-user")
	sess, err := as.New()
	if err != nil {
		return nil, err
	}
	client := iam.New(sess)

	var results *iam.ListAttachedUserPoliciesOutput
	results, err = client.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
		UserName: user,
	})
	if err != nil {
		return nil, err
	}
	kl.Log(results)
	return results, nil
}

func hasConsoleAccess(userName *string) (bool, error) {
	kl := kiam.Extend("has-console-access")
	sess, err := as.New()
	if err != nil {
		return false, err
	}
	client := iam.New(sess)

	var results *iam.GetLoginProfileOutput
	results, err = client.GetLoginProfile(&iam.GetLoginProfileInput{
		UserName: userName,
	})
	kl.Log(results)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeNoSuchEntityException:
				kl.Printf("%s does NOT have console access", aws.StringValue(userName))
				return false, nil
			default:
				kl.Printf("%s had the following error: %e", aws.StringValue(userName), aerr)
				return false, aerr
			}
		}
		kl.Printf("%s had the following error: %e", aws.StringValue(userName), err)
		return false, err
	}

	kl.Printf("%s has console access", aws.StringValue(userName))
	return true, nil
}

func ViewUserDetails(users []*User, userName string) error {
	kl := kiam.Extend("view-user-details")
	if userName == "" {
		// prompt for user selection
		var options []string
		for _, user := range users {
			options = append(options, user.UserName())
		}

		sort.Strings(options)

		prompt := &survey.Select{
			Message: "Pick a User:",
			Options: options,
		}
		err := survey.AskOne(prompt, &userName)
		if err != nil {
			return cli.Exit(err, 2)
		}

		kl.Log(userName)

		var user *User
		for _, u := range users {
			if u.UserName() == userName {
				user = u
			}
		}
		kl.Log(user)
		if user == nil {
			return fmt.Errorf("user not found: %s", userName)
		}

		err = ListPermissions(user.Permissions())
		if err != nil {
			return cli.Exit(err, 2)
		}

		detach := false
		toDetach := &survey.Confirm{
			Message: fmt.Sprintf("Detach permissions from %s?", userName),
		}
		err = survey.AskOne(toDetach, &detach)
		if err != nil {
			return cli.Exit(err, 2)
		}
		if detach {
			err = DetachPermissions(user.Permissions(), user.UserName())
			if err != nil {
				return cli.Exit(err, 2)
			}
		}
	}
	return nil
}
