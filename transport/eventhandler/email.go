package eventhandler

import (
	"fmt"
	"log"

	"github.com/kasbuunk/microservice/email"
	"github.com/kasbuunk/microservice/eventbus"
)

type EmailEventHandler struct {
	App             email.App
	EventSubscriber eventbus.EventSubscriber
}

// Handle listens for events that match the Stream or Subject and invokes the appropriate domain behaviour.
func (s EmailEventHandler) Handle() {
	fmt.Println("Email service listening for events.")
	for {
		// Starts process in loop, awaiting published messages.
		eventBus, err := s.EventSubscriber.Subscribe("AUTH", "USER_REGISTERED")
		if err != nil {
			log.Fatal(fmt.Errorf("subscribing: %w", err))
		}
		fmt.Println(<-eventBus)
		// TODO Switch statement based on message stream/subject.
		// Invoke behaviour that marks a user as PENDING ACTIVATION
	}
}

func NewEmailEventHandler(app email.App, eventSubscriber eventbus.EventSubscriber) EmailEventHandler {
	return EmailEventHandler{
		App:             app,
		EventSubscriber: eventSubscriber,
	}
}
