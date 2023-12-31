package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func GetSwitch(id string, exitOnError bool) *Switch {
	for _, s := range config.Switch {
		if s.Id == id {
			return s
		}
	}
	if exitOnError {
		log.Fatalf("Switch '%s' not found in configuration.yaml", id)
	} else {
		log.Printf("Switch '%s' not found in configuration.yaml", id)
	}
	return nil
}

func InitSwitches() {
	for _, s := range config.Switch {
		// some urls include env vars
		s.fixURLs()
		// update status if possible on startup
		s.getStatus()
	}
}

func (s *Switch) SetLastAction() {
	s.LastAction = time.Now()
	s.Disabled = 0
}

func (s *Switch) getStatus() string {
	if s.StatusUrl == "" {
		return s.State
	}

	r, respCode, err := SendRequest(s.StatusUrl, s.Insecure, 3, s.JsonPath != "", s.APITokenValue)
	if err != nil || respCode != 200 {
		log.Printf("Switch Offline: %s", s.Id)
		for _, n := range config.Notification {
			n.SendNotification("Currently unable to get status for switch '%s'. Please check the logs.", s.Id)
		}
	} else {
		// use resp to get status
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			log.Println(err)
			return s.State
		}
		state := strings.ToLower(fmt.Sprintf("%s", getJsonValue(string(b), s.JsonPath)))
		if state == "on" || state == "off" {
			s.State = state
		} else {
			log.Printf("Unknown status '%s' received for switch '%s'", state, s.Id)
			return s.State
		}
	}

	return s.State
}

func (s *Switch) setStatus(state string) {
	if s.StatusUrl == "" {
		s.State = state
	}
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
	if s.DisabledAt.Add(s.Disabled).Before(time.Now()) {
		return false
	}
	return true
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
	if s.getStatus() == "on" {
		return
	}
	// check for disable parameter
	if s.isDisabled() {
		Debug("Cannot turn on '%s' as it is currently disabled (%s)", s.Id, Reason)
		return
	}

	s.SetLastAction()
	if !config.DryRun {
		_, respCode, err := SendRequest(s.On, s.Insecure, 3, false, s.APITokenValue)
		if err != nil || respCode != 200 {
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
	if s.getStatus() == "off" {
		return
	}
	s.SetLastAction()
	if !config.DryRun {
		_, respCode, err := SendRequest(s.Off, s.Insecure, 3, false, s.APITokenValue)
		if err != nil || respCode != 200 {
			log.Printf("Switch Offline: %s", s.Id)
			for _, n := range config.Notification {
				n.SendNotification("Currently unable to turn off switch '%s'. Please check the logs.", s.Id)
			}
		}
	}
	s.setStatus("off")
	log.Printf("Switch Off: '%s' (%s)", s.Id, reason)
}

func (s *Switch) WriteStatus(w http.ResponseWriter) {
	writeResponse(w, s, false)
}
