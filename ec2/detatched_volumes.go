package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
	"os"
)

func ListDetachedVolumes(c *cli.Context) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(c.String("region")),
	}))
	client := ec2.New(sess)

	results, err := client.DescribeVolumes(&ec2.DescribeVolumesInput{})

	if err != nil {
		fmt.Println("Failed to list instances", err)
		return
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

	for _, volume := range results.Volumes {
		if len(volume.Attachments) <= 0 {
			volCnt++
			volParams := &ec2.DescribeVolumesInput{
				VolumeIds: []*string{volume.VolumeId},
			}
			volumes, err2 := client.DescribeVolumes(volParams)
			if err2 != nil {
				fmt.Println("Failed to list volumes", err2)
				return
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
				fmt.Println("Failed to list instances", err3)
				return
			}
			numSnaps := len(snapshots.Snapshots)
			snapCnt += int64(numSnaps)

			var lsnap int64
			if numSnaps > 0 {
				snapSize += aws.Int64Value(volumes.Volumes[0].Size)
				lsnap = aws.Int64Value(volumes.Volumes[0].Size)
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
}
