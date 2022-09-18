// Package gql exposes the Server interface that runs and listens for requests.
// It accepts configuration on where to run and how to store and retrieve entities.
package gql

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/kasbuunk/microservice/app/adapter/gql/graphql"
	"github.com/kasbuunk/microservice/app/adapter/gql/graphql/generated"
	"github.com/kasbuunk/microservice/app/auth"
)

type Server struct{}

type Config struct {
	Port     int
	Endpoint string
}

// New takes an endpoint  returns a new server
func New(endpoint string, auth auth.App) (Server, error) {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					Auth: auth,
				}}))

	http.Handle("/", playground.Handler("GraphQL playground", endpoint))
	http.Handle(endpoint, srv)

	return Server{}, nil
}

func (s Server) Serve(port int) error {
	log.Printf("connect to http://localhost:%d/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	return nil
}
