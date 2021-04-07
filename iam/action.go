package iam

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kyokomi/emoji/v2"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
)

// determineCheck converts threshold units to hours
func determineCheck(threshold int64, units string) (check int64) {
	switch units {
	case "hours":
		check = threshold
	case "days":
		check = threshold * 24
	case "weeks":
		check = threshold * 7 * 24
	case "months":
		check = threshold * 30 * 24
	default:
		// days
		check = threshold * 24
	}
	return
}

// keyActions determines which action to run based on CLI parameters
func keyActions(c *cli.Context) error {
	units := strings.ToLower(c.String("units"))
	allowedFilters := []string{"", "hours", "days", "weeks", "months"}
	if !funk.ContainsString(allowedFilters, units) {
		return cli.Exit(fmt.Sprintf("Invalid value for units. Must be one of: %v", allowedFilters), 3)
	}

	threshold := c.Int64("threshold")
	check := determineCheck(threshold, units)

	users, err := getAllUsers(&buildUserDataOptions{
		checkConsoleAccess: false,
		getPermissions:     false,
		getAccessKeys:      true,
	})
	if err != nil {
		return cli.Exit(err, 2)
	}

	var toAction []*AccessKey
	command := c.Command.Name
	for _, user := range users {
		for _, key := range user.accessKeys {
			switch command {
			case "deactivate":
				if markToDeactivate(key, check) {
					toAction = append(toAction, key)
				}
			case "delete":
				if markToDelete(key, check) {
					toAction = append(toAction, key)
				}
			case "unused":
				if markAsNeverUsed(key) {
					toAction = append(toAction, key)
				}
			case "recent":
				if markAsRecentlyUsed(key, check) {
					toAction = append(toAction, key)
				}
			}
		}
	}

	if len(toAction) == 0 {
		fmt.Println(emoji.Sprint(":check_mark_button: No Access Keys qualify."))
		return nil
	}

	switch command {
	case "deactivate":
		fmt.Println(emoji.Sprintf(":warning: Found %d Access Keys that qualify for deactivation :warning:", len(toAction)))
		renderUserAccessKeys(toAction, "name")
	case "delete":
		fmt.Println(emoji.Sprintf(":warning: Found %d Access Keys that qualify for deletion :warning:", len(toAction)))
		renderUserAccessKeys(toAction, "name")
	case "unused":
		fmt.Println(emoji.Sprintf(":doughnut: Found %d Access Keys that have NEVER been used :coffee:", len(toAction)))
		renderUserAccessKeys(toAction, "created")
	case "recent":
		fmt.Println(emoji.Sprintf(":peacock: Found %d Access Keys used in the last %d %s :whale:", len(toAction), threshold, units))
		renderUserAccessKeys(toAction, "activity")
	}

	switch command {
	case "deactivate":
		takeAction := false
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
	case "delete":
		takeAction := false
		p := &survey.Confirm{
			Message: "Delete Keys?",
		}
		err = survey.AskOne(p, &takeAction)
		if err != nil {
			return err
		}

		if takeAction {
			err = actionOnUserAccessKey(toAction, "DELETE")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
