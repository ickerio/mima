package main

import (
	"net/http"
	"flag"

	"github.com/gorilla/mux"
	"github.com/ickerio/mima/providors"
	"github.com/ickerio/mima/api"
	"github.com/ickerio/mima/util"
)

var prov VpsProvidor

var apiKey = flag.Int("apikey", "", "API key for either Vultr or DigitalOcean")

type Options struct {
	label string
}

func main() {
	prov = NewProvidor(Vultr, apiKey)

	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/info", handleInfo).Methods(http.MethodGet)
	api.HandleFunc("/start", handleStart).Methods(http.MethodGet)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
