package main

import (
	"log"
)

var (
	c Configuration
)

func main() {
	log.Println("main(): Starting...")
	c = LoadConfiguration()

	if DayTime() {
		log.Println("Day time")
	} else {
		log.Println("Night time")
	}

	FanInit()

	go MonitorLights()
	go MonitorTemperature()
	go MonitorHumidity()
	// go MonitorButtons()
	MonitorMisting()
}
