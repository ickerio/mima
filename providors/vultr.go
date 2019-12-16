package providors

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
	var s Server

	res, err := v.client.Server.List(context.Background())
	if err != nil {
		return s, err
	}

	for _, el := range res {
		if el.Label == v.name {
			s = Server{
				Name:             el.Label,
				Os:               el.Os,
				Memory:           el.RAM,
				Storage:          el.Disk,
				CPUCount:         el.VPSCpus,
				IP:               el.MainIP,
				CurrentBandwidth: el.CurrentBandwidth,
				AllowedBandwidth: el.AllowedBandwidth,
				Location:         el.Location,
				Cost:             el.Cost,
				Created:          el.Created,
				Password:         el.DefaultPassword,
			}
			return s, nil
		}
	}
	return s, errors.New("Server currently offline")
}

// CreateServer TOODODODODODODO
func (v Vultr) CreateServer() {
	vpsOptions := &govultr.ServerOptions{
		Label:                "awesome-go-app",
		Hostname:             "awesome-go.com",
		EnablePrivateNetwork: true,
		AutoBackups:          true,
		EnableIPV6:           true,
	}

	// RegionId, VpsPlanID, OsID can be grabbed from their respective API calls
	res, err := v.client.Server.Create(context.Background(), 1, 201, 1, vpsOptions)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}
