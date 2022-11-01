package main

import (
	"log"
	"time"
)

func MonitorLights() {
	for {
		if Sunrise() {
			SunriseLights()
		} else if Sunset() {
			SunsetLights()
		} else if DayTime() {
			DayTimeLights()
		} else {
			NightTimeLights()
		}
		time.Sleep(1 * time.Minute)
	}
}

func SunriseLights() {
	for _, l := range c.Switches {
		if l.Type == "light" && l.Sunrise == "on" {
			log.Printf("SunriseLights(): Turning on light '%s'", l.Name)
		}
	}
}

func SunsetLights() {
	for _, l := range c.Switches {
		if l.Type == "light" && l.Sunset == "off" {
			log.Printf("SunsetLights(): Turning off light '%s'", l.Name)
		}
	}
}

func DayTimeLights() {
	for _, l := range c.Switches {
		if l.Type == "light" {
			LightOn(l)
		}
	}
}

func NightTimeLights() {
	for _, l := range c.Switches {
		if l.Type == "light" {
			LightOff(l)
		}
	}
}

func LightOn(l Switch) {
	if GetSwitchState(l) == "off" {
		log.Printf("LightOn(): Turning on light '%s'", l.Name)
		SetSwitchState(l, "on")
	}
}

func LightOff(l Switch) {
	if GetSwitchState(l) == "on" {
		log.Printf("LightOff(): Turning off light '%s'", l.Name)
		SetSwitchState(l, "off")
	}
}
