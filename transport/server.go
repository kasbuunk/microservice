package transport

type Server interface {
	Serve(int) error
}
