package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func GetSensor(id string) *Sensor {
	for _, s := range config.Sensor {
		if s.Id == id {
			return s
		}
	}
	log.Fatalf("Sensor '%s' not found in configuration.yaml", id)
	return nil
}

func InitSensors() {
	for _, s := range config.Sensor {
		go s.monitor()
	}
	time.Sleep(5 * time.Second) // give abit of time for any sensors to collect data
}

func (s *Sensor) SetValue(value int) {
	s.Value = value
}

func (s *Sensor) GetValue() int {
	return s.Value
}

func (s *Sensor) getSensorValue() int {
	r, respCode, err := SendRequest(s.Url, s.Insecure)
	if err != nil {
		log.Println(err)
		return 0
	}
	if respCode != 200 {
		log.Printf("Unable to get sensor value, response code: %v", respCode)
		return 0
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Println(err)
		return 0
	}
	value := getJsonValue(string(b), s.JsonPath)
	intValue, err := strconv.Atoi(fmt.Sprintf("%.0f", value))
	s.SetValue(intValue)
	s.checkValue()
	return intValue
}

func (s *Sensor) checkValue() {
	if s.GetValue() == 0 {
		log.Printf("Sensor Offline: '%s'", s.Id)
		for _, n := range config.Notification {
			n.SendNotification("Currently unable to get value for sensor '%s'. Please check the logs.", s.Id)
		}
	}
}

func (s *Sensor) monitor() {
	val := s.getSensorValue()
	log.Printf("Monitoring sensor '%s' (%v%s)", s.Id, val, s.Unit)
	for {
		val = s.getSensorValue()
		Debug("%s: %v%s", strings.Title(s.Id), val, s.Unit)
		time.Sleep(1 * time.Minute)
	}
}
