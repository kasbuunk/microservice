package core

import (
	"fmt"

	"github.com/kasbuunk/microservice/app/email"
	"github.com/kasbuunk/microservice/app/port"
)

// App implements the API interface.
type App struct {
	EventBus    port.EventBus
	EmailClient port.Client
}

func (s App) Send() error {
	msg := port.Event{
		Stream:  "EMAIL",
		Subject: "ACTIVATION_REQUEST_SENT",
		Body:    port.Body("new user registered with email"),
	}
	err := s.EventBus.Publish(msg)
	if err != nil {
		return fmt.Errorf("publishing msg: %w", err)
	}
	return nil
}

func New(busClient port.EventBus, emailClient port.Client) email.App {
	return App{
		EventBus:    busClient,
		EmailClient: emailClient}
}
