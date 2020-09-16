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
    - [port](#port)
    - [amazon](#amazon)
    - [direct-ip-mapping, dim](#direct-ip-mapping-dim)
- [iam](#iam)
    - [user-report, ur](#user-report-ur)
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

**--all**: Process ALL Security Groups, not just attached

**--approved, -a**="": CIDR blocks that are approved (csv)

**--ignore-ports, -p**="": Ports that can be ignored (csv) (default: 80,443,3,4,3-4)

**--ignore-protocols**="": Protocols to ignore. Can be tcp,udp,icmp (csv)

**--warn, -w**="": CIDR blocks that will cause a warning (csv) (default: 204.0.0.0/8)

### port

generate a report comparing SG rules with input CIDR blocks on a specific port

```
$ gw-aws-audit sg ports --ports 22,3306 --allowed 10.176.0.0/16,10.175.0.0/16 --alert 174.0.0.0/8,1.2.3.4/32

This command will generate a report for a set of PORTS for attached Security Groups.

A list of Approved CIDRs is required. This is typically the CIDR block associated
with your VPC.
```

**--alert, -b**="": CIDR blocks that will cause an alert (csv) (default: 174.0.0.0/8)

**--all**: Process ALL Security Groups, not just attached

**--approved, -a**="": CIDR blocks that are approved (csv)

**--ignore-protocols**="": Protocols to ignore. Can be tcp,udp,icmp (csv)

**--ports, -p**="": Ports to generate report on (csv) (default: 22)

**--warn, -w**="": CIDR blocks that will cause a warning (csv) (default: 204.0.0.0/8)

### amazon

generate a report of allow SG with rules mapped to known AWS IPs

```
This method loads the current version of https://ip-ranges.amazonaws.com/ip-ranges.json
and compares the CIDR blocks against all Security Groups.
```

### direct-ip-mapping, dim

generate report of Security Groups with direct mappings to EC2 instances

```
This method will generate a report comparing all Security Groups with all 
EC2 instances to determine where you have a direct IP mapping.

This will note Internal and External IP usage as well.
```

## iam


### user-report, ur

generates report of IAM Users and Access Key Usage

```
This action will generate a report for all Users within an AWS account with the details
specific user authentication methods.

┌──────────────┬────────┬───────────┬─────────┬────────────┬─────────────────────────────────────────────────────────────────────────┐
│              │        │           │         │            │                           ACCESS KEY DETAILS                            │
│ USER         │ STATUS │       AGE │ CONSOLE │ LAST LOGIN │               KEY ID | STATUS | AGE | LAST USED | SERVICE               │
├──────────────┼────────┼───────────┼─────────┼────────────┼─────────────────────────────────────────────────────────────────────────┤
│ user12345    │   PASS │  123 days │      NO │       NONE │                               0 API Keys                                │
├──────────────┼────────┼───────────┼─────────┼────────────┼─────────────────────────────────────────────────────────────────────────┤
│ bot-user-123 │   WARN │  236 days │      NO │       NONE │                               2 API Keys                                │
│              │        │           │         │            │ AKIAIOSFODNN6EXAMPLE │ Active │ 229 days │   229 days 22 hours   │ s3   │
│              │        │           │         │            │ AKIAIOSFODNN5EXAMPLE │ Active │ 228 days │ 51 minutes 24 seconds │ sts  │
├──────────────┼────────┼───────────┼─────────┼────────────┼─────────────────────────────────────────────────────────────────────────┤
│ userAOK123   │   FAIL │   43 days │     YES │     5 days │                               1 API Key                                 │
│              │        │           │         │            │   AKIAIOSFODNN3EXAMPLE │ Active │ 43 days │ 22 hours 5 minutes │ ec2    │
└──────────────┴────────┴───────────┴─────────┴────────────┴─────────────────────────────────────────────────────────────────────────┘

USER [string]:
  - The user name

STATUS [enum]:
  - PASS: When a does NOT have Console Access and has NO Access Keys
  - FAIL: When a User has Console Access
  - WARN: When a User does NOT have Console Acces, but does have at least 1 Access Key
  - UNKNOWN: Catch all for cases not handled.

AGE [duration]:
  - Time since User was created

CONSOLE [bool]:
  - Does the User have Console Access? YES/NO

LAST LOGIN [duration]:
  - Time since User was created
  - NONE if the User does not have Console Access or if the User has NEVER logged in.

ACCESS KEY DETAILS [sub table]:
  - Primary header row is the number of Access Keys associated with the User

  KEY ID [string]:
    - The AWS_ACCESS_KEY_ID

  STATUS [enum]:
    - Active/Inactive

  LAST USED [duration]:
    - Time since the Access Key was last used.

  SERVICE [string]:
    - The last AWS Service that the Access Key was used to access at the "LAST USED" time.
```

## cw

CloudWatch related commands

### enhanced-monitoring

Produce report of Enhanced Monitoring enabled EC2 & RDS instances

## install-manpage

Generate and install man page

>NOTE: Windows is not supported

## version, v

Print version info

