package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetSwitch(t *testing.T) {
	config.Switch = []*Switch{
		{Id: "heating"},
		{Id: "fan"},
	}
	expected := config.Switch[0]
	result := GetSwitch("heating")

	if result != expected {
		t.Errorf("Expected '%+v' but got '%+v'", expected, result)
	}
}

// Switch set last action function tests
func TestSwitchSetLastAction(t *testing.T) {
	s := &Switch{}

	// Set last action and ensure it's set correctly
	s.SetLastAction()
	if s.LastAction.IsZero() {
		t.Errorf("Last action was not set")
	}
}

// Switch get status and set status function tests
func TestSwitchGetSetStatus(t *testing.T) {
	s := &Switch{}

	// Set status to on and ensure that it's set correctly
	s.setStatus("on")
	if s.Status != "on" {
		t.Errorf("Status was not set correctly")
	}

	// Get status and ensure that it's returned correctly
	status := s.getStatus()
	if status != s.Status {
		t.Errorf("getStatus did not return the correct status")
	}
}

// Switch set disable function tests
func TestSwitchDisable(t *testing.T) {
	s := &Switch{Id: "disable"}

	// Set disable and ensure that it's set correctly
	s.Disable("1m", "")
	if s.Disabled != 1*time.Minute {
		t.Errorf("Disable was not set correctly")
	}
}

func TestSwitchIsDisabled(t *testing.T) {
	// Test case 1: switch should not be disabled
	s1 := Switch{LastAction: time.Now()}
	s1.Disable("", "")
	s1.Enable("")
	if s1.isDisabled() != false {
		t.Errorf("Test case 1 failed: expected false but got true")
	}

	// Test case 2: switch should be disabled for 2 seconds
	s1.Disable("2s", "")
	time.Sleep(1 * time.Second) // sleep for 1 second
	if s1.isDisabled() != true {
		t.Errorf("Test case 2 failed: expected true but got false")
	}
	time.Sleep(3 * time.Second) // sleep for 3 seconds
	// should not be disabled anymore
	if s1.isDisabled() != false {
		t.Errorf("Test case 2 failed: expected false but got true")
	}
}

// Switch set on URL and set off URL function tests
func TestFixURLs(t *testing.T) {
	// set env variable
	os.Setenv("FIX_URL_VAL", "testing")

	s := &Switch{
		On:  "http://test.com/$FIX_URL_VAL/on",
		Off: "http://test.com/$FIX_URL_VAL/off",
	}

	// Fix URLs and ensure that they're fixed correctly
	s.fixURLs()
	if s.On != "http://test.com/testing/on" {
		t.Errorf("Expected environment variables were not replaced for On URL")
	}
	if s.Off != "http://test.com/testing/off" {
		t.Errorf("Expected environment variables were not replaced for Off URL")
	}

	// reset
	os.Unsetenv("FIX_URL_VAL")
}

// Switch turn on function tests
func TestSwitchTurnOn(t *testing.T) {
	// disable http calls when turning on/off switches
	config.DryRun = true

	s := &Switch{
		Id:  "on-test",
		On:  "http://test.com/on",
		Off: "http://test.com/off",
	}

	// Test case 1: Switch turns on
	s.TurnOff("") // ensure switch is off
	s.TurnOn("", "")

	if s.getStatus() != "on" {
		t.Errorf("Switch did not turn on as expected")
	}

	// Test case 2: Switch doesn't turn on again when already on with UseInMemoryStatus enabled
	var buf bytes.Buffer
	log.SetOutput(&buf)
	config.UseInMemoryStatus = true
	s.TurnOn("", "")

	if strings.Contains(buf.String(), "Switch On 'on-test'") {
		t.Errorf("Switch turned on again while 'USE_IN_MEMORY_STATUS' was set")
	}

	// Test case 3: Switch turns on when already on with UseInMemoryStatus disabled
	config.UseInMemoryStatus = false
	buf.Reset()
	s.TurnOn("", "")
	if !strings.Contains(buf.String(), "Switch On: 'on-test'") {
		t.Errorf("Switch did not turn on again while 'USE_IN_MEMORY_STATUS' was unset, got %v", buf.String())
	}

	// Test case 4: Switch turns off after 'for' duration
	s.TurnOff("")         // ensure switch is off
	go s.TurnOn("2s", "") // turn on for 2 seconds
	time.Sleep(1 * time.Second)
	if s.getStatus() != "on" {
		t.Errorf("Switch did not turn on for 2 seconds as expected")
	}
	time.Sleep(3 * time.Second)
	if s.getStatus() != "off" {
		t.Errorf("Switch did not turn off after 2 seconds as expected")
	}

	// Test case 5: Switch does not turn on if recently turned off and Disable is specified
	s.TurnOff("")
	s.Disable("10m", "")
	s.TurnOn("", "")
	if s.getStatus() == "on" {
		t.Errorf("Switch was disabled and should not have turned on")
	}

	//reset
	config.DryRun = false
	log.SetOutput(os.Stderr)
}

// Switch turn off function tests
func TestSwitchTurnOff(t *testing.T) {
	// disable http calls when turning on/off switches
	config.DryRun = true

	s := &Switch{
		Id:  "test",
		On:  "http://test.com/on",
		Off: "http://test.com/off",
	}

	// Test that the switch turns off
	s.TurnOn("", "") // ensure switch is on
	s.TurnOff("")

	if s.getStatus() != "off" {
		t.Errorf("Switch did not turn off as expected")
	}

	// reset
	config.DryRun = false
}
