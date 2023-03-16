package main

import (
	"time"
)

func isTimeBetween(t1, t2 time.Time) bool {
	now, _ := time.Parse("15:04", time.Now().Format("15:04"))

	if now.After(t1) && now.Before(t2) {
		return true
	} else if now.Equal(t1) {
		return true
	} else {
		return false
	}
}

func isDayTime() bool {
	return isTimeBetween(config.Day.StartTime, config.Night.StartTime)
}

func isSunrise() bool {
	return isTimeBetween(config.Sunrise.StartTime, config.Day.StartTime)
}

func isSunset() bool {
	return isTimeBetween(config.Sunset.StartTime, config.Night.StartTime)
}

func InitTime() {
	go monitorTime()
}

func monitorTime() {
	for {
		time.Sleep(sleepingAction)
		if isSunrise() {
			Debug("Setting Sunrise configuration")
			for _, a := range config.Sunrise.Action {
				RunAction(a)
			}
		} else if isSunset() {
			Debug("Setting Sunset configuration")
			for _, a := range config.Sunset.Action {
				RunAction(a)
			}
		} else if isDayTime() {
			Debug("Setting Day Time configuration")
			for _, a := range config.Day.Action {
				RunAction(a)
			}
		} else {
			Debug("Setting Night Time configuration")
			for _, a := range config.Night.Action {
				RunAction(a)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
