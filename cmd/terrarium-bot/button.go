package main

import (
	"log"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

func MonitorButtons() {
	for i, b := range c.GPIO {
		if b.Type == "button" {
			log.Printf("MonitorButtons(): Monitoring button '%s'", b.Action)
			go MonitorButton(i, b)
		}
	}
}

func MonitorButton(buttonIndex int, button GPIO) {
	pin := rpio.Pin(button.Pin)

	if err := rpio.Open(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rpio.Close()

	pin.Input()
	pin.PullUp()
	pin.Detect(rpio.FallEdge)

	for {
		if pin.EdgeDetected() {
			c.GPIO[buttonIndex].LastStateChange = time.Now()
		}
		time.Sleep(time.Second / 2)
	}

}
