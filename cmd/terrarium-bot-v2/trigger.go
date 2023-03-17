package main

import (
	"strconv"
	"strings"
	"time"
)

func InitTriggers() {
	for _, t := range config.Trigger {
		if t.Sensor != "" {
			go t.monitorSensor()
		}
	}
}

func GenerateReason(sensor string, value int, unit string, maxValue int) string {
	return strings.Title(sensor) + " currently " + strconv.Itoa(value) + unit + "/" + strconv.Itoa(maxValue) + unit
}

func (t *Trigger) monitorSensor() {
	s := GetSensor(t.Sensor)
	var runAction bool = false
	var reason string = ""
	for {
		value := s.GetValue()
		runAction = false
		reason = ""
		if value != 0 {
			// don't do anything if value is 0
			if isDayTime() {
				// day time
				if value > t.When.Day.Above && t.When.Day.Above != 0 {
					// trigger action
					runAction = true
					reason = GenerateReason(t.Sensor, value, s.Unit, t.When.Day.Above)
				}
				if value < t.When.Day.Below && t.When.Day.Below != 0 {
					// trigger action
					runAction = true
					reason = GenerateReason(t.Sensor, value, s.Unit, t.When.Day.Below)
				}
			} else {
				// night time
				if value > t.When.Night.Above && t.When.Night.Above != 0 {
					// trigger action
					runAction = true
					reason = GenerateReason(t.Sensor, value, s.Unit, t.When.Night.Above)
				}
				if value < t.When.Night.Below && t.When.Night.Below != 0 {
					// trigger action
					runAction = true
					reason = GenerateReason(t.Sensor, value, s.Unit, t.When.Night.Below)
				}
			}
			if runAction {
				t.doAction(reason)
			} else {
				if isDayTime() {
					// day time
					if t.When.Day.Above != 0 {
						reason = GenerateReason(t.Sensor, value, s.Unit, t.When.Day.Above)
					}
					if t.When.Day.Below != 0 {
						reason = GenerateReason(t.Sensor, value, s.Unit, t.When.Day.Below)
					}
				} else {
					// night time
					if t.When.Night.Above != 0 {
						reason = GenerateReason(t.Sensor, value, s.Unit, t.When.Night.Above)
					}
					if t.When.Night.Below != 0 {
						reason = GenerateReason(t.Sensor, value, s.Unit, t.When.Night.Below)
					}
				}

				t.doElseAction(reason)
			}
		}
		if isTesting {
			return
		}
		time.Sleep(1 * time.Minute)
	}
}

func isTriggerEndpoint(endpoint string) (bool, *Trigger) {
	for _, t := range config.Trigger {
		if t.Endpoint == "/"+endpoint {
			return true, t
		}
	}
	return false, &Trigger{}
}

func (t *Trigger) doAction(reason string) {
	for _, a := range t.Action {
		RunAction(a, reason)
	}
}

func (t *Trigger) doElseAction(reason string) {
	for _, a := range t.Else {
		RunAction(a, reason)
	}
}
