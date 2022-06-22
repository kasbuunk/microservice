package main

import (
	"github.com/kasbuunk/microservice/client/postmark"
	"log"

	"github.com/kasbuunk/microservice/api/auth"
	"github.com/kasbuunk/microservice/api/email"
	"github.com/kasbuunk/microservice/config"
	"github.com/kasbuunk/microservice/event"
	"github.com/kasbuunk/microservice/listener"
	"github.com/kasbuunk/microservice/repository"
	"github.com/kasbuunk/microservice/server"
	"github.com/kasbuunk/microservice/storage"
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
	userRepo := repository.New(db)
	emailClient := postmark.New(conf.Postmark)

	authAPI := auth.New(userRepo, bus)

	emailAPI := email.New(bus, emailClient)

	svc, err := server.New(conf.Server.GQLEndpoint, authAPI)
	if err != nil {
		log.Fatalf("Initialisation of server failed: %v", err)
	}

	authSubscriber := listener.NewAuth(authAPI, bus)
	go authSubscriber.Listen()

	emailSubscriber := listener.NewEmail(emailAPI, bus)
	go emailSubscriber.Listen()

	err = svc.Serve(conf.Server.Port)
	if err != nil {
		log.Fatalf("Serving failed: %v", err)
	}
}
