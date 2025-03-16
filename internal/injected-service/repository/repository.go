//go:generate mockery --name=Repository --output=../../generated/mocks --case=underscore
package repository

import (
	"libary-service/internal/domain"
)

// Repository aggregates all data access methods.
type Repository interface {
	Connect() error
	Disconnect() error

	GetBooks() ([]domain.Book, error)
	GetBookByID(id string) (domain.Book, error)
	CreateBook(book domain.Book) (domain.Book, error)
	UpdateBook(book domain.Book) (domain.Book, error)
	DeleteBook(id string) error

	GetUsers() ([]domain.User, error)
	GetUserByID(id string) (domain.User, error)
	CreateUser(user domain.User) (domain.User, error)
	UpdateUser(user domain.User) (domain.User, error)
	DeleteUser(id string) error

	GetLendings() ([]domain.Lending, error)
	GetLendingByID(id string) (domain.Lending, error)
	CreateLending(lending domain.Lending) (domain.Lending, error)
	UpdateLending(lending domain.Lending) (domain.Lending, error)
	DeleteLending(id string) error
}
