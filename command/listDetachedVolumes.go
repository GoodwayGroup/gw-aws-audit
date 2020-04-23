package command

import (
	"github.com/GoodwayGroup/gw-aws-audit/ec2"

	"github.com/yitsushi/go-commander"
)
type ListDetachedVolumesCommand struct {
}

// Execute is the main function. It will be called on ListDetachedVolumesCommand command
func (c *ListDetachedVolumesCommand) Execute(opts *commander.CommandHelper) {
	ec2.ListDetachedVolumes()
}

// NewListDetachedVolumesCommand creates a new ListDetachedVolumesCommand command
func NewListDetachedVolumesCommand(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &ListDetachedVolumesCommand{},
		Help: &commander.CommandDescriptor{
			Name:             "ec2-list-detached-volumes",
			ShortDescription: "List detached EBS volumes and snapshot counts",
		},
	}
}
