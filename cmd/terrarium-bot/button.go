package main

import (
	"log"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func MonitorButtons() {
	for i, b := range c.GPIO {
		if b.Type == "button" {
			if c.Debug {
				log.Printf("MonitorButtons(): Monitoring button '%s'", b.Action)
			}
			go MonitorButton(i, b)
		}
	}
}

func MonitorButton(buttonIndex int, button GPIO) {
	pin := rpio.Pin(button.Pin)

	if err := rpio.Open(); err != nil {
		log.Println(err)
		return
	}

	pin.Input()
	pin.PullUp()
	pin.Detect(rpio.FallEdge)

	for {
		if pin.EdgeDetected() {
			log.Printf("Button Press: %s", c.GPIO[buttonIndex].Action)
			c.GPIO[buttonIndex].LastStateChange = time.Now()
		}
		time.Sleep(750 * time.Millisecond)
	}

}
