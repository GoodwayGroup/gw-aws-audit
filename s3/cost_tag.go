package s3

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	as "github.com/clok/awssession"
	"github.com/remeh/sizedwaitgroup"
	"github.com/thoas/go-funk"
	"sync/atomic"
)

var (
	kact  = k.Extend("AddCostTag")
	ktag  = k.Extend("checkCostTag")
	kup   = k.Extend("updateTags")
	khtr  = k.Extend("handleGetTagsResponse")
	kproc = kact.Extend("processBucket")
)

// Add the s3-cost-name tag with bucket name as value to ALL S3 Buckets
func AddCostTag() error {
	metrics := lib.Metrics{}

	sess, err := as.New()
	if err != nil {
		return err
	}
	s3svc := s3.New(sess)

	var result *s3.ListBucketsOutput
	result, err = s3svc.ListBuckets(&s3.ListBucketsInput{})
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
	sess, _ := as.New()
	region, err := s3manager.GetBucketRegion(aws.BackgroundContext(), sess, *bucketName, "us-west-2")
	kproc.Printf("%s | region: %s", *bucketName, region)
	lib.HandleResponse(err)

	// Create session for region
	// Create region based s3 service
	s3svc := s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	})))

	if tagSet, hasCostTag := checkCostTag(s3svc, bucketName); !hasCostTag {
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

				if updateTags(s3svc, newTags) {
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

		if updateTags(s3svc, newTags) {
			return map[string]int{"Processed": 1, "Modified": 1}
		}
	}

	kproc.Printf("%s | done processing", *bucketName)
	return map[string]int{"Processed": 1, "Skipped": 1}
}

func updateTags(s3svc *s3.S3, newTags *s3.PutBucketTaggingInput) bool {
	_, err := s3svc.PutBucketTagging(newTags)
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

func checkCostTag(s3svc *s3.S3, bucketName *string) ([]*s3.Tag, bool) {
	result, err := s3svc.GetBucketTagging(&s3.GetBucketTaggingInput{
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
		} else {
			kproc.Printf("%s | s3-cost-name not found", *bucketName)
			return result.TagSet, false
		}
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
		} else {
			panic(err)
		}
	}
	return true
}
