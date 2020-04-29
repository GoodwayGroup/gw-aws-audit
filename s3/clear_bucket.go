package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cenkalti/backoff/v4"
	"github.com/remeh/sizedwaitgroup"
	"github.com/tcnksm/go-input"
	"os"
	"sync/atomic"
	"time"
)

func ClearBucketObjects(bucketName string) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	fmt.Println("-- WARNING -- PAY ATTENTION -- FOR REALS --")
	fmt.Printf("This will delete ALL objects in %s\n", bucketName)
	fmt.Println("-- THIS ACTION IS NOT REVERSIBLE --")
	query := fmt.Sprintf("Are you SUPER sure? [%s]", bucketName)
	_, inputErr := ui.Ask(query, &input.Options{
		Required: true,
		// Validate input
		ValidateFunc: func(s string) error {
			if s != bucketName {
				return fmt.Errorf("Input must be %s to coninue. Exiting.", bucketName)
			}

			return nil
		},
	})
	if inputErr != nil {
		panic(inputErr)
	}
	fmt.Printf("Proceeding with batch delete for bucket: %s\n", bucketName)

	var pageNum int64
	var listed int64
	var deleted int64
	var retries int64
	swg := sizedwaitgroup.New(7)
	startTime := time.Now()

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}))
	s3svc := s3.New(sess)

	err := s3svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket:  &bucketName,
		MaxKeys: aws.Int64(1000),
	},
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			atomic.AddInt64(&listed, int64(len(page.Contents)))
			if len(page.Contents) > 0 {
				go func() {
					defer swg.Done()
					backoff.Retry(func() error {
						var objects []*s3.ObjectIdentifier
						for _, obj := range page.Contents {
							objects = append(objects, &s3.ObjectIdentifier{Key: obj.Key})
						}

						delInput := s3.DeleteObjectsInput{
							Bucket: &bucketName,
							Delete: &s3.Delete{
								Objects: objects,
								Quiet:   aws.Bool(true),
							},
						}

						_, err := s3svc.DeleteObjects(&delInput)
						hasError := handleResponse(err, &retries)
						if hasError {
							dps := float64(deleted) / time.Since(startTime).Seconds()
							fmt.Printf("\rPages: %d Listed: %d Deleted: %d Retries: %d DPS: %.2f", pageNum, listed, deleted, retries, dps)
							return err
						}
						atomic.AddInt64(&deleted, int64(len(page.Contents)))
						return nil
					}, backoff.NewExponentialBackOff())
				}()

				swg.Add()
				atomic.AddInt64(&pageNum, 1)
				dps := float64(deleted) / time.Since(startTime).Seconds()
				fmt.Printf("\rPages: %d Listed: %d Deleted: %d Retries: %d DPS: %.2f", pageNum, listed, deleted, retries, dps)
			}
			return !lastPage
		})

	swg.Wait()

	handleResponse(err, &retries)
	fmt.Println("Process complete.")
	dps := float64(deleted) / time.Since(startTime).Seconds()
	fmt.Printf("Pages: %d Listed: %d Deleted: %d Retries: %d DPS: %.2f", pageNum, listed, deleted, retries, dps)
}

func handleResponse(err error, retries *int64) (hasError bool) {
	if err != nil {
		if _, ok := err.(awserr.Error); ok {
			// Get error details
			// fmt.Printf("\nError Code: %s Message: %s\nRetrying...\n", awsErr.Code(), awsErr.Message())
			atomic.AddInt64(retries, 1)
			return true
		} else {
			panic(err)
		}
	} else {
		return false
	}
}
