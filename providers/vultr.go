package providers

import (
	"context"
	"errors"
	"fmt"

	"github.com/vultr/govultr"
)

// Vultr Client for Vultr VPS
type Vultr struct {
	apiKey string
	client *govultr.Client
	name   string
}

// Info Retrieves all hosted VPS servers
func (v Vultr) Info() (Server, error) {
	var server Server

	res, err := v.client.Server.List(context.Background())
	if err != nil {
		return server, err
	}

	for _, element := range res {
		if element.Label == v.name {
			server = Server{
				Name:             element.Label,
				Os:               element.Os,
				Memory:           element.RAM,
				Storage:          element.Disk,
				CPUCount:         element.VPSCpus,
				IP:               element.MainIP,
				CurrentBandwidth: element.CurrentBandwidth,
				AllowedBandwidth: element.AllowedBandwidth,
				Location:         element.Location,
				Cost:             element.Cost,
				Created:          element.Created,
				Password:         element.DefaultPassword,
			}
			return server, nil
		}
	}
	return server, errors.New("Server currently offline")
}

// CreateServer TOODODODODODODO
func (v Vultr) CreateServer() {
	vpsOptions := &govultr.ServerOptions{
		Label: v.name,
	}

	// RegionId, VpsPlanID, OsID can be grabbed from their respective API calls
	res, err := v.client.Server.Create(context.Background(), 1, 201, 1, vpsOptions)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}
