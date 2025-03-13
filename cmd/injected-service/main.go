package main

import (
	"libary-service/internal/app"
	"libary-service/internal/repository/postgres"
	"libary-service/internal/router/gin"
	"libary-service/internal/validation/validator"
)

func main() {
	repository := postgresrepository.New()
	if err := repository.Connect(); err != nil {
		panic(err)
	}
	defer repository.Disconnect()
	validator := validator.New(repository)
	service := app.NewLibaryService(repository, validator)
	router := gin.NewGinRouter(service)
	router.Serve(":8080")
}
