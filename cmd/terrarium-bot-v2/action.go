package main

import (
	"log"
	"strings"
	"time"
)

var sleepingAction time.Duration // used by time.go to pause normal actions until after this duration

func RunAction(a string, reason string) {
	Debug("Running action '%s'", a)
	str := strings.Split(a, ".")

	if len(str) <= 1 {
		log.Printf("Invalid action '%s'", a)
		return
	}

	arg1 := str[0]
	arg2 := str[1]
	switch arg1 {
	case "sleep":
		sleepDuration, err := time.ParseDuration(arg2)
		if err != nil {
			log.Printf("Invalid sleep duration '%s'", arg2)
			return
		}
		sleepingAction = sleepDuration
		time.Sleep(sleepDuration)
		sleepingAction = 0
	default:
		s := GetSwitch(arg1)
		switch strings.ToLower(arg2) {
		case "on":
			s.TurnOn(reason)
		case "off":
			s.TurnOff(reason)
		case "disable":
			if len(str) == 3 {
				s.SetDisableCustom(str[2], reason)
				return
			}
			if len(str) == 2 {
				// disable 'forever'
				s.SetDisableCustom("87660h", reason) // 10 years :)
				return
			}
			log.Printf("Invalid parameters for disable action '%s'", a)
		case "enable":
			s.SetDisableCustom("0", reason)
			s.Disable = 0
		default:
			log.Printf("Unknown action '%s'", a)
		}
	}
}
