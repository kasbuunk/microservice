package main

import (
	"log"

	"github.com/kasbuunk/microservice/api/auth"
	"github.com/kasbuunk/microservice/api/client"
	"github.com/kasbuunk/microservice/api/email"
	"github.com/kasbuunk/microservice/client/email"
	"github.com/kasbuunk/microservice/client/eventbus"
	"github.com/kasbuunk/microservice/client/repository/storage"
	"github.com/kasbuunk/microservice/client/repository/user"
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
	userRepo := userrepo.New(db)
	emailClient := emailclient.New(conf.Postmark)
	busClient := eventbusclient.New([]client.Stream{
		"AUTH",
		"EMAIL",
	})

	// Initialise APIs that implement all core domain logic, injecting dependencies.
	authAPI := auth.New(userRepo, busClient)
	emailAPI := email.New(busClient, emailClient)

	// Initialise sources of input: servers and listeners.
	authSubscriber := authsubscriber.New(authAPI, busClient)
	emailSubscriber := emailsubscriber.New(emailAPI, busClient)

	// Start processes that listen for events.
	go authSubscriber.SubscribeToEvents()
	go emailSubscriber.SubscribeToEvents()

	// Initialise Graphql http server.
	authServer, err := gqlserver.New(conf.Server.GQLEndpoint, authAPI)
	if err != nil {
		log.Fatalf("Initialisation of server failed: %v", err)
	}

	// Start process that serves GraphQL requests.
	err = authServer.Serve(conf.Server.Port)
	if err != nil {
		log.Fatalf("Serving failed: %v", err)
	}
}
