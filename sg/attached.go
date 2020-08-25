package sg

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
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
	t.AppendHeader(table.Row{"ID", "Name", "Attachments"})

	sort.Slice(attached, func(i, j int) bool {
		return strings.ToLower(attached[i].name) < strings.ToLower(attached[j].name)
	})
	total := 0
	for _, sg := range attached {
		var attachments []string
		for t, cnt := range sg.attached {
			total += cnt
			attachments = append(attachments, fmt.Sprintf("%s: %d", t, cnt))
		}
		t.AppendRow([]interface{}{sg.id, sg.name, strings.Join(attachments, " ")})
	}
	t.AppendFooter(table.Row{
		"Total Attached",
		fmt.Sprintf("%d / %d (%.0f%%)", len(attached), len(sgs), 100*float64(len(attached))/float64(len(sgs))),
		fmt.Sprintf("%d usages", total),
	})
	t.Render()

	return nil
}
