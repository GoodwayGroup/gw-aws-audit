package iam

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsIAM "github.com/aws/aws-sdk-go/service/iam"
	"github.com/clok/kemba"
	"github.com/hako/durafmt"
	"github.com/logrusorgru/aurora/v3"
)

type Permission struct {
	ARN  string
	Name string
	Type string
}

type User struct {
	arn      *string
	userName *string
	userID   *string
	// TODO: Deprecate hasPassword
	hasPassword      bool
	hasConsoleAccess bool
	passwordLastUsed *time.Time
	createDate       *time.Time
	accessKeys       []*AccessKey
	permissions      []*Permission
}

var (
	nilTime = time.Time{}
	k       = kemba.New("iam:types")
)

func (u User) ARN() string {
	return aws.StringValue(u.arn)
}

func (u User) UserName() string {
	return aws.StringValue(u.userName)
}

func (u User) ID() string {
	return aws.StringValue(u.userID)
}

func (u *User) detectConsoleAccess() error {
	kl := kiam.Extend("User:detectConsoleAccess")
	client := session.GetIAMClient()

	var err error
	var results *awsIAM.GetLoginProfileOutput
	results, err = client.GetLoginProfile(&awsIAM.GetLoginProfileInput{
		UserName: u.userName,
	})
	kl.Log(results)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case awsIAM.ErrCodeNoSuchEntityException:
				kl.Printf("%s does NOT have console access", u.UserName())
				u.hasConsoleAccess = false
				return nil
			default:
				kl.Printf("%s had the following error: %e", u.UserName(), aerr)
				u.hasConsoleAccess = false
				return aerr
			}
		}
		kl.Printf("%s had the following error: %e", u.UserName(), err)
		u.hasConsoleAccess = false
		return err
	}

	kl.Printf("%s has console access", u.UserName())
	u.hasConsoleAccess = true
	return nil
}

func (u User) HasConsoleAccess() bool {
	return u.hasConsoleAccess
}

func (u User) LastLogin() time.Time {
	return aws.TimeValue(u.passwordLastUsed)
}

func (u User) LastLoginDuration() string {
	if u.LastLogin() == nilTime {
		return ""
	}
	return durafmt.Parse(time.Since(u.LastLogin())).LimitToUnit("days").LimitFirstN(1).String()
}

func (u User) CreatedDate() time.Time {
	return aws.TimeValue(u.createDate)
}

func (u User) CreatedDateDuration() string {
	if u.CreatedDate() == nilTime {
		return ""
	}
	return durafmt.Parse(time.Since(u.CreatedDate())).LimitToUnit("days").LimitFirstN(1).String()
}

func (u *User) GetAccessKeys() error {
	kl := kiam.Extend("User:GetAccessKeys")
	client := session.GetIAMClient()

	var err error
	var results *awsIAM.ListAccessKeysOutput
	results, err = client.ListAccessKeys(&awsIAM.ListAccessKeysInput{
		UserName: u.userName,
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return err
	}

	kl.Printf("found %d access keys", len(results.AccessKeyMetadata))
	var keys []*AccessKey
	for _, key := range results.AccessKeyMetadata {

		// Look up the last used
		var data *awsIAM.GetAccessKeyLastUsedOutput
		data, err = client.GetAccessKeyLastUsed(&awsIAM.GetAccessKeyLastUsedInput{
			AccessKeyId: key.AccessKeyId,
		})
		if err != nil {
			return err
		}

		keys = append(keys, &AccessKey{
			id:          key.AccessKeyId,
			createdDate: key.CreateDate,
			status:      key.Status,
			lastUsed:    data.AccessKeyLastUsed,
			userName:    u.userName,
		})
	}

	kl.Log(keys)
	u.accessKeys = keys

	return nil
}

func (u User) AccessKeys() []*AccessKey {
	return u.accessKeys
}

func (u User) HasAccessKeys() bool {
	return u.accessKeys != nil
}

func (u User) AccessKeysCount() int {
	return len(u.accessKeys)
}

func (u *User) GetPermissions() error {
	kl := kiam.Extend("getUserPermissions")
	groups, err := u.getGroups()
	if err != nil {
		return err
	}

	var permissions []*Permission
	for _, group := range groups.Groups {
		permissions = append(permissions, &Permission{
			Type: "GROUP",
			ARN:  aws.StringValue(group.Arn),
			Name: aws.StringValue(group.GroupName),
		})
	}

	policies, err := u.getAttachedPolicies()
	if err != nil {
		return err
	}

	for _, policy := range policies.AttachedPolicies {
		permissions = append(permissions, &Permission{
			Type: "POLICY",
			ARN:  aws.StringValue(policy.PolicyArn),
			Name: aws.StringValue(policy.PolicyName),
		})
	}

	inline, err := u.getInlinePolicies()
	if err != nil {
		return err
	}

	for _, policy := range inline.PolicyNames {
		permissions = append(permissions, &Permission{
			Type: "INLINE",
			Name: aws.StringValue(policy),
		})
	}

	kl.Log(permissions)
	u.permissions = permissions

	return nil
}

func (u User) getGroups() (*awsIAM.ListGroupsForUserOutput, error) {
	kl := kiam.Extend("User:getGroups")
	client := session.GetIAMClient()

	var err error
	var results *awsIAM.ListGroupsForUserOutput
	results, err = client.ListGroupsForUser(&awsIAM.ListGroupsForUserInput{
		UserName: u.userName,
		MaxItems: aws.Int64(1000),
	})
	if err != nil {
		return nil, err
	}
	kl.Log(results)
	return results, nil
}

func (u User) getAttachedPolicies() (*awsIAM.ListAttachedUserPoliciesOutput, error) {
	kl := kiam.Extend("User:getAttachedPolicies")
	client := session.GetIAMClient()

	var err error
	var results *awsIAM.ListAttachedUserPoliciesOutput
	results, err = client.ListAttachedUserPolicies(&awsIAM.ListAttachedUserPoliciesInput{
		UserName: u.userName,
	})
	if err != nil {
		return nil, err
	}
	kl.Log(results)
	return results, nil
}

func (u User) getInlinePolicies() (*awsIAM.ListUserPoliciesOutput, error) {
	kl := kiam.Extend("User:getInlinePolicies")
	client := session.GetIAMClient()

	var err error
	var results *awsIAM.ListUserPoliciesOutput
	results, err = client.ListUserPolicies(&awsIAM.ListUserPoliciesInput{
		UserName: u.userName,
	})
	if err != nil {
		return nil, err
	}
	kl.Log(results)
	return results, nil
}

func (u User) HasPermissions() bool {
	return u.permissions != nil
}

func (u User) Permissions() []*Permission {
	return u.permissions
}

func (u User) PermissionsCount() int {
	return len(u.permissions)
}

func (u User) Groups() []*Permission {
	return permissionsByType(u.permissions, "GROUP")
}

func (u User) Policies() []*Permission {
	return permissionsByType(u.permissions, "POLICY")
}

func (u User) InlinePolicies() []*Permission {
	return permissionsByType(u.permissions, "INLINE")
}

// TODO: Make the Status more robust
func (u User) CheckStatus() string {
	switch {
	case !u.HasConsoleAccess() && !u.HasAccessKeys():
		return "pass"
	case u.HasConsoleAccess():
		return "fail"
	case u.HasAccessKeys():
		status := "pass"
		for _, ak := range u.accessKeys {
			if status == "pass" && aws.StringValue(ak.status) != "Inactive" {
				status = "warn"
			}
		}
		return status
	default:
		return "unknown"
	}
}

func (u User) FormattedCheckStatus() string {
	switch u.CheckStatus() {
	case "pass":
		return aurora.Green("PASS").String()
	case "warn":
		return aurora.Yellow("WARN").String()
	case "fail":
		return aurora.Red("FAIL").String()
	default:
		return aurora.Gray(8, "UNKNOWN").String()
	}
}

func (u User) FormattedLastLoginDateDuration() string {
	if u.HasConsoleAccess() {
		return u.LastLoginDuration()
	}
	return aurora.Gray(8, "NONE").String()
}
