package iam

import (
	"github.com/hako/durafmt"
	"github.com/logrusorgru/aurora/v3"
	"time"
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
		return aurora.Green(v).String()
	case v == 1:
		return aurora.Yellow(v).String()
	default:
		return aurora.Red(v).String()
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
