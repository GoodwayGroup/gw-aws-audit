package rds

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	as "github.com/clok/awssession"
	"github.com/clok/kemba"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

// List RDS instances with CW Enhanced Monitoring enabled.
func ListMonitoringEnabled() error {
	k := kemba.New("gw-aws-audit:rds:ListMonitoringEnabled")
	sess, err := as.New()
	if err != nil {
		return err
	}
	client := rds.New(sess)
	cnt := 0

	var result *rds.DescribeDBInstancesOutput
	result, err = client.DescribeDBInstances(&rds.DescribeDBInstancesInput{})

	if err != nil {
		fmt.Println("Failed to list instances")
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"DB Instance", "Engine"})

	k.Printf("checking %d RDS instances", len(result.DBInstances))
	for _, db := range result.DBInstances {
		if aws.Int64Value(db.MonitoringInterval) != 0 {
			cnt++

			name := aws.StringValue(db.DBInstanceIdentifier)
			engine := aws.StringValue(db.Engine)

			t.AppendRow([]interface{}{name, engine})
		}
	}

	// There are a LOT of metrics to consider
	// See: https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Monitoring.OS.html
	t.AppendFooter(table.Row{"DB Instances", cnt})
	t.Render()
	return nil
}
