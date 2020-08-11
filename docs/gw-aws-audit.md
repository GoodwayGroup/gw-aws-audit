% gw-aws-audit 8
# NAME
gw-aws-audit - a collection of tools to audit AWS.
# SYNOPSIS
gw-aws-audit


# COMMAND TREE

- [s3](#s3)
    - [add-cost-tag](#add-cost-tag)
    - [metrics](#metrics)
    - [clear-bucket, exterminatus](#clear-bucket-exterminatus)
- [rds](#rds)
    - [enhanced-monitoring](#enhanced-monitoring)
- [ec2](#ec2)
    - [enhanced-monitoring](#enhanced-monitoring)
    - [detached-volumes](#detached-volumes)
    - [stopped-hosts](#stopped-hosts)
- [cw](#cw)
    - [enhanced-monitoring](#enhanced-monitoring)
- [install-manpage](#install-manpage)
- [version, v](#version-v)

**Usage**:
```
gw-aws-audit [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# COMMANDS

## s3

S3 related commands

### add-cost-tag

Add s3-cost-name to all S3 buckets

```
Idempotent action that will add the `s3-cost-name` tag to ALL S3 buckets for a
given account.

The value will be the Bucket name.
```

### metrics

Get usage metrics

```
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
```

### clear-bucket, exterminatus

Clear all Objects within a given Bucket

```
Efficiently delete all objects within a bucket.

This process will run multiple paged deletes in parallel. It will handle API
throttling from AWS with an exponential backoff with retry. 
```

**--bucket, -b**="": Bucket to clear

## rds

RDS related commands

### enhanced-monitoring

Produce report of Enhanced Monitoring enabled instances

## ec2

EC2 related commands

### enhanced-monitoring

Produce report of Enhanced Monitoring enabled instances

### detached-volumes

List detached EBS volumes and snapshot counts

### stopped-hosts

List stopped EC2 hosts and associated EBS volumes

## cw

CloudWatch related commands

### enhanced-monitoring

Produce report of Enhanced Monitoring enabled EC2 & RDS instances

## install-manpage

Generate and install man page

## version, v

Print version info

