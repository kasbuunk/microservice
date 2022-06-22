package main

import (
	"log"

	"github.com/kasbuunk/microservice/api/auth"
	"github.com/kasbuunk/microservice/api/email"
	"github.com/kasbuunk/microservice/client/postmark"
	"github.com/kasbuunk/microservice/client/repository/storage"
	"github.com/kasbuunk/microservice/client/repository/user"
	"github.com/kasbuunk/microservice/config"
	"github.com/kasbuunk/microservice/event"
	"github.com/kasbuunk/microservice/input/listener"
	"github.com/kasbuunk/microservice/input/server"
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

	streams := []event.Stream{
		"AUTH",
		"EMAIL",
	}
	bus := event.NewMessageBus(streams)

	// Initialise repositories and clients.
	userRepo := user.New(db)
	emailClient := postmark.New(conf.Postmark)

	// Initialise APIs
	authAPI := auth.New(userRepo, bus)
	emailAPI := email.New(bus, emailClient)

	// Initialise sources of input: servers and listeners.
	authSubscriber := listener.NewAuth(authAPI, bus)
	emailSubscriber := listener.NewEmail(emailAPI, bus)

	authServer, err := server.New(conf.Server.GQLEndpoint, authAPI)
	if err != nil {
		log.Fatalf("Initialisation of server failed: %v", err)
	}

	// Start process that listens for requests.
	err = authServer.Serve(conf.Server.Port)
	if err != nil {
		log.Fatalf("Serving failed: %v", err)
	}

	// Start processes that listen for events.
	go authSubscriber.Listen()
	go emailSubscriber.Listen()

}
