package eventbus

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/kasbuunk/microservice/port"
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
	orderStream := port.Stream("ORDER")
	invoiceStream := port.Stream("INVOICE")
	incomingEvents, err := s.EventBus.Subscribe(orderStream, "*")
	assert.NoError(s.T(), err)

	orderEvent := port.Event{
		Stream:  orderStream,
		Subject: "order placed",
	}
	invoiceEvent := port.Event{
		Stream:  invoiceStream,
		Subject: "invoice paid",
	}
	go func() {
		for {
			select {
			case incomingEvent := <-incomingEvents:
				s.Equal(orderEvent, incomingEvent)
				s.NotEqual(invoiceEvent.Stream, incomingEvent)
				return
			}
		}
	}()
	err = s.EventBus.Publish(orderEvent)
	assert.NoError(s.T(), err)
	time.Sleep(20 * time.Millisecond)
}

func TestEventBusTestSuite(t *testing.T) {
	suite.Run(t, new(EventBusTestSuite))
}
