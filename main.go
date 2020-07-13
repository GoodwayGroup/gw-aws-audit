package main

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/ec2"
	"github.com/GoodwayGroup/gw-aws-audit/info"
	"github.com/GoodwayGroup/gw-aws-audit/rds"
	"github.com/GoodwayGroup/gw-aws-audit/s3"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli/v2"
)

var version string

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s %s (%s/%s)\n", info.AppName, version, runtime.GOOS, runtime.GOARCH)
	}
}

func main() {
	app := &cli.App{
		Name:     info.AppName,
		Version:  version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Derek Smith",
				Email: "dsmith@goodwaygroup.com",
			},
		},
		Copyright:            "(c) 2020 Goodway Group",
		HelpName:             info.AppName,
		Usage:                "a collection of tools to audit AWS.",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "region",
				Usage:       "AWS region",
				Required:    false,
				Value:       "us-east-1",
				DefaultText: "us-east-1",
				EnvVars:     []string{"AWS_REGION"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version info",
				Action: func(c *cli.Context) error {
					fmt.Printf("%s %s (%s/%s)\n", info.AppName, version, runtime.GOOS, runtime.GOARCH)
					return nil
				},
			},
			{
				Name:  "s3",
				Usage: "S3 related commands",
				Subcommands: []*cli.Command{
					{
						Name:  "add-cost-tag",
						Usage: "Add s3-cost-name to all S3 buckets",
						Action: func(c *cli.Context) error {
							s3.AddCostTag(c)
							return nil
						},
					},
					{
						Name:  "metrics",
						Usage: "Get usage metrics",
						Action: func(c *cli.Context) error {
							s3.GetBucketMetrics(c)
							return nil
						},
					},
					{
						Name:    "clear-bucket",
						Aliases: []string{"exterminatus"},
						Usage:   "Clear all Objects within a given Bucket",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Aliases:  []string{"b"},
								Usage:    "Bucket to clear",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							s3.ClearBucketObjects(c)
							return nil
						},
					},
				},
			},
			{
				Name:  "rds",
				Usage: "RDS related commands",
				Subcommands: []*cli.Command{
					{
						Name:  "enhanced-monitoring",
						Usage: "Produce report of Enhanced Monitoring enabled instances",
						Action: func(c *cli.Context) error {
							rds.ListMonitoringEnabled(c)
							return nil
						},
					},
				},
			},
			{
				Name:  "ec2",
				Usage: "EC2 related commands",
				Subcommands: []*cli.Command{
					{
						Name:  "enhanced-monitoring",
						Usage: "Produce report of Enhanced Monitoring enabled instances",
						Action: func(c *cli.Context) error {
							ec2.ListMonitoringEnabled(c)
							return nil
						},
					},
					{
						Name:  "detached-volumes",
						Usage: "List detached EBS volumes and snapshot counts",
						Action: func(c *cli.Context) error {
							ec2.ListDetachedVolumes(c)
							return nil
						},
					},
					{
						Name:  "stopped-hosts",
						Usage: "List stopped EC2 hosts and associated EBS volumes",
						Action: func(c *cli.Context) error {
							ec2.ListStoppedHosts(c)
							return nil
						},
					},
				},
			},
			{
				Name:  "cw",
				Usage: "CloudWatch related commands",
				Subcommands: []*cli.Command{
					{
						Name:  "enhanced-monitoring",
						Usage: "Produce report of Enhanced Monitoring enabled EC2 & RDS instances",
						Action: func(c *cli.Context) error {
							fmt.Println("Enhanced Metrics can add a cost. See: https://aws.amazon.com/cloudwatch/pricing/")
							fmt.Printf("Checking for EC2 Enhanced Monitoring\n\n")
							ec2.ListMonitoringEnabled(c)
							fmt.Printf("\n\nChecking for RDS Enhanced Monitoring\n\n")
							rds.ListMonitoringEnabled(c)
							return nil
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
