package email

import (
	"fmt"

	"github.com/kasbuunk/microservice/api/client"
)

// API provides the interface that maps closely to however you wish to communicate with external components.
// It may be a one-to-one mapping to a graphql schema or grpc service.
// Other contexts, or 'domains', should communicate with each other through their APIs.
type API interface {
	Send() error
}

// Service implements the API interface.
type Service struct {
	BusClient   client.EventBusClient
	EmailClient client.EmailClient
}

func (s Service) Send() error {
	msg := client.Event{
		Stream:  "EMAIL",
		Subject: "ACTIVATION_REQUEST_SENT",
		Body:    client.Body("new user registered with email"),
	}
	err := s.BusClient.Publish(msg)
	if err != nil {
		return fmt.Errorf("publishing msg: %w", err)
	}
	return nil
}

func New(busCLient client.EventBusClient, emailClient client.EmailClient) API {
	return Service{
		BusClient:   busCLient,
		EmailClient: emailClient}
}
