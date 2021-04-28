package helpers

import (
	// Import embed for generic use of go:embed
	_ "embed"
	"io/ioutil"

	"github.com/clok/ghlabels/types"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

//go:embed embeds/defaults.yml
var defaultLabelsBytes []byte

var kdc = k.Extend("DetermineConfig")

func GetDefaultConfig() string {
	return string(defaultLabelsBytes)
}

func DetermineConfig(c *cli.Context) (*types.Config, error) {
	var defaultLabels types.Config
	err := yaml.Unmarshal(defaultLabelsBytes, &defaultLabels)
	if err != nil {
		return nil, err
	}

	var userConfig types.Config
	var config types.Config
	if c.String("config") != "" {
		yamlFile, err := ioutil.ReadFile(c.String("config"))
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(yamlFile, &userConfig)
		if err != nil {
			return nil, err
		}
		if c.Bool("merge-with-defaults") {
			defaultLabels.MergeLeft(userConfig)
			config = defaultLabels
		} else {
			config = userConfig
		}
	} else {
		config = defaultLabels
	}
	kdc.Log(config)

	return &config, nil
}
