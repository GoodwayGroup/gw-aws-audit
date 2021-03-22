package iam

import (
	"github.com/aws/aws-sdk-go/aws"
	awsIAM "github.com/aws/aws-sdk-go/service/iam"
	"github.com/hako/durafmt"
	"github.com/logrusorgru/aurora/v3"
	"time"
)

type User struct {
	arn      *string
	userName *string
	userID   *string
	// TODO: Deprecate hasPassword
	hasPassword      bool
	hasConsoleAccess bool
	passwordLastUsed *time.Time
	createDate       *time.Time
	accessKeys       []*accessKey
	permissions      []*permission
}

var (
	nilTime = time.Time{}
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

func (u User) HasAccessKeys() bool {
	return u.accessKeys != nil
}

func (u User) AccessKeysCount() int {
	return len(u.accessKeys)
}

func (u User) HasPermissions() bool {
	return u.permissions != nil
}

func (u User) Permissions() []*permission {
	return u.permissions
}

func (u User) PermissionsCount() int {
	return len(u.permissions)
}

func (u User) Groups() []*permission {
	return permissionsByType(u.permissions, "GROUP")
}

func (u User) Policies() []*permission {
	return permissionsByType(u.permissions, "GROUP")
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

type accessKey struct {
	id          *string
	createdDate *time.Time
	status      *string
	lastUsed    *awsIAM.AccessKeyLastUsed
}

type permission struct {
	ARN  string
	Name string
	Type string
}
