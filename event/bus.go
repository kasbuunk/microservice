package event

type Subject string
type Stream string
type Body string

type Message struct {
	Stream  Stream
	Subject Subject
	Body    Body
}

type Subscription struct {
	Connection chan Message
	Stream     Stream
	Subject    Subject
}

type Bus interface {
	Subscribe(Stream, Subject) (chan Message, error)
	Publish(Message) error
}

type Subscriber interface {
	Subscribe(Stream, Subject) (chan Message, error)
}

type Publisher interface {
	Publish(Message) error
}

type bus struct {
	// Holds Subscriptions in memory for now, might be delegated elsewhere
	// to remain stateless in case of horizontal scaling. Streams do not change at runtime.
	Streams       []Stream
	Subscriptions []Subscription
}

func (b *bus) Publish(msg Message) error {
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

func (b *bus) Subscribe(stream Stream, subject Subject) (chan Message, error) {
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
func NewMessageBus(streams []Stream) Bus {
	return &bus{
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
