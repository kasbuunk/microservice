package main

import (
	"log"

	"github.com/kasbuunk/microservice/api/auth"
	"github.com/kasbuunk/microservice/api/email"
	"github.com/kasbuunk/microservice/config"
	"github.com/kasbuunk/microservice/events"
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

	streams := []events.Stream{
		"AUTH",

		"EMAIL",
	}
	bus := events.NewMessageBus(streams)

	userRepo := repository.New(db)
	authService := auth.New(userRepo, bus)
	go authService.Subscribe()

	emailService := email.New(bus)
	go emailService.Subscribe()

	svc, err := server.New(conf.Server.GQLEndpoint, authService)
	if err != nil {
		log.Fatalf("Initialisation of server failed: %v", err)
	}

	err = svc.Serve(conf.Server.Port)
	if err != nil {
		log.Fatalf("Serving failed: %v", err)
	}

}
