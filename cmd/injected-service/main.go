package main

import (
	"libary-service/internal/injected-service/app"
	"libary-service/internal/injected-service/repository/postgres"
	"libary-service/internal/injected-service/router/gin"
	"libary-service/internal/injected-service/validation/validator"
	"log"
)

func main() {
	repository := postgresrepository.New()
	if err := repository.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repository.Disconnect()
	validator := validator.New(repository)
	service := app.NewLibaryService(repository, validator)
	router := gin.NewGinRouter(service)
	if err := router.Serve(":8080"); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
