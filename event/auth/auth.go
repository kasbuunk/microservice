package authhandler

import (
	"fmt"
	"log"

	"github.com/kasbuunk/microservice/app/auth"
	"github.com/kasbuunk/microservice/app/port"
)

type EventHandler struct {
	API       auth.App
	BusClient port.EventBus
}

// Handle listens for events that match the Stream or Subject and invokes the appropriate domain behaviour.
func (s EventHandler) Handle() {
	fmt.Println("Auth service listening for events.")
	for {
		// Starts process in loop, awaiting published messages.
		eventBus, err := s.BusClient.Subscribe("AUTH", "ACTIVATION_REQUEST_SENT")
		if err != nil {
			log.Fatal(fmt.Errorf("subscribing: %w", err))
		}
		fmt.Println(<-eventBus)
		// TODO Switch statement based on message stream/subject.
		// Invoke behaviour that marks a user as PENDING ACTIVATION
	}
}

func New(api auth.App, bus port.EventBus) EventHandler {
	return EventHandler{
		API:       api,
		BusClient: bus,
	}
}
