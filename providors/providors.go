package providors

import (
	"github.com/vultr/govultr"
)

/*
 * Vultr Enum for the Vultr VPS when using NewProvidor()
 * DigitalOcean Enum for the DigitalOcean VPS when using NewProvidor()
 */
const (
	Vultr        = iota
	DigitalOcean = iota
)

// VpsProvidor Interface for the Vultr, DigitalOcean
type VpsProvidor interface {
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

// NewProvidor Instanciates the desired VPS client
func NewProvidor(providor int, apiKey string) VpsProvidor {
	switch providor {
	case Vultr:
		v := govultr.NewClient(nil, apiKey)
		return VultrProvidor{apiKey: apiKey, client: v}
	case DigitalOcean:
		v := govultr.NewClient(nil, apiKey)
		return VultrProvidor{apiKey: apiKey, client: v}
	default:
		panic("No such providor")
	}
}
