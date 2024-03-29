package main

import (
	"fmt"
	"log"
	"time"
)

func GetTrigger(id string) *Trigger {
	for _, t := range config.Trigger {
		if t.Id == id {
			return t
		}
	}
	log.Fatalf("Trigger '%s' not found in configuration.yaml", id)
	return nil
}

func InitTriggers() {
	for _, t := range config.Trigger {
		go t.monitor()
		time.Sleep(1 * time.Second) // stop triggers clashing during startup
	}
}

func GenerateReason(value float64, unit string, maxValue float64) string {
	return fmt.Sprintf("%.2f", value) + unit + "/" + fmt.Sprintf("%.2f", maxValue) + unit
}

func (t *Trigger) monitor() {
	var (
		s             *Sensor
		runAction     bool
		reason        string
		value         float64
		previousValue float64
	)

	// get trigger sensor if one is set
	if t.Sensor != "" {
		s = GetSensor(t.Sensor)
	}

	// initialize the last triggered time
	t.LastTriggered = time.Now()

	for {
		runAction = false
		reason = ""
		previousValue = 0
		value = 0

		if t.isDisabled() {
			Debug("Trigger %s is currently disabled", t.Id)
			if isTesting {
				return
			}
			time.Sleep(1 * time.Minute)
			continue
		}

		if s != nil {
			// get value from the sensor
			previousValue = s.GetPreviousValue()
			value = s.GetValue()
		}
		valueSet := value != 0
		dayAbove := t.When.Day.Above != 0
		dayBelow := t.When.Day.Below != 0
		nightAbove := t.When.Night.Above != 0
		nightBelow := t.When.Night.Below != 0
		dayEvery := t.When.Day.Every != 0
		nightEvery := t.When.Night.Every != 0

		// has the value dropped/increased by x or more
		dayDroppedBy := t.When.Day.DroppedBy != 0 && value <= previousValue-t.When.Day.DroppedBy
		nightDroppedBy := t.When.Night.DroppedBy != 0 && value <= previousValue-t.When.Night.DroppedBy
		dayIncreasedBy := t.When.Day.IncreasedBy != 0 && value >= previousValue+t.When.Day.IncreasedBy
		nightIncreasedBy := t.When.Night.IncreasedBy != 0 && value >= previousValue+t.When.Night.IncreasedBy

		// check triggers based on time of day and value
		if isDayTime() {
			if valueSet && (dayDroppedBy || dayIncreasedBy) {
				runAction = true
				reason = fmt.Sprintf("%.2f", previousValue) + s.Unit + "->" + fmt.Sprintf("%.2f", value) + s.Unit
			} else if valueSet && dayAbove && value > t.When.Day.Above {
				runAction = true
				reason = GenerateReason(value, s.Unit, t.When.Day.Above)
			} else if valueSet && dayBelow && value < t.When.Day.Below {
				runAction = true
				reason = GenerateReason(value, s.Unit, t.When.Day.Below)
			} else if dayEvery && t.LastTriggered.Before(time.Now().Add(-t.When.Day.Every)) {
				reason = "Trigger '" + t.Id + "' scheduled every " + t.When.Day.Every.String()
				runAction = true
			}
		} else {
			if valueSet && (nightDroppedBy || nightIncreasedBy) {
				runAction = true
				reason = fmt.Sprintf("%.2f", previousValue) + s.Unit + "->" + fmt.Sprintf("%.2f", value) + s.Unit
			} else if valueSet && nightAbove && value > t.When.Night.Above {
				runAction = true
				reason = GenerateReason(value, s.Unit, t.When.Night.Above)
			} else if valueSet && nightBelow && value < t.When.Night.Below {
				runAction = true
				reason = GenerateReason(value, s.Unit, t.When.Night.Below)
			} else if nightEvery && t.LastTriggered.Before(time.Now().Add(-t.When.Night.Every)) {
				reason = "Trigger '" + t.Id + "' scheduled every " + t.When.Night.Every.String()
				runAction = true
			}
		}

		// run actions/else actions
		if runAction {
			t.doAction(reason)
			t.LastTriggered = time.Now() // update the last triggered time
		} else if valueSet {
			// do else action
			if isDayTime() {
				// day time
				if dayAbove {
					reason = GenerateReason(value, s.Unit, t.When.Day.Above)
				} else if dayBelow {
					reason = GenerateReason(value, s.Unit, t.When.Day.Below)
				}
			} else {
				// night time
				if nightAbove {
					reason = GenerateReason(value, s.Unit, t.When.Night.Above)
				}
				if nightBelow {
					reason = GenerateReason(value, s.Unit, t.When.Night.Below)
				}
			}
			t.doElseAction(reason)
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

func (t *Trigger) Enable(reason string) {
	t.Disabled = 0
	log.Printf("Trigger Enabled: '%s'", t.Id)
}

func (t *Trigger) Disable(duration string, reason string) {
	if duration == "" {
		// 10 years.. 'forever'
		duration = "87660h"
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Printf("Invalid disable duration '%s'", duration)
		return
	}
	t.DisabledAt = time.Now()
	t.Disabled = d
	if duration == "87660h" {
		log.Printf("Trigger Disabled: '%s'", t.Id)
	} else {
		log.Printf("Trigger Disabled: '%s' for %s", t.Id, d)
	}
}

func (t *Trigger) isDisabled() bool {
	if t.Disabled == 0 {
		return false
	}
	if t.DisabledAt.Add(t.Disabled).Before(time.Now()) {
		return false
	}
	return true
}
