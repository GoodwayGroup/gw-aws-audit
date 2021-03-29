package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/clok/kemba"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

var (
	k = kemba.New("gw-aws-audit:ec2")
)

// ListStoppedHosts will list all stopped EC2 hosts and attached EBS Volumes for those hosts for a given region.
func ListStoppedHosts() error {
	kl := k.Extend("ListStoppedHosts")
	client := session.GetEC2Client()

	var err error
	var results *ec2.DescribeInstancesOutput
	results, err = client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("stopped"),
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Failed to list instances")
		return err
	}

	var instCnt int64
	var volCnt int64
	var volumeSize int64
	var volumeCosts int64
	var snapCnt int64
	var snapSize int64

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Instance ID", "Name", "Volume", "Size (GB)", "Snapshots", "min Size (GB)", "Costs"})

	kl.Printf("found %d reservations", len(results.Reservations))
	for _, reservations := range results.Reservations {
		kl.Printf("%2s found %d instances", "└>", len(reservations.Instances))
		for _, instance := range reservations.Instances {
			instCnt++
			var name string
			for _, tag := range instance.Tags {
				if aws.StringValue(tag.Key) == "Name" {
					name = aws.StringValue(tag.Value)
				}
			}

			t.AppendRow([]interface{}{
				aws.StringValue(instance.InstanceId),
				name,
			})

			kl.Printf("%4s found %d volumes", "└>", len(instance.BlockDeviceMappings))
			for _, b := range instance.BlockDeviceMappings {
				volCnt++
				var volumes *ec2.DescribeVolumesOutput
				volumes, err = client.DescribeVolumes(&ec2.DescribeVolumesInput{
					VolumeIds: []*string{b.Ebs.VolumeId},
				})
				if err != nil {
					fmt.Println("Failed to list instances")
					return err
				}
				volumeSize += aws.Int64Value(volumes.Volumes[0].Size)

				var costs int64
				if aws.StringValue(volumes.Volumes[0].VolumeType) == "gp2" {
					costs = aws.Int64Value(volumes.Volumes[0].Size)
				} else {
					costs = int64(10*float32(aws.Int64Value(volumes.Volumes[0].Size))/2) + int64(float32(aws.Int64Value(volumes.Volumes[0].Iops))*0.65)
				}
				volumeCosts += costs

				var snapshots *ec2.DescribeSnapshotsOutput
				snapshots, err = client.DescribeSnapshots(&ec2.DescribeSnapshotsInput{
					Filters: []*ec2.Filter{
						{
							Name: aws.String("volume-id"),
							Values: []*string{
								b.Ebs.VolumeId,
							},
						},
					},
				})
				if err != nil {
					fmt.Println("Failed to list instances")
					return err
				}
				numSnaps := len(snapshots.Snapshots)
				snapCnt += int64(numSnaps)
				kl.Printf("%6s found %d snapshots", "└>", numSnaps)

				var lsnap int64
				if numSnaps > 0 {
					snapSize += aws.Int64Value(volumes.Volumes[0].Size)
					lsnap = aws.Int64Value(volumes.Volumes[0].Size)
				}

				t.AppendRow([]interface{}{
					"",
					"",
					aws.StringValue(b.Ebs.VolumeId),
					aws.Int64Value(volumes.Volumes[0].Size),
					numSnaps,
					lsnap,
					fmt.Sprintf("$%.2f", float32(costs)/10),
				})
			}
		}
	}

	t.AppendFooter(table.Row{
		"TOTALS",
		fmt.Sprintf("%d Instances", instCnt),
		fmt.Sprintf("%d Volumes", volCnt),
		fmt.Sprintf("%d GB", volumeSize),
		snapCnt,
		fmt.Sprintf("%d GB", snapSize),
		fmt.Sprintf("$%.2f", float32(volumeCosts)/10),
	})
	t.Render()
	return nil
}
