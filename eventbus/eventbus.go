// Package eventbus contains all interfaces the Apps need to implement their domain logic. The adapters
// implement the ports and are injected as dependencies upon initialisation of Apps.
package eventbus

// EventBus offers the caller the interface to Subscribe or Publish to the EventBus, encapsulating its
// technical implementation.
type EventBus interface {
	Subscribe(Stream, Subject) (chan Event, error)
	Publish(Event) error
}

type EventSubscriber interface {
	Subscribe(Stream, Subject) (chan Event, error)
}

type EventPublisher interface {
	Publish(Event) error
}

type Subject string
type Stream string
type Body string

type Event struct {
	Stream  Stream
	Subject Subject
	Body    Body
}

type Subscription struct {
	EventBus chan Event
	Stream   Stream
	Subject  Subject
}
