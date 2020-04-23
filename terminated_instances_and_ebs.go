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
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("stopped"),
				},
			},
		},
	}

	metrics := lib.SafeCounter{}

	results, err := client.DescribeInstances(params)

	for _, reservations := range results.Reservations {
		for _, instance := range reservations.Instances {
			if processInstance(instance, metrics, client) {
				return
			}
		}
	}

	if err != nil {
		fmt.Println("Failed to list instances", err)
		return
	}

	fmt.Printf("\nVolumes found: %d\n", metrics.Value("Volumes"))
	fmt.Printf("Total volume size (GB): %d\n\n", metrics.Value("SumVolumeSize"))
	fmt.Printf("Snapshots found: %d\n", metrics.Value("Snapshots"))
	fmt.Printf("Minimum snapshot size (GB): %d\n", metrics.Value("SumSnapshotSize"))
}

func processInstance(instance *ec2.Instance, metrics lib.SafeCounter, client *ec2.EC2) bool {
	fmt.Printf("Instance %s - %s\n", aws.StringValue(instance.InstanceId), aws.StringValue(instance.Tags[2].Value))
	for _, b := range instance.BlockDeviceMappings {
		metrics.Inc("Volumes")
		volParams := &ec2.DescribeVolumesInput{
			VolumeIds: []*string{b.Ebs.VolumeId},
		}
		volumes, err2 := client.DescribeVolumes(volParams)
		if err2 != nil {
			fmt.Println("Failed to list instances", err2)
			return true
		}
		metrics.Add("SumVolumeSize", int(aws.Int64Value(volumes.Volumes[0].Size)))

		snapParams := &ec2.DescribeSnapshotsInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("volume-id"),
					Values: []*string{
						b.Ebs.VolumeId,
					},
				},
			},
		}
		snapshots, err3 := client.DescribeSnapshots(snapParams)
		if err3 != nil {
			fmt.Println("Failed to list instances", err3)
			return true
		}
		numSnaps := len(snapshots.Snapshots)
		metrics.Add("Snapshots", numSnaps)
		if numSnaps > 0 {
			metrics.Add("SumSnapshotSize", int(aws.Int64Value(volumes.Volumes[0].Size)))
		}
		fmt.Printf("\tvolume: %s size: %d snaps: %d\n", aws.StringValue(b.Ebs.VolumeId), aws.Int64Value(volumes.Volumes[0].Size), numSnaps)
	}
	return false
}
