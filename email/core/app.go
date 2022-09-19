package core

import (
	"fmt"

	"github.com/kasbuunk/microservice/email"
	"github.com/kasbuunk/microservice/port"
)

// App implements the API interface.
type App struct {
	EventPublisher port.EventPublisher
	EmailClient    port.EmailClient
}

func (s App) Send() error {
	msg := port.Event{
		Stream:  "EMAIL",
		Subject: "ACTIVATION_REQUEST_SENT",
		Body:    port.Body("new user registered with email"),
	}
	err := s.EventPublisher.Publish(msg)
	if err != nil {
		return fmt.Errorf("publishing msg: %w", err)
	}
	return nil
}

func New(eventPublisher port.EventPublisher, emailClient port.EmailClient) email.App {
	return App{
		EventPublisher: eventPublisher,
		EmailClient:    emailClient}
}
