package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
	"sync/atomic"
)

func ClearBucketObjects(bucketName string) {
	var pageNum int64
	var listed int64
	var deleted int64
	var wg sync.WaitGroup

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}))
	s3svc := s3.New(sess)

	err := s3svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
		Bucket:  &bucketName,
		MaxKeys: aws.Int64(1000),
	},
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			if len(page.Contents) > 0 {
				wg.Add(1)
				go func() {
					atomic.AddInt64(&listed, int64(len(page.Contents)))

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
					handleListObjectResponse(err)
					atomic.AddInt64(&deleted, int64(len(page.Contents)))
					wg.Done()
				}()
				atomic.AddInt64(&pageNum, 1)
				fmt.Printf("Page: %d Objects: %d Listed: %d Deleted: %d\n", pageNum, len(page.Contents), listed, deleted)
			}
			return !lastPage
		})

	wg.Wait()

	handleListObjectResponse(err)
	fmt.Println("Process complete.")
	fmt.Printf("Pages: %d Objects: %d Deleted: %d\n", pageNum, listed, deleted)
}

func handleListObjectResponse(err error) (hasTags bool) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NoSuchTagSet" {
				return false
			} else {
				// Get error details
				// fmt.Printf("Error for bucket %s", aws.StringValue(bucketName))
				// fmt.Println("Error:", awsErr.Code(), awsErr.Message())
				return false
			}
		} else {
			panic(err)
		}
	} else {
		return true
	}
}
