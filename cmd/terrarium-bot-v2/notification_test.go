package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestGetNotification(t *testing.T) {
	config.Notification = []*Notification{
		{Id: "pushover"},
		{Id: "email"},
	}
	expected := config.Notification[0]
	result := GetNotification("pushover")

	if result != expected {
		t.Errorf("Expected '%+v' but got '%+v'", expected, result)
	}
}

func TestSetLastNotification(t *testing.T) {
	n := Notification{Id: "pushover", Device: "Device-1"}
	n.setLastNotification()

	if n.LastNotification.IsZero() {
		t.Error("Expected LastNotification to be set, but it was not")
	}
}

func TestSendNotification(t *testing.T) {
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
