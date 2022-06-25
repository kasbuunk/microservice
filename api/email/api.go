package email

import (
	"fmt"

	"github.com/kasbuunk/microservice/api/client/email"
	"github.com/kasbuunk/microservice/api/client/eventbus"
)

// API provides the interface that maps closely to however you wish to communicate with external components.
// It may be a one-to-one mapping to a graphql schema or grpc service.
// Other contexts, or 'domains', should communicate with each other through their APIs.
type API interface {
	Send() error
}

// Service implements the API interface.
type Service struct {
	EventBus    eventbus.Client
	EmailClient email.Client
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

func New(busCLient eventbus.Client, emailClient email.Client) API {
	return Service{
		EventBus:    busCLient,
		EmailClient: emailClient}
}
