package iam

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/logrusorgru/aurora/v3"
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

func viewUserDetails(user *User) {
	fmt.Println("User Details")
	renderUserDetails(user)

	fmt.Printf("\nPersmissions\n")
	renderPermissions(user.Permissions())

	fmt.Printf("\nAccess Keys\n")
	renderUserAccessKeys(user.AccessKeys())
}
