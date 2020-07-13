package lib

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"log"
	"os"
)

// Generic response AWS response handler
func HandleResponse(err error, silent bool) (hasError bool) {
	l := log.New(os.Stderr, "", 0)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Get error details
			if !silent {
				l.Printf("\nError Code: %s Message: %s\n", awsErr.Code(), awsErr.Message())
			}
			return true
		} else {
			panic(err)
		}
	}
	return false
}
