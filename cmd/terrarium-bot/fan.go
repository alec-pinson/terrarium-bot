package main

import (
	"log"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func FanInit() {
	for _, f := range c.GPIO {
		if f.Type == "fan" {
			FanOff(true)
		}
	}
}

func FanOn() {
	for i, f := range c.GPIO {
		if f.Type == "fan" {
			if f.LastStateChange.Add(f.Sleep).Before(time.Now()) && lastMistTime.Add(f.SleepPostMist).Before(time.Now()) && f.State != "on" {
				log.Printf("Fan On: %s", f.Name)
				c.GPIO[i].LastStateChange = time.Now()
				c.GPIO[i].State = "on"
				SetFan(f.Pin, f.Speed)
				time.Sleep(f.Length)
				FanOff()
			}
		}
	}
}

func FanOff(NoLog ...bool) {
	for i, f := range c.GPIO {
		if f.Type == "fan" {
			if f.LastStateChange.Add(f.Sleep).Before(time.Now()) && f.State != "off" {
				if len(NoLog) == 0 {
					log.Printf("Fan Off: %s", f.Name)
				}
				c.GPIO[i].LastStateChange = time.Now()
				c.GPIO[i].State = "off"
				SetFan(f.Pin, 0)
			}
		}
	}
}

func SetFan(pinNumber int, speed int) {
	if c.Debug {
		return
	}
	err := rpio.Open()
	if err != nil {
		os.Exit(1)
	}
	defer rpio.Close()

	pin := rpio.Pin(pinNumber)
	pin.Mode(rpio.Output)
	if speed == 0 {
		pin.Write(rpio.Low)
	} else {
		// pin.Mode(rpio.Pwm)
		// pin.Freq(25)
		// pin.DutyCycle(uint32(speed), uint32(speed))
		pin.Write(rpio.High)
	}
}
