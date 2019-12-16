package util

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

// Config is the structure of .mima.yml config file
type Config struct {
	Keys struct {
		Vultr        string `yaml:"vultr"`
		DigitalOcean string `yaml:"digitalocean"`
	} `yaml:"keys"`
	Servers []struct {
		Name     string `yaml:"name"`
		Provider string `yaml:"provider"`
		Plan     string `yaml:"plan"`
		Region   string `yaml:"region"`
	} `yaml:"servers"`
}

// GetConfig will open, parses and return the YAML config file
func GetConfig(fileName string) (Config, error) {
	var cfg Config

	f, err := os.Open(fileName)
	if err != nil {
		return cfg, errors.New("Could not find config file")
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, errors.New("Could not parse YAML config file")
	}

	return cfg, nil
}
