package main

import "time"

type TerrariumPiSensorResp struct {
	State struct {
		Sensors struct {
			Current float32 `json:"current"`
		} `json:"sensors"`
	} `json:"state"`
}

type HomeAssistantStateResp struct {
	State      string    `json:"state"`
	LastChange time.Time `json:"last_changed"`
}
