package main

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/lib"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsEast1RegionID),
	}))
	client := ec2.New(sess)

	metrics := lib.SafeCounter{}

	results, err := client.DescribeVolumes(&ec2.DescribeVolumesInput{})

	fmt.Println("Detached Volumes")
	for _, volume := range results.Volumes {
		if len(volume.Attachments) <= 0 {
			metrics.Inc("Volumes")
			volParams := &ec2.DescribeVolumesInput{
				VolumeIds: []*string{volume.VolumeId},
			}
			volumes, err2 := client.DescribeVolumes(volParams)
			if err2 != nil {
				fmt.Println("Failed to list volumes", err2)
				return
			}
			metrics.Add("SumVolumeSize", int(aws.Int64Value(volumes.Volumes[0].Size)))

			var costs int
			if aws.StringValue(volume.VolumeType) == "gp2" {
				costs = int(aws.Int64Value(volumes.Volumes[0].Size))
			} else {
				costs = int(10 * float32(aws.Int64Value(volumes.Volumes[0].Size)) / 2) + int(float32(aws.Int64Value(volumes.Volumes[0].Iops)) * 0.65)
			}
			metrics.Add("VolumeCosts", costs)

			snapParams := &ec2.DescribeSnapshotsInput{
				Filters: []*ec2.Filter{
					{
						Name: aws.String("volume-id"),
						Values: []*string{
							volume.VolumeId,
						},
					},
				},
			}
			snapshots, err3 := client.DescribeSnapshots(snapParams)
			if err3 != nil {
				fmt.Println("Failed to list instances", err3)
				return
			}
			numSnaps := len(snapshots.Snapshots)
			metrics.Add("Snapshots", numSnaps)
			if numSnaps > 0 {
				metrics.Add("SumSnapshotSize", int(aws.Int64Value(volumes.Volumes[0].Size)))
			}
			fmt.Printf("volume: %s size: %d snaps: %d costs: %.2f\n", aws.StringValue(volume.VolumeId), aws.Int64Value(volumes.Volumes[0].Size), numSnaps, float32(costs) / 10)
		}
	}

	if err != nil {
		fmt.Println("Failed to list instances", err)
		return
	}

	fmt.Printf("\nVolumes found: %d\n", metrics.Value("Volumes"))
	fmt.Printf("Total volume size (GB): %d\n", metrics.Value("SumVolumeSize"))
	fmt.Printf("Total volume costs: %.2f\n\n", float32(metrics.Value("VolumeCosts")) / 10)
	fmt.Printf("Snapshots found: %d\n", metrics.Value("Snapshots"))
	fmt.Printf("Minimum snapshot size (GB): %d\n", metrics.Value("SumSnapshotSize"))
}
