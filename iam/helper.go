package iam

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"time"

	"github.com/hako/durafmt"
	"github.com/logrusorgru/aurora/v3"
)

func formattedYesNo(v bool) string {
	switch {
	case v:
		return aurora.Red("YES").String()
	default:
		return aurora.Green("NO").String()
	}
}

func formattedKeyCount(v int) string {
	switch {
	case v == 0:
		return fmt.Sprintf("%s API Keys", aurora.Green(v).String())
	case v == 1:
		return fmt.Sprintf("%s API Key", aurora.Yellow(v).String())
	default:
		return fmt.Sprintf("%s API Keys", aurora.Red(v).String())
	}
}

func formattedStatus(v string) string {
	switch {
	case v == "Active":
		return aurora.Yellow(v).String()
	default:
		return aurora.Green(v).String()
	}
}

func dateDuration(d time.Time, n int) string {
	if d.Year() < 2000 {
		return ""
	}
	return durafmt.Parse(time.Since(d)).LimitToUnit("days").LimitFirstN(n).String()
}

func formattedDateDuration(v string) string {
	if v == "" {
		return aurora.Gray(8, "NONE").String()
	}
	return v
}

func findPermission(a []*Permission, arn string) *Permission {
	for _, n := range a {
		if arn == n.ARN {
			return n
		}
	}
	return nil
}

func permissionsByType(a []*Permission, t string) []*Permission {
	var permissions []*Permission
	for _, n := range a {
		if t == n.Type {
			permissions = append(permissions, n)
		}
	}
	return permissions
}

func findAccessKey(a []*AccessKey, id string) *AccessKey {
	for _, n := range a {
		if id == *n.id {
			return n
		}
	}
	return nil
}

func markToDeactivate(key *AccessKey, daysAgo int64) bool {
	if aws.StringValue(key.status) == "Inactive" {
		return false
	}

	// If the key has NEVER been used, check the create date
	if key.lastUsed.LastUsedDate == nil {
		created := int64(time.Since(aws.TimeValue(key.createdDate)).Hours())
		if created < daysAgo {
			return false
		}
	}

	// If key has not been used in the past N days
	delta := int64(time.Since(aws.TimeValue(key.lastUsed.LastUsedDate)).Hours())
	return delta >= daysAgo
}
