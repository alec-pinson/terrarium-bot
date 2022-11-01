package main

import (
	"log"
	"time"
)

var (
	c Configuration
)

func main() {
	log.Println("Starting...")
	c = LoadConfiguration()

	GetSunriseTime()
	GetSunsetTime()
	if DayTime() {
		log.Printf("Current Humidity: %v%%/%v%%", GetHumidity(), c.Humidity.Day.Maximum)
		log.Printf("Current Temperature: %vc/%vc", GetTemperature(), c.Temperature.Day.Maximum)
	} else {
		log.Printf("Current Humidity: %v%%/%v%%", GetHumidity(), c.Humidity.Night.Maximum)
		log.Printf("Current Temperature: %vc/%vc", GetTemperature(), c.Temperature.Night.Maximum)
	}

	FanInit()

	go MonitorLights()
	go MonitorTemperature()
	go MonitorHumidity()
	go MonitorButtons()
	go MonitorMisting()

	for {
		time.Sleep(1 * time.Minute)
	}
}
