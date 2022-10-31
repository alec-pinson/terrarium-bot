package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetSwitchState(Switch Switch) string {
	client := http.Client{}
	req, err := http.NewRequest("GET", c.HomeAssistant.URL+"/api/states/"+Switch.ID, nil)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}

	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + c.HomeAssistant.Token},
	}

	response, err := client.Do(req)
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var resp HomeAssistantStateResp
	json.Unmarshal(responseData, &resp)
	return resp.State
}
