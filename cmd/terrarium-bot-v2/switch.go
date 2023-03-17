package main

import (
	"log"
	"os"
	"strings"
	"time"
)

func GetSwitch(id string) *Switch {
	for _, s := range config.Switch {
		if s.Id == id {
			return s
		}
	}
	log.Fatalf("Switch '%s' not found in configuration.yaml", id)
	return &Switch{}
}

func InitSwitches() {
	for _, s := range config.Switch {
		// some urls include env vars
		s.fixURLs()
		if s.Every != 0 {
			// some switches are configured to run every x minutes/hours
			go s.monitor()
		}
	}
}

func (s *Switch) monitor() {
	log.Printf("Switch '%s' has been set to turn on every %s for %s (not during the night)", s.Id, s.Every, s.For)
	if s.Disable != 0 {
		log.Printf("Switch '%s' will not run more than once every %s", s.Id, s.Disable)
	}
	if !isTesting {
		// if testing, skip this
		time.Sleep(s.Every)
	}
	for {
		if isDayTime() { // maybe add something here so it doesn't mist straight after sunrise
			lastAction := s.GetLastAction()
			// has this action been ran in the past x minutes/hours
			if lastAction.Before(time.Now().Add(-s.Every)) {
				// nope
				s.TurnOn("Scheduled every " + s.Every.String())
				time.Sleep(s.Every)
			}
		}
		if isTesting {
			return
		}
		time.Sleep(30 * time.Second)
	}
}

func (s *Switch) SetLastAction() {
	s.LastAction = time.Now()
	s.DisableCustom = 0
}

func (s *Switch) GetLastAction() time.Time {
	return s.LastAction
}

func (s *Switch) getStatus() string {
	return s.Status
}

func (s *Switch) setStatus(status string) {
	s.Status = status
}

func (s *Switch) SetDisableCustom(duration string, reason string) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Printf("Invalid disable duration '%s'", duration)
		return
	}
	s.SetLastAction()
	s.DisableCustom = d
	if duration == "87660h" {
		// 10 years.. 'forever'
		log.Printf("Switch '%s' has been disabled", s.Id)
	} else {
		log.Printf("Switch '%s' has been disabled, this will last %s", s.Id, d)
	}
}

func (s *Switch) isDisabled() bool {
	if s.Disable == 0 && s.DisableCustom == 0 {
		return false
	}
	if s.DisableCustom == 0 {
		// disable custom not set, do normal
		if s.LastAction.Add(s.Disable).After(time.Now()) {
			return true
		}
	} else {
		// disable custom is set, use this instead
		if s.LastAction.Add(s.DisableCustom).After(time.Now()) {
			return true
		}
	}
	return false
}

func (s *Switch) setOnUrl(url string) {
	s.On = url
}

func (s *Switch) setOffUrl(url string) {
	s.Off = url
}

func (s *Switch) fixURLs() {
	if strings.Contains(s.On, "$") {
		s.setOnUrl(os.ExpandEnv(s.On))
	}
	if strings.Contains(s.Off, "$") {
		s.setOffUrl(os.ExpandEnv(s.Off))
	}
}

func (s *Switch) TurnOn(reason string) {
	if config.UseInMemoryStatus && s.getStatus() == "on" {
		return
	}
	// check for disable parameter
	if s.isDisabled() {
		log.Printf("Cannot turn on '%s' as it is currently disabled (%s)", s.Id, reason)
		return
	}

	s.SetLastAction()
	if !config.DryRun {
		SendRequest(s.On)
	}
	s.setStatus("on")
	if s.For != 0 {
		log.Printf("%s (Turning on '%s' for %v)", reason, s.Id, s.For)
		time.Sleep(s.For)
		s.TurnOff(s.For.String() + " has elapsed")
	} else {
		log.Printf("%s (Turning on '%s')", reason, s.Id)
	}
}

func (s *Switch) TurnOff(reason string) {
	if config.UseInMemoryStatus && s.getStatus() == "off" {
		return
	}
	s.SetLastAction()
	if !config.DryRun {
		SendRequest(s.Off)
	}
	s.setStatus("off")
	log.Printf("%s (Turning off '%s')", reason, s.Id)
}
