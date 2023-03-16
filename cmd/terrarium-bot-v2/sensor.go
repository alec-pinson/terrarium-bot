package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

func getSensor(id string) (*Sensor, error) {
	for _, s := range config.Sensor {
		if s.Id == id {
			return &s, nil
		}
	}
	return &Sensor{}, errors.New("Sensor '" + id + "' not found.")
}

func (s Sensor) getIdx() (int, error) {
	for idx, ss := range config.Sensor {
		if s.Id == ss.Id {
			return idx, nil
		}
	}
	return 0, errors.New("Sensor '" + s.Id + "' not found.")
}

func InitSensors() {
	for _, s := range config.Sensor {
		go s.monitor()
	}
}

func (s Sensor) SetValue(value int) {
	idx, err := s.getIdx()
	if err != nil {
		log.Println(err)
		return
	}
	config.Sensor[idx].Value = value
}

func (s Sensor) GetValue() int {
	idx, err := s.getIdx()
	if err != nil {
		log.Println(err)
		return 0
	}
	return config.Sensor[idx].Value
}

func (s Sensor) monitor() {
	log.Printf("Monitoring sensor '%s'", s.Id)
	for {
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
		log.Printf("%s: %v%s", strings.Title(s.Id), intValue, s.Unit)
		time.Sleep(1 * time.Minute)
	}
}
