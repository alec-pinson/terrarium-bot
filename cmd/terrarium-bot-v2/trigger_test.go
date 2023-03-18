package main

import (
	"testing"
	"time"
)

func TestMonitorSensor(t *testing.T) {
	// disable http calls when turning on/off switches
	config.DryRun = true
	// set testing mode so we exit the for loop
	isTesting = true

	// set up a trigger with a simulated sensor and switch
	trigger := Trigger{
		Sensor: "mock-sensor",
	}
	trigger.When.Day.Above = 50
	trigger.Action = append(trigger.Action, "switch.mock-switch.on")
	trigger.Else = append(trigger.Else, "switch.mock-switch.off")
	// create a simulated sensor with value 60
	s := Sensor{
		Id:    "mock-sensor",
		Value: 60,
	}
	ss := Switch{
		Id:  "mock-switch",
		On:  "mock-switch.com/turnOn",
		Off: "mock-switch.com/turnOff",
	}
	config.Sensor = append(config.Sensor, &s)
	config.Switch = append(config.Switch, &ss)

	// ensure it is day time
	config.Day.StartTime, _ = time.Parse("15:04", "00:00")
	config.Night.StartTime, _ = time.Parse("15:04", "23:59")

	// call monitor sensor and check it turns on the switch
	trigger.monitor()
	if ss.getStatus() != "on" {
		t.Errorf("Trigger should have turned on the mock-switch but didn't")
	}

	// set senor value to lower to try trigger the else action
	s.Value = 40
	trigger.monitor()
	if ss.getStatus() != "off" {
		t.Errorf("Trigger should have turned off the mock-switch but didn't")
	}

	// reset
	config.DryRun = false
	isTesting = false
}

func TestIsTriggerEndpoint(t *testing.T) {
	config.Trigger = append(config.Trigger, &Trigger{
		Endpoint: "/endpoint",
	})

	isTrigger, _ := isTriggerEndpoint("endpoint")

	if !isTrigger {
		t.Errorf("Expected endpoint to exist")
	}

	isTrigger, _ = isTriggerEndpoint("non-existent")

	if isTrigger {
		t.Errorf("Expected endpoint to not exist")
	}
}
