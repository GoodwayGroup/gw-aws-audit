<a name="unreleased"></a>
## [Unreleased]


<a name="v1.18.0"></a>
## [v1.18.0] - 2021-03-30
### Chore
- regenerate CHANGELOG.md
- update changelog config
- **docs:** updating docs for version v1.18.0

### Code Refactoring
- **aws client:** centralize client session management to reduce the number of instantiations of an AWS Session
- **iam:** cleanup render commands

### Features
- **iam:** add key deactivation bulk action


<a name="v1.17.0"></a>
## [v1.17.0] - 2021-03-26
### Bug Fixes
- **deps:** update module github.com/aws/aws-sdk-go to v1.38.6 ([#47](https://github.com/GoodwayGroup/gw-aws-audit/issues/47))

### Chore
- update README
- **docs:** updating docs for version v1.17.0

### Features
- **iam:** adding interactive management of users to the tool ([#48](https://github.com/GoodwayGroup/gw-aws-audit/issues/48))
- **release:** v1.17.0

### BREAKING CHANGE

The command `gw-aws-audit iam user-report` has been renamed to `gw-aws-audit iam report`


<a name="v1.16.0"></a>
## [v1.16.0] - 2021-03-22
### Chore
- **docs:** updating docs for version v1.16.0

### Features
- **iam:** user group and policy management
- **release:** v1.16.0


<a name="v1.15.2"></a>
## [v1.15.2] - 2021-03-17
### Bug Fixes
- **deps:** update github.com/hako/durafmt commit hash to 3a2c319 ([#45](https://github.com/GoodwayGroup/gw-aws-audit/issues/45))
- **deps:** update module github.com/aws/aws-sdk-go to v1.37.33 ([#46](https://github.com/GoodwayGroup/gw-aws-audit/issues/46))
- **goreleaser:** addressed type in binary naming

### Features
- **release:** v1.15.2


<a name="v1.15.1"></a>
## [v1.15.1] - 2021-03-17
### Chore
- **release:** update homebrew push

### Features
- **release:** v1.15.1


<a name="v1.15.0"></a>
## [v1.15.0] - 2021-03-15
### Chore
- **ci:** remove deprecated release commands
- **ci:** port to using golangci-lint and goreleaser github action
- **deps:** update kemba, awessession, cdocs, go-funk and aws-sdf-go
- **go.mod:** bump to go 1.16

### Features
- **release:** v1.15.0


<a name="v1.14.3"></a>
## [v1.14.3] - 2021-03-03
### Bug Fixes
- **deps:** update all non-major dependencies ([#40](https://github.com/GoodwayGroup/gw-aws-audit/issues/40))

### Chore
- **deps:** update module stretchr/testify to v1.7.0 ([#32](https://github.com/GoodwayGroup/gw-aws-audit/issues/32))
- **deps:** update module aws/aws-sdk-go to v1.36.28 ([#31](https://github.com/GoodwayGroup/gw-aws-audit/issues/31))
- **deps:** update module aws/aws-sdk-go to v1.36.18 ([#30](https://github.com/GoodwayGroup/gw-aws-audit/issues/30))
- **github actions:** add go proxy warming
- **renovate:** add extension for group:allNonMajor
- **renovate:** add gomodTidy option

### Features
- **release:** v1.14.3


<a name="v1.14.2"></a>
## [v1.14.2] - 2020-12-29
### Chore
- **deps:** update module aws/aws-sdk-go to v1.36.17
- **deps:** update module urfave/cli/v2 to v2.3.0
- **deps:** update awssession to v0.1.5 and cdocs to v0.2.3
- **deps:** update module aws/aws-sdk-go to v1.36.13 ([#27](https://github.com/GoodwayGroup/gw-aws-audit/issues/27))
- **deps:** update actions/setup-go action to v2 ([#26](https://github.com/GoodwayGroup/gw-aws-audit/issues/26))
- **deps:** update actions/checkout action to v2 ([#25](https://github.com/GoodwayGroup/gw-aws-audit/issues/25))
- **deps:** update module aws/aws-sdk-go to v1.35.17 ([#24](https://github.com/GoodwayGroup/gw-aws-audit/issues/24))

### Features
- **release:** v1.14.2


<a name="v1.14.1"></a>
## [v1.14.1] - 2020-10-13
### Chore
- **deps:** updated aws/aws-sdk-go, cenkalti/backoff/v4 and jedib0t/go-pretty/v6

### Features
- **release:** v1.14.1


<a name="v1.14.0"></a>
## [v1.14.0] - 2020-09-30
### Chore
- **deps:** update module aws/aws-sdk-go to v1.35.0 ([#19](https://github.com/GoodwayGroup/gw-aws-audit/issues/19))
- **docs:** updating docs for version v1.14.0

### Features
- **rds:** add public command to generate report of publicly exposed RDS instances ([#20](https://github.com/GoodwayGroup/gw-aws-audit/issues/20))
- **release:** v1.14.0


<a name="v1.13.1"></a>
## [v1.13.1] - 2020-09-24
### Chore
- **docs:** updating docs for version v1.13.1
- **tests:** update tests for types

### Features
- **iam:** updated check status pass conditions for user-report and added summary output
- **release:** v1.13.1


<a name="v1.13.0"></a>
## [v1.13.0] - 2020-09-15
### Chore
- **deps:** update module aws/aws-sdk-go to v1.34.24 ([#17](https://github.com/GoodwayGroup/gw-aws-audit/issues/17))
- **docs:** updating docs for version v1.13.0

### Features
- **iam:** added user-report action to report on console access and api keys ([#18](https://github.com/GoodwayGroup/gw-aws-audit/issues/18))
- **release:** v1.13.0


<a name="v1.12.0"></a>
## [v1.12.0] - 2020-09-11
### Features
- **release:** v1.12.0
- **sg dim:** added network interfaces to the direct IP mapping


<a name="v1.11.0"></a>
## [v1.11.0] - 2020-09-03
### Chore
- **docs:** updating docs for version v1.11.0

### Features
- **release:** v1.11.0
- **sg:** added direct-ip-mapping command


<a name="v1.10.0"></a>
## [v1.10.0] - 2020-09-01
### Chore
- **docs:** updating docs for version v1.10.0

### Features
- **protocol:** add flags to filter out specific protocols
- **release:** v1.10.0
- **sg:** add mapping of known AWS IP blocks to reports


<a name="v1.9.1"></a>
## [v1.9.1] - 2020-08-31
### Bug Fixes
- addressed bug that was causing SGs to be skipped in processing

### Chore
- **deps:** update module aws/aws-sdk-go to v1.34.14

### Features
- **release:** v1.9.1


<a name="v1.9.0"></a>
## [v1.9.0] - 2020-08-31
### Chore
- add github lint action
- **deps:** update module aws/aws-sdk-go to v1.34.13
- **docs:** updating docs for version v1.9.0
- **lint:** clean up some formatting concerns

### Features
- **release:** v1.9.0
- **sg:** added port command to generate report of exposure of a port across Security Groups


<a name="v1.8.1"></a>
## [v1.8.1] - 2020-08-27
### Chore
- **docs:** updating docs for version v1.8.1
- **sg:** add port 3 to ignored ports list

### Features
- **release:** v1.8.1


<a name="v1.8.0"></a>
## [v1.8.0] - 2020-08-26
### Chore
- **docs:** updating docs for version v1.8.0

### Features
- **ec2:** added task to generate a report of pem key usage
- **release:** v1.8.0
- **sg:** CIDR and port reporting for attached security groups


<a name="v1.7.0"></a>
## [v1.7.0] - 2020-08-26
### Chore
- linting

### Features
- **helper:** added audit.sh helper tool
- **release:** v1.7.0
- **sg:** add security group summary metrics


<a name="v1.6.2"></a>
## [v1.6.2] - 2020-08-25
### Bug Fixes
- **release:** fix typo in name of release file

### Chore
- **docs:** updating docs for version v1.6.2

### Features
- **release:** v1.6.2


<a name="v1.6.1"></a>
## [v1.6.1] - 2020-08-25
### Chore
- updated release script to include publish to github

### Features
- **release:** v1.6.1


<a name="v1.6.0"></a>
## [v1.6.0] - 2020-08-25
### Chore
- **deps:** bump version of kemba and aws-sdk-go
- **make:** don't update go.mod with gox

### Features
- **release:** v1.6.0
- **sg:** added Security Group attach/detach reports


<a name="v1.5.1"></a>
## [v1.5.1] - 2020-08-21
### Chore
- **deps:** udpate clok/kemba, clok/awssession, clok/cdocs, jedib0t/go-pretty/v6 and aws/aws-sdk-go
- **deps:** add renovate.json
- **deps:** update module aws/aws-sdk-go to v1.34.8 ([#6](https://github.com/GoodwayGroup/gw-aws-audit/issues/6))
- **deps:** update module thoas/go-funk to v0.7.0 ([#8](https://github.com/GoodwayGroup/gw-aws-audit/issues/8))
- **deps:** update module jedib0t/go-pretty to v6 ([#9](https://github.com/GoodwayGroup/gw-aws-audit/issues/9))
- **deps:** update module clok/awssession to v0.1.4 ([#7](https://github.com/GoodwayGroup/gw-aws-audit/issues/7))
- **renovate:** add renovate.json
- **renovate:** clean up dupe config

### Features
- **release:** v1.5.1


<a name="v1.5.0"></a>
## [v1.5.0] - 2020-08-13
### Chore
- **docs:** updating docs for version v1.4.0

### Features
- **release:** v1.5.0

### Fest
- **cdocs:** integrate cdocs library


<a name="v1.4.0"></a>
## [v1.4.0] - 2020-08-11
### Chore
- **docs:** update docs out put and manpage generation

### Features
- **release:** v1.4.0


<a name="v1.3.3"></a>
## [v1.3.3] - 2020-08-04
### Chore
- update dependencies

### Features
- **release:** v1.3.3


<a name="v1.3.2"></a>
## [v1.3.2] - 2020-07-24
### Chore
- updated release process to auto push branch and tag

### Features
- **release:** v1.3.2


<a name="v1.3.1"></a>
## [v1.3.1] - 2020-07-24
### Chore
- update release script

### Features
- **release:** v1.3.1
- **release:** v1.3.0

### Tech Debt
- **logging:** improved logging using kemba


<a name="v1.3.0"></a>
## [v1.3.0] - 2020-07-13
### Bug Fixes
- **metrics:** fix bug in order of columns of s3 metrics report and with the count of objects for a bucket ([#3](https://github.com/GoodwayGroup/gw-aws-audit/issues/3))

### Features
- **release:** v1.2.0


<a name="v1.2.0"></a>
## [v1.2.0] - 2020-06-03
### Chore
- **readme:** update README
- **table:** use more agreeable table rendering style
- **version:** bump version to v1.2.0
- **version:** bump version

### Tech Debt
- **changelog:** switch to using git-ghglog
- **release:** update realease.sh script


<a name="v1.1.0"></a>
## [v1.1.0] - 2020-05-18
### Chore
- update readme

### Features
- **region:** support passing region as flag or ENV var

### Pull Requests
- Merge pull request [#2](https://github.com/GoodwayGroup/gw-aws-audit/issues/2) from GoodwayGroup/release/v1.0.0


<a name="v1.0.0"></a>
## [v1.0.0] - 2020-05-14
### Features
- add failed tracking
- **cli:** port from go-commander to urfave/cli/v2 ([#1](https://github.com/GoodwayGroup/gw-aws-audit/issues/1))
- **s3:** Add support for all regions when adding s3-cost-name


<a name="0.5.0"></a>
## [0.5.0] - 2020-04-29
### Tech Debt
- refactor naming conventions to follow best practices


<a name="0.4.0"></a>
## [0.4.0] - 2020-04-28
### Bug Fixes
- tune wait group count to maximize in VPC call rate

### Chore
- bump version to 0.5.0

### Features
- added table render and enhanced monitoring listing
- added s3-bucket-metrics command


<a name="0.3.3"></a>
## [0.3.3] - 2020-04-24
### Bug Fixes
- addressed race condition in wait group

### Chore
- bump version


<a name="0.3.2"></a>
## [0.3.2] - 2020-04-24
### Bug Fixes
- scale back workers and update messaging

### Chore
- bump version


<a name="0.3.1"></a>
## [0.3.1] - 2020-04-24
### Bug Fixes
- lower worker count


<a name="0.3.0"></a>
## [0.3.0] - 2020-04-24
### Features
- **retries:** added retry with backoff and limited workers to 20


<a name="0.2.0"></a>
## [0.2.0] - 2020-04-23
### Features
- **yolo:** added prompt before performing Exterminatus on bucket
- **yolo:** added prompt before performing Exterminatus on bucket


<a name="0.1.0"></a>
## 0.1.0 - 2020-04-23
### Features
- **cli:** initial CLI implementation


[Unreleased]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.18.0...HEAD
[v1.18.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.17.0...v1.18.0
[v1.17.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.16.0...v1.17.0
[v1.16.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.15.2...v1.16.0
[v1.15.2]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.15.1...v1.15.2
[v1.15.1]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.15.0...v1.15.1
[v1.15.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.14.3...v1.15.0
[v1.14.3]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.14.2...v1.14.3
[v1.14.2]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.14.1...v1.14.2
[v1.14.1]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.14.0...v1.14.1
[v1.14.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.13.1...v1.14.0
[v1.13.1]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.13.0...v1.13.1
[v1.13.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.12.0...v1.13.0
[v1.12.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.11.0...v1.12.0
[v1.11.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.10.0...v1.11.0
[v1.10.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.9.1...v1.10.0
[v1.9.1]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.9.0...v1.9.1
[v1.9.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.8.1...v1.9.0
[v1.8.1]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.8.0...v1.8.1
[v1.8.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.7.0...v1.8.0
[v1.7.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.6.2...v1.7.0
[v1.6.2]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.6.1...v1.6.2
[v1.6.1]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.6.0...v1.6.1
[v1.6.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.5.1...v1.6.0
[v1.5.1]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.5.0...v1.5.1
[v1.5.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.4.0...v1.5.0
[v1.4.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.3.3...v1.4.0
[v1.3.3]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.3.2...v1.3.3
[v1.3.2]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.3.1...v1.3.2
[v1.3.1]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.3.0...v1.3.1
[v1.3.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.2.0...v1.3.0
[v1.2.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/0.5.0...v1.0.0
[0.5.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/0.4.0...0.5.0
[0.4.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/0.3.3...0.4.0
[0.3.3]: https://github.com/GoodwayGroup/gw-aws-audit/compare/0.3.2...0.3.3
[0.3.2]: https://github.com/GoodwayGroup/gw-aws-audit/compare/0.3.1...0.3.2
[0.3.1]: https://github.com/GoodwayGroup/gw-aws-audit/compare/0.3.0...0.3.1
[0.3.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/0.2.0...0.3.0
[0.2.0]: https://github.com/GoodwayGroup/gw-aws-audit/compare/0.1.0...0.2.0
