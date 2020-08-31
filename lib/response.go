package lib

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/clok/kemba"
)

var (
	k = kemba.New("gw-aws-audit:lib:HandleResponse")
)

// Generic response AWS response handler
func HandleResponse(err error) (hasError bool) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			k.Printf("AWS Error Code: %s", awsErr.Code())
			k.Printf("AWS Error Message: %s", awsErr.Message())
			k.Log(awsErr)
			return true
		}
		// TODO: refactor to remove panic
		panic(err)
	}
	return false
}
