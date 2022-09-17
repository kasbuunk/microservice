package core

import (
	"fmt"

	"github.com/kasbuunk/microservice/app/email"
	emailclient "github.com/kasbuunk/microservice/app/email/port"
	"github.com/kasbuunk/microservice/app/eventbus"
)

// App implements the API interface.
type App struct {
	EventBus    eventbus.EventBus
	EmailClient emailclient.Client
}

func (s App) Send() error {
	msg := eventbus.Event{
		Stream:  "EMAIL",
		Subject: "ACTIVATION_REQUEST_SENT",
		Body:    eventbus.Body("new user registered with email"),
	}
	err := s.EventBus.Publish(msg)
	if err != nil {
		return fmt.Errorf("publishing msg: %w", err)
	}
	return nil
}

func New(busClient eventbus.EventBus, emailClient emailclient.Client) email.App {
	return App{
		EventBus:    busClient,
		EmailClient: emailClient}
}
