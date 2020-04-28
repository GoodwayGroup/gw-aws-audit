package command

import (
	"github.com/GoodwayGroup/gw-aws-audit/s3"

	"github.com/yitsushi/go-commander"
)

type GetBucketMetricsCommand struct {
}

// Execute is the main function. It will be called on GetBucketMetricsCommand command
func (c *GetBucketMetricsCommand) Execute(opts *commander.CommandHelper) {
	s3.GetBucketMetrics()
}

// NewGetBucketMetricsCommand creates a new GetBucketMetricsCommand command
func NewGetBucketMetricsCommand(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &GetBucketMetricsCommand{},
		Help: &commander.CommandDescriptor{
			Name:             "s3-bucket-metrics",
			ShortDescription: "Print out bucket metrics to STDOUT",
		},
	}
}
