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

// // Switch monitor function tests
// func TestSwitchMonitor(t *testing.T) {
// 	// disable http calls when turning on/off switches
// 	config.DryRun = true
// 	// set testing mode so we exit the for loop
// 	isTesting = true

// 	s := &Switch{
// 		Id:            "test-monitor-switch",
// 		Every:         1 * time.Second,
// 		For:           0,
// 		Disable:       0,
// 		On:            "http://test.com/on",
// 		Off:           "http://test.com/off",
// 		DisableCustom: 0,
// 	}

// 	// Test case 1: Light should automatically turn on every 1 second during day time
// 	// ensure switch is off
// 	s.TurnOff("")

// 	// set last action to some time ago
// 	s.LastAction = time.Now().Add(-(5 * time.Minute))

// 	// ensure it is day time
// 	config.Day.StartTime, _ = time.Parse("15:04", "00:00")
// 	config.Night.StartTime, _ = time.Parse("15:04", "23:59")

// 	// begin monitor
// 	s.monitor()

// 	// check that the switch turns on
// 	if s.getStatus() != "on" {
// 		t.Errorf("Switch did not turn on as expected")
// 	}

// 	// Test case 2: Make sure switch cannot turn on during the night
// 	// ensure switch is off
// 	s.TurnOff("")

// 	// set last action to some time ago
// 	s.LastAction = time.Now().Add(-(5 * time.Minute))

// 	// ensure it is night time
// 	config.Day.StartTime, _ = time.Parse("15:04", "23:59")
// 	config.Night.StartTime, _ = time.Parse("15:04", "00:00")

// 	s.monitor()

// 	if s.getStatus() == "on" {
// 		t.Errorf("Switch should not turn on during the night")
// 	}

// 	// reset back to day time
// 	config.Day.StartTime, _ = time.Parse("15:04", "00:00")
// 	config.Night.StartTime, _ = time.Parse("15:04", "23:59")

// 	// Test case 3: Make sure switch does not turn on during the disabled duration
// 	// ensure switch is off
// 	s.TurnOff("")
// 	// set last action to 1 second ago
// 	s.LastAction = time.Now().Add(-(1 * time.Second))
// 	// Set disable duration to 20 seconds
// 	s.Disable = 2 * time.Second
// 	s.monitor()

// 	if s.getStatus() == "on" {
// 		t.Errorf("Switch turned on while disabled")
// 	}

// 	// Test case 4: Switch should turn on after disable duration passes
// 	time.Sleep(3 * time.Second)
// 	s.monitor()

// 	if s.getStatus() != "on" {
// 		t.Errorf("Switch should turn on again after disable duration passes")
// 	}

// 	// Test case 5: Make sure switch does not turn on during the custom disabled duration
// 	// ensure switch is off
// 	s.TurnOff("")
// 	// set last action to 1 second ago
// 	s.LastAction = time.Now().Add(-(1 * time.Second))
// 	// Set custom disable duration to 2 seconds
// 	s.SetDisableCustom("2s", "")
// 	s.monitor()

// 	if s.getStatus() == "on" {
// 		t.Errorf("Switch turned on while custom disable duration was set")
// 	}

// 	// Test case 6: Switch should turn on after custom disable duration passes
// 	time.Sleep(3 * time.Second)
// 	s.monitor()

// 	if s.getStatus() != "on" {
// 		t.Errorf("Switch should turn on again after custom disable duration passes")
// 	}

// 	// reset
// 	config.DryRun = false
// 	isTesting = false
// }

// Switch set last action function tests
func TestSwitchSetLastAction(t *testing.T) {
	s := &Switch{}

	// Set last action and ensure it's set correctly
	s.SetLastAction()
	if s.LastAction.IsZero() {
		t.Errorf("Last action was not set")
	}
}

// Switch get last action function tests
func TestSwitchGetLastAction(t *testing.T) {
	s := &Switch{
		LastAction: time.Now(),
	}

	// Get last action and ensure that it's returned correctly
	lastAction := s.GetLastAction()
	if lastAction != s.LastAction {
		t.Errorf("GetLastAction did not return the correct last action")
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
func TestSwitchSetOnOffUrl(t *testing.T) {
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

	// Set on URL and ensure that it's set correctly
	s.setOnUrl("http://test.com/on")
	if s.On != "http://test.com/on" {
		t.Errorf("On URL was not set correctly")
	}

	// Set off URL and ensure that it's set correctly
	s.setOffUrl("http://test.com/off")
	if s.Off != "http://test.com/off" {
		t.Errorf("Off URL was not set correctly")
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
