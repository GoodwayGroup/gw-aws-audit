package sg

import (
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
	"strings"
)

func portInList(ports []string, port string) bool {
	if strings.HasPrefix(port, "sg-") {
		return true
	}
	return funk.ContainsString(ports, port)
}

// GeneratePortReport will generate a report of PORT exposure from Security Groups.
func GeneratePortReport(c *cli.Context) error {
	kl := ksg.Extend("GeneratePortReport")

	kl.Log(c.String("ports"))
	checkPorts := strings.Split(c.String("ports"), ",")
	fmt.Printf("----\n Ports: %s\n----\n\n", checkPorts)

	err := generateReport(c, portInList, checkPorts)
	if err != nil {
		return err
	}

	return nil
}
