# ghlabels
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://github.com/clok/ghlabels/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/clok/ghlabels)](https://goreportcard.com/report/clok/ghlabels)
[![Coverage Status](https://coveralls.io/repos/github/clok/ghlabels/badge.svg)](https://coveralls.io/github/clok/ghlabels)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/clok/ghlabels?tab=overview)

Simple CLI tool to help manage labels across repos

```text
$ bin/ghlabels 
NAME:
   ghlabels - label sync for repos and organizations

USAGE:
   ghlabels [global options] command [command options] [arguments...]

COMMANDS:
   sync             sync labels - delete, rename, update
   dump-defaults    print default labels yaml to STDOUT
   stats            prints out repo stats
   install-manpage  Generate and install man page
   help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Development

1. Fork the [clok/ghlabels](https://github.com/clok/ghlabels) repo
1. Use `go >= 1.16`
1. Branch & Code
1. Run linters :broom: `golangci-lint run`
    - The project uses [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
1. Commit with a Conventional Commit
1. Open a PR