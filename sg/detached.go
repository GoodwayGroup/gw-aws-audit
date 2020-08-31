package sg

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
	"sort"
	"strings"
)

// ListDetachedSecurityGroups generates a report listing out all Security Groups
// that are NOT attached to a Network Interface
func ListDetachedSecurityGroups() error {
	sgs, err := getAnnotatedSecurityGroups()
	if err != nil {
		return err
	}

	var detached []*securityGroup
	for _, sg := range sgs {
		if sg.attached == nil {
			detached = append(detached, sg)
		}
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"ID", "Name"})

	sort.Slice(detached, func(i, j int) bool {
		return strings.ToLower(detached[i].name) < strings.ToLower(detached[j].name)
	})
	for _, sg := range detached {
		t.AppendRow([]interface{}{sg.id, sg.name})
	}
	t.AppendFooter(table.Row{"Total Detached", fmt.Sprintf("%d / %d (%.0f%%)", len(detached), len(sgs), 100*float64(len(detached))/float64(len(sgs)))})
	t.Render()

	return nil
}
