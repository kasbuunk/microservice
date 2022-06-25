// Package event defines the interface of a Subscriber as a source of input to the program.
package event

type Subscriber interface {
	SubscribeToEvents()
}
