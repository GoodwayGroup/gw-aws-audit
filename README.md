# GW AWS Audit Tool
> NOTE: This is a specialized tool to help with actions to take during an audit of AWS usage.

[![Go Report Card](https://goreportcard.com/badge/GoodwayGroup/gw-aws-audit)](https://goreportcard.com/report/GoodwayGroup/gw-aws-audit)

## Basic Usage

Please see [the docs for details on the commands.](./docs/gw-aws-audit.md)

Useful for clearing **large S3 buckets (many millions of objects)**, identifying egress EBS volumes and tracking S3 spend.

```
$ gw-aws-audit help
NAME:
   gw-aws-audit - a collection of tools to audit AWS.

USAGE:
   gw-aws-audit [global options] command [command options] [arguments...]

AUTHOR:
   Derek Smith <dsmith@goodwaygroup.com>

COMMANDS:
   s3               S3 related commands
   rds              RDS related commands
   ec2              EC2 related commands
   sg               Security Group related commands
   cw               CloudWatch related commands
   install-manpage  Generate and install man page
   version, v       Print version info
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)

COPYRIGHT:
   (c) 2020 Goodway Group
```

## Installation

### [`asdf` plugin](https://github.com/GoodwayGroup/asdf-gw-aws-audit)

Add plugin:

```
$ asdf plugin-add gw-aws-audit https://github.com/GoodwayGroup/asdf-gw-aws-audit.git
```

Install the latest version:

```
$ asdf install gw-aws-audit latest
```

### [Homebrew](https://brew.sh) (for macOS users)

```
brew tap GoodwayGroup/gw-aws-audit
brew install gw-aws-audit
```

### curl binary

```
$ curl https://i.jpillora.com/GoodwayGroup/gw-aws-audit! | bash
```

### man page

To install `man` page:

```
$ gw-aws-audit install-manpage
```

## Audit helper script

There is a bash helper script in the repo at [audit.sh](audit.sh). This tool is useful in running an audit across many regions at once.

```
$ ./audit.sh

audit.sh helper script for gw-aws-audit

Usage:
    audit.sh [gw-aws-audit commands]

Examples:
> This will run the 'gw-aws-audit sg detached' command for every region in the US (default)

    $ audit.sh sg detached

> This will run the 'gw-aws-audit ec2 stopped-hosts' for ONLY the us-west-2 region

    $ AWS_REGION=us-west-2 audit.sh ec2 stopped-hosts

> This will run the 'gw-aws-audit ec2 stopped-hosts' for every region in the EU

    $ REGION=eu audit.sh ec2 stopped-hosts

> This will run the 'gw-aws-audit cw monitoring' using a specific version of the tool.

    $ BIN_PATH=./bin/gw-aws-audit audit.sh cw monitoring

Note: REGION env values (default: US):
US: us-east-1 us-east-2 us-west-1 us-west-2
EU: eu-central-1 eu-west-1 eu-west-2 eu-west-3 eu-south-1 eu-north-1
AP: ap-east-1 ap-south-1 ap-northeast-3 ap-northeast-2 ap-southeast-1 ap-southeast-2 ap-northeast-1
CH: cn-north-1 cn-northwest-1
ROW: af-south-1 me-south-1 sa-east-2
ALL: All of the above combined

You can also set AWS_REGION and that will supersede the value of REGION
✔ Have fun!
```

### Command Categories

There are commands for `s3`, `ec2`, `rds`, `sg` and `cw`

**s3**
```
$ gw-aws-audit s3
...
COMMANDS:
   add-cost-tag                Add s3-cost-name to all S3 buckets
   metrics                     Get usage metrics
   clear-bucket, exterminatus  Clear all Objects within a given Bucket
```

**ec2**
```
$ gw-aws-audit ec2
...
COMMANDS:
   enhanced-monitoring  Produce report of Enhanced Monitoring enabled instances
   detached-volumes     List detached EBS volumes and snapshot counts
   stopped-hosts        List stopped EC2 hosts and associated EBS volumes
```

**rds**
```
$ gw-aws-audit rds
...
COMMANDS:
   enhanced-monitoring  Produce report of Enhanced Monitoring enabled instances
```

**cw**
```
$ gw-aws-audit cw
...
COMMANDS:
   enhanced-monitoring  Produce report of Enhanced Monitoring enabled EC2 & RDS instances
```

### Example Outputs

#### ec2 stopped-hosts
```
$ gw-aws-audit ec2 stopped-hosts
                INSTANCE ID          NAME            VOLUME                 SIZE (GB)  SNAPSHOTS  MIN SIZE (GB)  COSTS
                i-09e42474f22039e23  dummy-box-test
                                                     vol-0d4b4a7bc95a4b8e4          8          0              0  $0.80
                                                     vol-0cc0f6cd3c99bc1cc          8          0              0  $0.80
                TOTALS               1 INSTANCES     2 VOLUMES                  16 GB          0           0 GB  $1.60
```

#### ec2 detached-volumes
```
$ gw-aws-audit ec2 detached-volumes
           VOLUME                 SIZE (GB)  SNAPSHOTS  MIN SIZE (GB)  COSTS
           vol-0cc0f6cd3c99bc1cc          8          0              0  $0.80
   TOTALS  1 VOLUMES                   8 GB          0           0 GB  $0.80
```

#### cw enhanced-monitoring
```
$ gw-aws-audit cw enhanced-monitoring
Enhanced Metrics can add a cost. See: https://aws.amazon.com/cloudwatch/pricing/
Checking for EC2 Enhanced Monitoring

 NAME                                              INSTANCE ID
 master-us-east-1c.masters.us-east-1.gwdocker.com  i-041h4jk12jk23sd
 jumpbox12                                         i-02412sdfgsgdfgs
 analytics-prod                                    i-0a87d921n1rtasd
 EC2 INSTANCES                                     3


Checking for RDS Enhanced Monitoring

 DB INSTANCE                                ENGINE
 airflow-womp-ba-prod-v2                    postgres
 dashboard-reporting-production-instance-1  aurora-postgresql
 service-loloyol-production                 aurora-mysql
 DB INSTANCES                               3
```

#### s3 clear-bucket
```
$ gw-aws-audit s3 clear-bucket --bucket yolo
-- WARNING -- PAY ATTENTION -- FOR REALS --
This will delete ALL objects in yolo
-- THIS ACTION IS NOT REVERSIBLE --
Are you SUPER sure? [yolo]
Enter a value: yolo

Proceeding with batch delete for bucket: yolo
Pages: 198788 Listed: 198788000 Deleted: 198781000 Retries: 18373 DPS: 1921.64
```

#### s3 metrics
```
$ gw-aws-audit s3 metrics > out.csv
Starting metrics pull...
Bucket metric pull complete. Buckets: 207 Processed: 207
```

## Built With

* go v1.14+
* make
* [git-chglog](https://github.com/git-chglog/git-chglog)
* [goreleaser](https://goreleaser.com/install/)

## Deployment

Run `./release.sh $VERSION`

This will update docs, changelog, add the tag, push main and the tag to the repo. The `goreleaser` action will publish the binaries to the Github Release.

If you want to simulate the `goreleaser` process, run the following command:

```
$ curl -sL https://git.io/goreleaser | bash -s -- --rm-dist --skip-publish --snapshot
```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

1. Fork the [GoodwayGroup/gw-aws-audit](https://github.com/GoodwayGroup/gw-aws-audit) repo
1. Use `go >= 1.16`
1. Branch & Code
1. Run linters :broom: `golangci-lint run`
    - The project uses [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
1. Commit with a Conventional Commit
1. Open a PR

## Versioning

We employ [git-chglog](https://github.com/git-chglog/git-chglog) to manage the [CHANGELOG.md](CHANGELOG.md). For the versions available, see the [tags on this repository](https://github.com/GoodwayGroup/gw-aws-audit/tags).

## Authors

* **Derek Smith** - [@clok](https://github.com/clok)

See also the list of [contributors](https://github.com/GoodwayGroup/gwvault/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Sponsors

[![goodwaygroup][goodwaygroup]](https://goodwaygroup.com)

[goodwaygroup]: https://s3.amazonaws.com/gw-crs-assets/goodwaygroup/logos/ggLogo_sm.png "Goodway Group"
