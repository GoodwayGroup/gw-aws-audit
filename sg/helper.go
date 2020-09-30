package sg

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/urfave/cli/v2"
	"net"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func isCiderIn(input *net.IPNet, cidrs []*net.IPNet) bool {
	for _, cidr := range cidrs {
		if cidr.Contains(input.IP) {
			return true
		}
	}
	return false
}

func getFuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func generateReport(c *cli.Context, checkFxn func(a []string, b string) bool, ports []string, ignoredProtocols map[string]bool) error {
	kl := ksg.Extend("generateReport")
	kl.Printf("generating report with settings - checkFxn: %s ports: %# v", getFuncName(checkFxn), ports)
	sgs, err := getAnnotatedSecurityGroups()
	if err != nil {
		return err
	}

	var securityGroups []*SecurityGroup
	for _, sg := range sgs {
		if sg.attached != nil || c.Bool("all") {
			securityGroups = append(securityGroups, sg)
		}
	}

	groupedCIDRs, err := generateIPBlockRules(c)
	if err != nil {
		return err
	}

	mappedSGs, err := processSecurityGroups(securityGroups, groupedCIDRs, checkFxn, ports, ignoredProtocols)
	if err != nil {
		return err
	}

	if len(mappedSGs.amazon) > 0 {
		fmt.Println("AMAZON: CIDR rules that are mapping to external Amazon IPs")
		printAmazonTable(mappedSGs.amazon)
		fmt.Println("")
	}

	if len(mappedSGs.unknown) > 0 {
		fmt.Println("UNKNOWN: It is not known if these are allowed or not based on filters provided.")
		printTable(mappedSGs.unknown)
		fmt.Println("")
	}

	if len(mappedSGs.sg) > 0 {
		fmt.Println("SECURITY GROUP: CIDR rules that are mapping Security Groups to other Security Groups")
		printTable(mappedSGs.sg)
		fmt.Println("")
	}

	if len(mappedSGs.warning) > 0 {
		fmt.Println("WARNINGS: CIDR rules flagged as warnings")
		printTable(mappedSGs.warning)
		fmt.Println("")
	}

	if len(mappedSGs.alert) > 0 {
		fmt.Println("ALERTS: CIDR rules flagged as alerts. Probable they are Public IPs")
		printTable(mappedSGs.alert)
		fmt.Println("")
	}

	if len(mappedSGs.wideOpen) > 0 {
		fmt.Println("WIDE OPEN: CIDR rules that are wide open. Should be verified that this is intended.")
		printTable(mappedSGs.wideOpen)
	}
	return nil
}

func parseToken(token string) (port string, protocol string, sgIDs string) {
	parts := strings.Split(token, "::")
	return parts[0], parts[1], parts[2]
}

func processSecurityGroups(securityGroups []*SecurityGroup, groupedCIDRs *groupedIPBlockRules, checkFxn func(a []string, b string) bool, ports []string, ignoredProtocols map[string]bool) (*groupedSecurityGroups, error) {
	kl := ksg.Extend("processSecurityGroups")
	mappedSGs := newGroupedSecurityGroups()

	awsIPs, err := getAWSIPRanges()
	if err != nil {
		return nil, err
	}

	for _, sec := range securityGroups {
		for token, rule := range sec.rules {
			port, proto, sgIDs := parseToken(token)

			if _, ok := ignoredProtocols[proto]; ok {
				kl.Printf("Skipping protocol: %s", proto)
				continue
			}

			switch {
			case rule == nil:
				for _, sgID := range strings.Split(sgIDs, ",") {
					mappedSGs.addToSG(sec, &portToIP{
						port:  port,
						ip:    sgID,
						proto: proto,
					})
				}
			case checkFxn(ports, port):
				for _, ip := range rule {
					_, ipv4Net, err := net.ParseCIDR(aws.StringValue(ip.CidrIp))
					if err != nil {
						return mappedSGs, err
					}

					portToIPValue := &portToIP{
						port:  port,
						proto: proto,
						ip:    ipv4Net.String(),
					}

					if ipv4Net.IP.String() == "0.0.0.0" {
						kl.Printf("%s\t%s\t%s", "FULL", ipv4Net.String(), port)
						mappedSGs.addToWideOpen(sec, portToIPValue)
						continue
					}

					if prefix, ok := findAWSCidr(ipv4Net, awsIPs); ok {
						kl.Printf("%s\t%s\t%s\t%s\t%s\n", ip.String(), port, prefix.IPPrefix, prefix.Service, prefix.Region)
						portToIPValue.prefix = prefix
						mappedSGs.addToAmazon(sec, portToIPValue)
						continue
					}

					switch {
					case isCiderIn(ipv4Net, groupedCIDRs.approved):
						kl.Printf("%s\t%s\t%s", "APPROVED", ipv4Net.String(), port)
					case isCiderIn(ipv4Net, groupedCIDRs.warning):
						kl.Printf("%s\t%s\t%s", "WARN", ipv4Net.String(), port)
						mappedSGs.addToWarning(sec, portToIPValue)
					case isCiderIn(ipv4Net, groupedCIDRs.alert):
						kl.Printf("%s\t%s\t%s", "ALERT", ipv4Net.String(), port)
						mappedSGs.addToAlert(sec, portToIPValue)
					default:
						kl.Printf("%s\t%s\t%s", "UNKNOWN", ipv4Net.String(), port)
						mappedSGs.addToUnknown(sec, portToIPValue)
					}
				}
			default:
				kl.Printf("skipping port %s", port)
			}
		}
	}
	return mappedSGs, nil
}

func generateIPBlockRules(c *cli.Context) (*groupedIPBlockRules, error) {
	kl := ksg.Extend("generateIPBlockRules")
	groupedCIDRs := newGroupedIPBlockRules()

	// Get approved CIDR blocks
	for _, cidr := range strings.Split(c.String("approved"), ",") {
		_, ipv4Net, err := net.ParseCIDR(cidr)
		if err != nil {
			return groupedCIDRs, err
		}
		groupedCIDRs.addToApproved(ipv4Net)
	}
	kl.Log("approved cidrs", groupedCIDRs.approved)

	// Get warn CIDR blocks
	if c.String("warn") != "" {
		for _, cidr := range strings.Split(c.String("warn"), ",") {
			_, ipv4Net, err := net.ParseCIDR(cidr)
			if err != nil {
				return groupedCIDRs, err
			}
			groupedCIDRs.addToWarning(ipv4Net)
		}
	}
	kl.Log("warn cidrs", groupedCIDRs.warning)

	// Get approved CIDR blocks
	for _, cidr := range strings.Split(c.String("alert"), ",") {
		_, ipv4Net, err := net.ParseCIDR(cidr)
		if err != nil {
			return groupedCIDRs, err
		}
		groupedCIDRs.addToAlert(ipv4Net)
	}
	kl.Log("alert cidrs", groupedCIDRs.alert)
	return groupedCIDRs, nil
}

func printTable(data map[*SecurityGroup][]*portToIP) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"ID", "Name", "Attachments", "Port", "Proto", "IP"})

	for sec, portToIps := range data {
		for i, pti := range portToIps {
			id := ""
			name := ""
			usage := ""
			if i == 0 {
				id = sec.id
				name = sec.name
				usage = sec.GetAttachmentsAsString()
			}
			t.AppendRow([]interface{}{
				id,
				name,
				usage,
				pti.port,
				pti.proto,
				pti.ip,
			})
		}
		t.AppendSeparator()
	}
	t.Render()
}

func printAmazonTable(data map[*SecurityGroup][]*portToIP) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"ID", "Name", "Attachments", "Port", "Proto", "IP", "Block", "Service", "Region"})

	for sec, portToIps := range data {
		for i, pti := range portToIps {
			id := ""
			name := ""
			usage := ""
			if i == 0 {
				id = sec.id
				name = sec.name
				usage = sec.GetAttachmentsAsString()
			}
			t.AppendRow([]interface{}{
				id,
				name,
				usage,
				pti.port,
				pti.proto,
				pti.ip,
				pti.prefix.IPPrefix,
				pti.prefix.GetService(),
				pti.prefix.Region,
			})
		}
		t.AppendSeparator()
	}
	t.Render()
}

func getJSON(url string, target interface{}) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(target)
}
