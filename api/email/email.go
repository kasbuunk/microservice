package email

import (
	"fmt"

	"github.com/kasbuunk/microservice/event"
)

// API provides the interface that maps closely to however you wish to communicate with external components.
// It may be a one-to-one mapping to a graphql schema or grpc service.
// Other contexts, or 'domains', should communicate with each other through their APIs.
type API interface {
	Send() error
}

// Service implements the API interface.
type Service struct {
	Bus         event.Bus
	EmailClient Client
}

func (s Service) Send() error {
	msg := event.Message{
		Stream:  "EMAIL",
		Subject: "ACTIVATION_REQUEST_SENT",
		Body:    event.Body("new user registered with email"),
	}
	err := s.Bus.Publish(msg)
	if err != nil {
		return fmt.Errorf("publishing msg: %w", err)
	}
	return nil
}

func New(bus event.Bus, emailClient Client) API {
	return Service{
		Bus:         bus,
		EmailClient: emailClient}
}

type Address string

type Client interface {
	SendActivationLink(Address) error
}
