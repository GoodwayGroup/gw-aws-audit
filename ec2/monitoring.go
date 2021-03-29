package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

// ListMonitoringEnabled will list EC2 instances with CW Enhanced Monitoring enabled.
func ListMonitoringEnabled() error {
	kl := k.Extend("ListMonitoringEnabled")
	client := session.GetEC2Client()

	cnt := 0
	var err error
	var result *ec2.DescribeInstancesOutput
	result, err = client.DescribeInstances(&ec2.DescribeInstancesInput{
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
		fmt.Println("Failed to list instances")
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"Name", "Instance ID"})

	kl.Printf("found %d reservations", len(result.Reservations))
	for _, reserve := range result.Reservations {
		kl.Printf("%2s found %d instances", "└>", len(reserve.Instances))
		if len(reserve.Instances) > 0 {
			for _, instance := range reserve.Instances {
				cnt++

				id := aws.StringValue(instance.InstanceId)

				var name string
				for _, tag := range reserve.Instances[0].Tags {
					if aws.StringValue(tag.Key) == "Name" {
						name = aws.StringValue(tag.Value)
					}
				}
				t.AppendRow([]interface{}{name, id})
			}
		}
	}

	// See: https://aws.amazon.com/cloudwatch/pricing/
	t.AppendFooter(table.Row{"EC2 Instances", cnt})
	t.Render()
	return nil
}
