package sg

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/logrusorgru/aurora/v3"
	"github.com/thoas/go-funk"
	"net"
	"os"
)

// GenerateExternalAWSIPReport
func GenerateExternalAWSIPReport() error {
	awsIPs, err := getAWSIPRanges()
	if err != nil {
		return err
	}

	sgs, err := getAnnotatedSecurityGroups()
	if err != nil {
		return err
	}

	var securityGroups []*securityGroup
	for _, sg := range sgs {
		securityGroups = append(securityGroups, sg)
	}

	mappedSGs := newGroupedSecurityGroups()
	for _, sec := range securityGroups {
		for token, rule := range sec.rules {
			port, proto, _ := parseToken(token)
			for _, ip := range rule {
				_, ipv4Net, err := net.ParseCIDR(aws.StringValue(ip.CidrIp))
				if err != nil {
					return err
				}

				// Skip
				if ipv4Net.IP.String() == "0.0.0.0" {
					continue
				}

				if prefix, ok := findAWSCidr(ipv4Net, awsIPs); ok {
					mappedSGs.addToAmazon(sec, &portToIP{
						port:   port,
						ip:     ipv4Net.String(),
						proto:  proto,
						prefix: prefix,
					})
				}
			}
		}
	}

	printAmazonTable(mappedSGs.amazon)

	return nil
}

func getAWSIPRanges() (*AWSIPs, error) {
	res := &AWSIPRanges{}
	err := getJSON("https://ip-ranges.amazonaws.com/ip-ranges.json", res)
	if err != nil {
		return nil, err
	}

	awsIPs := &AWSIPs{
		list:  []*net.IPNet{},
		table: map[*net.IPNet]*Prefix{},
	}
	for _, prefix := range res.Prefixes {
		cidr, err := prefix.GetCIDR()
		if err != nil {
			return nil, err
		}

		awsIPs.list = append(awsIPs.list, cidr)
		awsIPs.table[cidr] = prefix
	}

	for _, prefix := range res.IPv6Prefixes {
		cidr, err := prefix.GetCIDR()
		if err != nil {
			return nil, err
		}

		awsIPs.list = append(awsIPs.list, cidr)
		awsIPs.table[cidr] = &Prefix{
			IPPrefix:           prefix.IPv6Prefix,
			Region:             prefix.Region,
			NetworkBorderGroup: prefix.NetworkBorderGroup,
			Service:            prefix.Service,
		}
	}
	return awsIPs, nil
}

func findAWSCidr(input *net.IPNet, awsMap *AWSIPs) (*Prefix, bool) {
	for _, cidr := range awsMap.list {
		if cidr.Contains(input.IP) {
			return awsMap.table[cidr], true
		}
	}
	return nil, false
}

func GenerateMappedEC2Report() error {
	kl := ksg.Extend("GenerateMappedEC2Report")
	sgs, err := getAllSecurityGroups()
	if err != nil {
		return err
	}

	// create reverse lookup table:
	// cidr -> port ->> SG
	lookup := map[string]map[string][]string{}
	var allIPs []string
	for sgID, secGroup := range sgs {
		for token, rule := range secGroup.rules {
			port, _, _ := parseToken(token)

			for _, ip := range rule {
				_, ipv4Net, err := net.ParseCIDR(aws.StringValue(ip.CidrIp))
				if err != nil {
					return err
				}
				ipStr := ipv4Net.String()

				if !funk.ContainsString(allIPs, ipStr) {
					allIPs = append(allIPs, ipStr)
				}

				if _, ok := lookup[ipStr]; !ok {
					lookup[ipStr] = map[string][]string{}
					lookup[ipStr][port] = []string{sgID}
				} else {
					if _, ok := lookup[ipStr][port]; !ok {
						lookup[ipStr][port] = []string{sgID}
					} else {
						lookup[ipStr][port] = append(lookup[ipStr][port], sgID)
					}
				}
			}
		}
	}

	ec2Info, err := ec2.GetEC2IPs()
	if err != nil {
		return err
	}

	nicInfo, err := ec2.GetInterfaceIPs()
	if err != nil {
		return err
	}

	allInfo := append(ec2Info, nicInfo...)

	kl.Log(allInfo)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"EC2 Instance", "", "", "", "", "Security Group", "", ""}, table.RowConfig{AutoMerge: true})
	t.AppendHeader(table.Row{"Type", "IP", "Name", "ID", "State", "ID", "Name", "Port"})

	kl.Log(allIPs)
	for _, info := range allInfo {
		if info.ExternalIP != "" {
			if funk.ContainsString(allIPs, fmt.Sprintf("%s/32", info.ExternalIP)) {
				_, ipv4Net, _ := net.ParseCIDR(fmt.Sprintf("%s/32", info.ExternalIP))
				ipStr := ipv4Net.String()
				kl.Printf("found mapped IP: %s -> %# v", info.ExternalIP, lookup[ipStr])

				for port, sgIDs := range lookup[ipStr] {
					for _, sgID := range sgIDs {
						secGroup := sgs[sgID]
						t.AppendRow([]interface{}{
							aurora.Red("PUBLIC"),
							info.ExternalIP,
							info.Name,
							info.ID,
							info.State,
							secGroup.id,
							secGroup.name,
							port,
						})
					}
				}
			}
		}
		if info.InternalIP != "" {
			if funk.ContainsString(allIPs, fmt.Sprintf("%s/32", info.InternalIP)) {
				_, ipv4Net, _ := net.ParseCIDR(fmt.Sprintf("%s/32", info.InternalIP))
				ipStr := ipv4Net.String()
				kl.Printf("found mapped IP: %s -> %# v", info.InternalIP, lookup[ipStr])

				for port, sgIDs := range lookup[ipStr] {
					for _, sgID := range sgIDs {
						secGroup := sgs[sgID]
						t.AppendRow([]interface{}{
							aurora.Green("INTERNAL"),
							info.InternalIP,
							info.Name,
							info.ID,
							info.State,
							secGroup.id,
							secGroup.name,
							port,
						})
					}
				}
			}
		}
	}
	t.SortBy([]table.SortBy{
		{Number: 3},
		{Number: 2},
		{Number: 6},
	})
	t.Render()

	return nil
}
