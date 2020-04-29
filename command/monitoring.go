package command

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/ec2"
	"github.com/GoodwayGroup/gw-aws-audit/rds"
	"github.com/yitsushi/go-commander"
)

type ListMonitoringEnabledCommand struct {
}

// Execute is the main function. It will be called on ListMonitoringEnabledCommand command
func (c *ListMonitoringEnabledCommand) Execute(opts *commander.CommandHelper) {
	fmt.Println("Enhanced Metrics can add a cost. See: https://aws.amazon.com/cloudwatch/pricing/")
	fmt.Printf("Checking for EC2 Enhanced Monitoring\n\n")
	ec2.ListMonitoringEnabled()
	fmt.Printf("\n\nChecking for RDS Enhanced Monitoring\n\n")
	rds.ListMonitoringEnabled()
}

// NewListMonitoringEnabledCommand creates a new ListMonitoringEnabledCommand command
func NewListMonitoringEnabledCommand(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &ListMonitoringEnabledCommand{},
		Help: &commander.CommandDescriptor{
			Name:             "monitoring-enabled",
			ShortDescription: "List EC2 & RDS hosts with Enhanced Monitoring enabled",
		},
	}
}
