package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

var cameraNightMode bool = true

func MonitorCamera() {
	CameraInit()
	for {
		if !mistMode {
			if DayTime() {
				SetCameraDayMode()
			} else {
				SetCameraNightMode()
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func CameraInit() {
	SetCameraDayMode()
}

func SetCameraDayMode() {
	if cameraNightMode {
		SendCameraCommand("toggle-rtsp-nightvision-off")
		SendCameraCommand("ir_led_off")
		SendCameraCommand("ir_cut_on")
		cameraNightMode = false
	}
}

func SetCameraNightMode() {
	if !cameraNightMode {
		SendCameraCommand("toggle-rtsp-nightvision-on")
		SendCameraCommand("ir_led_on")
		SendCameraCommand("ir_cut_off")
		cameraNightMode = true
	}
}

func SendCameraCommand(Command string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	_, err := client.Get("https://" + c.Camera.Username + ":" + c.Camera.Password + "@" + c.Camera.Hostname + "/cgi-bin/action.cgi?cmd=" + Command)
	if err != nil {
		log.Println(err.Error())
	}
}
