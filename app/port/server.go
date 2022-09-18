package port

type Server interface {
	Serve(int) error
}
