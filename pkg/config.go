package pkg

import (
	"log"
	"strings"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type Config struct {
	Enable              []string `hcl:"enable,optional"`
	SkipShellCompletion []string `hcl:"skipShellCompletion,optional"`
}

func ReadConfig(fileName string) (*Config, error) {
	var config struct {
		Config *Config `hcl:"Config,block"`
	}
	err := hclsimple.DecodeFile(fileName, nil, &config)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
	return config.Config, err

}
func (c *Config) FilterShellCompletions(moduleNames ...string) []string {
	if c == nil {
		return moduleNames
	}
	var filtered []string
	for _, s := range moduleNames {
		ok := false
		for _, s2 := range c.SkipShellCompletion {
			if strings.Contains(s, s2) {
				ok = true
			}
		}
		if !ok {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
