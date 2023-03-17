package main

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// set up environment variables
	os.Setenv("DEBUG", "true")
	os.Setenv("DRY_RUN", "true")
	os.Setenv("CONFIG_FILE", "test.yaml")
	os.Setenv("NOTIFICATION_USER_TOKEN", "user123")
	os.Setenv("NOTIFICATION_API_TOKEN", "api123")
	os.Setenv("USE_IN_MEMORY_STATUS", "false")

	// create test configuration
	expectedConfig := Config{
		Debug:   true,
		DryRun:  true,
		File:    "test.yaml",
		Day:     StartAction{Start: "08:00", Action: []string{"do this", "then that"}},
		Night:   StartAction{Start: "22:00", Action: []string{"do that", "then this"}},
		Sunrise: StartAction{Start: "06:00", Action: []string{"do something"}},
		Sunset:  StartAction{Start: "18:00", Action: []string{"do something else"}},
		Trigger: []*Trigger{{Sensor: "sensor1", Endpoint: "endpoint1", When: When{}, Action: []string{"do this"}}},
		Switch:  []*Switch{{Id: "switch1", On: "on", For: 5 * time.Minute, Every: 10 * time.Minute, Off: "off", Disable: 0}},
		Sensor:  []*Sensor{{Id: "sensor1", Url: "http://sensor1.com", JsonPath: "path1", Unit: "unit1"}},
		Notification: []*Notification{{
			Id:               "notif1",
			AntiSpam:         10 * time.Minute,
			Device:           "device1",
			UserToken:        "NOTIFICATION_USER_TOKEN",
			APIToken:         "NOTIFICATION_API_TOKEN",
			UserTokenValue:   "",
			APITokenValue:    "",
			LastNotification: time.Time{},
		}},
		Alert:             []*Alert{{Sensor: "sensor1", When: When{}, After: 30 * time.Minute, Notification: []string{"notif1"}}},
		UseInMemoryStatus: false,
	}

	// convert times
	expectedConfig.Day.StartTime, _ = time.Parse("15:04", expectedConfig.Day.Start)
	expectedConfig.Night.StartTime, _ = time.Parse("15:04", expectedConfig.Night.Start)
	expectedConfig.Sunrise.StartTime, _ = time.Parse("15:04", expectedConfig.Sunrise.Start)
	expectedConfig.Sunset.StartTime, _ = time.Parse("15:04", expectedConfig.Sunset.Start)

	// create the test yaml file
	data := `
day:
  start: 08:00
  action:
    - do this
    - then that
night:
  start: 22:00
  action:
    - do that
    - then this
sunrise:
  start: 06:00
  action:
    - do something
sunset:
  start: 18:00
  action:
    - do something else
trigger:
  - sensor: sensor1
    endpoint: endpoint1
    when:
      day:
        below: 50
        above: 100
      night:
        below: 20
        above: 80
    action:
      - do this
    else:
      - do that
switch:
  - id: switch1
    on: on
    for: 5m
    every: 10m
    off: off
    disable: 0
sensor:
  - id: sensor1
    url: http://sensor1.com
    jsonPath: path1
    unit: unit1
notification:
  - id: notif1
    antiSpam: 10m
    device: device1
    userToken: NOTIFICATION_USER_TOKEN
    apiToken: NOTIFICATION_API_TOKEN
    lastNotification: ""
alert:
  - sensor: sensor1
    when:
      day:
        below: 50
        above: 100
      night:
        below: 20
        above: 80
    after: 30m
    notification:
      - notif1
useInMemoryStatus: false`
	err := os.WriteFile("test.yaml", []byte(data), 0644)
	if err != nil {
		t.Errorf("Error creating test.yaml: %v", err)
	}

	// load the configuration
	loadedConfig := expectedConfig.Load()

	// ensure loaded configuration matches expected values
	if loadedConfig.Debug != expectedConfig.Debug {
		t.Errorf("Expected Debug %v, got %v", expectedConfig.Debug, loadedConfig.Debug)
	}
	if loadedConfig.DryRun != expectedConfig.DryRun {
		t.Errorf("Expected DryRun %v, got %v", expectedConfig.DryRun, loadedConfig.DryRun)
	}
	if loadedConfig.File != expectedConfig.File {
		t.Errorf("Expected File %v, got %v", expectedConfig.File, loadedConfig.File)
	}
	if loadedConfig.Day.Start != expectedConfig.Day.Start {
		t.Errorf("Expected Day.Start %v, got %v", expectedConfig.Day.Start, loadedConfig.Day.Start)
	}
	if !loadedConfig.Day.StartTime.Equal(expectedConfig.Day.StartTime) {
		t.Errorf("Expected Day.StartTime %v, got %v", expectedConfig.Day.StartTime, loadedConfig.Day.StartTime)
	}
	if len(loadedConfig.Day.Action) != len(expectedConfig.Day.Action) {
		t.Errorf("Expected %v actions, got %v", len(expectedConfig.Day.Action), len(loadedConfig.Day.Action))
	}
	if loadedConfig.Night.Start != expectedConfig.Night.Start {
		t.Errorf("Expected Night.Start %v, got %v", expectedConfig.Night.Start, loadedConfig.Night.Start)
	}
	if !loadedConfig.Night.StartTime.Equal(expectedConfig.Night.StartTime) {
		t.Errorf("Expected Night.StartTime %v, got %v", expectedConfig.Night.StartTime, loadedConfig.Night.StartTime)
	}
	if len(loadedConfig.Night.Action) != len(expectedConfig.Night.Action) {
		t.Errorf("Expected %v actions, got %v", len(expectedConfig.Night.Action), len(loadedConfig.Night.Action))
	}
	if loadedConfig.Sunrise.Start != expectedConfig.Sunrise.Start {
		t.Errorf("Expected Sunrise.Start %v, got %v", expectedConfig.Sunrise.Start, loadedConfig.Sunrise.Start)
	}
	if !loadedConfig.Sunrise.StartTime.Equal(expectedConfig.Sunrise.StartTime) {
		t.Errorf("Expected Sunrise.StartTime %v, got %v", expectedConfig.Sunrise.StartTime, loadedConfig.Sunrise.StartTime)
	}
	if len(loadedConfig.Sunrise.Action) != len(expectedConfig.Sunrise.Action) {
		t.Errorf("Expected %v actions, got %v", len(expectedConfig.Sunrise.Action), len(loadedConfig.Sunrise.Action))
	}
	if loadedConfig.Sunset.Start != expectedConfig.Sunset.Start {
		t.Errorf("Expected Sunset.Start %v, got %v", expectedConfig.Sunset.Start, loadedConfig.Sunset.Start)
	}
	if !loadedConfig.Sunset.StartTime.Equal(expectedConfig.Sunset.StartTime) {
		t.Errorf("Expected Sunset.StartTime %v, got %v", expectedConfig.Sunset.StartTime, loadedConfig.Sunset.StartTime)
	}
	if len(loadedConfig.Sunset.Action) != len(expectedConfig.Sunset.Action) {
		t.Errorf("Expected %v actions, got %v", len(expectedConfig.Sunset.Action), len(loadedConfig.Sunset.Action))
	}
	if len(loadedConfig.Trigger) != len(expectedConfig.Trigger) {
		t.Errorf("Expected %v triggers, got %v", len(expectedConfig.Trigger), len(loadedConfig.Trigger))
	}
	if len(loadedConfig.Switch) != len(expectedConfig.Switch) {
		t.Errorf("Expected %v switches, got %v", len(expectedConfig.Switch), len(loadedConfig.Switch))
	}
	if len(loadedConfig.Sensor) != len(expectedConfig.Sensor) {
		t.Errorf("Expected %v sensors, got %v", len(expectedConfig.Sensor), len(loadedConfig.Sensor))
	}
	if len(loadedConfig.Notification) != len(expectedConfig.Notification) {
		t.Errorf("Expected %v notifications, got %v", len(expectedConfig.Notification), len(loadedConfig.Notification))
	}
	if len(loadedConfig.Alert) != len(expectedConfig.Alert) {
		t.Errorf("Expected %v alerts, got %v", len(expectedConfig.Alert), len(loadedConfig.Alert))
	}
	if loadedConfig.UseInMemoryStatus != expectedConfig.UseInMemoryStatus {
		t.Errorf("Expected UseInMemoryStatus %v, got %v", expectedConfig.UseInMemoryStatus, loadedConfig.UseInMemoryStatus)
	}
	if loadedConfig.Notification[0].UserTokenValue != "user123" {
		t.Errorf("Expected UserTokenValue %v, got %v", "user123", loadedConfig.Notification[0].UserTokenValue)
	}
	if loadedConfig.Notification[0].APITokenValue != "api123" {
		t.Errorf("Expected APITokenValue %v, got %v", "api123", loadedConfig.Notification[0].APITokenValue)
	}

	// clean up test files
	os.Remove("test.yaml")
}
