package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	as "github.com/clok/awssession"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

// ListPemKeyUsage will generate a report of named pem keys used at creation of an EC2 host
func ListPemKeyUsage() error {
	kl := k.Extend("ListPemKeyUsage")
	sess, err := as.New()
	if err != nil {
		return err
	}
	client := ec2.New(sess)

	results, err := client.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					// pending | running | shutting-down | terminated | stopping | stopped
					aws.String("stopped"),
					aws.String("running"),
					aws.String("stopping"),
					aws.String("shutting-down"),
					aws.String("pending"),
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
	t.AppendHeader(table.Row{"Instance ID", "Name", "PEM Key"})

	kl.Printf("found %d reservations", len(results.Reservations))
	for _, reservations := range results.Reservations {
		kl.Printf("%2s found %d instances", "└>", len(reservations.Instances))
		for _, instance := range reservations.Instances {
			if instance.KeyName == nil {
				kl.Printf("Skipping host, no PEM Key found: %s", aws.StringValue(instance.InstanceId))
				break
			}

			var name string
			for _, tag := range instance.Tags {
				if aws.StringValue(tag.Key) == "Name" {
					name = aws.StringValue(tag.Value)
				}
			}

			t.AppendRow([]interface{}{
				aws.StringValue(instance.InstanceId),
				name,
				aws.StringValue(instance.KeyName),
			})
		}
	}

	t.Render()
	return nil
}
