package providors

import (
	"context"
	"fmt"

	"github.com/vultr/govultr"
)

// VultrProvidor Client for Vultr VPS
type VultrProvidor struct {
	apiKey string
	client *govultr.Client
}

// ListServers Retrieves all hosted VPS servers
func (v VultrProvidor) ListServers() []Server {
	res, _ := v.client.Server.List(context.Background())
	var list []Server
	for _, el := range res {
		list = append(list, Server{
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
		})
	}
	return list
}

// CreateServer TOODODODODODODO
func (v VultrProvidor) CreateServer() {
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
