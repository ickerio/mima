package printer

import (
	"fmt"

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
