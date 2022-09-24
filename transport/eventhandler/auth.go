package eventhandler

import (
	"fmt"
	"log"

	"github.com/kasbuunk/microservice/auth"
	"github.com/kasbuunk/microservice/eventbus"
)

type AuthEventHandler struct {
	App             auth.App
	EventSubscriber eventbus.EventSubscriber
}

// Handle listens for events that match the Stream or Subject and invokes the appropriate domain behaviour.
func (s AuthEventHandler) Handle() {
	fmt.Println("Auth service listening for events.")
	for {
		// Starts process in loop, awaiting published messages.
		eventBus, err := s.EventSubscriber.Subscribe("AUTH", "ACTIVATION_REQUEST_SENT")
		if err != nil {
			log.Fatal(fmt.Errorf("subscribing: %w", err))
		}
		fmt.Println(<-eventBus)
		// TODO Switch statement based on message stream/subject.
		// Invoke behaviour that marks a user as PENDING ACTIVATION
	}
}

func NewAuthEventHandler(app auth.App, eventSubscriber eventbus.EventSubscriber) AuthEventHandler {
	return AuthEventHandler{
		App:             app,
		EventSubscriber: eventSubscriber,
	}
}
