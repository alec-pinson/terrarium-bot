package main

import (
	"time"
)

func InitTriggers() {
	for _, t := range config.Trigger {
		if t.Sensor != "" {
			go t.monitorSensor()
		}
	}
}

func (t Trigger) monitorSensor() {
	s, _ := getSensor(t.Sensor)
	var runAction bool = false
	for {
		value := s.GetValue()
		runAction = false
		if value != 0 {
			// don't do anything if value is 0
			if isDayTime() {
				// day time
				if value > t.When.Day.Above && t.When.Day.Above != 0 {
					// trigger action
					runAction = true
				}
				if value < t.When.Day.Below && t.When.Day.Below != 0 {
					// trigger action
					runAction = true
				}
			} else {
				// night time
				if value > t.When.Night.Above && t.When.Night.Above != 0 {
					// trigger action
					runAction = true
				}
				if value < t.When.Night.Below && t.When.Night.Below != 0 {
					// trigger action
					runAction = true
				}
			}
			if runAction {
				t.doAction()
			} else {
				t.doElseAction()
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func isTriggerEndpoint(endpoint string) (bool, Trigger) {
	for _, t := range config.Trigger {
		if t.Endpoint == "/"+endpoint {
			return true, t
		}
	}
	return false, Trigger{}
}

func (t Trigger) doAction() {
	for _, a := range t.Action {
		RunAction(a)
	}
}

func (t Trigger) doElseAction() {
	for _, a := range t.Else {
		RunAction(a)
	}
}
