package localbus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/kasbuunk/microservice/eventbus"
)

type EventBusTestSuite struct {
	suite.Suite
	EventBus *EventBus
}

func (s *EventBusTestSuite) SetupTest() {
	s.EventBus = New([]string{"ORDER", "INVOICE", "USER"})
}

func (s *EventBusTestSuite) TestPubSub() {
	// Init some participants in the event bus. Some subscribers and publishers.
	orderStream := eventbus.Stream("ORDER")
	invoiceStream := eventbus.Stream("INVOICE")
	incomingEvents, err := s.EventBus.Subscribe(orderStream, "*")
	assert.NoError(s.T(), err)

	orderEvent := eventbus.Event{
		Stream:  orderStream,
		Subject: "order placed",
	}
	invoiceEvent := eventbus.Event{
		Stream:  invoiceStream,
		Subject: "invoice paid",
	}
	go func() {
		incomingEvent := <-incomingEvents
		s.Equal(orderEvent, incomingEvent)
		s.NotEqual(invoiceEvent.Stream, incomingEvent)
	}()

	err = s.EventBus.Publish(orderEvent)
	assert.NoError(s.T(), err)
	time.Sleep(20 * time.Millisecond)
}

func TestEventBusTestSuite(t *testing.T) {
	suite.Run(t, new(EventBusTestSuite))
}
