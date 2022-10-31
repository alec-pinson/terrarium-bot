package main

import (
	"log"
	"time"
)

var lastNotificationTime time.Time

func SendNotification(alertMessage string) {
	if lastNotificationTime.Add(c.Alerts.AntiSpam.Sleep).Before(time.Now()) {
		log.Println("SendNotification(): Sent alert: '" + alertMessage + "'")
		lastNotificationTime = time.Now()
	} else {
		log.Println("SendNotification(): Alert not sent: '" + alertMessage + "' (anti-spam)")
	}
}
