package main

import (
	"log"
	"time"

	"github.com/gregdel/pushover"
)

var lastNotificationTime time.Time

func SendNotification(alertMessage string) {
	if lastNotificationTime.Add(c.Alerts.AntiSpam.Sleep).Before(time.Now()) {
		log.Println("SendNotification(): Sent alert: '" + alertMessage + "'")
		lastNotificationTime = time.Now()
		// PushoverNotification(alertMessage)
	} else {
		log.Println("SendNotification(): Alert not sent: '" + alertMessage + "' (anti-spam)")
	}
}

func PushoverNotification(alertMessage string) {
	app := pushover.New(c.Alerts.Pushover.APIToken)
	recipient := pushover.NewRecipient(c.Alerts.Pushover.UserToken)

	message := &pushover.Message{
		Message:    alertMessage,
		DeviceName: c.Alerts.Pushover.Device,
	}

	response, err := app.SendMessage(message, recipient)
	if err != nil {
		log.Panic(err)
	}

	if c.Debug {
		log.Println(response)
	}
}
