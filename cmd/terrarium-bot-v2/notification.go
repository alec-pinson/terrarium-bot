package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gregdel/pushover"
)

func getNotification(id string) (*Notification, error) {
	for _, n := range config.Notification {
		if n.Id == id {
			return &n, nil
		}
	}
	return &Notification{}, errors.New("Notification '" + id + "' not found.")
}

func (n Notification) getIdx() (int, error) {
	for idx, nn := range config.Notification {
		if n.Id == nn.Id {
			return idx, nil
		}
	}
	return 0, errors.New("Notification '" + n.Id + "' not found.")
}

func (n Notification) setLastNotification() {
	idx, err := n.getIdx()
	if err != nil {
		log.Println(err)
		return
	}
	config.Notification[idx].LastNotification = time.Now()
}

func (n Notification) SendNotification(s string, v ...any) {
	var alertMessage string = fmt.Sprintf(s, v...)

	if n.LastNotification.Add(n.AntiSpam).Before(time.Now()) {
		// make sure we're not spamming
		if n.do(alertMessage) {
			n.setLastNotification()
			log.Println("Sent alert: " + alertMessage)
		}
	} else {
		Debug("Alert not sent: %s (anti-spam @ %s)", alertMessage, n.AntiSpam)
	}
}

func (n Notification) do(alertMessage string) bool {
	switch id := n.Id; {
	case id == "pushover":
		PushoverNotification(n, alertMessage)
	default:
		log.Printf("Unknown notification type '%s", id)
		return false
	}
	return true
}

func PushoverNotification(n Notification, alertMessage string) {
	if config.Debug {
		return
	}
	app := pushover.New(n.APITokenValue)
	recipient := pushover.NewRecipient(n.UserTokenValue)

	message := &pushover.Message{
		Message:    alertMessage,
		DeviceName: n.Device,
	}

	response, err := app.SendMessage(message, recipient)
	if err != nil {
		log.Println(err)
	}

	Debug("%v", response)
}
