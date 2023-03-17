package main

import (
	"log"
	"strings"
	"time"
)

func GetAlert(sensor string) *Alert {
	for _, a := range config.Alert {
		if a.Sensor == sensor {
			return a
		}
	}
	log.Fatalf("Unable to alert on sensor '%s', not found in configuration.yaml", sensor)
	return &Alert{}
}

func InitAlerting() {
	for _, a := range config.Alert {
		go a.monitor()
	}
}

func (a *Alert) monitor() {
	// maybe add a sleep here for startup, dont want alerts straight away
	for {
		s := GetSensor(a.Sensor)
		value := s.GetValue()
		if !isSunset() && !isSunrise() && value != 0 {
			// don't alert between sunset/sunrise or if value is 0
			if isDayTime() {
				// day time
				if value > a.When.Day.Above {
					a.Failing("%s is currently %v%s/%v%s", strings.Title(a.Sensor), value, s.Unit, a.When.Day.Above, s.Unit)
				} else if value < a.When.Day.Below {
					a.Failing("%s is currently %v%s/%v%s", strings.Title(a.Sensor), value, s.Unit, a.When.Day.Below, s.Unit)
				} else {
					// clear alerts
					a.Clear()
				}
			} else {
				// night time
				if value > a.When.Night.Above {
					a.Failing("%s is currently %v%s/%v%s", strings.Title(a.Sensor), value, s.Unit, a.When.Night.Above, s.Unit)
				} else if value < a.When.Night.Below {
					a.Failing("%s is currently %v%s/%v%s", strings.Title(a.Sensor), value, s.Unit, a.When.Night.Below, s.Unit)
				} else {
					// clear alerts
					a.Clear()
				}
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

func (a *Alert) getFailTime() time.Time {
	return a.FailedTime
}

func (a *Alert) setFailTime(t time.Time) {
	a.FailedTime = t
}

func (a *Alert) isFailing() bool {
	failTime := a.getFailTime()
	if time.Now().After(failTime.Add(a.After)) {
		a.Clear()
		return true
	}
	return false
}

func (a *Alert) Failing(s string, v ...any) {
	emptyTime := time.Time{}
	failTime := a.getFailTime()
	if failTime == emptyTime {
		a.setFailTime(time.Now())
	} else if a.isFailing() {
		a.sendNotification(s, v...)
	}
}

func (a *Alert) Clear() {
	a.setFailTime(time.Time{})
}

func (a *Alert) sendNotification(s string, v ...any) {
	for _, nId := range a.Notification {
		n := GetNotification(nId)
		n.SendNotification(s, v...)
	}
}
