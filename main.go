package main

import (
	"log"

	authapp "github.com/kasbuunk/microservice/auth/core"
	"github.com/kasbuunk/microservice/auth/port/repository/user"
	"github.com/kasbuunk/microservice/config"
	emailapp "github.com/kasbuunk/microservice/email/core"
	"github.com/kasbuunk/microservice/email/port/email"
	"github.com/kasbuunk/microservice/eventbus/localbus"
	"github.com/kasbuunk/microservice/internal/storage"
	"github.com/kasbuunk/microservice/transport/eventhandler"
	"github.com/kasbuunk/microservice/transport/gql"
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
	eventBus := localbus.New([]string{
		"AUTH",
		"EMAIL",
	})

	// Initialise Apps that implement all core domain logic, injecting dependencies.
	authApp := authapp.New(eventBus, userRepo)
	emailApp := emailapp.New(eventBus, emailClient)

	// Initialise sources of input: event handlers.
	go eventhandler.NewAuthEventHandler(authApp, eventBus).Handle()
	go eventhandler.NewEmailEventHandler(emailApp, eventBus).Handle()

	// Initialise Graphql http server.
	authServer := gql.New(authApp, conf.GQLServer.Endpoint)

	// Start process that serves GraphQL requests.
	err = authServer.Serve(conf.GQLServer.Port)
	if err != nil {
		log.Fatalf("Serving failed: %v", err)
	}
}
