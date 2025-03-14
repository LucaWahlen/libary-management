package router

import (
	"github.com/gin-gonic/gin"
	"libary-service/internal/direct-service/app"
	"net/http"
)

var r *gin.Engine

func NewGinRouter() {
	r = gin.Default()

	GET("/books", app.GetBooks)
	GET("/books/:id", app.GetBookByID)
	POST("/books", app.CreateBook)
	PUT("/books/:id", app.UpdateBook)
	DELETE("/books/:id", app.DeleteBook)

	GET("/users", app.GetUsers)
	GET("/users/:id", app.GetUserByID)
	POST("/users", app.CreateUser)
	PUT("/users/:id", app.UpdateUser)
	DELETE("/users/:id", app.DeleteUser)

	GET("/lendings", app.GetLendings)
	GET("/lendings/:id", app.GetLendingByID)
	POST("/lendings", app.CreateLending)
	PUT("/lendings/:id", app.UpdateLending)
	DELETE("/lendings/:id", app.DeleteLending)
}

func GET(path string, handler http.HandlerFunc) {
	r.GET(path, gin.WrapF(handler))
}

func POST(path string, handler http.HandlerFunc) {
	r.POST(path, gin.WrapF(handler))
}

func PUT(path string, handler http.HandlerFunc) {
	r.PUT(path, gin.WrapF(handler))
}

func DELETE(path string, handler http.HandlerFunc) {
	r.DELETE(path, gin.WrapF(handler))
}

func Serve(addr string) error {
	return r.Run(addr)
}
