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
		// if s.Every != 0 {
		// 	// some switches are configured to run every x minutes/hours
		// 	go s.monitor()
		// }
	}
}

// func (s *Switch) monitor() {
// 	log.Printf("Switch '%s' has been set to turn on every %s for %s (not during the night)", s.Id, s.Every, s.For)
// 	if s.Disable != 0 {
// 		log.Printf("Switch '%s' will not run more than once every %s", s.Id, s.Disable)
// 	}
// 	if !isTesting {
// 		// if testing, skip this
// 		time.Sleep(s.Every)
// 	}
// 	for {
// 		if isDayTime() { // maybe add something here so it doesn't mist straight after sunrise
// 			lastAction := s.GetLastAction()
// 			// has this action been ran in the past x minutes/hours
// 			if lastAction.Before(time.Now().Add(-s.Every)) {
// 				// nope
// 				s.TurnOn("Scheduled every " + s.Every.String())
// 				time.Sleep(s.Every)
// 			}
// 		}
// 		if isTesting {
// 			return
// 		}
// 		time.Sleep(30 * time.Second)
// 	}
// }

func (s *Switch) SetLastAction() {
	s.LastAction = time.Now()
	s.Disabled = 0
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
	s.SetLastAction()
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
	if s.LastAction.Add(s.Disabled).After(time.Now()) {
		return true
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
