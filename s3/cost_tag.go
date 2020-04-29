package s3

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/thoas/go-funk"
	"sync"
	"sync/atomic"
)

func AddCostTag() {
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

	var wg sync.WaitGroup
	for _, bucket := range result.Buckets {
		wg.Add(1)
		go func(bucketName *string) {
			defer wg.Done()

			details := processBucket(s3svc, bucketName)
			for key := range details {
				switch key {
				case "Processed":
					atomic.AddInt64(&metrics.Processed, 1)
				case "Modified":
					atomic.AddInt64(&metrics.Modified, 1)
				case "Skipped":
					atomic.AddInt64(&metrics.Skipped, 1)
				}
			}
			fmt.Printf("\rBuckets: %d Processed: %d Updated: %d Skipped %d", numBuckets, metrics.Processed, metrics.Modified, metrics.Skipped)
		}(bucket.Name)
	}

	wg.Wait()
	fmt.Printf("\nBucket tagging complete\n")
}

func processBucket(s3svc *s3.S3, bucketName *string) (details map[string]int) {
	result, err := s3svc.GetBucketTagging(&s3.GetBucketTaggingInput{
		Bucket: bucketName,
	})
	hasTags := handleGetTagsResponse(err)

	if hasTags {
		keys := make([]string, 0, len(result.TagSet))

		for _, tag := range result.TagSet {
			keys = append(keys, aws.StringValue(tag.Key))
		}

		if !funk.ContainsString(keys, "s3-cost-name") {
			result.TagSet = append(result.TagSet, &s3.Tag{
				Key:   aws.String("s3-cost-name"),
				Value: bucketName,
			})

			newTags := &s3.PutBucketTaggingInput{
				Bucket: bucketName,
				Tagging: &s3.Tagging{
					TagSet: result.TagSet,
				},
			}

			if updateTags(s3svc, newTags) {
				return map[string]int{"Processed": 1, "Modified": 1}
			}
			return map[string]int{"Processed": 1, "Skipped": 1}
		} else {
			return map[string]int{"Processed": 1, "Skipped": 1}
		}
	} else {
		newTags := &s3.PutBucketTaggingInput{
			Bucket: bucketName,
			Tagging: &s3.Tagging{
				TagSet: []*s3.Tag{
					{
						Key:   aws.String("s3-cost-name"),
						Value: bucketName,
					},
				},
			},
		}

		if updateTags(s3svc, newTags) {
			return map[string]int{"Processed": 1, "Modified": 1}
		}
		return map[string]int{"Processed": 1, "Skipped": 1}
	}
}

func updateTags(s3svc *s3.S3, newTags *s3.PutBucketTaggingInput) bool {
	_, err := s3svc.PutBucketTagging(newTags)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				// fmt.Println(aerr.Error())
				return false
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			// fmt.Println(err.Error())
			return false
		}

	}
	return true
}

func handleGetTagsResponse(err error) (hasTags bool) {
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
