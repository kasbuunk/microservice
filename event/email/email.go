package emailhandler

import (
	"fmt"
	"log"

	"github.com/kasbuunk/microservice/app/email"
	"github.com/kasbuunk/microservice/app/port"
)

type EventHandler struct {
	App       email.App
	BusClient port.EventBus
}

// Handle listens for events that match the Stream or Subject and invokes the appropriate domain behaviour.
func (s EventHandler) Handle() {
	fmt.Println("Email service listening for events.")
	for {
		// Starts process in loop, awaiting published messages.
		eventBus, err := s.BusClient.Subscribe("AUTH", "USER_REGISTERED")
		if err != nil {
			log.Fatal(fmt.Errorf("subscribing: %w", err))
		}
		fmt.Println(<-eventBus)
		// TODO Switch statement based on message stream/subject.
		// Invoke behaviour that marks a user as PENDING ACTIVATION
	}
}

func New(api email.App, bus port.EventBus) EventHandler {
	return EventHandler{
		App:       api,
		BusClient: bus,
	}
}
