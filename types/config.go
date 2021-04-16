package types

import "strings"

type Label struct {
	Name string `yaml:"name"`
	Color string `yaml:"color"`
	Description string `yaml:"description"`
}

func (l *Label) SetName(s string) {
	l.Name = s
}

type Rename struct {
	From string `yaml:"from"`
	To string `yaml:"to"`
}

type Config struct {
	Rename []Rename `yaml:"rename"`
	Remove []string `yaml:"remove"`
	Sync []Label `yaml:"sync"`
}

func (c *Config) FindSyncLabel(name string) *Label {
	lower := strings.ToLower(name)
	for _, l := range c.Sync {
		if strings.ToLower(l.Name) == lower {
			return &l
		}
	}
	return nil
}

// MergeLeft takes the input config and overwrites the values of the original
// when there is a a duplicate detected.
func (c *Config) MergeLeft(config Config) {
	// time to merge, always pick User over Default
	if len(config.Rename) != 0 {
		dl := c.Rename
		for _, ul := range config.Rename {
			for i, l := range dl {
				if l.From == ul.From {
					dl = append(dl[:i], dl[i+1:]...)
				}
			}
		}
		c.Rename = append(config.Rename, dl...)
	}

	if len(config.Remove) != 0 {
		dl := c.Remove
		for _, ul := range config.Remove {
			for i, l := range dl {
				if l == ul {
					dl = append(dl[:i], dl[i+1:]...)
				}
			}
		}
		c.Remove = append(config.Remove, dl...)
	}

	if len(config.Sync) != 0 {
		dl := c.Sync
		for _, ul := range config.Sync {
			for i, l := range dl {
				if l.Name == ul.Name {
					dl = append(dl[:i], dl[i+1:]...)
				}
			}
		}
		c.Sync = append(config.Sync, dl...)
	}
}