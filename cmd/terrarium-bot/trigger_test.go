package main

import (
	"testing"
	"time"
)

func TestMonitorTrigger(t *testing.T) {
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

	config.Debug = true

	// test dsiabled triggers
	trigger.Disable("2s", "")
	s.Value = 60
	time.Sleep(1 * time.Second)
	trigger.monitor()
	if ss.getStatus() != "off" {
		t.Errorf("Trigger should still be turned off but isn't")
	}
	time.Sleep(3 * time.Second)
	trigger.monitor()
	if ss.getStatus() != "on" {
		t.Errorf("Trigger should have turned on the mock-switch but didn't")
	}

	// reuse same trigger to test DroppedBy
	trigger.When.Day.Above = 0
	// set when dropped by to 5
	trigger.When.Day.DroppedBy = 5

	// set the sensor to 60
	s.SetValue(60)
	// call trigger monitor to make sure it turns off the switch
	trigger.monitor()
	if ss.getStatus() != "off" {
		t.Errorf("DropBy trigger should have turned off the mock-switch but didn't")
	}

	// set the sensor to 56
	s.SetValue(56)
	// call trigger monitor to make sure it turns off the switch
	trigger.monitor()
	if ss.getStatus() != "off" {
		t.Errorf("DropBy trigger should have turned off the mock switch but didn't. Current sensor val: %v, previous sensor val: %v", s.Value, s.PreviousValue)
	}

	// drop sensor value by 5
	s.SetValue(51)
	// call trigger monitor to make sure it turns off the switch
	trigger.monitor()
	if ss.getStatus() != "on" {
		t.Errorf("DropBy trigger should have turned on the mock switch but didn't. Current sensor val: %v, previous sensor val: %v", s.Value, s.PreviousValue)
	}

	// reset the switch
	s.SetValue(51)
	trigger.monitor()
	if ss.getStatus() != "off" {
		t.Errorf("Trigger should have turned off the mock switch but didn't. Current sensor val: %v, previous sensor val: %v", s.Value, s.PreviousValue)
	}

	// set when dropped by to 0
	trigger.When.Day.DroppedBy = 0
	// now to test increased by
	trigger.When.Day.IncreasedBy = 5

	// set the sensor to 53
	s.SetValue(53)
	// call trigger monitor to make sure it turns off the switch
	trigger.monitor()
	if ss.getStatus() != "off" {
		t.Errorf("IncreaseBy trigger should have turned off the mock-switch but didn't")
	}

	// increase sensor value by 5+
	s.SetValue(60)
	// call trigger monitor to make sure it turns off the switch
	trigger.monitor()
	if ss.getStatus() != "on" {
		t.Errorf("IncreaseBy trigger should have turned on the mock switch but didn't. Current sensor val: %v, previous sensor val: %v", s.Value, s.PreviousValue)
	}

	// reset
	config.Debug = false
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

func TestTriggerDisable(t *testing.T) {
	trigger := &Trigger{Id: "disable"}

	// Set disable and ensure that it's set correctly
	trigger.Disable("1m", "")
	if trigger.Disabled != 1*time.Minute {
		t.Errorf("Trigger Disable was not set correctly")
	}
}

func TestTriggerIsDisabled(t *testing.T) {
	// Test case 1: trigger should not be disabled
	trigger := Trigger{Disabled: 0}
	trigger.Disable("", "")
	trigger.Enable("")
	if trigger.isDisabled() {
		t.Errorf("Test case 1 failed: expected false but got true")
	}

	// Test case 2: trigger should be disabled for 2 seconds
	trigger.Disable("2s", "")
	time.Sleep(1 * time.Second) // sleep for 1 second
	if !trigger.isDisabled() {
		t.Errorf("Test case 2 failed: expected true but got false")
	}
	time.Sleep(3 * time.Second) // sleep for 3 seconds
	// should not be disabled anymore
	if trigger.isDisabled() {
		t.Errorf("Test case 2 failed: expected false but got true")
	}
}
