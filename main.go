package main

import (
	"log"

	"github.com/kasbuunk/microservice/auth"
	"github.com/kasbuunk/microservice/config"
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

	repo := repository.New(db)

	authService := auth.New(repo)

	svc, err := server.New(conf.Server, authService)
	if err != nil {
		log.Fatalf("Initialisation of server failed: %v", err)
	}

	err = svc.Serve(conf.Server)
	if err != nil {
		panic(err)
	}
}
