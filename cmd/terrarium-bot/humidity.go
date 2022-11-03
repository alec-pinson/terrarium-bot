package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	lastMistTime time.Time = time.Now()
	mistMode     bool      = false
)

func MonitorMisting() {
	for {
		for _, m := range c.Switches {
			if m.Type == "mister" {
				if lastMistTime.Add(m.Sleep).Before(time.Now()) && DayTime() && dayStartTime.Add(30*time.Minute).Before(time.Now()) { // allow 30 minutes before misting in the morning
					log.Printf("Not misted for %s, current humidity (%v%%/%v%%). Misting now.", m.Sleep, GetHumidity(), c.Humidity.Day.Maximum)
					Mist()
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func MonitorHumidity() {
	FanInit()
	for {
		if DayTime() {
			switch humidity := GetHumidity(); {
			case humidity >= c.Humidity.Day.Maximum+c.Alerts.Humidity.Threshold:
				FanOn()
				if lastMistTime.Add(30 * time.Minute).Before(time.Now()) {
					SendNotification("Humidity is very high: %v%%/%v%%", humidity, c.Humidity.Day.Maximum+c.Alerts.Humidity.Threshold)
				}
			case humidity >= c.Humidity.Day.Maximum:
				FanOn()
			case humidity <= c.Humidity.Day.Minumum-c.Alerts.Humidity.Threshold:
				FanOff()
				SendNotification("Humidity is very low: %v%%/%v%%", humidity, c.Humidity.Day.Minumum-c.Alerts.Humidity.Threshold)
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
				if lastMistTime.Add(30 * time.Minute).Before(time.Now()) {
					SendNotification("Humidity is very high: %v%%/%v%%", humidity, c.Humidity.Night.Maximum+c.Alerts.Humidity.Threshold)
				}
			case humidity >= c.Humidity.Night.Maximum:
				FanOn()
			case humidity <= c.Humidity.Night.Minumum-c.Alerts.Humidity.Threshold:
				FanOff()
				SendNotification("Humidity is very low: %v%%/%v%%", humidity, c.Humidity.Night.Minumum-c.Alerts.Humidity.Threshold)
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
	return int(resp.State.Sensors.Current)
}

func Mist() {
	for _, b := range c.GPIO {
		if b.Type == "button" && b.Action == "prevent mist" && b.LastStateChange.Add(b.Sleep).After(time.Now()) {
			log.Printf("Misting has been prevented via button press")
			return
		}
	}

	if lastMistTime.Add(30 * time.Minute).After(time.Now()) {
		log.Printf("Will not mist, hit hard coded 30 minute misting limit")
		return
	}

	log.Printf("Misting will begin shortly, current humidity (%v%%/%v%%)", GetHumidity(), c.Humidity.Day.Maximum)
	mistMode = true
	for _, l := range c.Switches {
		if l.Type == "light" {
			LightOff(l)
		}
	}
	time.Sleep(5 * time.Minute)
	for _, m := range c.Switches {
		if m.Type == "mister" {
			DoMist(m)
		}
	}

	for _, l := range c.Switches {
		if l.Type == "light" {
			LightOn(l)
		}
	}
	mistMode = false
}

func DoMist(Mister Switch) {
	log.Printf("Misting for %v seconds", Mister.Length)
	lastMistTime = time.Now()
	SetSwitchState(Mister, "on")
	time.Sleep(Mister.Length)
	SetSwitchState(Mister, "off")
	time.Sleep(1 * time.Second)
	SetSwitchState(Mister, "off")
	time.Sleep(1 * time.Second)
	SetSwitchState(Mister, "off")
	time.Sleep(1 * time.Second)
	SetSwitchState(Mister, "off")
}
