package iam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	as "github.com/clok/awssession"
	"github.com/clok/kemba"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/remeh/sizedwaitgroup"
	"os"
	"sort"
)

var (
	kiam = kemba.New("gw-aws-audit:iam")
)

func ListUsers() error {
	kl := kiam.Extend("list-users")
	data, err := getAllUsersWithAccessKeyData()
	if err != nil {
		return err
	}

	kl.Log(data)

	// sort user list
	sort.Slice(data, func(i, j int) bool {
		return data[i].UserName() < data[j].UserName()
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{
		"User",
		"Status",
		"Age",
		"Console",
		"Last Login",
		"Key Count",
		"Key ID",
		"Status",
		"Age",
		"Last Used",
		"Service",
	})

	for _, user := range data {
		t.AppendRow([]interface{}{
			user.UserName(),
			user.FormattedCheckStatus(),
			user.CreatedDateDuration(),
			formattedYesNo(user.HasConsoleAccess()),
			user.FormattedLastLoginDateDuration(),
			formattedKeyCount(user.AccessKeysCount()),
			"",
			"",
			"",
			"",
			"",
		})
		if len(user.accessKeys) > 0 {
			for _, key := range user.accessKeys {
				t.AppendRow([]interface{}{
					"",
					"",
					"",
					"",
					"",
					"",
					aws.StringValue(key.id),
					formattedStatus(aws.StringValue(key.status)),
					dateDuration(aws.TimeValue(key.createdDate), 1),
					formattedDateDuration(dateDuration(aws.TimeValue(key.lastUsed.LastUsedDate), 2)),
					aws.StringValue(key.lastUsed.ServiceName),
				})
			}
		}
		t.AppendSeparator()
	}
	t.Render()

	return nil
}

func getAllUsersWithAccessKeyData() ([]*iamUser, error) {
	kl := kiam.Extend("getAllUsersWithAccessKeyData")
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

	swg := sizedwaitgroup.New(20)

	users := make([]*iamUser, len(results.Users))
	kl.Printf("found %d users", len(results.Users))
	for i, user := range results.Users {
		go func(i int, user *iam.User) {
			defer swg.Done()
			kmeta.Printf("[%d] Executing goroutine for user %s", i, aws.StringValue(user.UserName))
			var iu *iamUser
			iu, err = buildUserData(user)
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

func buildUserData(user *iam.User) (*iamUser, error) {
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

	iu := &iamUser{
		arn:              user.Arn,
		hasPassword:      hasPassword,
		hasConsoleAccess: hasConsole,
		passwordLastUsed: user.PasswordLastUsed,
		createDate:       user.CreateDate,
		userName:         user.UserName,
		userID:           user.UserId,
		accessKeys:       keys,
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