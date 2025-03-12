package router

import (
	"net/http"
)

type Router interface {
	GET(path string, handler http.HandlerFunc)
	POST(path string, handler http.HandlerFunc)
	PUT(path string, handler http.HandlerFunc)
	DELETE(path string, handler http.HandlerFunc)
	Serve(addr string) error
}
