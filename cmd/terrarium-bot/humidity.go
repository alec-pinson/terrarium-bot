package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func MonitorHumidity() {
	for {
		if DayTime() {
			switch humidity := GetHumidity(); {
			case humidity >= c.Humidity.Day.Maximum+c.Alerts.Humidity.Threshold:
				FanOn()
				SendHumidityAlert("Really high humidity levels")
			case humidity >= c.Humidity.Day.Maximum:
				FanOn()
			case humidity <= c.Humidity.Day.Minumum-c.Alerts.Humidity.Threshold:
				FanOff()
				SendHumidityAlert("Way below humidity levels")
			case humidity >= c.Humidity.Day.Minumum && humidity < c.Humidity.Day.Maximum:
				FanOff()
			}
		} else {
			switch humidity := GetHumidity(); {
			case humidity >= c.Humidity.Night.Maximum+c.Alerts.Humidity.Threshold:
				FanOn()
				SendHumidityAlert("Really high humidity levels")
			case humidity >= c.Humidity.Night.Maximum:
				FanOn()
			case humidity <= c.Humidity.Night.Minumum-c.Alerts.Humidity.Threshold:
				FanOff()
				SendHumidityAlert("Way below humidity levels")
			case humidity >= c.Humidity.Night.Minumum && humidity < c.Humidity.Night.Maximum:
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

func FanOn() {
	log.Println("Turning on fan")
}

func FanOff() {
	log.Println("Turning off fan")
}

func SendHumidityAlert(alertMessage string) {
	log.Println("SendHumidityAlert(): Sent alert: '" + alertMessage + "'")
}
