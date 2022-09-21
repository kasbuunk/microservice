// Package gql implements the Server interface that runs and listens for requests.
// It accepts configuration on where to run and how to store and retrieve entities.
package gql

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"

	"github.com/kasbuunk/microservice/auth"
	"github.com/kasbuunk/microservice/transport/gql/graphql"
	"github.com/kasbuunk/microservice/transport/gql/graphql/generated"
)

type Config struct {
	Port     int
	Endpoint string
}

type Server struct{}

func (s Server) Serve(port int) error {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	return nil
}

// New takes an endpoint  returns a new server
func New(auth auth.App, endpoint string) Server {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					Auth: auth,
				}}))

	http.Handle(endpoint, srv)

	return Server{}
}
