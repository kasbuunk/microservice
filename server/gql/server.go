// Package server exposes the Server interface that runs and listens for requests.
// It accepts configuration on where to run and how to store and retrieve entities.
package gqlserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/kasbuunk/microservice/api/auth"
	"github.com/kasbuunk/microservice/server/gql/graphql"
	"github.com/kasbuunk/microservice/server/gql/graphql/generated"
)

type Server interface {
	Serve(Port) error
}

type Service struct{}

type Config struct {
	Port        Port
	GQLEndpoint GQLEndpoint
}

type Port int
type GQLEndpoint string

// New takes an endpoint  returns a new server
func New(endpoint GQLEndpoint, auth auth.API) (Server, error) {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					Auth: auth,
				}}))

	http.Handle("/", playground.Handler("GraphQL playground", string(endpoint)))
	http.Handle(string(endpoint), srv)

	return Service{}, nil
}

func (s Service) Serve(port Port) error {
	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	return nil
}
