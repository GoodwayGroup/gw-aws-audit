package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

// ListDetachedVolumes will list all EBS volumes that are in a Detached state, along with predicted associated cost.
func ListDetachedVolumes() error {
	kl := k.Extend("ListDetachedVolumes")
	client := session.GetEC2Client()

	var err error
	var results *ec2.DescribeVolumesOutput
	results, err = client.DescribeVolumes(&ec2.DescribeVolumesInput{})
	if err != nil {
		fmt.Println("Failed to list instances")
		return err
	}

	var volCnt int64
	var volumeSize int64
	var volumeCosts int64
	var snapCnt int64
	var snapSize int64

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"", "Volume", "Size (GB)", "Snapshots", "min Size (GB)", "Costs"})

	kl.Printf("found %d reservations", len(results.Volumes))
	for _, volume := range results.Volumes {
		if len(volume.Attachments) == 0 {
			volCnt++
			volParams := &ec2.DescribeVolumesInput{
				VolumeIds: []*string{volume.VolumeId},
			}
			volumes, err2 := client.DescribeVolumes(volParams)
			if err2 != nil {
				fmt.Println("Failed to list volumes")
				return err2
			}
			volumeSize += aws.Int64Value(volumes.Volumes[0].Size)

			var costs int64
			if aws.StringValue(volume.VolumeType) == "gp2" {
				costs = aws.Int64Value(volumes.Volumes[0].Size)
			} else {
				costs = int64(10*float32(aws.Int64Value(volumes.Volumes[0].Size))/2) + int64(float32(aws.Int64Value(volumes.Volumes[0].Iops))*0.65)
			}
			volumeCosts += costs

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
				fmt.Println("Failed to list instances")
				return err3
			}
			numSnaps := len(snapshots.Snapshots)
			snapCnt += int64(numSnaps)
			kl.Printf("%2s found %d snapshots", "└>", numSnaps)

			var lsnap int64
			if numSnaps > 0 {
				snapSize += aws.Int64Value(volumes.Volumes[0].Size)
				lsnap = aws.Int64Value(volumes.Volumes[0].Size)
				kl.Printf("%4s lsnap: %d snapSize: %d", "└>", lsnap, snapSize)
			}

			t.AppendRow([]interface{}{
				"",
				aws.StringValue(volume.VolumeId),
				aws.Int64Value(volumes.Volumes[0].Size),
				numSnaps,
				lsnap,
				fmt.Sprintf("$%.2f", float32(costs)/10),
			})
		}
	}

	t.AppendFooter(table.Row{
		"TOTALS",
		fmt.Sprintf("%d Volumes", volCnt),
		fmt.Sprintf("%d GB", volumeSize),
		snapCnt,
		fmt.Sprintf("%d GB", snapSize),
		fmt.Sprintf("$%.2f", float32(volumeCosts)/10),
	})
	t.Render()
	return nil
}
