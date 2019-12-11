package main

import (
	"flag"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ickerio/mima/api"
	"github.com/ickerio/mima/providors"
)

/*
type Options struct {
	label string
}
*/

func parseFlags() (string, string) {
	var apiKey = flag.String("key", "", "API key for either Vultr or DigitalOcean")
	var port = flag.String("port", "8080", "Port to host web server on")
	flag.Parse()

	if *apiKey == "" {
		flag.PrintDefaults()
		panic("No apikey set!")
	}

	return *apiKey, *port
}

func main() {
	apiKey, port := parseFlags()

	server := mux.NewRouter()
	handler := api.Handler{Providor: providors.NewProvidor(providors.Vultr, apiKey)}

	router := server.PathPrefix("/api").Subrouter()
	router.HandleFunc("/info", handler.HandleInfo).Methods(http.MethodGet)
	router.HandleFunc("/start", handler.HandleStart).Methods(http.MethodGet)

	server.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	if err := http.ListenAndServe(":"+port, server); err != nil {
		panic(err)
	}
}
