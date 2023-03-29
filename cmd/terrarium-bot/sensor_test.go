package main

import (
	"testing"
)

func TestGetSensor(t *testing.T) {
	config.Sensor = []*Sensor{
		{Id: "temperature"},
		{Id: "humidity"},
	}
	expected := config.Sensor[0]
	result := GetSensor("temperature")

	if result != expected {
		t.Errorf("Expected '%+v' but got '%+v'", expected, result)
	}
}

func TestSensorSetValue(t *testing.T) {
	s := &Sensor{}
	s.SetValue(42)
	if s.Value != 42 {
		t.Errorf("Expected sensor value to be 42, but got %d", s.Value)
	}
}

func TestSensorGetValue(t *testing.T) {
	s := &Sensor{Value: 42}
	v := s.GetValue()
	if v != 42 {
		t.Errorf("Expected sensor value to be 42, but got %d", v)
	}
}
