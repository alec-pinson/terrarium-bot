package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
				if dayEndTime.Add(1 * time.Hour).Before(time.Now()) { // allow an hour cool down
					SendNotification("It's very hot: %vc/%vc", temperature, c.Temperature.Night.Maximum+c.Alerts.Temperature.Threshold)
				}
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
		log.Println("Returning Temperature.Maximum - 1")
		if DayTime() {
			return c.Temperature.Day.Maximum - 1
		} else {
			return c.Temperature.Night.Maximum - 1
		}
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		log.Println("Returning Temperature.Maximum - 1")
		if DayTime() {
			return c.Temperature.Day.Maximum - 1
		} else {
			return c.Temperature.Night.Maximum - 1
		}
	}
	var resp TerrariumPiSensorResp
	json.Unmarshal(responseData, &resp)
	return int(resp.State.Sensors.Current)
}

func HeatingOn() {
	for _, h := range c.Switches {
		if h.Type == "heat" {
			if GetSwitchState(h) == "off" {
				if DayTime() {
					log.Printf("Heater On: %s (%vc/%vc)", h.Name, GetTemperature(), c.Temperature.Day.Maximum)
				} else {
					log.Printf("Heater On: %s (%vc/%vc)", h.Name, GetTemperature(), c.Temperature.Night.Maximum)
				}
				SetSwitchState(h, "on")
			}
		}
	}
}

func HeatingOff() {
	for _, h := range c.Switches {
		if h.Type == "heat" {
			if GetSwitchState(h) == "on" {
				if DayTime() {
					log.Printf("Heater Off: %s (%vc/%vc)", h.Name, GetTemperature(), c.Temperature.Day.Maximum)
				} else {
					log.Printf("Heater Off: %s (%vc/%vc)", h.Name, GetTemperature(), c.Temperature.Night.Maximum)
				}
				SetSwitchState(h, "off")
			}
		}
	}
}
