package rds

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/sg"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/clok/kemba"
	"github.com/jedib0t/go-pretty/v6/table"
	"net"
	"os"
)

// ListPublicInterfaces will list RDS instances with a public interface attached.
func ListPublicInterfaces() error {
	k := kemba.New("gw-aws-audit:rds:ListPublicInterfaces")
	client := session.GetRDSClient()

	var err error
	var result *rds.DescribeDBInstancesOutput
	result, err = client.DescribeDBInstances(&rds.DescribeDBInstancesInput{})

	if err != nil {
		fmt.Println("Failed to list instances")
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"DB Instance", "Engine", "Security Groups"})

	k.Printf("checking %d RDS instances", len(result.DBInstances))
	cnt := 0
	for _, db := range result.DBInstances {
		if aws.BoolValue(db.PubliclyAccessible) {
			cnt++

			var sgIDs []*string
			for _, sec := range db.VpcSecurityGroups {
				sgIDs = append(sgIDs, sec.VpcSecurityGroupId)
			}
			sgs, err := sg.GetSecurityGroups(sgIDs)
			if err != nil {
				return err
			}
			var securityGroups []*sg.SecurityGroup
			for _, sec := range sgs {
				securityGroups = append(securityGroups, sec)
			}
			k.Log(securityGroups)

			var ips []string
			var stub string
			for _, sec := range securityGroups {
				for token, rule := range sec.Rules() {
					port, _, _ := sec.ParseRuleToken(token)
					for _, ip := range rule {
						_, ipv4Net, _ := net.ParseCIDR(aws.StringValue(ip.CidrIp))
						ips = append(ips, ipv4Net.String())
					}
					stub = fmt.Sprintf("%s\t%s\t%s\n\n\t", sec.ID(), sec.Name(), port)
					for i, ip := range ips {
						if i != 0 && i%4 == 0 {
							stub = fmt.Sprintf("%s\n\t", stub)
						}
						stub = fmt.Sprintf("%s %20s", stub, ip)
					}
					stub = fmt.Sprintf("%s\n", stub)
				}
			}

			name := aws.StringValue(db.DBInstanceIdentifier)
			engine := aws.StringValue(db.Engine)

			t.AppendRow([]interface{}{name, engine, stub})
		}
	}

	// There are a LOT of metrics to consider
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Monitoring.OS.html
	t.AppendFooter(table.Row{"DB Instances", cnt})
	t.Render()
	return nil
}
