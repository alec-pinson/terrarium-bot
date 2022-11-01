package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var lastMistTime time.Time = time.Now()

func MonitorMisting() {
	for {
		for _, m := range c.Switches {
			if m.Type == "mister" {
				if lastMistTime.Add(m.Sleep).Before(time.Now()) && DayTime() {
					Mist()
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func MonitorHumidity() {
	for {
		if DayTime() {
			switch humidity := GetHumidity(); {
			case humidity >= c.Humidity.Day.Maximum+c.Alerts.Humidity.Threshold:
				FanOn()
				SendNotification("Really high humidity levels")
			case humidity >= c.Humidity.Day.Maximum:
				FanOn()
			case humidity <= c.Humidity.Day.Minumum-c.Alerts.Humidity.Threshold:
				FanOff()
				SendNotification("Way below humidity levels")
			case humidity <= c.Humidity.Day.Minumum:
				FanOff()
				Mist()
			case humidity > c.Humidity.Day.Minumum && humidity < c.Humidity.Day.Maximum:
				FanOff()
			}
		} else {
			switch humidity := GetHumidity(); {
			case humidity >= c.Humidity.Night.Maximum+c.Alerts.Humidity.Threshold:
				FanOn()
				SendNotification("Really high humidity levels")
			case humidity >= c.Humidity.Night.Maximum:
				FanOn()
			case humidity <= c.Humidity.Night.Minumum-c.Alerts.Humidity.Threshold:
				FanOff()
				SendNotification("Way below humidity levels")
			case humidity <= c.Humidity.Night.Minumum:
				FanOff()
			case humidity > c.Humidity.Night.Minumum && humidity < c.Humidity.Night.Maximum:
				FanOff()
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func GetHumidity() int {
	response, err := http.Get(c.Humidity.Url)
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
	log.Printf("GetHumidity(): Current Humidity: %v", int(resp.State.Sensors.Current))
	return int(resp.State.Sensors.Current)
}

func Mist() {
	for _, b := range c.GPIO {
		if b.Type == "button" && b.Action == "prevent mist" && b.LastStateChange.Add(b.Sleep).After(time.Now()) {
			log.Printf("Misting has been prevented via button press")
			return
		}
	}

	log.Printf("Misting will begin shortly")
	for _, l := range c.Switches {
		if l.Type == "light" {
			LightOff(l)
		}
	}
	time.Sleep(5 * time.Minute)
	for _, m := range c.Switches {
		if m.Type == "mister" {
			log.Printf("Misting for %v seconds", m.Length)
			lastMistTime = time.Now()
		}
	}
	time.Sleep(5 * time.Second)
	for _, l := range c.Switches {
		if l.Type == "light" {
			LightOn(l)
		}
	}
}
