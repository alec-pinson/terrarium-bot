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
	- alert.humidity.disable
	- alert.humidity.disable.1h
	- alert.humidity.enable
	- echo.Misting will begin in 5m
*/

func RunAction(a string, reason string) bool {
	Debug("Running action '%s'", a)
	args := strings.Split(a, ".")

	if len(args) <= 1 {
		log.Printf("Invalid action '%s'", a)
		return false
	}

	switch args[0] {
	case "sleep":
		return runSleepAction(args, reason)
	case "trigger":
		return runTriggerAction(args, reason)
	case "switch":
		return runSwitchAction(args, reason)
	case "alert":
		return runAlertAction(args, reason)
	case "echo":
		return runEchoAction(args, reason)
	default:
		log.Printf("Unknown action '%s'", strings.Join(args, "."))
		return false
	}
}

func runEchoAction(args []string, reason string) bool {
	message := strings.Join(args, ".")
	log.Println(message[5:]) // print without echo.
	return true
}

func runSleepAction(args []string, reason string) bool {
	sleepDuration, err := time.ParseDuration(args[1])
	if err != nil {
		log.Printf("Invalid sleep duration '%s'", args[1])
		return false
	}
	time.Sleep(sleepDuration)
	return true
}

func runTriggerAction(args []string, reason string) bool {
	t := GetTrigger(args[1])
	switch strings.ToLower(args[2]) {
	case "disable":
		if len(args) == 4 {
			t.Disable(args[3], reason)
			return true
		}
		if len(args) == 3 {
			// disable 'forever'
			t.Disable("", reason)
			return true
		}
		log.Printf("Invalid parameters for disable action '%s'", strings.Join(args, "."))
		return false
	case "enable":
		t.Enable(reason)
		return true
	default:
		log.Printf("Unknown trigger action '%s'", strings.Join(args, "."))
		return false
	}
}

func runSwitchAction(args []string, reason string) bool {
	s := GetSwitch(args[1], true)
	switch strings.ToLower(args[2]) {
	case "on":
		if len(args) == 4 {
			// set on for duration
			s.TurnOn(args[3], reason)
			return true
		} else {
			// turn on forever
			s.TurnOn("", reason)
			return true
		}
	case "off":
		s.TurnOff(reason)
		return true
	case "disable":
		if len(args) == 4 {
			s.Disable(args[3], reason)
			return true
		}
		if len(args) == 3 {
			// disable 'forever'
			s.Disable("", reason)
			return true
		}
		log.Printf("Invalid parameters for disable action '%s'", strings.Join(args, "."))
		return false
	case "enable":
		s.Enable(reason)
		return true
	default:
		log.Printf("Unknown switch action '%s'", strings.Join(args, "."))
		return false
	}
}

func runAlertAction(args []string, reason string) bool {
	a := GetAlert(args[1])
	switch strings.ToLower(args[2]) {
	case "disable":
		if len(args) == 4 {
			a.Disable(args[3], reason)
			return true
		}
		if len(args) == 3 {
			// disable 'forever'
			a.Disable("", reason)
			return true
		}
		log.Printf("Invalid parameters for disable action '%s'", strings.Join(args, "."))
		return false
	case "enable":
		a.Enable(reason)
		return true
	default:
		log.Printf("Unknown alert action '%s'", strings.Join(args, "."))
		return false
	}
}
