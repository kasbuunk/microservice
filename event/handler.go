// Package event defines the interface of a Handler as a source of input to the program.
package event

type Handler interface {
	// Handle listens for events that match the Stream or Subject and invokes the appropriate domain behaviour.
	Handle()
}
