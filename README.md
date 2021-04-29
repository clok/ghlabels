# ghlabels
[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://github.com/clok/ghlabels/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/clok/ghlabels)](https://goreportcard.com/report/clok/ghlabels)
[![Coverage Status](https://coveralls.io/repos/github/clok/ghlabels/badge.svg)](https://coveralls.io/github/clok/ghlabels)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/clok/ghlabels?tab=overview)

Simple CLI tool to help manage labels across repos

> Please see [the docs for details on the commands.](./docs/ghlabels.md)

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

## Configuration

To view the default configuration, run the `dump-defaults` command.

```text
$ ghlabels dump-defaults
```

The order of operations is:

1. Rename
2. Sync
3. Delete

A custom configuration file can be provided using the `--config, -c` flag. The
file passed to this config option must use the following structure. You do
not need to have all 3 Top Level sections in the config file for the
configuration to be valid.

By default, the `--config, -c` flag will overwrite the default configuration.
You can merge the provided configuration with the default using the 
`--merge-with-defaults, -m` boolean flag. This will take the default configuration
and merge in the user provided configuration, with the user config taking precedence.

### Configuration Schema

| Top Level 	| Type   	| Description                                                                      	| Structure                                                   	|
|-----------	|--------	|----------------------------------------------------------------------------------	|-------------------------------------------------------------	|
| `rename`  	| `List` 	| List of Label names to rename `from` a given name `to` a new name.               	| `{ from: "string", to: "string" }`                          	|
| `remove`  	| `List` 	| List of Labels to be deleted from a Repo.                                        	| `string`                                                    	|
| `sync`    	| `List` 	| List of Label configuration that will be used to create or update a given Label. 	| `{ name: "string", color: "sting", description: "string" }` 	|

## Development

1. Fork the [clok/ghlabels](https://github.com/clok/ghlabels) repo
1. Use `go >= 1.16`
1. Branch & Code
1. Run linters :broom: `golangci-lint run`
    - The project uses [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
1. Commit with a Conventional Commit
1. Open a PR