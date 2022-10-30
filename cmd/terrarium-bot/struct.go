package main

type TerrariumPiSensorResp struct {
	State struct {
		Sensors struct {
			Current float32 `json:"current"`
		} `json:"sensors"`
	} `json:"state"`
}
