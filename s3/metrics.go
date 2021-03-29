package s3

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/remeh/sizedwaitgroup"
	"log"
	"os"
	"sync/atomic"
	"time"
)

// Get ALL S3 Bucket metrics for a given region
func GetBucketMetrics() error {
	// create logger to STDERR
	l := log.New(os.Stderr, "", 0)
	l.Println("Starting metrics pull...")

	client := session.GetS3Client()
	result, err := client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		l.Println("Failed to list buckets")
		return err
	}

	swg := sizedwaitgroup.New(10)
	metrics := lib.Metrics{}

	fmt.Println("Bucket,Objects,Size (Bytes),Size (GB),Size (TB),Bytes per Object,MB per Object,Has Cost Tag")
	kmetrics.Printf("number of buckets: %d", len(result.Buckets))
	for _, bucket := range result.Buckets {
		go func(bucketName *string) {
			defer swg.Done()

			details := processBucketMetrics(bucketName)
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
		}(bucket.Name)
		swg.Add()
	}

	swg.Wait()
	l.Printf("Bucket metric pull complete. Buckets: %d Processed: %d\n", len(result.Buckets), metrics.Processed)
	return nil
}

func processBucketMetrics(bucketName *string) (details map[string]int) {
	kpbm.Printf("%s | processing", *bucketName)

	region, err := s3manager.GetBucketRegion(aws.BackgroundContext(), session.Session(), *bucketName, "us-west-2")
	kpbm.Printf("%s | region: %s", *bucketName, region)
	lib.HandleResponse(err)

	// Create session for region
	client := s3.New(awsSession.Must(awsSession.NewSession(&aws.Config{
		Region: aws.String(region),
	})))
	cwsvc := session.GetCloudWatchClient()

	// Check for s3-cost-name tag
	_, hasCostTag := checkCostTag(client, bucketName)

	st := time.Now().AddDate(0, 0, -2)
	et := time.Now().AddDate(0, 0, -1)
	kpbm.Printf("%s | Start: %s End: %s", *bucketName, st, et)

	// Pull bucket bytes size
	sizeResult, err := cwsvc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		StartTime: aws.Time(st),
		EndTime:   aws.Time(et),
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

	lib.HandleResponse(err)

	var sizeInBytes float64
	if len(sizeResult.Datapoints) > 0 {
		sizeInBytes = aws.Float64Value(sizeResult.Datapoints[0].Average)
	}

	// Pull bucket object counts
	countResult, err := cwsvc.GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		StartTime: aws.Time(st),
		EndTime:   aws.Time(et),
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

	lib.HandleResponse(err)

	var objectCount float64
	if len(countResult.Datapoints) > 0 {
		objectCount = aws.Float64Value(countResult.Datapoints[0].Average)
	}

	sizeInGB := sizeInBytes / 1000 / 1000 / 1000
	sizeInTB := sizeInBytes / 1000 / 1000 / 1000 / 1000
	bytesPerObject := sizeInBytes / objectCount
	megabytesPerObject := sizeInBytes / objectCount / 1000 / 1000
	if objectCount == 0 {
		bytesPerObject = 0
		megabytesPerObject = 0
	}
	hasTag := "no"
	if hasCostTag {
		hasTag = "yes"
	}

	// Print to STDOUT
	fmt.Printf("%s,%.0f,%.0f,%.2f,%.2f,%.2f,%.2f,%s\n", aws.StringValue(bucketName), objectCount, sizeInBytes, sizeInGB, sizeInTB, bytesPerObject, megabytesPerObject, hasTag)

	kpbm.Printf("%s | done processing", *bucketName)
	return map[string]int{"Processed": 1}
}
