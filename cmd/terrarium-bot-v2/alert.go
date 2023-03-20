package main

import (
	"log"
	"time"
)

func GetAlert(id string) *Alert {
	for _, a := range config.Alert {
		if a.Id == id {
			if GetSensor(a.Sensor) == nil {
				log.Fatalf("Unable to alert on '%s' as sensor '%s', not found in configuration.yaml", id, a.Sensor)
				return nil
			}
			return a
		}
	}
	log.Fatalf("Alert for '%s' not found in configuration.yaml", id)
	return nil
}

func InitAlerting() {
	for _, a := range config.Alert {
		go a.monitor()
	}
}

func (a *Alert) monitor() {
	// maybe add a sleep here for startup, dont want alerts straight away
	s := GetSensor(a.Sensor)
	for {
		value := s.GetValue()
		if !isSunset() && !isSunrise() && value != 0 {
			// don't alert between sunset/sunrise or if value is 0
			if isDayTime() {
				// day time
				if value > a.When.Day.Above {
					a.Failing("%v%s/%v%s", value, s.Unit, a.When.Day.Above, s.Unit)
				} else if value < a.When.Day.Below {
					a.Failing("%v%s/%v%s", value, s.Unit, a.When.Day.Below, s.Unit)
				} else {
					// clear alerts
					a.Clear()
				}
			} else {
				// night time
				if value > a.When.Night.Above {
					a.Failing("%v%s/%v%s", value, s.Unit, a.When.Night.Above, s.Unit)
				} else if value < a.When.Night.Below {
					a.Failing("%v%s/%v%s", value, s.Unit, a.When.Night.Below, s.Unit)
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

func (a *Alert) Enable(reason string) {
	a.Disabled = 0
	log.Printf("Alert Enabled: '%s'", a.Id)
}

func (a *Alert) Disable(duration string, reason string) {
	if duration == "" {
		// 10 years.. 'forever'
		duration = "87660h"
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Printf("Invalid alert disable duration '%s'", duration)
		return
	}
	a.DisabledAt = time.Now()
	a.Disabled = d
	if duration == "87660h" {
		log.Printf("Alert Disabled: '%s'", a.Id)
	} else {
		log.Printf("Alert Disabled: '%s' for %s", a.Id, d)
	}
}

func (a *Alert) isDisabled() bool {
	if a.Disabled == 0 {
		return false
	}
	if a.DisabledAt.Add(a.Disabled).After(time.Now()) {
		return true
	}
	return false
}
