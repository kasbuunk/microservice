package listener

import (
	"fmt"
	"log"

	"github.com/kasbuunk/microservice/api/email"
	"github.com/kasbuunk/microservice/event"
)

type Email struct {
	API email.API
	Bus event.Bus
}

// Listen listens for messages that match the Stream or Subject.
func (s Email) Listen() {
	fmt.Println("API service listening for messages.")
	for {
		// starts process in loop, in goroutine that awaits published messages and invokes api calls
		messageChannel, err := s.Bus.Subscribe("AUTH", "USER_REGISTERED")
		if err != nil {
			log.Fatal(fmt.Errorf("subscribing: %w", err))
		}
		fmt.Println(<-messageChannel)
		// TODO Switch statement based on message stream/subject.
		// Invoke behaviour that marks a user as PENDING ACTIVATION
	}
}

func NewEmail(api email.API, bus event.Bus) Listener {
	return Email{
		API: api,
		Bus: bus,
	}
}
