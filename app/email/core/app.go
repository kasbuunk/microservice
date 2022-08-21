package core

import (
	"fmt"
	emailclient "github.com/kasbuunk/microservice/app/dependency/email"
	"github.com/kasbuunk/microservice/app/dependency/eventbus"
	"github.com/kasbuunk/microservice/app/email"
)

// Service implements the API interface.
type Service struct {
	EventBus    eventbus.Client
	EmailClient emailclient.Client
}

func (s Service) Send() error {
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

func New(busClient eventbus.Client, emailClient emailclient.Client) email.App {
	return Service{
		EventBus:    busClient,
		EmailClient: emailClient}
}
