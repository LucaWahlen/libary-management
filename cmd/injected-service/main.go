package main

import (
	"libary-service/internal/app"
	"libary-service/internal/repository/postgres"
	"libary-service/internal/router"
	"libary-service/internal/validation/validator"
)

func main() {
	repository := postgresrepository.New()

	repository.Connect()
	defer repository.Disconnect()

	validator := validator.NewValidator(repository)

	service := app.NewService(repository, validator)

	router := router.NewGinRouter(service)

	router.Serve(":8080")
}
