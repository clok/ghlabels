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

sync labels - delete, rename, update

### all

Sync labels across ALL repos within an org or for a user

**--config, -c**="": Path to config file withs labels to sync.

**--merge-with-defaults, -m**: Merge provided config with defaults, otherwise only use the provided config.

**--org, -o**="": GitHub Organization to view. Cannot be used with User flag.

**--user, -u**="": GitHub User to view. Cannot be used with Organization flag.

### repo

Sync labels for a single repo

**--config, -c**="": Path to config file withs labels to sync.

**--merge-with-defaults, -m**: Merge provided config with defaults, otherwise only use the provided config.

**--repo, -r**="": Repo name including owner. Examlple: clok/ghlabels

## dump-defaults

print default labels yaml to STDOUT

## stats

prints out repo stats

**--org, -o**="": GitHub Organization to view. Cannot be used with User flag.

**--user, -u**="": GitHub User to view. Cannot be used with Organization flag.

## install-manpage

Generate and install man page

>NOTE: Windows is not supported

