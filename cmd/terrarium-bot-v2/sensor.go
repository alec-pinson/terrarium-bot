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
	return &Sensor{}
}

func InitSensors() {
	for _, s := range config.Sensor {
		go s.monitor()
	}
}

func (s *Sensor) SetValue(value int) {
	s.Value = value
}

func (s *Sensor) GetValue() int {
	return s.Value
}

func (s *Sensor) getSensorValue() int {
	r, err := SendRequest(s.Url)
	if err != nil {
		log.Println(err)
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Println(err)
	}
	value := getJsonValue(string(b), s.JsonPath)
	intValue, err := strconv.Atoi(fmt.Sprintf("%.0f", value))
	s.SetValue(intValue)
	return intValue
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
