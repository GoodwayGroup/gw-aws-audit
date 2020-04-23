package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/tcnksm/go-input"
	"os"
)

func BatchDeleteObjects(bucketName string) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	fmt.Println("-- WARNING -- PAY ATTENTION --")
	fmt.Printf("This will delete ALL objects in %s\n", bucketName)
	fmt.Println("-- THIS ACTION IS NOT REVERABLE --")
	query := "Are you SUPER sure? [YES/no]"
	name, err := ui.Ask(query, &input.Options{
		Required: true,
		// Validate input
		ValidateFunc: func(s string) error {
			if s != "YES" && s != "no" {
				return fmt.Errorf("input must be YES or no")
			}

			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Response: %s\n", name)
	fmt.Printf("Proceeding with batch delete for bucket: %s\n", bucketName)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}))
	svc := s3.New(sess)

	input := &s3.ListObjectsInput{
		Bucket:  &bucketName,
	}
	// Create a delete list objects iterator
	iter := s3manager.NewDeleteListIterator(svc, input)

	// Create the BatchDelete client
	batcher := s3manager.NewBatchDeleteWithClient(svc)

	ctx := aws.BackgroundContext()
	if err := batcher.Delete(ctx, iter); err != nil {
		panic(err)
	}
	fmt.Printf("Batch delete of objects in %s complete.\n", bucketName)
	fmt.Println(ctx)
}
