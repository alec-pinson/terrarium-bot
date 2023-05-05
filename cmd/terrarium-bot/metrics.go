package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct{}

func (metrics Metrics) Start() {
	log.Println("Starting Metrics server...")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8081", nil)
	time.Sleep(1 * time.Second) // give some time for metrics server to start before moving on
	log.Println("Metrics Server started...")
}
