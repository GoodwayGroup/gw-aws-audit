package rds

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
	"os"
)

// List RDS instances with CW Enhanced Monitoring enabled.
func ListMonitoringEnabled(c *cli.Context) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(c.String("region")),
	}))
	client := rds.New(sess)
	cnt := 0

	result, err := client.DescribeDBInstances(&rds.DescribeDBInstancesInput{})

	if err != nil {
		fmt.Println("Failed to list instances", err)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.AppendHeader(table.Row{"DB Instance", "Engine"})

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
}
