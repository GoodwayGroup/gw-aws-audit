package iam

import (
	"fmt"
	"os"
	"sort"
	"sync/atomic"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	awsIAM "github.com/aws/aws-sdk-go/service/iam"
	"github.com/kyokomi/emoji/v2"
	"github.com/logrusorgru/aurora/v3"
	"github.com/remeh/sizedwaitgroup"
)

func getUser(user string, opts *buildUserDataOptions) (*User, error) {
	kl := kiam.Extend("getUser")
	client := session.GetIAMClient()

	var err error
	var result *awsIAM.GetUserOutput
	result, err = client.GetUser(&awsIAM.GetUserInput{
		UserName: &user,
	})
	if err != nil {
		return nil, err
	}
	kl.Log(result)

	var in *User
	in, err = buildUserData(result.User, &buildUserDataOptions{
		checkConsoleAccess: opts.checkConsoleAccess,
		getPermissions:     opts.getPermissions,
		getAccessKeys:      opts.getAccessKeys,
	})
	if err != nil {
		return nil, err
	}

	return in, nil
}

func getAllUsers(opts *buildUserDataOptions) ([]*User, error) {
	kl := kiam.Extend("getAllUsers")
	kmeta := kl.Extend("meta")
	client := session.GetIAMClient()

	var err error
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
	for i, user := range results.Users {
		go func(i int, user *awsIAM.User) {
			defer swg.Done()
			kmeta.Printf("[%d] Executing goroutine for user %s", i, aws.StringValue(user.UserName))
			var iu *User
			iu, err = buildUserData(user, &buildUserDataOptions{
				checkConsoleAccess: opts.checkConsoleAccess,
				getPermissions:     opts.getPermissions,
				getAccessKeys:      opts.getAccessKeys,
			})
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

type buildUserDataOptions struct {
	checkConsoleAccess bool
	getPermissions     bool
	getAccessKeys      bool
}

func buildUserData(user *awsIAM.User, opts *buildUserDataOptions) (*User, error) {
	kbud.Printf("user '%s' with opts: %# v", *user.Arn, opts)
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

	if opts.checkConsoleAccess {
		err := iu.detectConsoleAccess()
		if err != nil {
			return nil, err
		}
	}

	if opts.getPermissions {
		err := iu.GetPermissions()
		if err != nil {
			return nil, err
		}
	}

	if opts.getAccessKeys {
		err := iu.GetAccessKeys()
		if err != nil {
			return nil, err
		}
	}

	kbuduser.Log(iu)

	return iu, nil
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
	var opts []string
	for _, perm := range permissions {
		if perm.Type != "INLINE" {
			opts = append(opts, perm.ARN)
		}
	}

	sort.Strings(opts)

	var toDelete []string
	prompt := &survey.MultiSelect{
		Message: "Select Permissions to Detach from User:",
		Options: opts,
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
	client := session.GetIAMClient()

	var err error
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
	client := session.GetIAMClient()

	var err error
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
	var opts []string
	for _, key := range keys {
		opts = append(opts, *key.id)
	}

	sort.Strings(opts)

	var toAction []string
	prompt := &survey.MultiSelect{
		Message: fmt.Sprintf("Select Access Keys to %s:", aurora.Red(action)),
		Options: opts,
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
