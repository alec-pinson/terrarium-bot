package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type APIServer struct{}

func (apiServer APIServer) Start() {
	var wg sync.WaitGroup
	log.Println("Starting API server...")
	http.HandleFunc("/", apiServer.Endpoint)
	wg.Add(1)
	go http.ListenAndServe(":8080", nil)
	log.Println("API Server started...")
	wg.Wait()
}

func (apiServer APIServer) Endpoint(w http.ResponseWriter, r *http.Request) {
	switch path := r.URL.Path[1:]; {
	case path == "health/live" || path == "health/ready":
		fmt.Fprintf(w, "ok")
	default:
		if path == "favicon.ico" {
			// ignore
			return
		}
		// check if this is a trigger endpoint, if so do action
		yes, t := isTriggerEndpoint(path)
		if yes {
			t.doAction("Triggered by endpoint '" + path + "'")
		}
	}
}
