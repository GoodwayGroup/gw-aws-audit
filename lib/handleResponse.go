package lib

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"log"
	"os"
)

func HandleResponse(err error) (hasError bool) {
	l := log.New(os.Stderr, "", 0)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			// Get error details
			l.Printf("\nError Code: %s Message: %s\n", awsErr.Code(), awsErr.Message())
			return true
		} else {
			panic(err)
		}
	} else {
		return false
	}
}
