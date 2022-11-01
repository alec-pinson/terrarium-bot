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
				SendNotification("It's very hot: %vc/%vc", temperature, c.Temperature.Day.Maximum+c.Alerts.Temperature.Threshold)
			case temperature >= c.Temperature.Day.Maximum:
				HeatingOff()
			case temperature <= c.Temperature.Day.Minumum-c.Alerts.Temperature.Threshold:
				HeatingOn()
				SendNotification("It's very cold: %vc/%vc", temperature, c.Temperature.Day.Minumum-c.Alerts.Temperature.Threshold)
			case temperature <= c.Temperature.Day.Minumum:
				HeatingOn()
			case temperature > c.Temperature.Day.Minumum && temperature < c.Temperature.Day.Maximum:
				HeatingOn()
			}
		} else {
			switch temperature := GetTemperature(); {
			case temperature >= c.Temperature.Night.Maximum+c.Alerts.Temperature.Threshold:
				HeatingOff()
				SendNotification("It's very hot: %vc/%vc", temperature, c.Temperature.Night.Maximum+c.Alerts.Temperature.Threshold)
			case temperature >= c.Temperature.Night.Maximum:
				HeatingOff()
			case temperature <= c.Temperature.Night.Minumum-c.Alerts.Temperature.Threshold:
				HeatingOn()
				SendNotification("It's very cold: %vc/%vc", temperature, c.Temperature.Night.Minumum-c.Alerts.Temperature.Threshold)
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
	return int(resp.State.Sensors.Current)
}

func HeatingOn() {
	for _, h := range c.Switches {
		if h.Type == "heat" {
			if GetSwitchState(h) == "off" {
				log.Printf("Heater On: %s", h.Name)
				SetSwitchState(h, "on")
			}
		}
	}
}

func HeatingOff() {
	for _, h := range c.Switches {
		if h.Type == "heat" {
			if GetSwitchState(h) == "on" {
				log.Printf("Heater Off: %s", h.Name)
				SetSwitchState(h, "off")
			}
		}
	}
}
