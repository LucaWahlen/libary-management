//go:generate mockery --name=Service --output=../../generated/mocks --case=underscore
package app

import (
	"net/http"
)

type Service interface {
	GetBooks(w http.ResponseWriter, r *http.Request)
	GetBookByID(w http.ResponseWriter, r *http.Request)
	CreateBook(w http.ResponseWriter, r *http.Request)
	UpdateBook(w http.ResponseWriter, r *http.Request)
	DeleteBook(w http.ResponseWriter, r *http.Request)

	GetUsers(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)

	GetLendings(w http.ResponseWriter, r *http.Request)
	GetLendingByID(w http.ResponseWriter, r *http.Request)
	CreateLending(w http.ResponseWriter, r *http.Request)
	UpdateLending(w http.ResponseWriter, r *http.Request)
	DeleteLending(w http.ResponseWriter, r *http.Request)
}
