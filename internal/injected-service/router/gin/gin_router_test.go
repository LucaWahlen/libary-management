package gin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"libary-service/generated/mocks"
)

func TestRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.Service)
	routes := []struct {
		method        string
		path          string
		serviceMethod string
		status        int
		response      string
	}{
		{"GET", "/books", "GetBooks", http.StatusOK, "mocked GetBooks"},
		{"GET", "/books/123", "GetBookByID", http.StatusOK, "mocked GetBookByID"},
		{"POST", "/books", "CreateBook", http.StatusCreated, "mocked CreateBook"},
		{"PUT", "/books/123", "UpdateBook", http.StatusOK, "mocked UpdateBook"},
		{"DELETE", "/books/123", "DeleteBook", http.StatusNoContent, ""},
		{"GET", "/users", "GetUsers", http.StatusOK, "mocked GetUsers"},
		{"GET", "/users/123", "GetUserByID", http.StatusOK, "mocked GetUserByID"},
		{"POST", "/users", "CreateUser", http.StatusCreated, "mocked CreateUser"},
		{"PUT", "/users/123", "UpdateUser", http.StatusOK, "mocked UpdateUser"},
		{"DELETE", "/users/123", "DeleteUser", http.StatusNoContent, ""},
		{"GET", "/lendings", "GetLendings", http.StatusOK, "mocked GetLendings"},
		{"GET", "/lendings/123", "GetLendingByID", http.StatusOK, "mocked GetLendingByID"},
		{"POST", "/lendings", "CreateLending", http.StatusCreated, "mocked CreateLending"},
		{"PUT", "/lendings/123", "UpdateLending", http.StatusOK, "mocked UpdateLending"},
		{"DELETE", "/lendings/123", "DeleteLending", http.StatusNoContent, ""},
	}
	for _, route := range routes {
		if route.serviceMethod == "DeleteBook" || route.serviceMethod == "DeleteUser" || route.serviceMethod == "DeleteLending" {
			mockService.On(route.serviceMethod, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(route.status)
			}).Once()
		} else {
			mockService.On(route.serviceMethod, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
				w := args.Get(0).(http.ResponseWriter)
				w.WriteHeader(route.status)
				if route.response != "" {
					w.Write([]byte(route.response))
				}
			}).Once()
		}
	}
	r := NewGinRouter(mockService)
	for _, route := range routes {
		req, err := http.NewRequest(route.method, route.path, nil)
		assert.NoError(t, err)
		recorder := httptest.NewRecorder()
		r.Engine.ServeHTTP(recorder, req)
		assert.Equal(t, route.status, recorder.Code)
		assert.Equal(t, route.response, recorder.Body.String())
	}
	mockService.AssertExpectations(t)
}

func TestServe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.Service)
	r := NewGinRouter(mockService)
	err := r.Serve("invalid")
	assert.Error(t, err)
}
