package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

func (s *Sensor) SetValue(value float64) {
	s.Value = value
}

func (s *Sensor) GetValue() float64 {
	return s.Value
}

func (s *Sensor) getSensorValue() float64 {
	r, respCode, err := SendRequest(s.Url, s.Insecure, 3, s.JsonPath != "")
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
	floatVal, err := strconv.ParseFloat(fmt.Sprint(value), 64)
	if err != nil {
		log.Println(err)
		return 0
	}
	roundedValue := math.Round(floatVal*100) / 100
	s.SetValue(roundedValue)
	s.checkValue()
	return roundedValue
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
	var sensorMetrics = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "terrarium_bot_sensor_" + s.Id,
		Help: "The current value of the " + s.Id + " sensor",
	})

	val := s.getSensorValue()
	log.Printf("Monitoring sensor '%s' (%v%s)", s.Id, val, s.Unit)
	for {
		val = s.getSensorValue()
		sensorMetrics.Set(float64(val))
		Debug("%s: %v%s", strings.Title(s.Id), val, s.Unit)
		time.Sleep(1 * time.Minute)
	}
}
