package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gregdel/pushover"
)

func GetNotification(id string) *Notification {
	for _, n := range config.Notification {
		if n.Id == id {
			return n
		}
	}
	log.Fatalf("Notification '%s' not found in configuration.yaml", id)
	return &Notification{}
}

func InitNotifications() {
	// send a startup notification (useful if it keeps crashing without you knowing)
	for _, n := range config.Notification {
		log.Println("Started: Terrarium bot")
		n.SendNotification("Terrarium bot started")
		// we don't want to block alerts because of this
		n.LastNotification = time.Now().Add(-24 * time.Hour)
	}
}

func (n *Notification) setLastNotification() {
	n.LastNotification = time.Now()
}

func (n *Notification) SendNotification(s string, v ...any) {
	var alertMessage string = fmt.Sprintf(s, v...)

	if n.LastNotification.Add(n.AntiSpam).Before(time.Now()) {
		// make sure we're not spamming
		if n.do(alertMessage) {
			n.setLastNotification()
			log.Println("Alert: " + alertMessage)
		}
	} else {
		Debug("Alert not sent: %s (anti-spam @ %s)", alertMessage, n.AntiSpam)
	}
}

func (n *Notification) do(alertMessage string) bool {
	switch id := n.Id; {
	case id == "pushover":
		PushoverNotification(*n, alertMessage)
	default:
		log.Printf("Unknown notification type '%s", id)
		return false
	}
	return true
}

func PushoverNotification(n Notification, alertMessage string) {
	if config.DryRun {
		return
	}
	app := pushover.New(n.APITokenValue)
	recipient := pushover.NewRecipient(n.UserTokenValue)

	message := &pushover.Message{
		Message:    alertMessage,
		DeviceName: n.Device,
		Sound:      n.Sound,
	}

	response, err := app.SendMessage(message, recipient)
	if err != nil {
		log.Println(err)
	}

	Debug("%v", response)
}
