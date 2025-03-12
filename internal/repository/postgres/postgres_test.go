package postgresrepository

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"libary-service/internal/domain"
)

var repo *PostgresRepository

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_URL", "postgres://postgres:password@localhost:5432/libraryDB?sslmode=disable")

	repo = New()
	if err := repo.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	_, err := repo.db.Exec(context.Background(), "TRUNCATE lendings, books, users CASCADE")
	if err != nil {
		log.Fatalf("Failed to truncate tables: %v", err)
	}

	code := m.Run()
	os.Exit(code)
}

// resetDB clears all data from the tables before each test.
func resetDB(t *testing.T) {
	_, err := repo.db.Exec(context.Background(), "TRUNCATE lendings, books, users RESTART IDENTITY CASCADE;")
	if err != nil {
		t.Fatalf("Failed to reset DB: %v", err)
	}
}

func TestConnectAndDisconnect(t *testing.T) {
	// Connect is already tested in TestMain.
	// Test Disconnect (as implemented, it always returns an error)
	err := repo.Disconnect()
	if err == nil || err.Error() != "error closing database connection" {
		t.Errorf("Disconnect expected error 'error closing database connection', got: %v", err)
	}
}

func TestBookMethods(t *testing.T) {
	resetDB(t)

	// CreateBook
	book := domain.Book{
		ID:     uuid.NewString(),
		Title:  "Test Book",
		Author: "Test Author",
	}
	createdBook, err := repo.CreateBook(book)
	if err != nil {
		t.Fatalf("CreateBook failed: %v", err)
	}
	if createdBook != book {
		t.Errorf("CreateBook: got %+v, want %+v", createdBook, book)
	}

	// GetBooks should return at least the created book
	books, err := repo.GetBooks()
	if err != nil {
		t.Fatalf("GetBooks failed: %v", err)
	}
	if len(books) != 1 {
		t.Errorf("GetBooks: expected 1 book, got %d", len(books))
	}

	// GetBookByID (found)
	gotBook, err := repo.GetBookByID(book.ID)
	if err != nil {
		t.Fatalf("GetBookByID failed: %v", err)
	}
	if gotBook.Title != book.Title || gotBook.Author != book.Author {
		t.Errorf("GetBookByID: got %+v, want %+v", gotBook, book)
	}

	// GetBookByID (not found)
	_, err = repo.GetBookByID(uuid.NewString())
	if err == nil || err.Error() != "book not found" {
		t.Errorf("GetBookByID with unknown ID: expected 'book not found' error, got %v", err)
	}

	// UpdateBook (successful)
	updatedBook := domain.Book{
		ID:     book.ID,
		Title:  "Updated Title",
		Author: "Updated Author",
	}
	b, err := repo.UpdateBook(updatedBook)
	if err != nil {
		t.Fatalf("UpdateBook failed: %v", err)
	}
	if b.Title != "Updated Title" || b.Author != "Updated Author" {
		t.Errorf("UpdateBook: got %+v, want %+v", b, updatedBook)
	}

	// UpdateBook (not found)
	nonexistent := domain.Book{
		ID:     uuid.NewString(),
		Title:  "No Title",
		Author: "No Author",
	}
	_, err = repo.UpdateBook(nonexistent)
	if err == nil || err.Error() != "book not found" {
		t.Errorf("UpdateBook for non-existent book: expected 'book not found', got %v", err)
	}

	// DeleteBook (successful)
	err = repo.DeleteBook(book.ID)
	if err != nil {
		t.Fatalf("DeleteBook failed: %v", err)
	}
	// Ensure the book is gone
	_, err = repo.GetBookByID(book.ID)
	if err == nil || err.Error() != "book not found" {
		t.Errorf("After DeleteBook, expected 'book not found', got %v", err)
	}

	// DeleteBook (not found)
	err = repo.DeleteBook(book.ID)
	if err == nil || err.Error() != "book not found" {
		t.Errorf("DeleteBook for non-existent book: expected 'book not found', got %v", err)
	}
}

func TestUserMethods(t *testing.T) {
	resetDB(t)

	// CreateUser
	user := domain.User{
		ID:    uuid.NewString(),
		Name:  "John Doe",
		Email: "john@example.com",
	}
	createdUser, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if createdUser != user {
		t.Errorf("CreateUser: got %+v, want %+v", createdUser, user)
	}

	// GetUsers should return the created user
	users, err := repo.GetUsers()
	if err != nil {
		t.Fatalf("GetUsers failed: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("GetUsers: expected 1 user, got %d", len(users))
	}

	// GetUserByID (found)
	gotUser, err := repo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if gotUser.Name != user.Name || gotUser.Email != user.Email {
		t.Errorf("GetUserByID: got %+v, want %+v", gotUser, user)
	}

	// GetUserByID (not found)
	_, err = repo.GetUserByID(uuid.NewString())
	if err == nil || err.Error() != "user not found" {
		t.Errorf("GetUserByID with unknown ID: expected 'user not found', got %v", err)
	}

	// UpdateUser (successful)
	updatedUser := domain.User{
		ID:    user.ID,
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}
	u, err := repo.UpdateUser(updatedUser)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}
	if u.Name != "Jane Doe" || u.Email != "jane@example.com" {
		t.Errorf("UpdateUser: got %+v, want %+v", u, updatedUser)
	}

	// UpdateUser (not found)
	nonexistent := domain.User{
		ID:    uuid.NewString(),
		Name:  "Nobody",
		Email: "nobody@example.com",
	}
	_, err = repo.UpdateUser(nonexistent)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("UpdateUser for non-existent user: expected 'user not found', got %v", err)
	}

	// DeleteUser (successful)
	err = repo.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	// Ensure the user is gone
	_, err = repo.GetUserByID(user.ID)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("After DeleteUser, expected 'user not found', got %v", err)
	}

	// DeleteUser (not found)
	err = repo.DeleteUser(user.ID)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("DeleteUser for non-existent user: expected 'user not found', got %v", err)
	}
}

func TestLendingMethods(t *testing.T) {
	resetDB(t)

	// First create a book and a user for the lending record.
	book := domain.Book{
		ID:     uuid.NewString(),
		Title:  "Lending Book",
		Author: "Author",
	}
	if _, err := repo.CreateBook(book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	user := domain.User{
		ID:    uuid.NewString(),
		Name:  "Lending User",
		Email: "lending@example.com",
	}
	if _, err := repo.CreateUser(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// CreateLending with a zero ReturnDate.
	now := time.Now().UTC().Truncate(time.Second)
	lending := domain.Lending{
		ID:       uuid.NewString(),
		BookID:   book.ID,
		UserID:   user.ID,
		LendDate: now,
		// ReturnDate is zero (i.e. not set)
	}
	createdLending, err := repo.CreateLending(lending)
	if err != nil {
		t.Fatalf("CreateLending failed: %v", err)
	}
	if createdLending.ReturnDate.IsZero() == false {
		t.Errorf("CreateLending: expected zero ReturnDate, got %v", createdLending.ReturnDate)
	}

	// GetLendings should include our new lending.
	lendings, err := repo.GetLendings()
	if err != nil {
		t.Fatalf("GetLendings failed: %v", err)
	}
	if len(lendings) != 1 {
		t.Errorf("GetLendings: expected 1 lending, got %d", len(lendings))
	}

	// GetLendingByID (found)
	gotLending, err := repo.GetLendingByID(lending.ID)
	if err != nil {
		t.Fatalf("GetLendingByID failed: %v", err)
	}
	if gotLending.BookID != lending.BookID || gotLending.UserID != lending.UserID {
		t.Errorf("GetLendingByID: got %+v, want %+v", gotLending, lending)
	}

	// GetLendingByID (not found)
	_, err = repo.GetLendingByID(uuid.NewString())
	if err == nil || err.Error() != "lending not found" {
		t.Errorf("GetLendingByID with unknown ID: expected 'lending not found', got %v", err)
	}

	// UpdateLending (successful)
	// Now set a ReturnDate
	returnTime := now.Add(24 * time.Hour)
	lending.ReturnDate = returnTime
	updatedLending, err := repo.UpdateLending(lending)
	if err != nil {
		t.Fatalf("UpdateLending failed: %v", err)
	}
	if !updatedLending.ReturnDate.Equal(returnTime) {
		t.Errorf("UpdateLending: expected ReturnDate %v, got %v", returnTime, updatedLending.ReturnDate)
	}

	// UpdateLending (not found)
	nonexistentLending := domain.Lending{
		ID:       uuid.NewString(),
		BookID:   book.ID,
		UserID:   user.ID,
		LendDate: now,
	}
	_, err = repo.UpdateLending(nonexistentLending)
	if err == nil || err.Error() != "lending not found" {
		t.Errorf("UpdateLending for non-existent lending: expected 'lending not found', got %v", err)
	}

	// DeleteLending (successful)
	err = repo.DeleteLending(lending.ID)
	if err != nil {
		t.Fatalf("DeleteLending failed: %v", err)
	}
	// Ensure the lending is gone
	_, err = repo.GetLendingByID(lending.ID)
	if err == nil || err.Error() != "lending not found" {
		t.Errorf("After DeleteLending, expected 'lending not found', got %v", err)
	}

	// DeleteLending (not found)
	err = repo.DeleteLending(lending.ID)
	if err == nil || err.Error() != "lending not found" {
		t.Errorf("DeleteLending for non-existent lending: expected 'lending not found', got %v", err)
	}
}
