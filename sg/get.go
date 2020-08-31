package sg

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	as "github.com/clok/awssession"
	"github.com/clok/kemba"
	"strings"
)

var (
	ksg = kemba.New("gw-aws-audit:sg")
)

func getAllSecurityGroups() (map[string]*securityGroup, error) {
	kl := ksg.Extend("get-all-sg")
	sess, err := as.New()
	if err != nil {
		return nil, err
	}
	client := ec2.New(sess)

	var results *ec2.DescribeSecurityGroupsOutput
	results, err = client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		MaxResults: aws.Int64(1000),
	})
	if err != nil {
		fmt.Println("Failed to list Security Groups")
		return nil, err
	}

	kl.Printf("found %d security groups", len(results.SecurityGroups))
	secGroups := make(map[string]*securityGroup, len(results.SecurityGroups))
	for _, sec := range results.SecurityGroups {
		rules := map[string][]*ec2.IpRange{}
		for _, rule := range sec.IpPermissions {
			var pstr string
			if rule.FromPort != nil {
				if aws.Int64Value(rule.FromPort) == 0 {
					pstr = "ALL"
				} else {
					pstr = fmt.Sprintf("%d", aws.Int64Value(rule.FromPort))
				}
			} else {
				pstr = "ALL"
			}
			if _, ok := rules[pstr]; !ok {
				rules[pstr] = rule.IpRanges
			} else {
				rules[pstr] = append(rules[pstr], rule.IpRanges...)
			}
		}

		secGroups[aws.StringValue(sec.GroupId)] = &securityGroup{
			id:    aws.StringValue(sec.GroupId),
			name:  aws.StringValue(sec.GroupName),
			rules: rules,
		}
	}

	return secGroups, nil
}

func detectAttachedSecurityGroups(sgs map[string]*securityGroup) error {
	kl := ksg.Extend("detect-attached")
	sess, err := as.New()
	if err != nil {
		return err
	}
	client := ec2.New(sess)

	var results *ec2.DescribeNetworkInterfacesOutput
	results, err = client.DescribeNetworkInterfaces(&ec2.DescribeNetworkInterfacesInput{
		MaxResults: aws.Int64(1000),
	})
	if err != nil {
		fmt.Println("Failed to list instances")
		return err
	}

	kl.Printf("found %d Network Interfaces", len(results.NetworkInterfaces))
	for _, network := range results.NetworkInterfaces {
		kl.Printf("%2s found %d security groups", "└>", len(network.Groups))
		if len(network.Groups) == 0 {
			kl.Extend("no-groups").Log(network)
		}
		for i, sec := range network.Groups {
			id := aws.StringValue(sec.GroupId)

			var owner string
			if network.Attachment != nil {
				switch ownerID := aws.StringValue(network.Attachment.InstanceOwnerId); {
				case ownerID == "amazon-aws":
					owner = aws.StringValue(network.InterfaceType)
				case strings.HasPrefix(ownerID, "amazon-"):
					owner = strings.Split(ownerID, "amazon-")[1]
				case strings.Contains(strings.ToLower(aws.StringValue(network.Description)), "eks"):
					owner = "eks"
				case strings.Contains(strings.ToLower(aws.StringValue(network.Description)), "efs"):
					owner = "efs"
				case strings.HasPrefix(aws.StringValue(network.Attachment.InstanceId), "i-"):
					owner = "ec2"
				case network.Attachment.InstanceId == nil:
					owner = aws.StringValue(network.InterfaceType)
				default:
					owner = aws.StringValue(network.Attachment.InstanceOwnerId)
				}
			} else {
				owner = "unknown"
			}

			if _, ok := sgs[id]; !ok {
				fmt.Printf("Found SG that was not in original list: %# v", sec)
			} else {
				if sgs[id].attached == nil {
					sgs[id].attached = map[string]int{}
				}
				sgs[id].attached[owner]++
			}

			stub := "├─"
			if i == len(network.Groups)-1 {
				stub = "└>"
			}
			kl.Printf("%6s %s -> %s", stub, id, aws.StringValue(sec.GroupName))
		}
	}

	kl.Extend("dump").Log(sgs)

	return nil
}

func getAnnotatedSecurityGroups() (map[string]*securityGroup, error) {
	// get all sgs in a region
	sgs, err := getAllSecurityGroups()
	if err != nil {
		return nil, err
	}

	err = detectAttachedSecurityGroups(sgs)
	if err != nil {
		return nil, err
	}
	return sgs, nil
}
