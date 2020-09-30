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

// GetSecurityGroups will retrieve a list of Security Group IDs with mapped ports
func GetSecurityGroups(sgIDs []*string) (map[string]*SecurityGroup, error) {
	kl := ksg.Extend("get-sg")
	sess, err := as.New()
	if err != nil {
		return nil, err
	}
	client := ec2.New(sess)

	kl.Printf("retrieving SG IDs: %# v", sgIDs)
	var results *ec2.DescribeSecurityGroupsOutput
	results, err = client.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupIds: sgIDs,
	})
	if err != nil {
		fmt.Println("Failed to list Security Groups")
		return nil, err
	}

	kl.Printf("found %d security groups", len(results.SecurityGroups))
	secGroups := processSecurityGroupsResponse(results)

	return secGroups, nil
}

func processSecurityGroupsResponse(results *ec2.DescribeSecurityGroupsOutput) map[string]*SecurityGroup {
	secGroups := make(map[string]*SecurityGroup, len(results.SecurityGroups))
	for _, sec := range results.SecurityGroups {
		rules := map[string][]*ec2.IpRange{}
		for _, rule := range sec.IpPermissions {
			var t = make([]string, len(rule.UserIdGroupPairs))
			for i, u := range rule.UserIdGroupPairs {
				t[i] = aws.StringValue(u.GroupId)
			}

			var token string
			switch {
			case rule.FromPort == nil && rule.IpRanges == nil:
				token = buildPortToken("ALL", "", rule.IpProtocol, t)
			case aws.Int64Value(rule.FromPort) != aws.Int64Value(rule.ToPort):
				token = buildPortToken(fmt.Sprintf("%d", aws.Int64Value(rule.FromPort)), fmt.Sprintf("%d", aws.Int64Value(rule.ToPort)), rule.IpProtocol, t)
			case aws.Int64Value(rule.FromPort) == 0:
				token = buildPortToken("ALL", "", rule.IpProtocol, t)
			case aws.Int64Value(rule.FromPort) == -1:
				token = buildPortToken("ALL", "", rule.IpProtocol, t)
			case rule.IpRanges == nil:
				token = buildPortToken(fmt.Sprintf("%d", aws.Int64Value(rule.FromPort)), "", rule.IpProtocol, t)
			default:
				token = buildPortToken(fmt.Sprintf("%d", aws.Int64Value(rule.FromPort)), "", rule.IpProtocol, nil)
			}

			if _, ok := rules[token]; !ok {
				rules[token] = rule.IpRanges
			} else {
				rules[token] = append(rules[token], rule.IpRanges...)
			}
		}

		secGroups[aws.StringValue(sec.GroupId)] = &SecurityGroup{
			id:    aws.StringValue(sec.GroupId),
			name:  aws.StringValue(sec.GroupName),
			rules: rules,
		}
	}
	return secGroups
}

func getAllSecurityGroups() (map[string]*SecurityGroup, error) {
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
	secGroups := processSecurityGroupsResponse(results)

	return secGroups, nil
}

func buildPortToken(fromPort string, toPort string, proto *string, securityGroups []string) string {
	var parts = make([]string, 3)

	switch {
	case fromPort != "" && toPort != "":
		parts[0] = fmt.Sprintf("%s-%s", fromPort, toPort)
	case fromPort != "":
		parts[0] = fromPort
	default:
		parts[0] = "UNKNOWN"
	}

	switch {
	case proto == nil:
		parts[1] = "NONE"
	case aws.StringValue(proto) == "-1":
		parts[1] = "ALL"
	default:
		parts[1] = aws.StringValue(proto)
	}

	if securityGroups != nil {
		parts[2] = strings.Join(securityGroups, ",")
	}

	return strings.Join(parts, "::")
}

func detectAttachedSecurityGroups(sgs map[string]*SecurityGroup) error {
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

func getAnnotatedSecurityGroups() (map[string]*SecurityGroup, error) {
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
