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
			log.Printf("MonitorButtons(): Monitoring button '%s'", b.Name)
			go MonitorButton(i, b)
		}
	}
}

func MonitorButton(buttonIndex int, button GPIO) {
	pin := rpio.Pin(button.Pin)

	// Open and map memory to access gpio, check for errors
	if err := rpio.Open(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// Unmap gpio memory when done
	defer rpio.Close()

	pin.Input()
	pin.PullUp()
	pin.Detect(rpio.FallEdge) // enable falling edge event detection

	log.Println("press a button")

	for i := 0; i < 2; {
		if pin.EdgeDetected() { // check if event occured
			log.Println("button pressed")
			c.GPIO[buttonIndex].LastStateChange = time.Now()
			i++
		}
		time.Sleep(time.Second / 2)
	}
	pin.Detect(rpio.NoEdge) // disable edge event detection

}
