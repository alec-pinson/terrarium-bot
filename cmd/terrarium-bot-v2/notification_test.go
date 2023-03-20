package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
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

	// configure notification
	n := Notification{
		Id:       "pushover",
		Device:   "some-device",
		AntiSpam: 1 * time.Second,
	}

	// catch log output
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Test case 1: Send notification
	n.SendNotification("Test notification")
	if !strings.Contains(buf.String(), "Alert:") {
		t.Errorf("Alert should have sent, log: %q", buf.String())
	}

	// Test case 2: Test anti-spam prevents sending again
	buf.Reset()
	n.SendNotification("Test notification")
	if strings.Contains(buf.String(), "Alert:") {
		t.Errorf("Notification should not have sent again due to anti-spam, log: %q", buf.String())
	}

	// Test case 3: Test anti-spam has ended
	buf.Reset()
	time.Sleep(3 * time.Second)
	n.SendNotification("Test notification")
	if !strings.Contains(buf.String(), "Alert:") {
		t.Errorf("Notification should have sent now anti-spam time has expired, log: %q", buf.String())
	}

	// reset
	config.DryRun = false
	log.SetOutput(os.Stderr)
}
