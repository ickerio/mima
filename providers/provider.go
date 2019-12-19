package providers

import (
	"errors"

	"github.com/ickerio/mima/parsers"
	"github.com/vultr/govultr"
)

// Provider is the interface for the Vultr, DigitalOcean structs
type Provider interface {
	Info() (Server, error)
	Start() error
	Stop() error
	Plans() ([]Plan, error)
	Regions() ([]Region, error)
	OS() ([]OS, error)
}

// Server details the information of a VPS server
type Server struct {
	ID               string
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

// Plan details the information of a VPS plan
type Plan struct {
	ID          int
	Description string
}

// Region details the information of a VPS region
type Region struct {
	ID   int
	Name string
}

// OS details the information of a VPS operating system
type OS struct {
	ID   int
	Name string
}

// GetNoAuth returns a provider without a key
func GetNoAuth(service string) (Provider, error) {
	switch service {
	case "Vultr":
		v := govultr.NewClient(nil, "")
		return Vultr{client: v}, nil
	case "DigitalOcean":
		v := govultr.NewClient(nil, "")
		return Vultr{client: v}, nil
	default:
		return nil, errors.New("Invalid provider used (Options: Vultr, DigitalOcean)")
	}

}

// GetFromConfig returns a provider from the config and user input
func GetFromConfig(conf parsers.Config, name string) (Provider, error) {
	for i := range conf.Servers {
		if conf.Servers[i].Name == name {
			switch conf.Servers[i].Provider {
			case "Vultr":
				v := govultr.NewClient(nil, conf.Keys.Vultr)
				return Vultr{client: v, name: name, region: conf.Servers[i].Region, plan: conf.Servers[i].Plan, os: conf.Servers[i].OS}, nil
			case "DigitalOcean":
				v := govultr.NewClient(nil, conf.Keys.DigitalOcean)
				return Vultr{client: v, name: name}, nil
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
