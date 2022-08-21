package main

import (
	"log"

	authapp "github.com/kasbuunk/microservice/app/auth/core"
	emailapp "github.com/kasbuunk/microservice/app/email/core"
	"github.com/kasbuunk/microservice/app/eventbus"
	"github.com/kasbuunk/microservice/config"
	"github.com/kasbuunk/microservice/event/auth"
	"github.com/kasbuunk/microservice/event/email"
	"github.com/kasbuunk/microservice/server/gql"
	"github.com/kasbuunk/microservice/side-effect/email"
	"github.com/kasbuunk/microservice/side-effect/eventbus"
	"github.com/kasbuunk/microservice/side-effect/repository/storage"
	"github.com/kasbuunk/microservice/side-effect/repository/user"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatalf("Loading environment configuration failed: %v", err)
	}

	db, err := storage.Connect(conf.DB)
	if err != nil {
		log.Fatalf("Connection to storage failed: %v", err)
	}

	// Initialise clients.
	userRepo := userrepo.New(db)
	emailClient := emailclient.New(conf.Postmark)
	eventBus := eventbusclient.New([]eventbus.Stream{
		"AUTH",
		"EMAIL",
	})

	// Initialise APIs that implement all core domain logic, injecting dependencies.
	authAPI := authapp.New(userRepo, eventBus)
	emailAPI := emailapp.New(eventBus, emailClient)

	// Initialise sources of input: event handlers.
	go authhandler.New(authAPI, eventBus).Handle()
	go emailhandler.New(emailAPI, eventBus).Handle()

	// Initialise Graphql http server.
	authServer, err := gqlserver.New(conf.GQLServer.Endpoint, authAPI)
	if err != nil {
		log.Fatalf("Initialisation of server failed: %v", err)
	}

	// Start process that serves GraphQL requests.
	err = authServer.Serve(conf.GQLServer.Port)
	if err != nil {
		log.Fatalf("Serving failed: %v", err)
	}
}
