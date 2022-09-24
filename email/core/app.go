package core

import (
	"fmt"

	"github.com/kasbuunk/microservice/email"
	"github.com/kasbuunk/microservice/email/port"
	"github.com/kasbuunk/microservice/eventbus"
)

// App implements the API interface.
type App struct {
	EventPublisher eventbus.EventPublisher
	EmailClient    port.EmailClient
}

func (s App) Send() error {
	msg := eventbus.Event{
		Stream:  "EMAIL",
		Subject: "ACTIVATION_REQUEST_SENT",
		Body:    eventbus.Body("new user registered with email"),
	}
	err := s.EventPublisher.Publish(msg)
	if err != nil {
		return fmt.Errorf("publishing msg: %w", err)
	}
	return nil
}

func New(eventPublisher eventbus.EventPublisher, emailClient port.EmailClient) email.App {
	return App{
		EventPublisher: eventPublisher,
		EmailClient:    emailClient}
}
