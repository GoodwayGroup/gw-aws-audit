package command

import (
	"github.com/GoodwayGroup/gw-aws-audit/s3"
	"github.com/yitsushi/go-commander"
)

type BatchDeleteObjectsCommand struct {
}

// Execute is the main function. It will be called on BatchDeleteObjectsCommand command
func (c *BatchDeleteObjectsCommand) Execute(opts *commander.CommandHelper) {
	bucketName := opts.Arg(0)
	if bucketName == "" {
		panic("Bucket name is required")
	}
	s3.BatchDeleteObjects(bucketName)
}

// NewBatchDeleteObjectsCommand creates a new BatchDeleteObjectsCommand command
func NewBatchDeleteObjectsCommand(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &BatchDeleteObjectsCommand{},
		Help: &commander.CommandDescriptor{
			Name:             "s3-batch-delete",
			ShortDescription: "Deletes all objects in bucket",
			Arguments: "<bucket>",
			Examples: []string{
				"athena-express-akia6cuah7-2019",
			},
		},
	}
}
