package iam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	as "github.com/clok/awssession"
	"github.com/clok/kemba"
	"time"
)

var (
	kiam = kemba.New("gw-aws-audit:iam")
)

type iamUser struct {
	arn *string
	hasPassword bool
	passwordLastUsed *time.Time
	createDate *time.Time
	userName *string
	userId *string
}

func ListUsers() error {
	kl := kiam.Extend("list-users")
	sess, err := as.New()
	if err != nil {
		return err
	}
	client := iam.New(sess)

	var results *iam.ListUsersOutput
	results, err = client.ListUsers(&iam.ListUsersInput{
		MaxItems:   aws.Int64(1000),
	})
	if err != nil {
		return err
	}

	users := []*iamUser{}
	kl.Printf("found %d users", len(results.Users))
	for _, user := range results.Users {
		var hasPassword bool
		if user.PasswordLastUsed != nil {
			hasPassword = true
		}
		users = append(users, &iamUser{
			arn:              user.Arn,
			hasPassword:      hasPassword,
			passwordLastUsed: user.PasswordLastUsed,
			createDate:       user.CreateDate,
			userName:         user.UserName,
			userId:           user.UserId,
		})
	}

	kl.Log(users)

	return nil
}