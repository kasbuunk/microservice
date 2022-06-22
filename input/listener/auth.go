package listener

import (
	"fmt"
	"log"

	"github.com/kasbuunk/microservice/api/auth"
	"github.com/kasbuunk/microservice/event"
)

type Auth struct {
	API        auth.API
	Subscriber event.Subscriber
}

// Listen listens for messages that match the Stream or Subject.
func (s Auth) Listen() {
	fmt.Println("API service listening for messages.")
	for {
		// starts process in loop, in goroutine that awaits published messages and invokes api calls
		messageChannel, err := s.Subscriber.Subscribe("EMAIL", "ACTIVATION_REQUEST_SENT")
		if err != nil {
			log.Fatal(fmt.Errorf("subscribing: %w", err))
		}
		fmt.Println(<-messageChannel)
		// TODO Switch statement based on message stream/subject.
		// Invoke behaviour that marks a user as PENDING ACTIVATION
	}
}

func NewAuth(api auth.API, sub event.Subscriber) Listener {
	return Auth{
		API:        api,
		Subscriber: sub,
	}
}
