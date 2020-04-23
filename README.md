# GW AWS Audit Tool
> NOTE: This is not perfect. It is a specialized tool to help with actions to take during an audit of AWS usage.

## Basic Usage

Use in place of `ansible-vault`. All commands are reimplemented. The tool will default to asking for your Vault password.

```
$ gw-aws-audit 
s3-add-cost-tag              Add 's3-cost-name' tag to all buckets
s3-clear-bucket <bucket>     Clear ALL objects from a Bucket
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
$ curl https://i.jpillora.com/GoodwayGroup/gwvault! | bash
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
