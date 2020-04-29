package command

import (
	"github.com/GoodwayGroup/gw-aws-audit/s3"

	"github.com/yitsushi/go-commander"
)

type AddCostTagCommand struct {
}

// Execute is the main function. It will be called on AddCostTagCommand command
func (c *AddCostTagCommand) Execute(opts *commander.CommandHelper) {
	s3.AddCostTag()
}

// NewAddCostTagCommand creates a new AddCostTagCommand command
func NewAddCostTagCommand(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &AddCostTagCommand{},
		Help: &commander.CommandDescriptor{
			Name:             "s3-add-cost-tag",
			ShortDescription: "Add 's3-cost-name' tag to all buckets",
		},
	}
}
