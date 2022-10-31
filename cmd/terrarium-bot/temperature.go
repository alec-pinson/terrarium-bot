package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func MonitorTemperature() {
	for {
		if DayTime() {
			switch temperature := GetTemperature(); {
			case temperature >= c.Temperature.Day.Maximum+c.Alerts.Temperature.Threshold:
				HeatingOff()
				SendNotification("Really hot")
			case temperature >= c.Temperature.Day.Maximum:
				HeatingOff()
			case temperature <= c.Temperature.Day.Minumum-c.Alerts.Temperature.Threshold:
				HeatingOn()
				SendNotification("Really cold")
			case temperature <= c.Temperature.Day.Minumum:
				HeatingOn()
			case temperature > c.Temperature.Day.Minumum && temperature < c.Temperature.Day.Maximum:
				HeatingOn()
			}
		} else {
			switch temperature := GetTemperature(); {
			case temperature >= c.Temperature.Night.Maximum+c.Alerts.Temperature.Threshold:
				HeatingOff()
				SendNotification("Really hot")
			case temperature >= c.Temperature.Night.Maximum:
				HeatingOff()
			case temperature <= c.Temperature.Night.Minumum-c.Alerts.Temperature.Threshold:
				HeatingOn()
				SendNotification("Really cold")
			case temperature <= c.Temperature.Night.Minumum:
				HeatingOn()
			case temperature > c.Temperature.Night.Minumum && temperature < c.Temperature.Night.Maximum:
				HeatingOn()
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func GetTemperature() int {
	response, err := http.Get(c.Temperature.Url)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var resp TerrariumPiSensorResp
	json.Unmarshal(responseData, &resp)
	log.Printf("GetTemperature(): Current Temperature: %v", int(resp.State.Sensors.Current))
	return int(resp.State.Sensors.Current)
}

func HeatingOn() {
	for _, h := range c.Switches {
		if h.Type == "heat" {
			if GetSwitchState(h) == "off" {
				log.Printf("HeatingOn(): Turning on heater '%s'", h.Name)
			}
		}
	}
}

func HeatingOff() {
	for _, h := range c.Switches {
		if h.Type == "heat" {
			if GetSwitchState(h) == "on" {
				log.Printf("HeatingOff(): Turning off heater '%s'", h.Name)
			}
		}
	}
}
