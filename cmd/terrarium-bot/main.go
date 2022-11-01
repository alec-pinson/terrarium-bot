package main

import (
	"log"
)

var (
	c Configuration
)

func main() {
	log.Println("Starting...")
	c = LoadConfiguration()

	GetSunriseTime()
	GetSunsetTime()

	log.Printf("Current Humidity: %v", int(GetHumidity()))
	log.Printf("Current Temperature: %v", int(GetTemperature()))

	FanInit()

	go MonitorLights()
	go MonitorTemperature()
	go MonitorHumidity()
	// go MonitorButtons()
	MonitorMisting()
}
