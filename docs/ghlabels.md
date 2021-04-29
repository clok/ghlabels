% ghlabels 8
# NAME
ghlabels - label sync for repos and organizations
# SYNOPSIS
ghlabels


# COMMAND TREE

- [sync](#sync)
    - [all](#all)
    - [repo](#repo)
- [dump-defaults](#dump-defaults)
- [stats](#stats)
- [install-manpage](#install-manpage)

**Usage**:
```
ghlabels [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# COMMANDS

## sync

sync labels - rename, sync, delete

### all

Sync labels across ALL repos within an org or for a user

```
This process will pull all repositories associated with a User or Organization
and perform the sync label action on all qualifying repositories.

A repository qualifies for label sync if is is NOT archived and NOT a fork.

To view the default configuration, run the `dump-defaults` command.

$ ghlabels dump-defaults

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
```

**--config, -c**="": Path to config file withs labels to sync.

**--merge-with-defaults, -m**: Merge provided config with defaults, otherwise only use the provided config.

**--org, -o**="": GitHub Organization to view. Cannot be used with User flag.

**--user, -u**="": GitHub User to view. Cannot be used with Organization flag.

### repo

Sync labels for a single repo

```
This process will perform the sync label action on all a single repository
provided via the `--repo, r` flag.

To view the default configuration, run the `dump-defaults` command.

$ ghlabels dump-defaults

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
```

**--config, -c**="": Path to config file withs labels to sync.

**--merge-with-defaults, -m**: Merge provided config with defaults, otherwise only use the provided config.

**--repo, -r**="": Repo name including owner. Examlple: clok/ghlabels

## dump-defaults

print default labels yaml to STDOUT

## stats

prints out repository stats

**--org, -o**="": GitHub Organization to view. Cannot be used with User flag.

**--user, -u**="": GitHub User to view. Cannot be used with Organization flag.

## install-manpage

Generate and install man page

>NOTE: Windows is not supported

