package gin

import (
	"github.com/gin-gonic/gin"
	"libary-service/internal/app"
	"net/http"
)

type GinRouter struct {
	Engine *gin.Engine
}

func NewGinRouter(service app.Service) *GinRouter {
	r := GinRouter{Engine: gin.Default()}

	r.GET("/books", service.GetBooks)
	r.GET("/books/:id", service.GetBookByID)
	r.POST("/books", service.CreateBook)
	r.PUT("/books/:id", service.UpdateBook)
	r.DELETE("/books/:id", service.DeleteBook)

	r.GET("/users", service.GetUsers)
	r.GET("/users/:id", service.GetUserByID)
	r.POST("/users", service.CreateUser)
	r.PUT("/users/:id", service.UpdateUser)
	r.DELETE("/users/:id", service.DeleteUser)

	r.GET("/lendings", service.GetLendings)
	r.GET("/lendings/:id", service.GetLendingByID)
	r.POST("/lendings", service.CreateLending)
	r.PUT("/lendings/:id", service.UpdateLending)
	r.DELETE("/lendings/:id", service.DeleteLending)

	return &r
}

func (r *GinRouter) GET(path string, handler http.HandlerFunc) {
	r.Engine.GET(path, gin.WrapF(handler))
}

func (r *GinRouter) POST(path string, handler http.HandlerFunc) {
	r.Engine.POST(path, gin.WrapF(handler))
}

func (r *GinRouter) PUT(path string, handler http.HandlerFunc) {
	r.Engine.PUT(path, gin.WrapF(handler))
}

func (r *GinRouter) DELETE(path string, handler http.HandlerFunc) {
	r.Engine.DELETE(path, gin.WrapF(handler))
}

func (r *GinRouter) Serve(addr string) error {
	return r.Engine.Run(addr)
}
