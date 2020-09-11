package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	as "github.com/clok/awssession"
)

func getActiveInstances() (*ec2.DescribeInstancesOutput, error) {
	kl := k.Extend("getActiveInstances")
	sess, err := as.New()
	if err != nil {
		return nil, err
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
		return nil, err
	}

	kl.Log(results)

	return results, nil
}

type Info struct {
	ID         string
	Name       string
	State      string
	InternalIP string
	ExternalIP string
}

func GetEC2IPs() ([]*Info, error) {
	kl := k.Extend("GetEC2IPs")
	results, err := getActiveInstances()
	if err != nil {
		return nil, err
	}

	kl.Printf("found %d reservations", len(results.Reservations))
	var info []*Info
	for _, reservations := range results.Reservations {
		kl.Printf("%2s found %d instances", "└>", len(reservations.Instances))
		for _, instance := range reservations.Instances {
			var name string
			for _, tag := range instance.Tags {
				if aws.StringValue(tag.Key) == "Name" {
					name = aws.StringValue(tag.Value)
				}
			}

			info = append(info, &Info{
				ID:         aws.StringValue(instance.InstanceId),
				Name:       name,
				State:      aws.StringValue(instance.State.Name),
				InternalIP: aws.StringValue(instance.PrivateIpAddress),
				ExternalIP: aws.StringValue(instance.PublicIpAddress),
			})
		}
	}
	return info, nil
}

func GetInterfaceIPs() ([]*Info, error) {
	kl := k.Extend("GetInterfaceIPs")
	sess, err := as.New()
	if err != nil {
		return nil, err
	}
	client := ec2.New(sess)

	results, err := client.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{})

	if err != nil {
		fmt.Println("Failed to list interfaces")
		return nil, err
	}

	kl.Log(results)

	kl.Printf("found %d interfaces", len(results.NetworkInterfaces))
	var info []*Info
	for _, nic := range results.NetworkInterfaces {
		var name string
		for _, tag := range nic.TagSet {
			if aws.StringValue(tag.Key) == "Name" {
				name = aws.StringValue(tag.Value)
			}
		}

		if name == "" {
			name = aws.StringValue(nic.Description)
		}

		var externalIP string
		if nic.Association != nil {
			externalIP = aws.StringValue(nic.Association.PublicIp)
		}

		info = append(info, &Info{
			ID:         aws.StringValue(nic.NetworkInterfaceId),
			Name:       name,
			State:      aws.StringValue(nic.Status),
			InternalIP: aws.StringValue(nic.PrivateIpAddress),
			ExternalIP: externalIP,
		})
	}
	kl.Log(info)

	return info, nil
}
