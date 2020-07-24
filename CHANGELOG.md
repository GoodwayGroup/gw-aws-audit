<a name="unreleased"></a>
## [Unreleased]


<a name="v1.3.2"></a>
## [v1.3.2] - 2020-07-24
### Chore
- updated release process to auto push branch and tag


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
### Features
- **release:** v1.2.0

### Pull Requests
- fix(metrics): fix bug in order of columns of s3 metrics report and with the count of objects for a bucket ([#3](https://github.com/GoodwayGroup/gw-aws-audit/issues/3))




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


###### Squashed Commits:
```
v1.0.0
```



<a name="v1.0.0"></a>
## [v1.0.0] - 2020-05-14
### Features
- add failed tracking
- **s3:** Add support for all regions when adding s3-cost-name

### Pull Requests
- feat(cli): port from go-commander to urfave/cli/v2 ([#1](https://github.com/GoodwayGroup/gw-aws-audit/issues/1))


###### Squashed Commits:
```
* feat(cli): port from go-commander to urfave/cli/v2

* chore: updated readme
```



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


[Unreleased]: https://github.com/GoodwayGroup/gw-aws-audit/compare/v1.3.2...HEAD
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
