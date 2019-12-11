package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handleInfo(w http.ResponseWriter, r *http.Request) {
	res := prov.ListServers()
	json.NewEncoder(w).Encode(res)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is start")
}
