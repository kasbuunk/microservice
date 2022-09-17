package main

import (
	"log"

	"github.com/kasbuunk/microservice/adapter/eventbus"
	"github.com/kasbuunk/microservice/app/auth/adapter/repository/storage"
	"github.com/kasbuunk/microservice/app/auth/adapter/repository/user"
	authapp "github.com/kasbuunk/microservice/app/auth/core"
	"github.com/kasbuunk/microservice/app/email/adapter/email"
	emailapp "github.com/kasbuunk/microservice/app/email/core"
	"github.com/kasbuunk/microservice/config"
	"github.com/kasbuunk/microservice/event/auth"
	"github.com/kasbuunk/microservice/event/email"
	"github.com/kasbuunk/microservice/server/gql"
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

	// Initialise adapters.
	userRepo := userrepo.New(db)
	emailClient := emailclient.New(conf.Postmark)
	eventBus := eventbus.New([]string{
		"AUTH",
		"EMAIL",
	})

	// Initialise Apps that implement all core domain logic, injecting dependencies.
	authApp := authapp.New(userRepo, eventBus)
	emailApp := emailapp.New(eventBus, emailClient)

	// Initialise sources of input: event handlers.
	go authhandler.New(authApp, eventBus).Handle()
	go emailhandler.New(emailApp, eventBus).Handle()

	// Initialise Graphql http server.
	authServer, err := gql.New(conf.GQLServer.Endpoint, authApp)
	if err != nil {
		log.Fatalf("Initialisation of server failed: %v", err)
	}

	// Start process that serves GraphQL requests.
	err = authServer.Serve(conf.GQLServer.Port)
	if err != nil {
		log.Fatalf("Serving failed: %v", err)
	}
}
