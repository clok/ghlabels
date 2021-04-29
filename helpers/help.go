package helpers

var customConfigText = `
To view the default configuration, run the ` + "`dump-defaults`" + ` command.

$ ghlabels dump-defaults

The order of operations is:

1. Rename
2. Sync
3. Delete

A custom configuration file can be provided using the ` + "`--config, -c`" + ` flag. The
file passed to this config option must use the following structure. You do
not need to have all 3 Top Level sections in the config file for the
configuration to be valid.

By default, the ` + "`--config, -c`" + ` flag will overwrite the default configuration.
You can merge the provided configuration with the default using the 
` + "`--merge-with-defaults, -m`" + ` boolean flag. This will take the default configuration
and merge in the user provided configuration, with the user config taking precedence.

### Configuration Schema

| Top Level 	| Type   	| Description                                                                      	| Structure                                                   	|
|-----------	|--------	|----------------------------------------------------------------------------------	|-------------------------------------------------------------	|
| ` + "`rename`" + `  	| ` + "`List`" + ` 	| List of Label names to rename ` + "`from`" + ` a given name ` + "`to`" + ` a new name.               	| ` + "`{ from: \"string\", to: \"string\" }`" + `                          	|
| ` + "`remove`" + `  	| ` + "`List`" + ` 	| List of Labels to be deleted from a Repo.                                        	| ` + "`string`" + `                                                    	|
| ` + "`sync`" + `    	| ` + "`List`" + ` 	| List of Label configuration that will be used to create or update a given Label. 	| ` + "`{ name: \"string\", color: \"sting\", description: \"string\" }`" + ` 	|
`

var SyncAllUsageText = `
This process will pull all repositories associated with a User or Organization
and perform the sync label action on all qualifying repositories.

A repository qualifies for label sync if is is NOT archived and NOT a fork.
` + customConfigText

var SyncRepoUsageText = `
This process will perform the sync label action on all a single repository
provided via the ` + "`--repo, r`" + ` flag.
` + customConfigText
