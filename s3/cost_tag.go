package s3

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/remeh/sizedwaitgroup"
	"github.com/thoas/go-funk"
	"sync/atomic"
)

// Add the s3-cost-name tag with bucket name as value to ALL S3 Buckets
func AddCostTag() error {
	metrics := lib.Metrics{}
	client := session.GetS3Client()

	var err error
	var result *s3.ListBucketsOutput
	result, err = client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Failed to list buckets")
		return err
	}

	// create and start new bar
	numBuckets := len(result.Buckets)
	kact.Printf("number of buckets: %d", numBuckets)

	swg := sizedwaitgroup.New(10)
	for _, bucket := range result.Buckets {
		go func(bucketName *string) {
			defer swg.Done()

			details := processBucket(bucketName)
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
		swg.Add()
	}

	swg.Wait()
	fmt.Printf("\nBucket tagging complete\n")
	return nil
}

func processBucket(bucketName *string) (details map[string]int) {
	kproc.Printf("%s | processing", *bucketName)
	region, err := s3manager.GetBucketRegion(aws.BackgroundContext(), session.Session(), *bucketName, "us-west-2")
	kproc.Printf("%s | region: %s", *bucketName, region)
	lib.HandleResponse(err)

	// Create session for region
	// Create region based s3 service
	client := s3.New(awsSession.Must(awsSession.NewSession(&aws.Config{
		Region: aws.String(region),
	})))

	if tagSet, hasCostTag := checkCostTag(client, bucketName); !hasCostTag {
		if tagSet != nil {
			kproc.Printf("%s | hasTags: %# v", *bucketName, tagSet)
			keys := make([]string, 0, len(tagSet))

			for _, tag := range tagSet {
				keys = append(keys, aws.StringValue(tag.Key))
			}

			if !funk.ContainsString(keys, "s3-cost-name") {
				tagSet = append(tagSet, &s3.Tag{
					Key:   aws.String("s3-cost-name"),
					Value: bucketName,
				})

				newTags := &s3.PutBucketTaggingInput{
					Bucket: bucketName,
					Tagging: &s3.Tagging{
						TagSet: tagSet,
					},
				}

				if updateTags(client, newTags) {
					return map[string]int{"Processed": 1, "Modified": 1}
				}
				return map[string]int{"Processed": 1, "Skipped": 1}
			}
			kproc.Printf("%s | s3-cost-name found", *bucketName)
			return map[string]int{"Processed": 1, "Skipped": 1}
		}

		kproc.Printf("%s | no tags found", *bucketName)
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

		if updateTags(client, newTags) {
			return map[string]int{"Processed": 1, "Modified": 1}
		}
	}

	kproc.Printf("%s | done processing", *bucketName)
	return map[string]int{"Processed": 1, "Skipped": 1}
}

func updateTags(client *s3.S3, newTags *s3.PutBucketTaggingInput) bool {
	_, err := client.PutBucketTagging(newTags)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			kup.Printf("AWS Error Code: %s", awsErr.Code())
			kup.Printf("AWS Error Message: %s", awsErr.Message())
			kup.Log(awsErr)
			return false
		}

		kup.Log(err)
		return false
	}
	return true
}

func checkCostTag(client *s3.S3, bucketName *string) ([]*s3.Tag, bool) {
	result, err := client.GetBucketTagging(&s3.GetBucketTaggingInput{
		Bucket: bucketName,
	})
	hasTags := handleGetTagsResponse(err)

	if hasTags {
		ktag.Printf("%s | hasTags: %# v", *bucketName, result.TagSet)
		keys := make([]string, 0, len(result.TagSet))

		for _, tag := range result.TagSet {
			keys = append(keys, aws.StringValue(tag.Key))
		}

		if funk.ContainsString(keys, "s3-cost-name") {
			kproc.Printf("%s | s3-cost-name found", *bucketName)
			return result.TagSet, true
		}
		kproc.Printf("%s | s3-cost-name not found", *bucketName)
		return result.TagSet, false
	}
	return nil, false
}

func handleGetTagsResponse(err error) (hasTags bool) {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NoSuchTagSet" {
				return false
			}
			khtr.Printf("AWS Error Code: %s", awsErr.Code())
			khtr.Printf("AWS Error Message: %s", awsErr.Message())
			khtr.Log(awsErr)
			return false
		}
		// TODO: refactor to remove panic
		panic(err)
	}
	return true
}
