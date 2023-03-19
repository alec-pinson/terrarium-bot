package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetAlert(t *testing.T) {
	config.Sensor = []*Sensor{
		{Id: "temperature"},
		{Id: "humidity"},
	}
	config.Alert = []*Alert{
		{Id: "test-alert-1", Sensor: "temperature"},
		{Id: "test-alert-2", Sensor: "humidity"},
	}

	alert := GetAlert("test-alert-1")
	assert.NotNil(t, alert)

	alert = GetAlert("test-alert-2")
	assert.NotNil(t, alert)
}

func TestSetFailTime(t *testing.T) {
	a := &Alert{}
	a.setFailTime(time.Now())

	assert.NotNil(t, a.getFailTime())
}

func TestClear(t *testing.T) {
	a := &Alert{
		FailedTime: time.Now()}
	a.Clear()

	assert.Equal(t, time.Time{}, a.getFailTime())
}

func TestIsFailing(t *testing.T) {
	a := &Alert{
		FailedTime: time.Now(),
		After:      5 * time.Minute}

	assert.False(t, a.isFailing())

	a = &Alert{
		FailedTime: time.Now().Add(-10 * time.Minute),
		After:      5 * time.Minute}

	assert.True(t, a.isFailing())
}

func TestSendAlertNotification(t *testing.T) {
	// do not send a real notification
	config.DryRun = true
	config.Notification = []*Notification{
		{
			Id:     "pushover",
			Device: "Some-device",
		},
	}

	a := &Alert{Notification: []string{"pushover"}}
	var buf bytes.Buffer
	log.SetOutput(&buf)
	a.sendNotification("Test notification")
	if !strings.Contains(buf.String(), "Alert:") {
		t.Errorf("Log should contain 'Alert:' but doesn't, log: %q", buf.String())
	}

	// reset
	config.DryRun = false
	log.SetOutput(os.Stderr)
}
