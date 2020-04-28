package s3

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/thoas/go-funk"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func GetBucketMetrics() {
	var wg sync.WaitGroup
	// create logger to STDERR
	l := log.New(os.Stderr, "", 0)
	metrics := lib.Metrics{}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}))
	s3svc := s3.New(sess)
	cwsvc := cloudwatch.New(sess)

	l.Println("Starting metrics pull...")

	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		l.Println("Failed to list buckets", err)
		return
	}

	// create and start new bar
	numBuckets := len(result.Buckets)

	fmt.Println("Bucket,Objects,Size (Bytes),Size (Gigabytes),Bytes per Object,MB per Object,Has Cost Tag")
	for _, bucket := range result.Buckets {
		wg.Add(1)
		go func(bucketName *string) {
			details := processBucketMetrics(s3svc, cwsvc, bucketName)
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
	l.Printf("Bucket metric pull complete. Buckets: %d Processed: %d\n", numBuckets, metrics.Processed)
}

func processBucketMetrics(s3svc *s3.S3, cwsvc *cloudwatch.CloudWatch, bucketName *string) (details map[string]int) {
	// Check for s3-cost-name tag
	hasCostTag := checkCostTag(s3svc, bucketName)

	// Pull bucket bytes size
	sizeResult, err := cwsvc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		StartTime: aws.Time(time.Now().AddDate(0, 0, -2)),
		EndTime:   aws.Time(time.Now().AddDate(0, 0, -1)),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("BucketName"),
				Value: bucketName,
			},
			{
				Name:  aws.String("StorageType"),
				Value: aws.String("StandardStorage"),
			},
		},
		MetricName: aws.String("BucketSizeBytes"),
		Namespace:  aws.String("AWS/S3"),
		Period:     aws.Int64(86400),
		Statistics: []*string{aws.String("Average")},
	})

	lib.HandleResponse(err, true)

	var sizeInBytes float64
	if len(sizeResult.Datapoints) > 0 {
		sizeInBytes = aws.Float64Value(sizeResult.Datapoints[0].Average)
	}

	// Pull bucket object counts
	countResult, err := cwsvc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		StartTime: aws.Time(time.Now().AddDate(0, 0, -2)),
		EndTime:   aws.Time(time.Now().AddDate(0, 0, -1)),
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("BucketName"),
				Value: bucketName,
			},
			{
				Name:  aws.String("StorageType"),
				Value: aws.String("AllStorageTypes"),
			},
		},
		MetricName: aws.String("NumberOfObjects"),
		Namespace:  aws.String("AWS/S3"),
		Period:     aws.Int64(86400),
		Statistics: []*string{aws.String("Average")},
	})

	lib.HandleResponse(err, true)

	var objectCount float64
	if len(sizeResult.Datapoints) > 0 {
		objectCount = aws.Float64Value(countResult.Datapoints[0].Average)
	}

	gigabytes := sizeInBytes / 1000 / 1000 / 1000
	bytesPerObject := sizeInBytes / objectCount
	megabytesPerObject := sizeInBytes / objectCount / 1000 / 1000
	if objectCount == 0 {
		bytesPerObject = 0
		megabytesPerObject = 0
	}
	ynTag := "no"
	if hasCostTag {
		ynTag = "yes"
	}

	// Print to STDOUT
	fmt.Printf("%s,%.0f,%.0f,%.2f,%.2f,%.2f,%s\n", aws.StringValue(bucketName), sizeInBytes, objectCount, gigabytes, bytesPerObject, megabytesPerObject, ynTag)

	return map[string]int{"Processed": 1}
}

func checkCostTag(s3svc *s3.S3, bucketName *string) (hasCostTag bool) {
	result, err := s3svc.GetBucketTagging(&s3.GetBucketTaggingInput{
		Bucket: bucketName,
	})
	hasTags := handleGetTagsResponse(err)

	if hasTags {
		keys := make([]string, 0, len(result.TagSet))

		for _, tag := range result.TagSet {
			keys = append(keys, aws.StringValue(tag.Key))
		}

		if funk.ContainsString(keys, "s3-cost-name") {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
