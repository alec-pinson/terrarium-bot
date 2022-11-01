package main

import (
	"log"
	"time"
)

func FanInit() {
	for _, f := range c.GPIO {
		if f.Type == "fan" {
			FanOff()
		}
	}
}

func FanOn() {
	for i, f := range c.GPIO {
		if f.Type == "fan" {
			if f.LastStateChange.Add(f.Sleep).Before(time.Now()) && lastMistTime.Add(f.SleepPostMist).Before(time.Now()) && f.State != "on" {
				log.Printf("FanOn(): Turning on fan '%s'", f.Name)
				c.GPIO[i].LastStateChange = time.Now()
				c.GPIO[i].State = "on"
				SetFan(f.Pin, f.Speed)
				time.Sleep(f.Length)
				FanOff()
			}
		}
	}
}

func FanOff() {
	for i, f := range c.GPIO {
		if f.Type == "fan" {
			if f.LastStateChange.Add(f.Sleep).Before(time.Now()) && f.State != "off" {
				log.Printf("FanOff(): Turning off fan '%s'", f.Name)
				c.GPIO[i].LastStateChange = time.Now()
				c.GPIO[i].State = "off"
				SetFan(f.Pin, 0)
			}
		}
	}
}

func SetFan(pinNumber int, speed int) {
	// uSpeed := uint32(speed)
	// err := rpio.Open()
	// if err != nil {
	// 	os.Exit(1)
	// }
	// defer rpio.Close()

	// pin := rpio.Pin(pinNumber)
	// pin.Mode(rpio.Pwm)
	// pin.Freq(25)
	// pin.DutyCycle(uSpeed, 32)
}
