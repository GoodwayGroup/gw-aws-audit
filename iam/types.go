package iam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/hako/durafmt"
	"github.com/logrusorgru/aurora/v3"
	"time"
)

type iamUser struct {
	arn      *string
	userName *string
	userID   *string
	// TODO: Deprecate hasPassword
	hasPassword      bool
	hasConsoleAccess bool
	passwordLastUsed *time.Time
	createDate       *time.Time
	accessKeys       []*accessKey
}

var (
	nilTime = time.Time{}
)

func (u iamUser) ARN() string {
	return aws.StringValue(u.arn)
}

func (u iamUser) UserName() string {
	return aws.StringValue(u.userName)
}

func (u iamUser) ID() string {
	return aws.StringValue(u.userID)
}

func (u iamUser) HasConsoleAccess() bool {
	return u.hasConsoleAccess
}

func (u iamUser) LastLogin() time.Time {
	return aws.TimeValue(u.passwordLastUsed)
}

func (u iamUser) LastLoginDuration() string {
	if u.LastLogin() == nilTime {
		return ""
	}
	return durafmt.Parse(time.Since(u.LastLogin())).LimitToUnit("days").LimitFirstN(1).String()
}

func (u iamUser) CreatedDate() time.Time {
	return aws.TimeValue(u.createDate)
}

func (u iamUser) CreatedDateDuration() string {
	if u.CreatedDate() == nilTime {
		return ""
	}
	return durafmt.Parse(time.Since(u.CreatedDate())).LimitToUnit("days").LimitFirstN(1).String()
}

func (u iamUser) HasAccessKeys() bool {
	return u.accessKeys != nil
}

func (u iamUser) AccessKeysCount() int {
	return len(u.accessKeys)
}

// TODO: Make the Status more robust
func (u iamUser) CheckStatus() string {
	switch {
	case !u.HasConsoleAccess() && !u.HasAccessKeys():
		return "pass"
	case u.HasConsoleAccess():
		return "fail"
	case u.HasAccessKeys():
		return "warn"
	default:
		return "unknown"
	}
}

func (u iamUser) FormattedCheckStatus() string {
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

func (u iamUser) FormattedLastLoginDateDuration() string {
	if u.HasConsoleAccess() {
		return u.LastLoginDuration()
	}
	return aurora.Gray(8, "NONE").String()
}

type accessKey struct {
	id          *string
	createdDate *time.Time
	status      *string
	lastUsed    *iam.AccessKeyLastUsed
}
