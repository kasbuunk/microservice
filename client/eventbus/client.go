// Package eventbusclient implements how events are published, transmitted and subscribed to.
// Hence, the domain core remains agnostic of how its events are distributed amongst services that
// subscribe. Here, the implementation of the client can be freely changed to connect to an external
// event store, such as Apache Kafka or NATS JetStream.
//
// At this moment, the implementation is in-memory, such that no network call is necessary. This
// suffices for further development in the domain core, while keeping the implementation abstracted
// away.
package eventbusclient

import "github.com/kasbuunk/microservice/api/client"

// eventBusClient implements the EventBusClient interface through which the caller can Subscribe to and Publish events.
type eventBusClient struct {
	// Holds Subscriptions in memory for now, might be delegated elsewhere
	// to remain stateless in case of horizontal scaling. Streams do not change at runtime.
	Streams       []client.Stream
	Subscriptions []client.Subscription
}

func (b *eventBusClient) Publish(msg client.Event) error {
	// For all subscribers that match the msg,
	for _, subscription := range b.Subscriptions { // b.Subscriptions() when delegated state.
		if subscribed(subscription, msg) {
			// send the msg to the sub
			subscription.EventBus <- msg
		}
	}
	// Never return an error, until the pubsub system is delegated to an external service.
	return nil
}

func (b *eventBusClient) Subscribe(stream client.Stream, subject client.Subject) (client.EventBus, error) {
	eventBus := make(client.EventBus)

	subscription := client.Subscription{
		EventBus: eventBus,
		Stream:   stream,
		Subject:  subject,
	}

	b.Subscriptions = append(b.Subscriptions, subscription)

	// Never return an error, until the pubsub system is delegated to an external service.
	return eventBus, nil
}

// New is initialised with a predetermined set of streams. Its subscriptions
// should be added after initialisation, upon passing it to the services. The services
// themselves are responsible for calling the method that adds their subscription.
func New(streams []client.Stream) client.EventBusClient {
	return &eventBusClient{
		Streams:       streams,
		Subscriptions: []client.Subscription{},
	}
}

func subscribed(sub client.Subscription, msg client.Event) bool {
	// TODO: allow regex matching.
	if sub.Subject == msg.Subject || sub.Stream == msg.Stream {
		return true
	}
	return false
}
