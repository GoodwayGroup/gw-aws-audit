package sg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/thoas/go-funk"
	"github.com/urfave/cli/v2"
	"net"
	"os"
	"strings"
)

type portToIp struct {
	port string
	ip   string
}

func GenerateCIDRReport(c *cli.Context) error {
	kl := ksg.Extend("GenerateCIDRReport")
	sgs, err := getAnnotatedSegurityGroups()
	if err != nil {
		return err
	}

	var attached []*securityGroup
	for _, sg := range sgs {
		if sg.attached != nil {
			attached = append(attached, sg)
		}
	}

	kl.Log(c.String("ignore-ports"))
	ignoredPorts := strings.Split(c.String("ignore-ports"), ",")
	fmt.Printf("----\nIgnored ports: %s\n----\n\n", ignoredPorts)

	// Get approved CIDR blocks
	var approvedCidrs []*net.IPNet
	for _, cidr := range strings.Split(c.String("approved"), ",") {
		_, ipv4Net, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}
		approvedCidrs = append(approvedCidrs, ipv4Net)
	}
	kl.Log("approved cidrs", approvedCidrs)

	// Get warn CIDR blocks
	var warnCidrs []*net.IPNet
	if c.String("warn") != "" {
		for _, cidr := range strings.Split(c.String("warn"), ",") {
			_, ipv4Net, err := net.ParseCIDR(cidr)
			if err != nil {
				return err
			}
			warnCidrs = append(warnCidrs, ipv4Net)
		}
	}
	kl.Log("warn cidrs", warnCidrs)

	// Get approved CIDR blocks
	var alertCidrs []*net.IPNet
	for _, cidr := range strings.Split(c.String("alert"), ",") {
		_, ipv4Net, err := net.ParseCIDR(cidr)
		if err != nil {
			return err
		}
		alertCidrs = append(alertCidrs, ipv4Net)
	}
	kl.Log("alert cidrs", alertCidrs)

	alerts := map[*securityGroup][]*portToIp{}
	warns := map[*securityGroup][]*portToIp{}
	unknowns := map[*securityGroup][]*portToIp{}
	wide := map[*securityGroup][]*portToIp{}
	for _, sec := range attached {
		for port, rule := range sec.rules {
			if !funk.ContainsString(ignoredPorts, port) {
				for _, ip := range rule {
					_, ipv4Net, err := net.ParseCIDR(aws.StringValue(ip.CidrIp))
					if err != nil {
						return err
					}

					if ipv4Net.IP.String() == "0.0.0.0" {
						kl.Printf("%s\t%s\t%s", "FULL", ipv4Net.String(), port)
						if _, ok := alerts[sec]; !ok {
							wide[sec] = []*portToIp{
								{
									port: port,
									ip:   ipv4Net.String(),
								},
							}
						} else {
							wide[sec] = append(wide[sec], &portToIp{
								port: port,
								ip:   ipv4Net.String(),
							})
						}
						break
					}

					if isCiderIn(ipv4Net, approvedCidrs) {
						kl.Printf("%s\t%s\t%s", "APPROVED", ipv4Net.String(), port)
					} else if isCiderIn(ipv4Net, warnCidrs) {
						kl.Printf("%s\t%s\t%s", "WARN", ipv4Net.String(), port)
						if _, ok := warns[sec]; !ok {
							warns[sec] = []*portToIp{
								{
									port: port,
									ip:   ipv4Net.String(),
								},
							}
						} else {
							warns[sec] = append(warns[sec], &portToIp{
								port: port,
								ip:   ipv4Net.String(),
							})
						}
					} else if isCiderIn(ipv4Net, alertCidrs) {
						kl.Printf("%s\t%s\t%s", "ALERT", ipv4Net.String(), port)
						if _, ok := alerts[sec]; !ok {
							alerts[sec] = []*portToIp{
								{
									port: port,
									ip:   ipv4Net.String(),
								},
							}
						} else {
							alerts[sec] = append(alerts[sec], &portToIp{
								port: port,
								ip:   ipv4Net.String(),
							})
						}
					} else {
						kl.Printf("%s\t%s\t%s", "UNKNOWN", ipv4Net.String(), port)
						if _, ok := unknowns[sec]; !ok {
							unknowns[sec] = []*portToIp{
								{
									port: port,
									ip:   ipv4Net.String(),
								},
							}
						} else {
							unknowns[sec] = append(unknowns[sec], &portToIp{
								port: port,
								ip:   ipv4Net.String(),
							})
						}
					}
				}
			} else {
				kl.Printf("ignoring port %s", port)
			}
		}
	}

	fmt.Println("Unknown if these are allowed or not based on filters provided.")
	printTable(unknowns)
	fmt.Println("")
	fmt.Println("WARNINGS: CIDR rules flagged as warnings")
	printTable(warns)
	fmt.Println("")
	fmt.Println("ALERTS: CIDR rules flagged as alerts. Probable they are Public IPs")
	printTable(alerts)
	fmt.Println("")
	fmt.Println("WIDE OPEN: CIDR rules that are wide open. Should be verified that this is intended.")
	printTable(wide)
	return nil
}

func printTable(data map[*securityGroup][]*portToIp) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"ID", "Name", "Port", "IP"})

	for sec, portToIps := range data {
		for i, pti := range portToIps {
			id := ""
			name := ""
			if i == 0 {
				id = sec.id
				name = sec.name
			}
			t.AppendRow([]interface{}{
				id,
				name,
				pti.port,
				pti.ip,
			})
		}
	}
	t.Render()
}

func isCiderIn(input *net.IPNet, cidrs []*net.IPNet) bool {
	for _, cidr := range cidrs {
		if cidr.Contains(input.IP) {
			return true
		}
	}
	return false
}
