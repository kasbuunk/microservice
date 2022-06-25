package main

import (
	"log"

	"github.com/kasbuunk/microservice/api/auth"
	"github.com/kasbuunk/microservice/api/email"
	"github.com/kasbuunk/microservice/client/email/postmark"
	"github.com/kasbuunk/microservice/client/eventbus"
	"github.com/kasbuunk/microservice/client/eventbus/eventbus"
	"github.com/kasbuunk/microservice/client/storage"
	"github.com/kasbuunk/microservice/client/userrepo/userrepoclient"
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

	// Initialise clients.
	userRepo := userrepoclient.New(db)
	emailClient := emailclient.New(conf.Postmark)
	eventBusClient := eventbusclient.New([]eventbus.Stream{
		"AUTH",
		"EMAIL",
	})

	// Initialise APIs that implement all core domain logic, injecting dependencies.
	authAPI := auth.New(userRepo, eventBusClient)
	emailAPI := email.New(eventBusClient, emailClient)

	// Initialise sources of input: servers and listeners.
	authEventHandler := authhandler.New(authAPI, eventBusClient)
	emailEventHandler := emailhandler.New(emailAPI, eventBusClient)

	// Start processes that listen for events.
	go authEventHandler.Handle()
	go emailEventHandler.Handle()

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
