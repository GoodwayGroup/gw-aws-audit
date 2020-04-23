package command

import (
	"github.com/yitsushi/go-commander"
	"github.com/GoodwayGroup/gw-aws-audit/ec2"
)
type ListStoppedHostsCommand struct {
}

// Execute is the main function. It will be called on ListStoppedHostsCommand command
func (c *ListStoppedHostsCommand) Execute(opts *commander.CommandHelper) {
	ec2.ListStoppedHosts()
}

// NewListStoppedHostsCommand creates a new ListStoppedHostsCommand command
func NewListStoppedHostsCommand(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &ListStoppedHostsCommand{},
		Help: &commander.CommandDescriptor{
			Name:             "ec2-list-stopped-hosts",
			ShortDescription: "List stopped EC2 hosts and associated EBS volumes",
		},
	}
}
