package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ickerio/mima/providors"
)

// Handler HTTP json api interface handler
type Handler struct {
	Providor providors.VpsProvidor
}

// HandleInfo Handles api/info
func (api Handler) HandleInfo(w http.ResponseWriter, r *http.Request) {
	res := api.Providor.ListServers()
	json.NewEncoder(w).Encode(res)
}

// HandleStart Handles api/start
func (api Handler) HandleStart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is start")
}
