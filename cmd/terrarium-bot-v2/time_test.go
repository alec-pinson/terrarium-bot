package main

import (
	"testing"
	"time"
)

func TestIsTimeBetween(t *testing.T) {
	t1, _ := time.Parse("15:04", "00:00")
	t2, _ := time.Parse("15:04", "23:59")

	if !isTimeBetween(t1, t2) {
		t.Error("Expected time to be between t1 and t2")
	}

	t1, _ = time.Parse("15:04", "23:59")
	t2, _ = time.Parse("15:04", "00:00")

	if isTimeBetween(t1, t2) {
		t.Error("Expected time to be outside of t1 and t2")
	}
}

func TestIsDayTime(t *testing.T) {
	config.Day.StartTime, _ = time.Parse("15:04", "23:59")
	config.Night.StartTime, _ = time.Parse("15:04", "00:00")

	if isDayTime() {
		t.Error("Expected it to be night")
	}

	config.Day.StartTime, _ = time.Parse("15:04", "00:00")
	config.Night.StartTime, _ = time.Parse("15:04", "23:59")

	if !isDayTime() {
		t.Error("Expected it to be daytime")
	}
}

func TestIsSunrise(t *testing.T) {
	config.Sunrise.StartTime, _ = time.Parse("15:04", "00:00")
	config.Day.StartTime, _ = time.Parse("15:04", "23:59")

	if !isSunrise() {
		t.Error("Expected it to be sunrise")
	}

	config.Sunrise.StartTime, _ = time.Parse("15:04", "23:59")
	config.Day.StartTime, _ = time.Parse("15:04", "00:00")

	if isSunrise() {
		t.Error("Expected it to be morning")
	}
}

func TestIsSunset(t *testing.T) {
	config.Sunset.StartTime, _ = time.Parse("15:04", "00:00")
	config.Night.StartTime, _ = time.Parse("15:04", "23:59")

	if !isSunset() {
		t.Error("Expected it to be sunset")
	}

	config.Sunset.StartTime, _ = time.Parse("15:04", "23:59")
	config.Night.StartTime, _ = time.Parse("15:04", "00:00")

	if isSunset() {
		t.Error("Expected it to be night")
	}
}
