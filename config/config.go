package config

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
		Providor string `yaml:"providor"`
		Plan     string `yaml:"plan"`
	} `yaml:"servers"`
}

// Get will open, parses and return the YAML config file
func Get(fileName string) (Config, error) {
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
