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

func ListMonitoringEnabled(c *cli.Context) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(c.String("region")),
	}))
	client := ec2.New(sess)
	cnt := 0

	result, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("monitoring-state"),
				Values: []*string{
					aws.String("enabled"),
				},
			},
		},
	})

	if err != nil {
		fmt.Println("Failed to list instances", err)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Name", "Instance ID"})

	for _, reserve := range result.Reservations {
		if len(reserve.Instances) > 0 {
			cnt++

			id := aws.StringValue(reserve.Instances[0].InstanceId)

			var name string
			for _, tag := range reserve.Instances[0].Tags {
				if aws.StringValue(tag.Key) == "Name" {
					name = aws.StringValue(tag.Value)
				}
			}
			t.AppendRow([]interface{}{name, id})
		}
	}

	// See: https://aws.amazon.com/cloudwatch/pricing/
	t.AppendFooter(table.Row{"EC2 Instances", cnt})
	t.Render()
}
