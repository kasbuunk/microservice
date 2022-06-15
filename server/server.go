// Package server exposes the Server interface that runs and listens for requests.
// It accepts configuration on where to run and how to store and retrieve entities.
package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/kasbuunk/microservice/graph"
	"github.com/kasbuunk/microservice/graph/generated"
	"github.com/kasbuunk/microservice/repository"
)

type Server interface {
	Serve(Config) error
}

type Service struct{}

type Config struct {
	Port        int
	GQLEndpoint string
}

// New takes an endpoint  returns a new server
func New(conf Config, repo repository.Repository) (Server, error) {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: &graph.Resolver{
					UserRepository: repo,
				}}))

	http.Handle("/", playground.Handler("GraphQL playground", conf.GQLEndpoint))
	http.Handle(conf.GQLEndpoint, srv)

	return Service{}, nil
}

func (s Service) Serve(conf Config) error {
	log.Printf("connect to http://localhost:%d/ for GraphQL playground", conf.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil))
	return nil
}
