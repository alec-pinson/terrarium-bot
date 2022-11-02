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
				if DayTime() {
					log.Printf("Fan On: %s (%v%%/%v%%)", f.Name, GetHumidity(), c.Humidity.Day.Maximum)
				} else {
					log.Printf("Fan On: %s (%v%%/%v%%)", f.Name, GetHumidity(), c.Humidity.Night.Maximum)
				}
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
					if DayTime() {
						log.Printf("Fan Off: %s (%v%%/%v%%)", f.Name, GetHumidity(), c.Humidity.Day.Maximum)
					} else {
						log.Printf("Fan Off: %s (%v%%/%v%%)", f.Name, GetHumidity(), c.Humidity.Night.Maximum)
					}
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
		log.Println(err)
		os.Exit(1)
	}

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
