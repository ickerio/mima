package providors

import (
	"errors"

	"github.com/ickerio/mima/config"
	"github.com/vultr/govultr"
)

// Providor is the interface for the Vultr, DigitalOcean structs
type Providor interface {
	Info() (Server, error)
}

// Server details the information of a VPS server
type Server struct {
	Name             string
	Os               string
	Memory           string
	Storage          string
	CPUCount         string
	IP               string
	CurrentBandwidth float64
	AllowedBandwidth string
	Location         string
	Cost             string
	Created          string
	Password         string
}

// Get returns a providor from the config and user input
func Get(config config.Config, name string) (Providor, error) {
	for i := range config.Servers {
		if config.Servers[i].Name == name {
			switch config.Servers[i].Providor {
			case "Vultr":
				v := govultr.NewClient(nil, config.Keys.Vultr)
				return Vultr{apiKey: config.Keys.Vultr, client: v, name: name}, nil
			case "DigitalOcean":
				v := govultr.NewClient(nil, config.Keys.DigitalOcean)
				return Vultr{apiKey: config.Keys.DigitalOcean, client: v, name: name}, nil
			default:
				return nil, errors.New("Invalid providor in config (Must be: Vultr or DigitalOcean)")
			}
		}
	}

	if name == "" {
		return nil, errors.New("Please enter an server name")
	}
	return nil, errors.New("Server name did not match any in config file")
}
