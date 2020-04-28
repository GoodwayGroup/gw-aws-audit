# GW AWS Audit Tool
> NOTE: This is not perfect. It is a specialized tool to help with actions to take during an audit of AWS usage.

## Basic Usage

Useful for clearing large S3 buckets, identifying egress EBS volumes and tracking S3 spend.

```
$ gw-aws-audit 
s3-add-cost-tag              Add 's3-cost-name' tag to all buckets
s3-clear-bucket <bucket>     Clear ALL objects from a Bucket
s3-bucket-metrics            Print out bucket metrics to STDOUT
monitoring-enabled           List EC2 & RDS hosts with Enhanced Monitoring enabled
version                      Print current version of this application
ec2-list-stopped-hosts       List stopped EC2 hosts and associated EBS volumes
ec2-list-detached-volumes    List detached EBS volumes and snapshot counts
help [command]               Display this help or a command specific help
```

```
$ gw-aws-audit help s3-clear-bucket
Usage: gw-aws-audit s3-clear-bucket <bucket>

Examples:
  gw-aws-audit s3-clear-bucket athena-results-ASDF1337
```

## Installation

```
$ curl https://i.jpillora.com/GoodwayGroup/gw-aws-audit! | bash
```

### Example Outputs

#### ec2-list-stopped-hosts
```
$ gw-aws-audit ec2-list-stopped-hosts
                INSTANCE ID          NAME            VOLUME                 SIZE (GB)  SNAPSHOTS  MIN SIZE (GB)  COSTS
                i-09e42474f22039e23  dummy-box-test
                                                     vol-0d4b4a7bc95a4b8e4          8          0              0  $0.80
                                                     vol-0cc0f6cd3c99bc1cc          8          0              0  $0.80
                TOTALS               1 INSTANCES     2 VOLUMES                  16 GB          0           0 GB  $1.60
```

#### ec2-list-detached-volumes
```
$ gw-aws-audit ec2-list-detached-volumes
           VOLUME                 SIZE (GB)  SNAPSHOTS  MIN SIZE (GB)  COSTS
           vol-0cc0f6cd3c99bc1cc          8          0              0  $0.80
   TOTALS  1 VOLUMES                   8 GB          0           0 GB  $0.80
```

#### monitoring-enabled
```
$ gw-aws-audit monitoring-enabled
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

#### s3-clear-bucket
```
$ gw-aws-audit s3-clear-bucket yolo
-- WARNING -- PAY ATTENTION -- FOR REALS --
This will delete ALL objects in yolo
-- THIS ACTION IS NOT REVERSIBLE --
Are you SUPER sure? [yolo]
Enter a value: yolo

Proceeding with batch delete for bucket: yolo
Pages: 198788 Listed: 198788000 Deleted: 198781000 Retries: 18373 DPS: 1921.64
```

#### s3-bucket-metrics
```
$ gw-aws-audit s3-bucket-metrics > out.csv
Starting metrics pull...
Bucket metric pull complete. Buckets: 207 Processed: 207
```

## Built With

* go v1.14+
* make
* [github.com/mitchellh/gox](https://github.com/mitchellh/gox)

## Deployment

Run `./release.sh $VERSION`

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We employ [auto-changelog](https://www.npmjs.com/package/auto-changelog) to manage the [CHANGELOG.md](CHANGELOG.md). For the versions available, see the [tags on this repository](https://github.com/GoodwayGroup/gwvault/tags).

## Authors

* **Derek Smith** - [@clok](https://github.com/clok)

See also the list of [contributors](https://github.com/GoodwayGroup/gwvault/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details

## Sponsors

[![goodwaygroup][goodwaygroup]](https://goodwaygroup.com)

[goodwaygroup]: https://s3.amazonaws.com/gw-crs-assets/goodwaygroup/logos/ggLogo_sm.png "Goodway Group"
