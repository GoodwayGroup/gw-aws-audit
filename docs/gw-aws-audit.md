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
    - [pem-keys](#pem-keys)
- [sg](#sg)
    - [detached](#detached)
    - [attached](#attached)
    - [cidr](#cidr)
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

### pem-keys

List instances and PEM key used at time of creation

## sg

Security Group related commands

### detached

generate a report of all Security Groups that are NOT attached to an instance

```
This command will scan the EC2 NetworkInterfaces to determine what
Security Groups are NOT attached/assigned in AWS.

```

### attached

generate a report of all Security Groups that are attached to an instance

```
This command will scan the EC2 NetworkInterfaces to determine what
Security Groups are attached/assigned in AWS.
```

### cidr

generate a report comparing SG rules with input CIDR blocks

```
$ gw-aws-audit sg cidr --allowed 10.176.0.0/16,10.175.0.0/16 --alert 174.0.0.0/8,1.2.3.4/32

This command will generate a report detecting the port to CIDR mapping rules 
for attached Security Groups. 

A list of Approved CIDRs is required. This is typically the CIDR block associated
with your VPC.
```

**--alert, -b**="": CIDR blocks that will cause an alert (csv) (default: 174.0.0.0/8)

**--approved, -a**="": CIDR blocks that are approved (csv)

**--ignore-ports, -p**="": Ports that can be ignored (csv) (default: 80,443)

**--warn, -w**="": CIDR blocks that will cause a warning (csv) (default: 204.0.0.0/8)

## cw

CloudWatch related commands

### enhanced-monitoring

Produce report of Enhanced Monitoring enabled EC2 & RDS instances

## install-manpage

Generate and install man page

>NOTE: Windows is not supported

## version, v

Print version info

