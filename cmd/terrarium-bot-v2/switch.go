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
	return nil
}

func InitSwitches() {
	for _, s := range config.Switch {
		// some urls include env vars
		s.fixURLs()
	}
}

func (s *Switch) SetLastAction() {
	s.LastAction = time.Now()
	s.Disabled = 0
}

func (s *Switch) getStatus() string {
	return s.Status
}

func (s *Switch) setStatus(status string) {
	s.Status = status
}

func (s *Switch) Enable(reason string) {
	s.Disabled = 0
	log.Printf("Switch Enabled: '%s'", s.Id)
}

func (s *Switch) Disable(duration string, reason string) {
	if duration == "" {
		// 10 years.. 'forever'
		duration = "87660h"
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		log.Printf("Invalid switch disable duration '%s'", duration)
		return
	}
	s.DisabledAt = time.Now()
	s.Disabled = d
	if duration == "87660h" {
		log.Printf("Switch Disabled: '%s'", s.Id)
	} else {
		log.Printf("Switch Disabled: '%s' for %s", s.Id, d)
	}
}

func (s *Switch) isDisabled() bool {
	if s.Disabled == 0 {
		return false
	}
	if s.DisabledAt.Add(s.Disabled).After(time.Now()) {
		return true
	}
	return false
}
func (s *Switch) fixURLs() {
	if strings.Contains(s.On, "$") {
		s.On = (os.ExpandEnv(s.On))
	}
	if strings.Contains(s.Off, "$") {
		s.Off = os.ExpandEnv(s.Off)
	}
}

func (s *Switch) TurnOn(For string, Reason string) {
	if config.UseInMemoryStatus && s.getStatus() == "on" {
		return
	}
	// check for disable parameter
	if s.isDisabled() {
		Debug("Cannot turn on '%s' as it is currently disabled (%s)", s.Id, Reason)
		return
	}

	s.SetLastAction()
	if !config.DryRun {
		_, err := SendRequest(s.On)
		if err != nil {
			log.Printf("Switch Offline: %s", s.Id)
			for _, n := range config.Notification {
				n.SendNotification("Currently unable to turn on switch '%s'. Please check the logs.", s.Id)
			}
		}
	}
	s.setStatus("on")
	if For != "" {
		onFor, _ := time.ParseDuration(For)
		log.Printf("Switch On: '%s' for %v (%s)", s.Id, For, Reason)
		time.Sleep(onFor)
		s.TurnOff(For + " has elapsed")
	} else {
		log.Printf("Switch On: '%s' (%s)", s.Id, Reason)
	}
}

func (s *Switch) TurnOff(reason string) {
	if config.UseInMemoryStatus && s.getStatus() == "off" {
		return
	}
	s.SetLastAction()
	if !config.DryRun {
		_, err := SendRequest(s.Off)
		if err != nil {
			log.Printf("Switch Offline: %s", s.Id)
			for _, n := range config.Notification {
				n.SendNotification("Currently unable to turn off switch '%s'. Please check the logs.", s.Id)
			}
		}
	}
	s.setStatus("off")
	log.Printf("Switch Off: '%s' (%s)", s.Id, reason)
}
