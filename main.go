package main

import (
	"fmt"
	"github.com/GoodwayGroup/gw-aws-audit/ec2"
	"github.com/GoodwayGroup/gw-aws-audit/info"
	"github.com/GoodwayGroup/gw-aws-audit/rds"
	"github.com/GoodwayGroup/gw-aws-audit/s3"
	"github.com/GoodwayGroup/gw-aws-audit/sg"
	"github.com/clok/cdocs"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
	"time"
)

var version string

func init() {
	cli.VersionPrinter = func(c *cli.Context) {
		_, _ = fmt.Fprintf(c.App.Writer, "%s %s (%s/%s)\n", info.AppName, version, runtime.GOOS, runtime.GOARCH)
	}
}

func main() {
	im, err := cdocs.InstallManpageCommand(&cdocs.InstallManpageCommandInput{
		AppName: info.AppName,
	})
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:     info.AppName,
		Version:  version,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
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
					{
						Name:  "pem-keys",
						Usage: "List instances and PEM key used at time of creation",
						Action: func(c *cli.Context) error {
							err := ec2.ListPemKeyUsage()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
				},
			},
			{
				Name:  "sg",
				Usage: "Security Group related commands",
				Subcommands: []*cli.Command{
					{
						Name:  "detached",
						Usage: "generate a report of all Security Groups that are NOT attached to an instance",
						UsageText: `
This command will scan the EC2 NetworkInterfaces to determine what
Security Groups are NOT attached/assigned in AWS.

`,
						Action: func(c *cli.Context) error {
							err := sg.ListDetachedSecurityGroups()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
					{
						Name:  "attached",
						Usage: "generate a report of all Security Groups that are attached to an instance",
						UsageText: `
This command will scan the EC2 NetworkInterfaces to determine what
Security Groups are attached/assigned in AWS.
`,
						Action: func(c *cli.Context) error {
							err := sg.ListAttachedSecurityGroups()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
					{
						Name:  "cidr",
						Usage: "generate a report comparing SG rules with input CIDR blocks",
						UsageText: `
$ gw-aws-audit sg cidr --allowed 10.176.0.0/16,10.175.0.0/16 --alert 174.0.0.0/8,1.2.3.4/32

This command will generate a report detecting the port to CIDR mapping rules 
for attached Security Groups. 

A list of Approved CIDRs is required. This is typically the CIDR block associated
with your VPC.
`,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "approved",
								Aliases:  []string{"a"},
								Usage:    "CIDR blocks that are approved (csv)",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "warn",
								Aliases: []string{"w"},
								Usage:   "CIDR blocks that will cause a warning (csv)",
								Value:   "204.0.0.0/8",
							},
							&cli.StringFlag{
								Name:    "alert",
								Aliases: []string{"b"},
								Usage:   "CIDR blocks that will cause an alert (csv)",
								Value:   "174.0.0.0/8",
							},
							&cli.StringFlag{
								Name:    "ignore-ports",
								Aliases: []string{"p"},
								Usage:   "Ports that can be ignored (csv)",
								Value:   "80,443,3,4,3-4",
							},
							&cli.StringFlag{
								Name:  "ignore-protocols",
								Usage: "Protocols to ignore. Can be tcp,udp,icmp (csv)",
							},
							&cli.BoolFlag{
								Name:  "all",
								Usage: "Process ALL Security Groups, not just attached",
								Value: false,
							},
						},
						Action: func(c *cli.Context) error {
							err := sg.GenerateCIDRReport(c)
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
					{
						Name:  "port",
						Usage: "generate a report comparing SG rules with input CIDR blocks on a specific port",
						UsageText: `
$ gw-aws-audit sg ports --ports 22,3306 --allowed 10.176.0.0/16,10.175.0.0/16 --alert 174.0.0.0/8,1.2.3.4/32

This command will generate a report for a set of PORTS for attached Security Groups.

A list of Approved CIDRs is required. This is typically the CIDR block associated
with your VPC.
`,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "ports",
								Aliases: []string{"p"},
								Usage:   "Ports to generate report on (csv)",
								Value:   "22",
							},
							&cli.StringFlag{
								Name:     "approved",
								Aliases:  []string{"a"},
								Usage:    "CIDR blocks that are approved (csv)",
								Required: true,
							},
							&cli.StringFlag{
								Name:    "warn",
								Aliases: []string{"w"},
								Usage:   "CIDR blocks that will cause a warning (csv)",
								Value:   "204.0.0.0/8",
							},
							&cli.StringFlag{
								Name:    "alert",
								Aliases: []string{"b"},
								Usage:   "CIDR blocks that will cause an alert (csv)",
								Value:   "174.0.0.0/8",
							},
							&cli.StringFlag{
								Name:  "ignore-protocols",
								Usage: "Protocols to ignore. Can be tcp,udp,icmp (csv)",
							},
							&cli.BoolFlag{
								Name:  "all",
								Usage: "Process ALL Security Groups, not just attached",
								Value: false,
							},
						},
						Action: func(context *cli.Context) error {
							err := sg.GeneratePortReport(context)
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
					{
						Name:  "amazon",
						Usage: "generate a report of allow SG with rules mapped to known AWS IPs",
						UsageText: `
This method loads the current version of https://ip-ranges.amazonaws.com/ip-ranges.json
and compares the CIDR blocks against all Security Groups.
`,
						Action: func(c *cli.Context) error {
							err := sg.GenerateExternalAWSIPReport()
							if err != nil {
								return cli.NewExitError(err, 2)
							}
							return nil
						},
					},
					{
						Name:    "direct-ip-mapping",
						Aliases: []string{"dim"},
						Usage:   "generate report of Security Groups with direct mappings to EC2 instances",
						UsageText: `
This method will generate a report comparing all Security Groups with all 
EC2 instances to determine where you have a direct IP mapping.

This will note Internal and External IP usage as well.
`,
						Action: func(c *cli.Context) error {
							err := sg.GenerateMappedEC2Report()
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
			im,
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
		docs, err := cdocs.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := cdocs.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
