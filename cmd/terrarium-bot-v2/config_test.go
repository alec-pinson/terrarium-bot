package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// set up environment variables
	os.Setenv("DEBUG", "true")
	os.Setenv("DRY_RUN", "true")
	os.Setenv("CONFIG_FILE", "test.yaml")
	os.Setenv("NOTIFICATION_USER_TOKEN", "user123")
	os.Setenv("NOTIFICATION_API_TOKEN", "api123")
	os.Setenv("USE_IN_MEMORY_STATUS", "true")

	// create the test yaml file
	configFileYaml := `
day:
  start: "06:00"
  action: ["doSomething"]
night:
  start: "22:00"
  action: ["doSomething"]
sunrise:
  start: "06:22"
  action: ["doSomething"]
sunset:
  start: "21:30"
  action: ["doSomething"]
trigger:
  - id: "trigger1"
    sensor: "sensor1"
    endpoint: "http://localhost:8080"
    when:
      day:
        below: 5
      night:
        below: 5
    action: ["doSomething"]
    else: ["doSomethingElse"]
  - id: "trigger2"
    sensor: "sensor2"
    endpoint: "http://anotherhost:8080"
    when:
      day:
        above: 100
        every: 10s
    action: ["doSomething", "doSomethingElse"]
    else: ["doSomethingElse", "doSomethingElseElse"]
switch:
  - id: "switch1"
    on: "http://localhost:8081/on"
    off: "http://localhost:8081/off"
    insecure: true
  - id: "switch2"
    on: "http://localhost:8081/on"
    off: "http://localhost:8081/off"
    insecure: false
sensor:
  - id: "sensor1"
    url: "http://localhost:8888"
    jsonPath: "json.path"
    unit: "unit"
  - id: "sensor2"
    url: "http://localhost:9999"
    jsonPath: "json.path"
    unit: "unit"
notification:
  - id: "notification1"
    antiSpam: "30s"
    device: "my_device"
    userToken: "NOTIFICATION_USER_TOKEN"
    apiToken: "NOTIFICATION_API_TOKEN"
alert:
  - id: "alert1"
    sensor: "sensor1"
    when:
      day:
        below: 10
      night:
        below: 5
    after: "20m"
    notification: ["notification1"]
useInMemoryStatus: true`
	err := os.WriteFile("test.yaml", []byte(configFileYaml), 0644)
	if err != nil {
		t.Errorf("Error creating test.yaml: %v", err)
	}

	// load the configuration
	loadedConfig := config.Load()

	// ensure loaded configuration matches expected values
	assert.True(t, loadedConfig.Debug)
	assert.True(t, loadedConfig.DryRun)
	assert.Equal(t, "test.yaml", loadedConfig.File)

	assert.Equal(t, "06:00", loadedConfig.Day.Start)
	assert.Equal(t, []string{"doSomething"}, loadedConfig.Day.Action)
	dayStartTime, err := time.Parse("15:04", "06:00")
	assert.NoError(t, err)
	assert.Equal(t, dayStartTime, loadedConfig.Day.StartTime)

	assert.Equal(t, "22:00", loadedConfig.Night.Start)
	assert.Equal(t, []string{"doSomething"}, loadedConfig.Night.Action)
	nightStartTime, err := time.Parse("15:04", "22:00")
	assert.NoError(t, err)
	assert.Equal(t, nightStartTime, loadedConfig.Night.StartTime)

	assert.Equal(t, "06:22", loadedConfig.Sunrise.Start)
	assert.Equal(t, []string{"doSomething"}, loadedConfig.Sunrise.Action)
	sunriseStartTime, err := time.Parse("15:04", "06:22")
	assert.NoError(t, err)
	assert.Equal(t, sunriseStartTime, loadedConfig.Sunrise.StartTime)

	assert.Equal(t, "21:30", loadedConfig.Sunset.Start)
	assert.Equal(t, []string{"doSomething"}, loadedConfig.Sunset.Action)
	sunsetStartTime, err := time.Parse("15:04", "21:30")
	assert.NoError(t, err)
	assert.Equal(t, sunsetStartTime, loadedConfig.Sunset.StartTime)

	assert.Len(t, loadedConfig.Trigger, 2)
	assert.Equal(t, "trigger1", loadedConfig.Trigger[0].Id)
	assert.Equal(t, "sensor1", loadedConfig.Trigger[0].Sensor)
	assert.Equal(t, "http://localhost:8080", loadedConfig.Trigger[0].Endpoint)
	assert.Equal(t, 5, loadedConfig.Trigger[0].When.Day.Below)
	assert.Equal(t, 5, loadedConfig.Trigger[0].When.Night.Below)
	assert.Equal(t, []string{"doSomething"}, loadedConfig.Trigger[0].Action)
	assert.Equal(t, []string{"doSomethingElse"}, loadedConfig.Trigger[0].Else)

	assert.Len(t, loadedConfig.Switch, 2)
	assert.Equal(t, "switch1", loadedConfig.Switch[0].Id)
	assert.Equal(t, "http://localhost:8081/on", loadedConfig.Switch[0].On)
	assert.Equal(t, "http://localhost:8081/off", loadedConfig.Switch[0].Off)
	assert.True(t, loadedConfig.Switch[0].Insecure)

	assert.Len(t, loadedConfig.Sensor, 2)
	assert.Equal(t, "sensor1", loadedConfig.Sensor[0].Id)
	assert.Equal(t, "http://localhost:8888", loadedConfig.Sensor[0].Url)
	assert.False(t, loadedConfig.Sensor[0].Insecure)
	assert.Equal(t, "json.path", loadedConfig.Sensor[0].JsonPath)
	assert.Equal(t, "unit", loadedConfig.Sensor[0].Unit)

	assert.Len(t, loadedConfig.Notification, 1)
	assert.Equal(t, "notification1", loadedConfig.Notification[0].Id)
	assert.Equal(t, 30*time.Second, loadedConfig.Notification[0].AntiSpam)
	assert.Equal(t, "my_device", loadedConfig.Notification[0].Device)
	assert.Equal(t, "NOTIFICATION_USER_TOKEN", loadedConfig.Notification[0].UserToken)
	assert.Equal(t, "NOTIFICATION_API_TOKEN", loadedConfig.Notification[0].APIToken)
	assert.Equal(t, "user123", loadedConfig.Notification[0].UserTokenValue)
	assert.Equal(t, "api123", loadedConfig.Notification[0].APITokenValue)

	assert.Len(t, loadedConfig.Alert, 1)
	assert.Equal(t, "alert1", loadedConfig.Alert[0].Id)
	assert.Equal(t, "sensor1", loadedConfig.Alert[0].Sensor)
	assert.Equal(t, 10, loadedConfig.Alert[0].When.Day.Below)
	assert.Equal(t, 5, loadedConfig.Alert[0].When.Night.Below)
	assert.Equal(t, 20*time.Minute, loadedConfig.Alert[0].After)
	assert.Equal(t, []string{"notification1"}, loadedConfig.Alert[0].Notification)

	assert.True(t, loadedConfig.UseInMemoryStatus)

	// clean up
	os.Remove("test.yaml")
}
