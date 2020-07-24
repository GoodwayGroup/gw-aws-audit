package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cenkalti/backoff/v4"
	as "github.com/clok/awssession"
	"github.com/clok/kemba"
	"github.com/remeh/sizedwaitgroup"
	"github.com/tcnksm/go-input"
	"github.com/urfave/cli/v2"
	"os"
	"sync/atomic"
	"time"
)

var (
	k    = kemba.New("gw-aws-audit:s3")
	kcbo = k.Extend("ClearBucketObjects")
	khr  = kcbo.Extend("handleResponse")
)

// Perform an asynchronous paged bulk delete of ALL objects within a Bucket.
func ClearBucketObjects(c *cli.Context) error {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	bucketName := c.String("bucket")

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
	var failed int64
	swg := sizedwaitgroup.New(7)
	startTime := time.Now()

	sess, err := as.New()
	if err != nil {
		return err
	}
	s3svc := s3.New(sess)

	err = s3svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket:  &bucketName,
		MaxKeys: aws.Int64(1000),
	},
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			atomic.AddInt64(&listed, int64(len(page.Contents)))
			if len(page.Contents) > 0 {
				go func() {
					defer swg.Done()
					err := backoff.Retry(func() error {
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
							fmt.Printf("\rPages: %d Listed: %d Deleted: %d Retries: %d Failed: %d DPS: %.2f", pageNum, listed, deleted, retries, failed, dps)
							return err
						}
						atomic.AddInt64(&deleted, int64(len(page.Contents)))
						return nil
					}, backoff.NewExponentialBackOff())
					if err != nil {
						atomic.AddInt64(&failed, int64(len(page.Contents)))
						return
					}
				}()

				swg.Add()
				atomic.AddInt64(&pageNum, 1)
				kcbo.Printf("%d Objects: %d lastPage: %t", pageNum, len(page.Contents), lastPage)
				dps := float64(deleted) / time.Since(startTime).Seconds()
				fmt.Printf("\rPages: %d Listed: %d Deleted: %d Retries: %d Failed: %d DPS: %.2f", pageNum, listed, deleted, retries, failed, dps)
			}
			return !lastPage
		})

	swg.Wait()

	handleResponse(err, &retries)
	fmt.Println("Process complete.")
	dps := float64(deleted) / time.Since(startTime).Seconds()
	fmt.Printf("Pages: %d Listed: %d Deleted: %d Retries: %d Failed: %d DPS: %.2f\n", pageNum, listed, deleted, retries, failed, dps)

	if err != nil {
		return err
	}
	return nil
}

func handleResponse(err error, retries *int64) (hasError bool) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			khr.Printf("AWS Error Code: %s", awsErr.Code())
			khr.Printf("AWS Error Message: %s", awsErr.Message())
			atomic.AddInt64(retries, 1)
			return true
		} else {
			panic(err)
		}
	} else {
		return false
	}
}
