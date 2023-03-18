package main

import (
	"testing"
	"time"
)

func TestRunAction(t *testing.T) {
	// disable http calls when turning on/off switches
	config.DryRun = true

	// Test sleep action
	start := time.Now()
	RunAction("sleep.1s", "")
	sleepTime := time.Since(start).Seconds()
	if sleepTime < 1 || sleepTime > 2 {
		t.Errorf("Expected sleepingAction to be 2s, but got %v", sleepTime)
	}

	s := &Switch{Id: "test-action-switch"}
	config.Switch = append(config.Switch, s)

	// Test the switch action turning on
	RunAction("switch.test-action-switch.on", "")
	if s.getStatus() != "on" {
		t.Errorf("RunAction('switch.test-action-switch.on') failed, got %v, want %v", s.getStatus(), "on")
	}

	// Test the switch action turning off and then disabling
	s.setStatus("on")
	RunAction("switch.test-action-switch.off", "")
	if s.getStatus() != "off" {
		t.Errorf("RunAction('switch.test-action-switch.off') failed, got %v, want %v", s.getStatus(), "off")
	}
	RunAction("switch.test-action-switch.disable.10m", "")
	if s.isDisabled() != true {
		t.Errorf("RunAction('switch.test-action-switch.disable.10m') failed, got %v, want %v", s.isDisabled(), true)
	}
	RunAction("switch.test-action-switch.enable", "")
	if s.isDisabled() != false {
		t.Errorf("RunAction('switch.test-action-switch.enable') failed, got %v, want %v", s.isDisabled(), false)
	}
	RunAction("switch.test-action-switch.disable", "")
	if s.isDisabled() != true {
		t.Errorf("RunAction('switch.test-action-switch.disable') failed, got %v, want %v", s.isDisabled(), true)
	}

	// reset
	config.DryRun = false
}
