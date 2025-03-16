package inmemoryrepository

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"libary-service/internal/domain"
	"testing"
	"time"
)

func TestConnect(t *testing.T) {
	assert.Nil(t, New().Connect())
}

func TestDisconnect(t *testing.T) {
	assert.Nil(t, New().Disconnect())
}

func TestCreateBook(t *testing.T) {
	repo := New()
	book := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Fellowship of the Ring",
		Author: "J. R. R. Tolkien",
	}

	result, err := repo.CreateBook(book)
	assert.NoError(t, err)
	assert.Equal(t, book, result)
	storedBook, err := repo.GetBookByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, book, storedBook)
}

func TestGetBooks(t *testing.T) {
	repo := New()
	book1 := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Two Towers",
		Author: "J. R. R. Tolkien",
	}
	book2 := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Return of the King",
		Author: "J. R. R. Tolkien",
	}

	_, err := repo.CreateBook(book1)
	assert.NoError(t, err)
	_, err = repo.CreateBook(book2)
	assert.NoError(t, err)

	books, err := repo.GetBooks()
	assert.NoError(t, err)
	assert.Len(t, books, 2)
	assert.Contains(t, books, book1)
	assert.Contains(t, books, book2)
}

func TestGetBookByID(t *testing.T) {
	repo := New()
	book := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Fellowship of the Ring",
		Author: "J. R. R. Tolkien",
	}

	_, err := repo.CreateBook(book)
	assert.NoError(t, err)

	result, err := repo.GetBookByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, book, result)

	_, err = repo.GetBookByID(uuid.New().String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "book not found")
}

func TestUpdateBook(t *testing.T) {
	repo := New()
	book := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Two Towers",
		Author: "J. R. R. Tolkien",
	}

	_, err := repo.CreateBook(book)
	assert.NoError(t, err)

	updated := book
	updated.Title = "The Two Towers (Revised Edition)"
	result, err := repo.UpdateBook(updated)
	assert.NoError(t, err)
	assert.Equal(t, updated, result)

	stored, err := repo.GetBookByID(book.ID)
	assert.NoError(t, err)
	assert.Equal(t, "The Two Towers (Revised Edition)", stored.Title)

	nonExistent := domain.Book{ID: uuid.New().String()}
	_, err = repo.UpdateBook(nonExistent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "book not found")
}

func TestDeleteBook(t *testing.T) {
	repo := New()
	book := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Return of the King",
		Author: "J. R. R. Tolkien",
	}

	_, err := repo.CreateBook(book)
	assert.NoError(t, err)

	err = repo.DeleteBook(book.ID)
	assert.NoError(t, err)

	_, err = repo.GetBookByID(book.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "book not found")

	err = repo.DeleteBook(uuid.New().String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "book not found")
}

func TestCreateUser(t *testing.T) {
	repo := New()
	user := domain.User{
		ID:    uuid.New().String(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}

	result, err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
	storedUser, err := repo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, storedUser)
}

func TestGetUsers(t *testing.T) {
	repo := New()
	user1 := domain.User{
		ID:    uuid.New().String(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}
	user2 := domain.User{
		ID:    uuid.New().String(),
		Name:  "Erika Mustermann",
		Email: "erika@mustermann.de",
	}

	_, err := repo.CreateUser(user1)
	assert.NoError(t, err)
	_, err = repo.CreateUser(user2)
	assert.NoError(t, err)

	users, err := repo.GetUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Contains(t, users, user1)
	assert.Contains(t, users, user2)
}

func TestGetUserByID(t *testing.T) {
	repo := New()
	user := domain.User{
		ID:    uuid.New().String(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}

	_, err := repo.CreateUser(user)
	assert.NoError(t, err)

	result, err := repo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user, result)

	_, err = repo.GetUserByID(uuid.New().String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUpdateUser(t *testing.T) {
	repo := New()
	user := domain.User{
		ID:    uuid.New().String(),
		Name:  "Erika Mustermann",
		Email: "erika@mustermann.de",
	}

	_, err := repo.CreateUser(user)
	assert.NoError(t, err)

	updated := user
	updated.Name = "Max Mustermann"
	updated.Email = "max@mustermann.de"
	result, err := repo.UpdateUser(updated)
	assert.NoError(t, err)
	assert.Equal(t, updated, result)

	stored, err := repo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "max@mustermann.de", stored.Email)
	assert.Equal(t, "Max Mustermann", stored.Name)

	nonExistent := domain.User{ID: uuid.New().String()}
	_, err = repo.UpdateUser(nonExistent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestDeleteUser(t *testing.T) {
	repo := New()
	user := domain.User{
		ID:    uuid.New().String(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}

	_, err := repo.CreateUser(user)
	assert.NoError(t, err)

	err = repo.DeleteUser(user.ID)
	assert.NoError(t, err)

	_, err = repo.GetUserByID(user.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	err = repo.DeleteUser(uuid.New().String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestCreateLending(t *testing.T) {
	repo := New()
	book := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Fellowship of the Ring",
		Author: "J. R. R. Tolkien",
	}
	user := domain.User{
		ID:    uuid.New().String(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}

	_, err := repo.CreateBook(book)
	assert.NoError(t, err)
	_, err = repo.CreateUser(user)
	assert.NoError(t, err)

	lending := domain.Lending{
		ID:         uuid.New().String(),
		BookID:     book.ID,
		UserID:     user.ID,
		LendDate:   time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 14),
	}

	result, err := repo.CreateLending(lending)
	assert.NoError(t, err)
	assert.Equal(t, lending, result)
	storedLending, err := repo.GetLendingByID(lending.ID)
	assert.NoError(t, err)
	assert.Equal(t, lending, storedLending)
}

func TestGetLendings(t *testing.T) {
	repo := New()
	book1 := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Two Towers",
		Author: "J. R. R. Tolkien",
	}
	book2 := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Return of the King",
		Author: "J. R. R. Tolkien",
	}
	user1 := domain.User{
		ID:    uuid.New().String(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}
	user2 := domain.User{
		ID:    uuid.New().String(),
		Name:  "Erika Mustermann",
		Email: "erika@mustermann.de",
	}

	_, err := repo.CreateBook(book1)
	assert.NoError(t, err)
	_, err = repo.CreateBook(book2)
	assert.NoError(t, err)
	_, err = repo.CreateUser(user1)
	assert.NoError(t, err)
	_, err = repo.CreateUser(user2)
	assert.NoError(t, err)

	lending1 := domain.Lending{
		ID:         uuid.New().String(),
		BookID:     book1.ID,
		UserID:     user1.ID,
		LendDate:   time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 14),
	}
	lending2 := domain.Lending{
		ID:         uuid.New().String(),
		BookID:     book2.ID,
		UserID:     user2.ID,
		LendDate:   time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 7),
	}

	_, err = repo.CreateLending(lending1)
	assert.NoError(t, err)
	_, err = repo.CreateLending(lending2)
	assert.NoError(t, err)

	lendings, err := repo.GetLendings()
	assert.NoError(t, err)
	assert.Len(t, lendings, 2)
	assert.Contains(t, lendings, lending1)
	assert.Contains(t, lendings, lending2)
}

func TestGetLendingByID(t *testing.T) {
	repo := New()
	book := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Fellowship of the Ring",
		Author: "J. R. R. Tolkien",
	}
	user := domain.User{
		ID:    uuid.New().String(),
		Name:  "Erika Mustermann",
		Email: "erika@mustermann.de",
	}

	_, err := repo.CreateBook(book)
	assert.NoError(t, err)
	_, err = repo.CreateUser(user)
	assert.NoError(t, err)

	lending := domain.Lending{
		ID:         uuid.New().String(),
		BookID:     book.ID,
		UserID:     user.ID,
		LendDate:   time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 14),
	}

	_, err = repo.CreateLending(lending)
	assert.NoError(t, err)

	result, err := repo.GetLendingByID(lending.ID)
	assert.NoError(t, err)
	assert.Equal(t, lending, result)

	_, err = repo.GetLendingByID(uuid.New().String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lending not found")
}

func TestUpdateLending(t *testing.T) {
	repo := New()
	book := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Two Towers",
		Author: "J. R. R. Tolkien",
	}
	user := domain.User{
		ID:    uuid.New().String(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}

	_, err := repo.CreateBook(book)
	assert.NoError(t, err)
	_, err = repo.CreateUser(user)
	assert.NoError(t, err)

	lending := domain.Lending{
		ID:         uuid.New().String(),
		BookID:     book.ID,
		UserID:     user.ID,
		LendDate:   time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 14),
	}

	_, err = repo.CreateLending(lending)
	assert.NoError(t, err)

	updated := lending
	updated.ReturnDate = time.Now().AddDate(0, 0, 21)
	result, err := repo.UpdateLending(updated)
	assert.NoError(t, err)
	assert.Equal(t, updated, result)

	stored, err := repo.GetLendingByID(lending.ID)
	assert.NoError(t, err)
	assert.Equal(t, updated.ReturnDate.Day(), stored.ReturnDate.Day())

	nonExistent := domain.Lending{ID: uuid.New().String()}
	_, err = repo.UpdateLending(nonExistent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lending not found")
}

func TestDeleteLending(t *testing.T) {
	repo := New()
	book := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Return of the King",
		Author: "J. R. R. Tolkien",
	}
	user := domain.User{
		ID:    uuid.New().String(),
		Name:  "Erika Mustermann",
		Email: "erika@mustermann.de",
	}

	_, err := repo.CreateBook(book)
	assert.NoError(t, err)
	_, err = repo.CreateUser(user)
	assert.NoError(t, err)

	lending := domain.Lending{
		ID:         uuid.New().String(),
		BookID:     book.ID,
		UserID:     user.ID,
		LendDate:   time.Now(),
		ReturnDate: time.Now().AddDate(0, 0, 14),
	}

	_, err = repo.CreateLending(lending)
	assert.NoError(t, err)

	err = repo.DeleteLending(lending.ID)
	assert.NoError(t, err)

	_, err = repo.GetLendingByID(lending.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lending not found")

	err = repo.DeleteLending(uuid.New().String())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lending not found")
}
