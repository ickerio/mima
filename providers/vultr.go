package providers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/vultr/govultr"
)

// Vultr Client for Vultr VPS
type Vultr struct {
	client *govultr.Client
	name   string
	region int
	plan   int
	os     int
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

// Start the desired server
func (v Vultr) Start() {
	vpsOptions := &govultr.ServerOptions{
		Label: v.name,
	}

	// RegionId, VpsPlanID, OsID can be grabbed from their respective API calls
	res, err := v.client.Server.Create(context.Background(), v.region, v.plan, v.os, vpsOptions)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

// Stop the desired server
func (v Vultr) Stop() {

}

// Plans will grab all regions available from the provider
func (v Vultr) Plans() ([]Plan, error) {
	var plans []Plan
	pla, err := v.client.Plan.GetVc2List(context.Background())

	for _, element := range pla {
		planID, _ := strconv.Atoi(element.PlanID)
		plans = append(plans, Plan{
			ID:          planID,
			Description: element.Name,
		})
	}

	return plans, err
}

// Regions will grab all regions available from the provider
func (v Vultr) Regions() ([]Region, error) {
	var regions []Region
	reg, err := v.client.Region.List(context.Background())

	for _, element := range reg {
		regionID, _ := strconv.Atoi(element.RegionID)
		regions = append(regions, Region{
			ID:   regionID,
			Name: element.Name,
		})
	}

	return regions, err
}

// OS will grab all operating systems available from the provider
func (v Vultr) OS() ([]OS, error) {
	var systems []OS
	sys, err := v.client.OS.List(context.Background())

	for _, element := range sys {
		systems = append(systems, OS{
			ID:   element.OsID,
			Name: element.Name,
		})
	}

	return systems, err
}
