package helpers

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func ValidateOrgUserArgs(c *cli.Context) error {
	if c.String("org") != "" && c.String("user") != "" {
		return fmt.Errorf("cannot pass both organization and user flag")
	}

	if c.String("org") == "" && c.String("user") == "" {
		return fmt.Errorf("must pass either an organization or user value")
	}

	return nil
}
