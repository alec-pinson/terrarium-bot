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
	config.Alert = []*Alert{
		{Sensor: "temperature"},
		{Sensor: "humidity"},
	}

	alert := GetAlert("temperature")
	assert.NotNil(t, alert)

	alert = GetAlert("humidity")
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
	config.Debug = true
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
	if !strings.Contains(buf.String(), "Sent alert:") {
		t.Errorf("Log should contain 'Sent alert:' but doesn't, log: %q", buf.String())
	}

	// reset
	config.Debug = false
	log.SetOutput(os.Stderr)
}
