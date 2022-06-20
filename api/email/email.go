package email

import (
	"fmt"
	"log"

	"github.com/kasbuunk/microservice/events"
)

// Email provides the api that maps closely to however you wish to communicate with external components.
// It may be a one-to-one mapping to a graphql schema or grpc service.
// Other contexts, or 'domains', should communicate with each other through their APIs.
type Email interface {
	Subscribe()
	Send() error
}

// Service implements the Email interface.
type Service struct {
	Bus events.MessageBus
}

// Subscribe listens for messages that match the Stream or Subject.
func (s Service) Subscribe() {
	fmt.Println("Email service listening for messages")
	// starts process in loop, in goroutine that awaits published messages and invokes api calls
	messageChannel, err := s.Bus.Subscribe("AUTH", "USER_REGISTERED")
	if err != nil {
		log.Fatal(fmt.Errorf("subscribing: %w", err))
	}
	fmt.Println(<-messageChannel)
}

func New(bus events.MessageBus) Email {
	return Service{
		Bus: bus,
	}
}

func (s Service) Send() error {
	msg := events.Message{
		Stream:  "EMAIL",
		Subject: "ACTIVATION_REQUEST_SENT",
		Body:    fmt.Sprintf("new user registered with email"),
	}
	err := s.Bus.Publish(msg)
	if err != nil {
		return fmt.Errorf("publishing msg: %w", err)
	}
	return nil
}
