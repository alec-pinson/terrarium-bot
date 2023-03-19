package main

import (
	"testing"
	"time"
)

func TestRunAction(t *testing.T) {
	// disable http calls when turning on/off switches
	config.DryRun = true

	//
	// Sleep action
	//
	// Test case 1: sleep runs for 1 second
	start := time.Now()
	RunAction("sleep.1s", "")
	sleepTime := time.Since(start).Seconds()
	if sleepTime < 1 || sleepTime > 2 {
		t.Errorf("Expected to sleep for 1 second, but slept for %v", sleepTime)
	}
	// Test case 2: sleep errors with invalid duration
	if RunAction("sleep.XYZs", "") {
		t.Errorf("Expected to false, but got true")
	}

	//
	// Switch action
	//
	s := &Switch{Id: "test-action-switch"}
	config.Switch = append(config.Switch, s)
	// Test case 1: switch turns on
	s.TurnOff("")
	RunAction("switch.test-action-switch.on", "")
	if s.getStatus() != "on" {
		t.Errorf("RunAction('switch.test-action-switch.on') failed, got %v, want %v", s.getStatus(), "on")
	}
	// Test case 2: switch turns on with duration
	s.TurnOff("")
	go RunAction("switch.test-action-switch.on.2s", "")
	time.Sleep(1 * time.Second)
	if s.getStatus() != "on" {
		t.Errorf("RunAction('switch.test-action-switch.on') failed, got %v, want %v", s.getStatus(), "on")
	}
	time.Sleep(1 * time.Second)
	if s.getStatus() != "off" {
		t.Errorf("RunAction('switch.test-action-switch.on') failed, got %v, want %v", s.getStatus(), "off")
	}
	// Test case 3: switch turning off
	s.TurnOn("", "")
	RunAction("switch.test-action-switch.off", "")
	if s.getStatus() != "off" {
		t.Errorf("RunAction('switch.test-action-switch.off') failed, got %v, want %v", s.getStatus(), "off")
	}
	// Test case 4: disable switch
	RunAction("switch.test-action-switch.disable", "")
	if s.isDisabled() != true {
		t.Errorf("RunAction('switch.test-action-switch.disable') failed, got %v, want %v", s.isDisabled(), true)
	}
	// Test case 5: enabling switch
	RunAction("switch.test-action-switch.enable", "")
	if s.isDisabled() != false {
		t.Errorf("RunAction('switch.test-action-switch.enable') failed, got %v, want %v", s.isDisabled(), false)
	}
	// Test case 6: disabling switch with duration
	RunAction("switch.test-action-switch.disable.1s", "")
	if s.isDisabled() != true {
		t.Errorf("RunAction('switch.test-action-switch.disable.1s') failed, got %v, want %v", s.isDisabled(), true)
	}
	time.Sleep(1 * time.Second)
	if s.isDisabled() != false {
		t.Errorf("RunAction('switch.test-action-switch.disable.1s') failed, got %v, want %v", s.isDisabled(), false)
	}

	//
	// Trigger action
	//
	trigger := &Trigger{Id: "test-action-trigger"}
	config.Trigger = append(config.Trigger, trigger)
	// Test case 1: disable trigger
	RunAction("trigger.test-action-trigger.disable", "")
	if trigger.isDisabled() != true {
		t.Errorf("RunAction('trigger.test-action-trigger.disable') failed, got %v, want %v", trigger.isDisabled(), true)
	}
	// Test case 2: enabling trigger
	RunAction("trigger.test-action-trigger.enable", "")
	if trigger.isDisabled() != false {
		t.Errorf("RunAction('trigger.test-action-trigger.enable') failed, got %v, want %v", trigger.isDisabled(), false)
	}
	// Test case 3: disabling trigger with duration
	RunAction("trigger.test-action-trigger.disable.1s", "")
	if trigger.isDisabled() != true {
		t.Errorf("RunAction('trigger.test-action-trigger.disable.1s') failed, got %v, want %v", trigger.isDisabled(), true)
	}
	time.Sleep(1 * time.Second)
	if trigger.isDisabled() != false {
		t.Errorf("RunAction('trigger.test-action-trigger.disable.1s') failed, got %v, want %v", trigger.isDisabled(), false)
	}

	// Alert action
	sensor := &Sensor{Id: "humidity"}
	config.Sensor = append(config.Sensor, sensor)
	alert := &Alert{Id: "test-action-alert", Sensor: "humidity"}
	config.Alert = append(config.Alert, alert)
	// Test case 1: disable alert
	RunAction("alert.test-action-alert.disable", "")
	if alert.isDisabled() != true {
		t.Errorf("RunAction('alert.test-action-alert.disable') failed, got %v, want %v", alert.isDisabled(), true)
	}
	// Test case 2: enabling alert
	RunAction("alert.test-action-alert.enable", "")
	if alert.isDisabled() != false {
		t.Errorf("RunAction('alert.test-action-alert.enable') failed, got %v, want %v", alert.isDisabled(), false)
	}
	// Test case 3: disabling alert with duration
	RunAction("alert.test-action-alert.disable.1s", "")
	if alert.isDisabled() != true {
		t.Errorf("RunAction('alert.test-action-alert.disable.1s') failed, got %v, want %v", alert.isDisabled(), true)
	}
	time.Sleep(1 * time.Second)
	if alert.isDisabled() != false {
		t.Errorf("RunAction('alert.test-action-alert.disable.1s') failed, got %v, want %v", alert.isDisabled(), false)
	}

	// reset
	config.DryRun = false
}
