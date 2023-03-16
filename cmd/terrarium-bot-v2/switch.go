package main

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"
)

func getSwitch(id string) (Switch, error) {
	for _, s := range config.Switch {
		if s.Id == id {
			return s, nil
		}
	}
	return Switch{}, errors.New("Switch '" + id + "' not found.")
}

func (s Switch) getIdx() (int, error) {
	for idx, ss := range config.Switch {
		if s.Id == ss.Id {
			return idx, nil
		}
	}
	return 0, errors.New("Switch '" + s.Id + "' not found.")
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

func (s Switch) monitor() {
	log.Printf("Switch '%s' has been set to turn on every %s for %s (not during the night)", s.Id, s.Every, s.For)
	if s.Disable != 0 {
		log.Printf("Switch '%s' will not run more than once every %s", s.Id, s.Disable)
	}
	time.Sleep(s.Every)
	for {
		if isDayTime() { // maybe add something here so it doesn't mist straight after sunrise
			lastAction, err := s.GetLastAction()
			if err == nil {
				// has this action been ran in the past x minutes/hours
				if lastAction.Before(time.Now().Add(-s.Every)) {
					// nope
					s.TurnOn()
					time.Sleep(s.Every)
				}
			}
		}
		time.Sleep(30 * time.Second)
	}
}

func (s Switch) SetLastAction() {
	idx, err := s.getIdx()
	if err != nil {
		log.Println(err)
		return
	}
	config.Switch[idx].LastAction = time.Now()
	config.Switch[idx].DisableCustom = 0
}

func (s Switch) GetLastAction() (time.Time, error) {
	idx, err := s.getIdx()
	if err != nil {
		log.Println(err)
		return time.Time{}, err
	}
	return config.Switch[idx].LastAction, nil
}

func (s Switch) getStatus() string {
	idx, err := s.getIdx()
	if err != nil {
		log.Println(err)
		return ""
	}
	return config.Switch[idx].Status
}

func (s Switch) setStatus(status string) {
	idx, err := s.getIdx()
	if err != nil {
		log.Println(err)
		return
	}
	config.Switch[idx].Status = status
}

func (s Switch) SetDisableCustom(d time.Duration) {
	idx, err := s.getIdx()
	if err != nil {
		log.Println(err)
		return
	}
	s.SetLastAction()
	config.Switch[idx].DisableCustom = d
}

func (s Switch) isDisabled() bool {
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

func (s Switch) setOnUrl(url string) {
	idx, err := s.getIdx()
	if err != nil {
		log.Println(err)
		return
	}
	config.Switch[idx].On = url
}

func (s Switch) setOffUrl(url string) {
	idx, err := s.getIdx()
	if err != nil {
		log.Println(err)
		return
	}
	config.Switch[idx].On = url
}

func (s Switch) fixURLs() {
	if strings.Contains(s.On, "$") {
		s.setOnUrl(os.ExpandEnv(s.On))
	}
	if strings.Contains(s.Off, "$") {
		s.setOffUrl(os.ExpandEnv(s.Off))
	}
}

func (s Switch) TurnOn() {
	s, err := getSwitch(s.Id)
	if err != nil {
		log.Println(err)
		return
	}
	if config.UseInMemoryStatus && s.getStatus() == "on" {
		return
	}
	// check for disable parameter
	if s.isDisabled() {
		log.Printf("Currently disabled, cannot turn on '%s'", s.Id)
		return
	}

	s.SetLastAction()
	if !config.DryRun {
		SendRequest(s.On)
	}
	s.setStatus("on")
	log.Printf("Turning on '%s' (%s)", s.Id, s.On)
	if s.For != 0 {
		time.Sleep(s.For)
		if !config.DryRun {
			SendRequest(s.Off)
		}
		s.setStatus("off")
		log.Printf("Turning off '%s' (%s)", s.Id, s.Off)
	}
}

func (s Switch) TurnOff() {
	s, err := getSwitch(s.Id)
	if err != nil {
		log.Println(err)
		return
	}
	if config.UseInMemoryStatus && s.getStatus() == "off" {
		return
	}
	if !config.DryRun {
		SendRequest(s.Off)
	}
	s.setStatus("off")
	log.Printf("Turning off '%s' (%s)", s.Id, s.Off)
}
