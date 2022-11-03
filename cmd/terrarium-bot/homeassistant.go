package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetSwitchState(Switch Switch) string {
	client := http.Client{}
	req, err := http.NewRequest("GET", c.HomeAssistant.URL+"/api/states/"+Switch.ID, nil)
	if err != nil {
		log.Print(err.Error())
		return "unknown"
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + c.HomeAssistant.Token},
	}

	response, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
		return "unknown"
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
		return "unknown"
	}
	var resp HomeAssistantStateResp
	json.Unmarshal(responseData, &resp)
	return resp.State
}

func SetSwitchState(Switch Switch, State string) {
	if c.Debug {
		return
	}
	var command string
	switch strings.ToLower(State) {
	case "on":
		command = "turn_on"
	case "off":
		command = "turn_off"
	}

	var body = []byte(`{"entity_id":"` + Switch.ID + `"}`)

	client := http.Client{}
	req, err := http.NewRequest("POST", c.HomeAssistant.URL+"/api/services/switch/"+command, bytes.NewBuffer(body))
	if err != nil {
		log.Print(err.Error())
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + c.HomeAssistant.Token},
	}

	response, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
	}
	var resp HomeAssistantStateResp
	json.Unmarshal(responseData, &resp)
}
