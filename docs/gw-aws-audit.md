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
    - [public](#public)
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
    - [report](#report)
    - [modify](#modify)
    - [permissions, p](#permissions-p)
    - [keys](#keys)
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

### public

Produce report of instances that have public interfaces attached

```
Produces a report that displays a list RDS servers that are configured as Publicly Accessible.

The report contains:

DB INSTANCE:
    - Name of the instance

ENGINE:
    - RDS DB engine

SECURITY GROUPS:
    - Security Group ID
    - Security Group Name
    - Inbound Port
    - CIDR rules applied to the Port
```

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

IAM related commands

### report

generates report of IAM Users and Access Key Usage

```
This action will generate a report for all Users within an AWS account with the details
specific user authentication methods.

Interactive mode will allow you to search for an IAM User and take actions once an IAM User is
selected.

USER [string]:
  - The user name

STATUS [enum]:
  - PASS: When a does NOT have Console Access and has NO Access Keys or only INACTIVE Access Keys
  - FAIL: When an IAM User has Console Access
  - WARN: When an IAM User does NOT have Console Access, but does have at least 1 ACTIVE Access Key
  - UNKNOWN: Catch all for cases not handled.

AGE [duration]:
  - Time since User was created

CONSOLE [bool]:
  - Does the User have Console Access? YES/NO

LAST LOGIN [duration]:
  - Time since User was created
  - NONE if the User does not have Console Access or if the User has NEVER logged in.

PERMISSIONS [struct]:
  - G: n -> Groups that the User belongs to
  - P: n -> Policies that are attached to the User

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

**--interactive, -i**: after generating the report, prompt for digging into a user

**--show-only**="": filter results to show only pass, warn or fail

### modify

modify an IAM User within AWS

```
This action allows you to take actions to modify a user's Permissions (Groups and Policies)
and the state of their Access Keys (Active, Inactive, Delete).
```

**--show-only**="": filter results to show only pass, warn or fail

**--user, -u**="": user name to look for

### permissions, p

view permissions that are associated with an IAM User

```
Produces a table of Groups and Policies that are attached to an IAM User.

Interactive mode allows for you to detach a permission from an IAM User.
```

**--interactive, -i**: interactive mode that allows for removal of permissions

**--user, -u**="": user name to look for

### keys

view Access Keys associated with an IAM User

```
Produces a table of Access Keys that are associated with an IAM User.

Interactive mode allows for you to Activate, Deactivate and Delete Access Keys.
```

**--interactive, -i**: interactive mode that allows status changes of keys

**--user, -u**="": user name to look for

## cw

CloudWatch related commands

### enhanced-monitoring

Produce report of Enhanced Monitoring enabled EC2 & RDS instances

## install-manpage

Generate and install man page

>NOTE: Windows is not supported

## version, v

Print version info

