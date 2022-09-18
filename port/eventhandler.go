package port

type EventHandler interface {
	// Handle listens for events that match the Stream or Subject and invokes the appropriate domain behaviour.
	Handle()
}
