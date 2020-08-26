package sg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/thoas/go-funk"
	"os"
	"sort"
	"strings"
)

// ListAttachedSecurityGroups generates a report listing out all Security Groups
// that are attached to a Network Interface
func ListAttachedSecurityGroups() error {
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

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"ID", "Name", "Attachments", "Ports", "CIDRs"})

	sort.Slice(attached, func(i, j int) bool {
		return strings.ToLower(attached[i].name) < strings.ToLower(attached[j].name)
	})
	total := 0
	var totalPorts []string
	var totalCidrsBlocks []string
	for _, sg := range attached {
		var attachments []string
		for t, cnt := range sg.attached {
			total += cnt
			attachments = append(attachments, fmt.Sprintf("%s: %d", t, cnt))
		}

		var ports []string
		var cidrsBlocks []string
		for port, cidrs := range sg.rules {
			ports = append(ports, port)
			for _, cidr := range cidrs {
				cidrsBlocks = append(cidrsBlocks, aws.StringValue(cidr.CidrIp))
			}
		}
		uniqPorts := funk.UniqString(ports)
		uniqCidrs := funk.UniqString(cidrsBlocks)
		totalPorts = append(totalPorts, uniqPorts...)
		totalCidrsBlocks = append(totalCidrsBlocks, uniqCidrs...)
		t.AppendRow([]interface{}{sg.id, sg.name, strings.Join(attachments, " "), len(uniqPorts), len(uniqCidrs)})
	}
	t.AppendFooter(table.Row{
		"Summary",
		fmt.Sprintf("%d / %d (%.0f%%)", len(attached), len(sgs), 100*float64(len(attached))/float64(len(sgs))),
		fmt.Sprintf("%d", total),
		fmt.Sprintf("%d", len(funk.UniqString(totalPorts))),
		fmt.Sprintf("%d", len(funk.UniqString(totalCidrsBlocks))),
	})
	t.AppendFooter(table.Row{
		"",
		"Attached / Total",
		"Total Usage",
		"Unique Ports",
		"Unique CIDRs",
	})
	t.Render()

	return nil
}
