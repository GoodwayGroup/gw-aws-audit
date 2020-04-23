package command

import (
	"github.com/GoodwayGroup/gw-aws-audit/s3"
	"github.com/yitsushi/go-commander"
)

type ClearBucketObjectsCommand struct {
}

// Execute is the main function. It will be called on ClearBucketObjectsCommand command
func (c *ClearBucketObjectsCommand) Execute(opts *commander.CommandHelper) {
	bucketName := opts.Arg(0)
	if bucketName == "" {
		panic("Bucket name is required")
	}
	s3.ClearBucketObjects(bucketName)
}

// NewClearBucketObjectsCommand creates a new ClearBucketObjectsCommand command
func NewClearBucketObjectsCommand(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &ClearBucketObjectsCommand{},
		Help: &commander.CommandDescriptor{
			Name:             "s3-clear-bucket",
			ShortDescription: "Clear ALL objects from a Bucket",
			Arguments:        "<bucket>",
			Examples: []string{
				"athena-results-ASDF1337",
			},
		},
	}
}
