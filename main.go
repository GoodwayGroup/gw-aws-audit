package main

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/ec2"
	"github.com/GoodwayGroup/gw-aws-audit/info"
	"github.com/GoodwayGroup/gw-aws-audit/rds"
	"github.com/GoodwayGroup/gw-aws-audit/s3"
	"io/ioutil"
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
		Commands: []*cli.Command{
			{
				Name:  "s3",
				Usage: "S3 related commands",
				Subcommands: []*cli.Command{
					{
						Name:  "add-cost-tag",
						Usage: "Add s3-cost-name to all S3 buckets",
						UsageText: `
Idempotent action that will add the ` + "`s3-cost-name`" + ` tag to ALL S3 buckets for a
given account.

The value will be the Bucket name.
`,
						Action: func(c *cli.Context) error {
							err := s3.AddCostTag()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
					{
						Name:  "metrics",
						Usage: "Get usage metrics",
						UsageText: `
Prints out a CSV report to STDOUT to help track usage across all buckets for a
given account.

Metrics per Bucket:

Objects (count)
Size (Bytes)
Size (GB)
Size (TB)
Bytes per Object
MB per Object
Has Cost Tag
`,
						Action: func(c *cli.Context) error {
							err := s3.GetBucketMetrics()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
					{
						Name:    "clear-bucket",
						Aliases: []string{"exterminatus"},
						Usage:   "Clear all Objects within a given Bucket",
						UsageText: `
Efficiently delete all objects within a bucket.

This process will run multiple paged deletes in parallel. It will handle API
throttling from AWS with an exponential backoff with retry. 
`,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "bucket",
								Aliases:  []string{"b"},
								Usage:    "Bucket to clear",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							err := s3.ClearBucketObjects(c)
							if err != nil {
								return cli.NewExitError(err, 2)
							}
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
							err := rds.ListMonitoringEnabled()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
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
							err := ec2.ListMonitoringEnabled()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
					{
						Name:  "detached-volumes",
						Usage: "List detached EBS volumes and snapshot counts",
						Action: func(c *cli.Context) error {
							err := ec2.ListDetachedVolumes()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
					{
						Name:  "stopped-hosts",
						Usage: "List stopped EC2 hosts and associated EBS volumes",
						Action: func(c *cli.Context) error {
							err := ec2.ListStoppedHosts()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
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
							err := ec2.ListMonitoringEnabled()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							fmt.Printf("\n\nChecking for RDS Enhanced Monitoring\n\n")
							err = rds.ListMonitoringEnabled()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "install-manpage",
				Usage: "Generate and install man page",
				Action: func(c *cli.Context) error {
					mp, _ := info.ToMan(c.App)
					err := ioutil.WriteFile("/usr/local/share/man/man8/gw-aws-audit.8", []byte(mp), 0644)
					if err != nil {
						return cli.NewExitError(fmt.Sprintf("Unable to install man page: %e", err), 2)
					}
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Print version info",
				Action: func(c *cli.Context) error {
					fmt.Printf("%s %s (%s/%s)\n", info.AppName, version, runtime.GOOS, runtime.GOARCH)
					return nil
				},
			},
		},
	}

	if os.Getenv("DOCS_MD") != "" {
		docs, err := info.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := info.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
