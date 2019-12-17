package providers

import (
	"context"
	"errors"
	"strings"

	"github.com/vultr/govultr"
)

// Vultr Client for Vultr VPS
type Vultr struct {
	client *govultr.Client
	name   string
}

// Info Retrieves all hosted VPS servers
func (v Vultr) Info() (Server, error) {
	var server Server

	servers, err := v.client.Server.List(context.Background())
	if err != nil {
		return server, err
	}

	for _, element := range servers {
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

// Regions will grab all regions available from the provider
func (v Vultr) Regions() ([]Region, error) {
	var regions []Region
	reg, err := v.client.Region.List(context.Background())

	for _, element := range reg {
		regions = append(regions, Region{
			ID:   element.RegionID,
			Name: element.Name,
		})
	}

	return regions, err
}

// Plans will grab all regions available from the provider
func (v Vultr) Plans() ([]Plan, error) {
	var plans []Plan
	plan, err := v.client.Plan.GetVc2List(context.Background())

	for _, element := range plan {
		plans = append(plans, Plan{
			ID:          element.PlanID,
			Description: strings.Replace(element.Name, ",", ", ", -1),
		})
	}

	return plans, err
}
