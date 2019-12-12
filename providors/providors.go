package providors

import (
	"errors"

	"github.com/ickerio/mima/config"
	"github.com/vultr/govultr"
)

// Providor Interface for the Vultr, DigitalOcean
type Providor interface {
	ListServers() []Server
}

// Server Details the information of a VPS server
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

// Get Returns a providor from the config and the user input
func Get(config config.Config, name string) (Providor, error) {
	for i := range config.Servers {
		if config.Servers[i].Name == name {
			switch config.Servers[i].Providor {
			case "Vultr":
				v := govultr.NewClient(nil, config.Keys.Vultr)
				return Vultr{apiKey: config.Keys.Vultr, client: v, instanceName: name}, nil
			case "DigitalOcean":
				v := govultr.NewClient(nil, config.Keys.DigitalOcean)
				return Vultr{apiKey: config.Keys.DigitalOcean, client: v, instanceName: name}, nil
			default:
				return nil, errors.New("Invalid providor in config")
			}
		}
	}
	return nil, errors.New("Could not match instance name to config file")
}
