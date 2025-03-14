package main

import (
	"libary-service/internal/direct-service/repository"
	"libary-service/internal/direct-service/router"
	"log"
)

func main() {

	if err := repository.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repository.Disconnect()
	router.NewGinRouter()
	if err := router.Serve(":8080"); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
