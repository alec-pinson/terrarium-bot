package main

import (
	"log"
	"strings"
	"time"
)

var sleepingAction time.Duration // used by time.go to pause normal actions until after this duration

func RunAction(a string) {
	Debug("Running action '%s'", a)
	str := strings.Split(a, ".")
	if str[0] == "sleep" {
		// do sleep command
		sleep, _ := time.ParseDuration(str[1])
		sleepingAction = sleep
		time.Sleep(sleep)
		sleepingAction = 0
	} else {
		s, _ := getSwitch(str[0])

		switch action := strings.ToLower(str[1]); {
		case action == "on":
			s.TurnOn()
		case action == "off":
			s.TurnOff()
		case action == "disable":
			d, _ := time.ParseDuration(str[2])
			s.SetDisableCustom(d)
			log.Printf("Switch '%s' has been disabled, this will last %s", s.Id, d)
		default:
			log.Printf("Unknown action '%s'", a)
		}
	}
}
