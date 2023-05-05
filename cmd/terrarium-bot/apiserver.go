package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type APIServer struct{}

var metricHttpRequestsReceived = promauto.NewCounter(prometheus.CounterOpts{
	Name: "terrarium_bot_http_requests_received_total",
	Help: "The total number of http requests received by terrarium bot",
})

func (apiServer APIServer) Start() {
	log.Println("Starting API server...")
	http.HandleFunc("/", apiServer.Endpoint)
	go http.ListenAndServe(":8080", nil)
	time.Sleep(1 * time.Second) // give some time for api server to start before moving on
	log.Println("API Server started...")
}

func (apiServer APIServer) Endpoint(w http.ResponseWriter, r *http.Request) {
	metricHttpRequestsReceived.Inc()
	switch path := r.URL.Path[1:]; {
	case path == "health/live" || path == "health/ready":
		fmt.Fprintf(w, "ok")
	case path == "switch":
		writeResponse(w, config.Switch, false)
	case strings.HasPrefix(path, "switch/"):
		GetSwitch(strings.Split(path, "/")[1], false).WriteStatus(w)
	default:
		if path == "favicon.ico" {
			// ignore
			return
		}
		// check if this is a trigger endpoint, if so do action
		yes, t := isTriggerEndpoint(path)
		if yes {
			fmt.Fprintf(w, "ok")
			t.doAction("Triggered by endpoint '" + path + "'")
		}
	}
}
