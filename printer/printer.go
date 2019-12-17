package printer

import (
	"encoding/csv"
	"fmt"
	"sort"
	"strings"

	"github.com/ickerio/mima/providers"
)

func PrintInfo(server providers.Server) {
	fmt.Printf(""+
		"┌─────────────────────────────────────────────┐\n"+
		"│ %-27v %15v │\n"+
		"├───────────────┬─────────────────────────────┤\n"+
		"│ OS            │ %27v │\n"+
		"│ Location      │ %27v │\n"+
		"│ Memory        │ %27v │\n"+
		"│ Storage       │ %27v │\n"+
		"│ CPU Count     │ %27v │\n"+
		"└───────────────┴─────────────────────────────┘",
		server.Name, server.IP, server.Os, server.Location, server.Memory, server.Storage, server.CPUCount,
	)
}

func PrintPlans(plans []providers.Plan) {
	var output string
	output += "" +
		"┌────────┬──────────────────────────────────────────┐\n" +
		"│ ID     │ Description                              │\n" +
		"├────────┼──────────────────────────────────────────┤\n"

	for _, plan := range plans {
		reader := csv.NewReader(strings.NewReader(plan.Description))
		record, err := reader.Read()
		if err != nil {
			fmt.Printf("An error occurred when printing plans.")
			return
		}

		var description string
		for i, item := range record {
			if i == len(record)-1 {
				description += item
			} else {
				description += item + ", "
			}
		}
		output += fmt.Sprintf("│ %-6v │ %40v │\n", plan.ID, description)
	}

	output += "└────────┴──────────────────────────────────────────┘"
	fmt.Print(output)
}

func PrintRegions(regions []providers.Region) {
	var output string
	output += "" +
		"┌────────┬──────────────────────┐\n" +
		"│ ID     │ Name                 │\n" +
		"├────────┼──────────────────────┤\n"

	sort.SliceStable(regions, func(i, j int) bool {
		return regions[i].Name[0] < regions[j].Name[0]
	})

	for _, region := range regions {
		output += fmt.Sprintf("│ %-6v │ %20v │\n", region.ID, region.Name)
	}

	output += "└────────┴──────────────────────┘"
	fmt.Print(output)
}
