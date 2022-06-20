package events

type Subject string
type Stream string

type Message struct {
	Stream  Stream
	Subject Subject
	Body    string
}

type Subscription struct {
	//Subscriber Subscriber
	Connection chan Message
	Stream     Stream
	Subject    Subject
}

type MessageBus interface {
	Subscribe(Stream, Subject) (chan Message, error)
	Publish(Message) error
}

type Bus struct {
	// Holds Subscriptions in memory for now, should be delegated elsewhere
	// to remain stateless in case of horizontal scaling. Streams do not change at runtime.
	Streams       []Stream
	Subscriptions []Subscription
}

func (b *Bus) Publish(msg Message) error {
	// For all subscribers that match the msg,
	for _, subscription := range b.Subscriptions { // b.Subscriptions() when delegated state.
		if subscribed(subscription, msg) {
			// send the msg to the sub
			subscription.Connection <- msg
		}
	}
	// Never return messages for now, only log. Perhaps return errors when the pubsub system is delegated.
	return nil
}

func (b *Bus) Subscribe(stream Stream, subject Subject) (chan Message, error) {
	c := make(chan Message)

	subscription := Subscription{
		Connection: c,
		Stream:     stream,
		Subject:    subject,
	}

	b.Subscriptions = append(b.Subscriptions, subscription)

	return c, nil
}

// NewMessageBus is initialised with a predetermined set of streams. Its subscriptions
// should be added after initialisation, upon passing it to the services. The services
// themselves are responsible for calling the method that adds their subscription.
func NewMessageBus(streams []Stream) MessageBus {
	return &Bus{
		Streams:       streams,
		Subscriptions: []Subscription{},
	}
}

func subscribed(sub Subscription, msg Message) bool {
	// TODO: allow regex matching.
	if sub.Subject == msg.Subject || sub.Stream == msg.Stream {
		return true
	}
	return false
}
