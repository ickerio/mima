package providers

import (
	"errors"

	"github.com/ickerio/mima/util"
	"github.com/vultr/govultr"
)

// Provider is the interface for the Vultr, DigitalOcean structs
type Provider interface {
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

// Get returns a provider from the config and user input
func Get(conf util.Config, name string) (Provider, error) {
	for i := range conf.Servers {
		if conf.Servers[i].Name == name {
			switch conf.Servers[i].Provider {
			case "Vultr":
				v := govultr.NewClient(nil, conf.Keys.Vultr)
				return Vultr{apiKey: conf.Keys.Vultr, client: v, name: name}, nil
			case "DigitalOcean":
				v := govultr.NewClient(nil, conf.Keys.DigitalOcean)
				return Vultr{apiKey: conf.Keys.DigitalOcean, client: v, name: name}, nil
			default:
				return nil, errors.New("Invalid provider in config (Must be: Vultr or DigitalOcean)")
			}
		}
	}

	if name == "" {
		return nil, errors.New("Please enter an server name")
	}
	return nil, errors.New("Server name did not match any in config file")
}
