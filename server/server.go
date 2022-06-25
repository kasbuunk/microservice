// Package server defines the interface of a server as a source of input to the program.
package server

type Port int

type Server interface {
	Serve(Port) error
}
