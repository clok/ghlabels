<a name="unreleased"></a>
## [Unreleased]


<a name="v1.1.0"></a>
## [v1.1.0] - 2021-05-24
### Chore
- create FUNDING.yml
- **lint:** address lint issues

### Feat
- **defaults:** added default security label

### Fix
- **deps:** update golang.org/x/oauth2 commit hash to f6687ab ([#5](https://github.com/clok/ghlabels/issues/5))


<a name="v1.0.0"></a>
## [v1.0.0] - 2021-05-13
### Chore
- update readme
- **docs:** updating docs for version v1.0.0
- **release:** v1.0.0

### Feat
- **config:** added config generation from a repo


<a name="v0.1.3"></a>
## [v0.1.3] - 2021-05-13
### Chore
- **deps:** update alpine docker tag to v3.13.5 ([#2](https://github.com/clok/ghlabels/issues/2))
- **release:** v0.1.3
- **renovate:** use defaults from clok/renovate-condig

### Fix
- **deps:** update golang.org/x/oauth2 commit hash to 81ed05c ([#1](https://github.com/clok/ghlabels/issues/1))
- **deps:** update all non-major dependencies ([#4](https://github.com/clok/ghlabels/issues/4))


<a name="v0.1.2"></a>
## [v0.1.2] - 2021-04-29
### Chore
- **docs:** updating docs for version v0.1.2
- **release:** v0.1.2
- **release:** updated release process to include release script


<a name="v0.1.1"></a>
## [v0.1.1] - 2021-04-29
### Chore
- update changelog for v0.1.1
- **docs:** update readme and command help text
- **release:** added homebrew installer


<a name="v0.1.0"></a>
## v0.1.0 - 2021-04-28
### Chore
- update changelog for v0.1.0
- generate docs
- update README
- lint
- **CHANGELOG:** generate fresh changelog
- **CHANGELOG:** added git-chglog config
- **actions:** add github actions for CI/CD
- **release:** update Makefile release command

### Feat
- **goreleaser:** add Dockerfile
- **sync:** support all and repo, with rename on labels that already are assigned
- **sync:** added sync repo and all commands
- **tracer bullet:** initial work with embedded defaults

### Fix
- **test:** update test command

### Refactor
- **lint:** updated linting rules and refactored commands
- **sync:** moved to functional pattern


[Unreleased]: https://github.com/clok/ghlabels/compare/v1.1.0...HEAD
[v1.1.0]: https://github.com/clok/ghlabels/compare/v1.0.0...v1.1.0
[v1.0.0]: https://github.com/clok/ghlabels/compare/v0.1.3...v1.0.0
[v0.1.3]: https://github.com/clok/ghlabels/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/clok/ghlabels/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/clok/ghlabels/compare/v0.1.0...v0.1.1
