package iam

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kyokomi/emoji/v2"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
)

var (
	ActionUserKeys = &cli.Command{
		Name:  "keys",
		Usage: "view Access Keys associated with an IAM User",
		UsageText: `
Produces a table of Access Keys that are associated with an IAM User.

Interactive mode allows for you to Activate, Deactivate and Delete Access Keys.
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "user name to look for",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "interactive mode that allows status changes of keys",
			},
		},
		Action: func(c *cli.Context) error {
			un := c.String("user")
			user := &User{
				userName: &un,
			}
			err := user.GetAccessKeys()
			if err != nil {
				return cli.Exit(err, 2)
			}

			renderUserAccessKeys(user.AccessKeys())

			if c.Bool("interactive") {
				err = promptForAccessKeyAction(user)
				if err != nil {
					return cli.Exit(err, 2)
				}
			}

			return nil
		},
	}
	ActionUserPermissions = &cli.Command{
		Name:    "permissions",
		Aliases: []string{"p"},
		Usage:   "view permissions that are associated with an IAM User",
		UsageText: `
Produces a table of Groups and Policies that are attached to an IAM User.

Interactive mode allows for you to detach a permission from an IAM User.
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "user",
				Aliases:  []string{"u"},
				Usage:    "user name to look for",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "interactive mode that allows for removal of permissions",
			},
		},
		Action: func(c *cli.Context) error {
			un := c.String("user")
			user := &User{
				userName: &un,
			}
			err := user.GetPermissions()
			if err != nil {
				return cli.Exit(err, 2)
			}

			renderPermissions(user.Permissions())

			if c.Bool("interactive") {
				err = promptForPermissionsAction(user)
				if err != nil {
					return cli.Exit(err, 2)
				}
			}
			return nil
		},
	}
	ActionDeprecatedUserReport = &cli.Command{
		Name:      "user-report",
		Usage:     "DEPRECATED: Please use the `iam user report` command",
		UsageText: "DEPRECATED: Please use the `iam user report` command",
		Action: func(c *cli.Context) error {
			return cli.Exit(fmt.Errorf("deprecated: please user the `iam report` command"), 2)
		},
		Hidden: true,
	}
	ActionUserReport = &cli.Command{
		Name:  "report",
		Usage: "generates report of IAM Users and Access Key Usage",
		UsageText: `
This action will generate a report for all Users within an AWS account with the details
specific user authentication methods.

Interactive mode will allow you to search for an IAM User and take actions once an IAM User is
selected.

USER [string]:
  - The user name

STATUS [enum]:
  - PASS: When a does NOT have Console Access and has NO Access Keys or only INACTIVE Access Keys
  - FAIL: When an IAM User has Console Access
  - WARN: When an IAM User does NOT have Console Access, but does have at least 1 ACTIVE Access Key
  - UNKNOWN: Catch all for cases not handled.

AGE [duration]:
  - Time since User was created

CONSOLE [bool]:
  - Does the User have Console Access? YES/NO

LAST LOGIN [duration]:
  - Time since User was created
  - NONE if the User does not have Console Access or if the User has NEVER logged in.

PERMISSIONS [struct]:
  - G: n -> Groups that the User belongs to
  - P: n -> Policies that are attached to the User

ACCESS KEY DETAILS [sub table]:
  - Primary header row is the number of Access Keys associated with the User

  KEY ID [string]:
    - The AWS_ACCESS_KEY_ID

  STATUS [enum]:
    - Active/Inactive

  LAST USED [duration]:
    - Time since the Access Key was last used.

  SERVICE [string]:
    - The last AWS Service that the Access Key was used to access at the "LAST USED" time.
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "show-only",
				Usage: "filter results to show only pass, warn or fail",
			},
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "after generating the report, prompt for digging into a user",
			},
		},
		Action: func(c *cli.Context) error {
			kl := kiam.Extend("ActionUserReport")
			showOnly := ""
			if c.String("show-only") != "" {
				showOnly = strings.ToLower(c.String("show-only"))
			}
			allowedFilters := []string{"", "pass", "warn", "fail"}
			if !funk.ContainsString(allowedFilters, showOnly) {
				return cli.Exit(fmt.Sprintf("Invalid value for show-only. Must be one of: %v", allowedFilters), 3)
			}

			users, err := getAllUsers(&buildUserDataOptions{
				checkConsoleAccess: true,
				getPermissions:     true,
				getAccessKeys:      true,
			})
			if err != nil {
				return cli.Exit(err, 2)
			}

			err = renderUsersReport(users, showOnly)
			if err != nil {
				return cli.Exit(err, 2)
			}

			if c.Bool("interactive") {
				// prompt for user selection
				var options []string
				for _, user := range users {
					options = append(options, user.UserName())
				}

				sort.Strings(options)

				var passedUser string
				prompt := &survey.Select{
					Message: "Pick an IAM User:",
					Options: options,
				}
				err := survey.AskOne(prompt, &passedUser)
				if err != nil {
					return cli.Exit(err, 2)
				}
				kl.Log(passedUser)

				var user *User
				for _, u := range users {
					if u.UserName() == passedUser {
						user = u
					}
				}
				kl.Log(user)

				if user == nil {
					return fmt.Errorf("user not found: %s", passedUser)
				}

				viewUserDetails(user)

				err = modifyUser(user)
				if err != nil {
					return cli.Exit(err, 2)
				}
			}
			return nil
		},
	}
	ActionUserModify = &cli.Command{
		Name:  "modify",
		Usage: "modify an IAM User within AWS",
		UsageText: `
This action allows you to take actions to modify a user's Permissions (Groups and Policies)
and the state of their Access Keys (Active, Inactive, Delete).
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "user name to look for",
			},
			&cli.StringFlag{
				Name:  "show-only",
				Usage: "filter results to show only pass, warn or fail",
			},
		},
		Action: func(c *cli.Context) error {
			kl := kiam.Extend("ActionUserModify")
			showOnly := ""
			if c.String("show-only") != "" {
				showOnly = strings.ToLower(c.String("show-only"))
			}
			allowedFilters := []string{"", "pass", "warn", "fail"}
			if !funk.ContainsString(allowedFilters, showOnly) {
				return cli.Exit(fmt.Sprintf("Invalid value for show-only. Must be one of: %v", allowedFilters), 3)
			}

			var err error
			var user *User
			passedUser := c.String("user")
			if passedUser == "" {
				users, err := getAllUsers(&buildUserDataOptions{})
				if err != nil {
					return cli.Exit(err, 2)
				}

				// prompt for user selection
				var options []string
				for _, user := range users {
					options = append(options, user.UserName())
				}

				sort.Strings(options)

				prompt := &survey.Select{
					Message: "Pick an IAM User:",
					Options: options,
				}
				err = survey.AskOne(prompt, &passedUser)
				if err != nil {
					return cli.Exit(err, 2)
				}
				kl.Log(passedUser)
			}

			user, err = getUser(passedUser, &buildUserDataOptions{})
			if err != nil {
				return cli.Exit(err, 2)
			}

			kl.Log(user)
			if user == nil {
				return fmt.Errorf("user not found: %s", passedUser)
			}

			viewUserDetails(user)

			err = modifyUser(user)
			if err != nil {
				return cli.Exit(err, 2)
			}

			return nil
		},
	}
	ActionKeysDeactivate = &cli.Command{
		Name:  "deactivate",
		Usage: "bulk deactivate Access Keys",
		UsageText: `
This action will check ALL Access Keys to determine if they meet the criteria
to be marked as INACTIVE within IAM.

Current rules are:

- If a keys HAS been used, the last usage was not within the last n(threshold) days
- If a key has NEVER been used, that the key was created at least n(threshold) days ago
`,
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:  "threshold",
				Usage: "number of days to pass as check for qualification",
				Value: 180,
			},
		},
		Action: func(c *cli.Context) error {
			users, err := getAllUsers(&buildUserDataOptions{
				checkConsoleAccess: false,
				getPermissions:     false,
				getAccessKeys:      true,
			})
			if err != nil {
				return cli.Exit(err, 2)
			}

			// sort user list
			sort.Slice(users, func(i, j int) bool {
				return strings.ToLower(users[i].UserName()) < strings.ToLower(users[j].UserName())
			})

			var toAction []*AccessKey
			thresholdCheck := c.Int64("threshold") * 24
			for _, user := range users {
				for _, key := range user.accessKeys {
					if markToDeactivate(key, thresholdCheck) {
						toAction = append(toAction, key)
					}
				}
			}

			if len(toAction) == 0 {
				fmt.Println(emoji.Sprint(":check_mark_button: No Access Keys qualify."))
				return nil
			}

			fmt.Println(emoji.Sprintf(":warning: Found %d Access Keys that qualify for deactivation :warning:", len(toAction)))

			takeAction := false
			prompt := &survey.Confirm{
				Message: "View keys that qualify?",
			}
			err = survey.AskOne(prompt, &takeAction)
			if err != nil {
				return err
			}

			if takeAction {
				renderUserAccessKeys(toAction)
			}

			takeAction = false
			p := &survey.Confirm{
				Message: "Deactivate Keys?",
			}
			err = survey.AskOne(p, &takeAction)
			if err != nil {
				return err
			}

			if takeAction {
				err = actionOnUserAccessKey(toAction, "DEACTIVATE")
				if err != nil {
					return err
				}
			}

			return nil
		},
	}
	ActionKeysRecent = &cli.Command{
		Name:      "recent",
		Usage:     "list Access Keys that have been recently used",
		UsageText: "This action will check ALL Access Keys to determine if they have been used within the threshold time.",
		Flags: []cli.Flag{
			&cli.Int64Flag{
				Name:    "threshold",
				Aliases: []string{"t"},
				Usage:   "number of Units to check for qualification",
				Value:   7,
			},
			&cli.StringFlag{
				Name:    "units",
				Aliases: []string{"u"},
				Usage:   "hours, days, weeks, months",
				Value:   "days",
			},
		},
		Action: func(c *cli.Context) error {
			units := strings.ToLower(c.String("units"))
			allowedFilters := []string{"hours", "days", "weeks", "months"}
			if !funk.ContainsString(allowedFilters, units) {
				return cli.Exit(fmt.Sprintf("Invalid value for units. Must be one of: %v", allowedFilters), 3)
			}

			// convert units to hours
			threshold := c.Int64("threshold")
			var check int64
			switch units {
			case "hours":
				check = threshold
			case "days":
				check = threshold * 24
			case "weeks":
				check = threshold * 7 * 24
			case "months":
				check = threshold * 30 * 24
			}

			users, err := getAllUsers(&buildUserDataOptions{
				checkConsoleAccess: false,
				getPermissions:     false,
				getAccessKeys:      true,
			})
			if err != nil {
				return cli.Exit(err, 2)
			}

			var toAction []*AccessKey
			for _, user := range users {
				for _, key := range user.accessKeys {
					if markAsRecentlyUsed(key, check) {
						toAction = append(toAction, key)
					}
				}
			}

			if len(toAction) == 0 {
				fmt.Println(emoji.Sprint(":check_mark_button: No Access Keys qualify."))
				return nil
			}
			fmt.Println(emoji.Sprintf(":peacock: Found %d Access Keys used in the last %d %s :whale:", len(toAction), threshold, units))

			// sort Access Keys list from most recent
			sort.Slice(toAction, func(i, j int) bool {
				return toAction[i].LastUsedDate().After(toAction[j].LastUsedDate())
			})

			renderUserAccessKeys(toAction)

			return nil
		},
	}
	ActionKeysUnused = &cli.Command{
		Name:  "unused",
		Usage: "list Access Keys that have NEVER been used",
		Action: func(c *cli.Context) error {
			users, err := getAllUsers(&buildUserDataOptions{
				checkConsoleAccess: false,
				getPermissions:     false,
				getAccessKeys:      true,
			})
			if err != nil {
				return cli.Exit(err, 2)
			}

			var toAction []*AccessKey
			for _, user := range users {
				for _, key := range user.accessKeys {
					if markAsNeverUsed(key) {
						toAction = append(toAction, key)
					}
				}
			}

			if len(toAction) == 0 {
				fmt.Println(emoji.Sprint(":check_mark_button: No Access Keys qualify."))
				return nil
			}
			fmt.Println(emoji.Sprintf(":doughnut: Found %d Access Keys that have NEVER been used :coffee:", len(toAction)))

			// sort Access Keys list from most recent
			sort.Slice(toAction, func(i, j int) bool {
				return toAction[i].CreatedDate().Before(toAction[j].CreatedDate())
			})

			renderUserAccessKeys(toAction)

			return nil
		},
	}
)
