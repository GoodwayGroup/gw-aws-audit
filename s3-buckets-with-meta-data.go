package main

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"sync"
	"sync/atomic"
)

func main() {
	var ops uint64
	var wg sync.WaitGroup
	metrics := lib.Metrics{}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}))
	s3svc := s3.New(sess)
	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Failed to list buckets", err)
		return
	}

	// create and start new bar
	numBuckets := len(result.Buckets)

	for _, bucket := range result.Buckets {
		wg.Add(1)
		go func(bucketName *string) {
			details := processBucketDetails(s3svc, bucketName)
			atomic.AddUint64(&ops, 1)
			for key, _ := range details {
				switch key {
				case "Processed":
					atomic.AddInt64(&metrics.Processed, 1)
				case "Modified":
					atomic.AddInt64(&metrics.Modified, 1)
				case "Skipped":
					atomic.AddInt64(&metrics.Skipped, 1)
				}
			}

			wg.Done()
		}(bucket.Name)
	}

	wg.Wait()
	fmt.Printf("Bucket tagging complete. Buckets: %d Processed: %d Updated: %d Skipped %d", numBuckets, metrics.Processed, metrics.Modified, metrics.Skipped)
}

func processBucketDetails(s3svc *s3.S3, bucketName *string) (details map[string]int) {
	pageNum := 0
	err := s3svc.ListObjectsV2Pages(&s3.ListObjectsV2Input{
			Bucket:  bucketName,
			MaxKeys: aws.Int64(2),
		},
		func(page *s3.ListObjectsV2Output, lastPage bool) bool {
			pageNum++
			fmt.Println(page)
			return pageNum <= 1
		})

	handleListObjectResponse(err)
	return map[string]int{"Processed": 1, "Modified": 1}
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

