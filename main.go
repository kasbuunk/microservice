package main

import (
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

	userRepo := repository.New(db)

	authAPI := auth.New(userRepo, bus)
	authSubscriber := listener.NewAuth(authAPI, bus)
	go authSubscriber.Listen()

	emailAPI := email.New(bus)
	emailSubscriber := listener.NewEmail(emailAPI, bus)
	go emailSubscriber.Listen()

	svc, err := server.New(conf.Server.GQLEndpoint, authAPI)
	if err != nil {
		log.Fatalf("Initialisation of server failed: %v", err)
	}

	err = svc.Serve(conf.Server.Port)
	if err != nil {
		log.Fatalf("Serving failed: %v", err)
	}
}
