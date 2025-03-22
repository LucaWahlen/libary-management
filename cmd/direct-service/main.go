package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"libary-service/internal/direct-service/app"
	"libary-service/internal/direct-service/repository"
	"log"
	"os"
)

func main() {
	//Connect to DB
	dbURL := os.Getenv("DATABASE_URL")
	var err error
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully connected to database")
	repository.DB = conn
	defer conn.Close(context.Background())

	//HTTP Router
	router := gin.Default()

	router.GET("/books", gin.WrapF(app.GetBooks))
	router.GET("/books/:id", gin.WrapF(app.GetBookByID))
	router.POST("/books", gin.WrapF(app.CreateBook))
	router.PUT("/books/:id", gin.WrapF(app.UpdateBook))
	router.DELETE("/books/:id", gin.WrapF(app.DeleteBook))

	router.GET("/users", gin.WrapF(app.GetUsers))
	router.GET("/users/:id", gin.WrapF(app.GetUserByID))
	router.POST("/users", gin.WrapF(app.CreateUser))
	router.PUT("/users/:id", gin.WrapF(app.UpdateUser))
	router.DELETE("/users/:id", gin.WrapF(app.DeleteUser))

	router.GET("/lendings", gin.WrapF(app.GetLendings))
	router.GET("/lendings/:id", gin.WrapF(app.GetLendingByID))
	router.POST("/lendings", gin.WrapF(app.CreateLending))
	router.PUT("/lendings/:id", gin.WrapF(app.UpdateLending))
	router.DELETE("/lendings/:id", gin.WrapF(app.DeleteLending))

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
