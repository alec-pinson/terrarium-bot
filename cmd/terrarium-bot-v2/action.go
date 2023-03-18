package main

import (
	"log"
	"strings"
	"time"
)

/*
example actions list
	- sleep.5s
	- switch.fan.on
	- switch.fan.on.10m
	- switch.fan.off
	- switch.fan.disable
	- swith.fan.enable
	- trigger.mist.disable
	- trigger.mist.disable.40m
	- trigger.mist.enable
*/

func RunAction(a string, reason string) {
	Debug("Running action '%s'", a)
	args := strings.Split(a, ".")

	if len(args) <= 1 {
		log.Printf("Invalid action '%s'", a)
		return
	}

	switch args[0] {
	case "sleep":
		runSleepAction(args, reason)
	case "trigger":
		runTriggerAction(args, reason)
	case "switch":
		runSwitchAction(args, reason)
	default:
		log.Printf("Unknown action '%s'", strings.Join(args, "."))
	}
}

func runSleepAction(args []string, reason string) {
	sleepDuration, err := time.ParseDuration(args[1])
	if err != nil {
		log.Printf("Invalid sleep duration '%s'", args[1])
		return
	}
	time.Sleep(sleepDuration)
}

func runTriggerAction(args []string, reason string) {
	t := GetTrigger(args[1])
	switch strings.ToLower(args[2]) {
	case "disable":
		if len(args) == 4 {
			t.Disable(args[3], reason)
			return
		}
		if len(args) == 3 {
			// disable 'forever'
			t.Disable("", reason)
			return
		}
		log.Printf("Invalid parameters for disable action '%s'", strings.Join(args, "."))
	case "enable":
		t.Enable(reason)
	default:
		log.Printf("Unknown switch action '%s'", strings.Join(args, "."))
	}
}

func runSwitchAction(args []string, reason string) {
	s := GetSwitch(args[1])
	switch strings.ToLower(args[2]) {
	case "on":
		if len(args) == 4 {
			// set on for duration
			s.TurnOn(args[3], reason)
		} else {
			// turn on forever
			s.TurnOn("", reason)
		}
	case "off":
		s.TurnOff(reason)
	case "disable":
		if len(args) == 4 {
			s.Disable(args[3], reason)
			return
		}
		if len(args) == 3 {
			// disable 'forever'
			s.Disable("", reason)
			return
		}
		log.Printf("Invalid parameters for disable action '%s'", strings.Join(args, "."))
	case "enable":
		s.Enable(reason)
	default:
		log.Printf("Unknown switch action '%s'", strings.Join(args, "."))
	}
}
