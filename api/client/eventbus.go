// Package client contains all interfaces the APIs need to implement their domain logic. The implementation of these
// clients are injected as dependencies upon initialisation of APIs.
package client

// EventBusClient offers the caller the interface to Subscribe or Publish to the EventBus, encapsulating its
// technical implementation.
type EventBusClient interface {
	Subscribe(Stream, Subject) (EventBus, error)
	Publish(Event) error
}

type EventBus chan Event

type Subject string
type Stream string
type Body string

type Event struct {
	Stream  Stream
	Subject Subject
	Body    Body
}

type Subscription struct {
	EventBus EventBus
	Stream   Stream
	Subject  Subject
}
