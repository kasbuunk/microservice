package emailsubscriber

import (
	"fmt"
	"log"

	"github.com/kasbuunk/microservice/api/client"
	"github.com/kasbuunk/microservice/api/email"
	"github.com/kasbuunk/microservice/event"
)

type Subscriber struct {
	API       email.API
	BusClient client.EventBusClient
}

// SubscribeToEvents listens for events that match the Stream or Subject.
func (s Subscriber) SubscribeToEvents() {
	fmt.Println("API service listening for messages.")
	for {
		// starts process in loop, in goroutine that awaits published messages and invokes api calls
		eventBus, err := s.BusClient.Subscribe("AUTH", "USER_REGISTERED")
		if err != nil {
			log.Fatal(fmt.Errorf("subscribing: %w", err))
		}
		fmt.Println(<-eventBus)
		// TODO Switch statement based on message stream/subject.
		// Invoke behaviour that marks a user as PENDING ACTIVATION
	}
}

func New(api email.API, bus client.EventBusClient) event.Subscriber {
	return Subscriber{
		API:       api,
		BusClient: bus,
	}
}
