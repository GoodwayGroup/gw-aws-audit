package sg

import (
	"fmt"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
	"strings"
)

func portNotInList(ports []string, port string) bool {
	if strings.HasPrefix(port, "sg-") {
		return true
	}
	return !funk.ContainsString(ports, port)
}

// GenerateCIDRReport will generate a report of CIDR block exposure from Security Groups.
func GenerateCIDRReport(c *cli.Context) error {
	kl := ksg.Extend("GenerateCIDRReport")

	kl.Log(c.String("ignore-ports"))
	ignoredPorts := strings.Split(c.String("ignore-ports"), ",")
	fmt.Printf("----\nIgnored ports: %s\n----\n\n", ignoredPorts)

	kl.Log(c.String("ignore-protocols"))
	var ignoredProtocols = make(map[string]bool)
	for _, v := range strings.Split(c.String("ignore-protocols"), ",") {
		ignoredProtocols[v] = true
	}

	err := generateReport(c, portNotInList, ignoredPorts, ignoredProtocols)
	if err != nil {
		return err
	}

	return nil
}
