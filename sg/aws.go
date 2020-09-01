package sg

import (
	"github.com/aws/aws-sdk-go/aws"
	"net"
)

// GenerateExternalAWSIPReport
func GenerateExternalAWSIPReport() error {
	awsIPs, err := getAWSIPRanges()
	if err != nil {
		return err
	}

	sgs, err := getAnnotatedSecurityGroups()
	if err != nil {
		return err
	}

	var securityGroups []*securityGroup
	for _, sg := range sgs {
		securityGroups = append(securityGroups, sg)
	}

	mappedSGs := newGroupedSecurityGroups()
	for _, sec := range securityGroups {
		for token, rule := range sec.rules {
			port, proto, _ := parseToken(token)
			for _, ip := range rule {
				_, ipv4Net, err := net.ParseCIDR(aws.StringValue(ip.CidrIp))
				if err != nil {
					return err
				}

				// Skip
				if ipv4Net.IP.String() == "0.0.0.0" {
					continue
				}

				if prefix, ok := findAWSCidr(ipv4Net, awsIPs); ok {
					mappedSGs.addToAmazon(sec, &portToIP{
						port:   port,
						ip:     ipv4Net.String(),
						proto:  proto,
						prefix: prefix,
					})
				}
			}
		}
	}

	printAmazonTable(mappedSGs.amazon)

	return nil
}

func getAWSIPRanges() (*AWSIPs, error) {
	res := &AWSIPRanges{}
	err := getJSON("https://ip-ranges.amazonaws.com/ip-ranges.json", res)
	if err != nil {
		return nil, err
	}

	awsIPs := &AWSIPs{
		list:  []*net.IPNet{},
		table: map[*net.IPNet]*Prefix{},
	}
	for _, prefix := range res.Prefixes {
		cidr, err := prefix.GetCIDR()
		if err != nil {
			return nil, err
		}

		awsIPs.list = append(awsIPs.list, cidr)
		awsIPs.table[cidr] = prefix
	}

	for _, prefix := range res.IPv6Prefixes {
		cidr, err := prefix.GetCIDR()
		if err != nil {
			return nil, err
		}

		awsIPs.list = append(awsIPs.list, cidr)
		awsIPs.table[cidr] = &Prefix{
			IPPrefix:           prefix.IPv6Prefix,
			Region:             prefix.Region,
			NetworkBorderGroup: prefix.NetworkBorderGroup,
			Service:            prefix.Service,
		}
	}
	return awsIPs, nil
}

func findAWSCidr(input *net.IPNet, awsMap *AWSIPs) (*Prefix, bool) {
	for _, cidr := range awsMap.list {
		if cidr.Contains(input.IP) {
			return awsMap.table[cidr], true
		}
	}
	return nil, false
}
