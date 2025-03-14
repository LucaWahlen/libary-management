package postgresrepository

import (
	"context"
	"github.com/google/uuid"
	"libary-service/internal/domain"
	"log"
	"os"
	"testing"
	"time"
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

func resetDB(t *testing.T) {
	_, err := repo.db.Exec(context.Background(), "TRUNCATE lendings, books, users CASCADE;")
	if err != nil {
		t.Fatalf("Failed to reset DB: %v", err)
	}
}

func TestConnectAndDisconnect(t *testing.T) {
	r := New()
	if err := r.Connect(); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	err := r.Disconnect()
	if err == nil || err.Error() != "error closing database connection" {
		t.Errorf("Disconnect expected error 'error closing database connection', got: %v", err)
	}
}

func TestConnectFailure(t *testing.T) {
	oldURL := os.Getenv("DATABASE_URL")
	defer os.Setenv("DATABASE_URL", oldURL)
	os.Setenv("DATABASE_URL", "invalid://postgres:password@localhost:5432/libraryDB?sslmode=disable")
	r := New()
	err := r.Connect()
	if err == nil {
		t.Error("Expected error when connecting with an invalid URL, got nil")
	}
}

func TestBookMethods(t *testing.T) {
	resetDB(t)
	book := domain.Book{
		ID:     uuid.NewString(),
		Title:  "The Fellowship of the Ring",
		Author: "J.R.R. Tolkien",
	}
	createdBook, err := repo.CreateBook(book)
	if err != nil {
		t.Fatalf("CreateBook failed: %v", err)
	}
	if createdBook != book {
		t.Errorf("CreateBook: got %+v, want %+v", createdBook, book)
	}
	books, err := repo.GetBooks()
	if err != nil {
		t.Fatalf("GetBooks failed: %v", err)
	}
	if len(books) != 1 {
		t.Errorf("GetBooks: expected 1 book, got %d", len(books))
	}
	gotBook, err := repo.GetBookByID(book.ID)
	if err != nil {
		t.Fatalf("GetBookByID failed: %v", err)
	}
	if gotBook.Title != book.Title || gotBook.Author != book.Author {
		t.Errorf("GetBookByID: got %+v, want %+v", gotBook, book)
	}
	_, err = repo.GetBookByID(uuid.NewString())
	if err == nil || err.Error() != "book not found" {
		t.Errorf("GetBookByID with unknown ID: expected 'book not found' error, got %v", err)
	}
	updatedBook := domain.Book{
		ID:     book.ID,
		Title:  "The Two Towers",
		Author: "J.R.R. Tolkien",
	}
	b, err := repo.UpdateBook(updatedBook)
	if err != nil {
		t.Fatalf("UpdateBook failed: %v", err)
	}
	if b.Title != "The Two Towers" || b.Author != "J.R.R. Tolkien" {
		t.Errorf("UpdateBook: got %+v, want %+v", b, updatedBook)
	}
	nonexistent := domain.Book{
		ID:     uuid.NewString(),
		Title:  "The Return of the King",
		Author: "J.R.R. Tolkien",
	}
	_, err = repo.UpdateBook(nonexistent)
	if err == nil || err.Error() != "book not found" {
		t.Errorf("UpdateBook for non-existent book: expected 'book not found', got %v", err)
	}
	err = repo.DeleteBook(book.ID)
	if err != nil {
		t.Fatalf("DeleteBook failed: %v", err)
	}
	_, err = repo.GetBookByID(book.ID)
	if err == nil || err.Error() != "book not found" {
		t.Errorf("After DeleteBook, expected 'book not found', got %v", err)
	}
	err = repo.DeleteBook(book.ID)
	if err == nil || err.Error() != "book not found" {
		t.Errorf("DeleteBook for non-existent book: expected 'book not found', got %v", err)
	}
}

func TestUserMethods(t *testing.T) {
	resetDB(t)
	user := domain.User{
		ID:    uuid.NewString(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}
	createdUser, err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}
	if createdUser != user {
		t.Errorf("CreateUser: got %+v, want %+v", createdUser, user)
	}
	users, err := repo.GetUsers()
	if err != nil {
		t.Fatalf("GetUsers failed: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("GetUsers: expected 1 user, got %d", len(users))
	}
	gotUser, err := repo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}
	if gotUser.Name != user.Name || gotUser.Email != user.Email {
		t.Errorf("GetUserByID: got %+v, want %+v", gotUser, user)
	}
	_, err = repo.GetUserByID(uuid.NewString())
	if err == nil || err.Error() != "user not found" {
		t.Errorf("GetUserByID with unknown ID: expected 'user not found', got %v", err)
	}
	updatedUser := domain.User{
		ID:    user.ID,
		Name:  "Erika Mustermann",
		Email: "erika@mustermann.de",
	}
	u, err := repo.UpdateUser(updatedUser)
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}
	if u.Name != "Erika Mustermann" || u.Email != "erika@mustermann.de" {
		t.Errorf("UpdateUser: got %+v, want %+v", u, updatedUser)
	}
	nonexistent := domain.User{
		ID:    uuid.NewString(),
		Name:  "Erika Mustermann",
		Email: "erika@mustermann.de",
	}
	_, err = repo.UpdateUser(nonexistent)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("UpdateUser for non-existent user: expected 'user not found', got %v", err)
	}
	err = repo.DeleteUser(user.ID)
	if err != nil {
		t.Fatalf("DeleteUser failed: %v", err)
	}
	_, err = repo.GetUserByID(user.ID)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("After DeleteUser, expected 'user not found', got %v", err)
	}
	err = repo.DeleteUser(user.ID)
	if err == nil || err.Error() != "user not found" {
		t.Errorf("DeleteUser for non-existent user: expected 'user not found', got %v", err)
	}
}

func TestLendingMethods(t *testing.T) {
	resetDB(t)
	book := domain.Book{
		ID:     uuid.NewString(),
		Title:  "The Fellowship of the Ring",
		Author: "J.R.R. Tolkien",
	}
	if _, err := repo.CreateBook(book); err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	user := domain.User{
		ID:    uuid.NewString(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}
	if _, err := repo.CreateUser(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	lendDate := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	lending := domain.Lending{
		ID:       uuid.NewString(),
		BookID:   book.ID,
		UserID:   user.ID,
		LendDate: lendDate,
	}
	createdLending, err := repo.CreateLending(lending)
	if err != nil {
		t.Fatalf("CreateLending without ReturnDate failed: %v", err)
	}
	if !createdLending.ReturnDate.IsZero() {
		t.Errorf("CreateLending: expected zero ReturnDate, got %v", createdLending.ReturnDate)
	}
	lendings, err := repo.GetLendings()
	if err != nil {
		t.Fatalf("GetLendings failed: %v", err)
	}
	if len(lendings) != 1 {
		t.Errorf("GetLendings: expected 1 lending, got %d", len(lendings))
	}
	gotLending, err := repo.GetLendingByID(lending.ID)
	if err != nil {
		t.Fatalf("GetLendingByID failed: %v", err)
	}
	if gotLending.BookID != lending.BookID || gotLending.UserID != lending.UserID {
		t.Errorf("GetLendingByID: got %+v, want %+v", gotLending, lending)
	}
	_, err = repo.GetLendingByID(uuid.NewString())
	if err == nil || err.Error() != "lending not found" {
		t.Errorf("GetLendingByID with unknown ID: expected 'lending not found', got %v", err)
	}
	returnTime := lendDate.Add(24 * time.Hour)
	lending.ReturnDate = returnTime
	updatedLending, err := repo.UpdateLending(lending)
	if err != nil {
		t.Fatalf("UpdateLending failed: %v", err)
	}
	if !updatedLending.ReturnDate.Equal(returnTime) {
		t.Errorf("UpdateLending: expected ReturnDate %v, got %v", returnTime, updatedLending.ReturnDate)
	}
	nonexistentLending := domain.Lending{
		ID:       uuid.NewString(),
		BookID:   book.ID,
		UserID:   user.ID,
		LendDate: lendDate,
	}
	_, err = repo.UpdateLending(nonexistentLending)
	if err == nil || err.Error() != "lending not found" {
		t.Errorf("UpdateLending for non-existent lending: expected 'lending not found', got %v", err)
	}
	err = repo.DeleteLending(lending.ID)
	if err != nil {
		t.Fatalf("DeleteLending failed: %v", err)
	}
	_, err = repo.GetLendingByID(lending.ID)
	if err == nil || err.Error() != "lending not found" {
		t.Errorf("After DeleteLending, expected 'lending not found', got %v", err)
	}
	err = repo.DeleteLending(lending.ID)
	if err == nil || err.Error() != "lending not found" {
		t.Errorf("DeleteLending for non-existent lending: expected 'lending not found', got %v", err)
	}
	lendDate = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	returnTime = lendDate.Add(48 * time.Hour)
	newLending := domain.Lending{
		ID:         uuid.NewString(),
		BookID:     book.ID,
		UserID:     user.ID,
		LendDate:   lendDate,
		ReturnDate: returnTime,
	}
	created, err := repo.CreateLending(newLending)
	if err != nil {
		t.Fatalf("CreateLending with ReturnDate failed: %v", err)
	}
	if !created.ReturnDate.Equal(returnTime) {
		t.Errorf("Expected ReturnDate %v, got %v", returnTime, created.ReturnDate)
	}
	lendings, err = repo.GetLendings()
	if err != nil {
		t.Fatalf("GetLendings failed: %v", err)
	}
	if len(lendings) != 1 {
		t.Errorf("GetLendings: expected 1 lending, got %d", len(lendings))
	}
	gotLending, err = repo.GetLendingByID(newLending.ID)
	if err != nil {
		t.Fatalf("GetLendingByID failed: %v", err)
	}
	if gotLending.ID != created.ID || gotLending.BookID != created.BookID || gotLending.UserID != created.UserID || !gotLending.LendDate.Equal(created.LendDate) || !gotLending.ReturnDate.Equal(created.ReturnDate) {
		t.Errorf("GetLendingByID: expected %+v, got %+v", created, gotLending)
	}
}

func TestMethodsAfterDisconnect(t *testing.T) {
	r := New()
	if err := r.Connect(); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	r.Disconnect()
	if _, err := r.GetBooks(); err == nil {
		t.Error("Expected error from GetBooks on disconnected connection")
	}
	if _, err := r.GetBookByID(uuid.NewString()); err == nil {
		t.Error("Expected error from GetBookByID on disconnected connection")
	}
	if _, err := r.CreateBook(domain.Book{ID: uuid.NewString(), Title: "X", Author: "Y"}); err == nil {
		t.Error("Expected error from CreateBook on disconnected connection")
	}
	if _, err := r.UpdateBook(domain.Book{ID: uuid.NewString(), Title: "X", Author: "Y"}); err == nil {
		t.Error("Expected error from UpdateBook on disconnected connection")
	}
	if err := r.DeleteBook(uuid.NewString()); err == nil {
		t.Error("Expected error from DeleteBook on disconnected connection")
	}
	if _, err := r.GetUsers(); err == nil {
		t.Error("Expected error from GetUsers on disconnected connection")
	}
	if _, err := r.GetUserByID(uuid.NewString()); err == nil {
		t.Error("Expected error from GetUserByID on disconnected connection")
	}
	if _, err := r.CreateUser(domain.User{ID: uuid.NewString(), Name: "Test", Email: "test@example.com"}); err == nil {
		t.Error("Expected error from CreateUser on disconnected connection")
	}
	if _, err := r.UpdateUser(domain.User{ID: uuid.NewString(), Name: "Test", Email: "test@example.com"}); err == nil {
		t.Error("Expected error from UpdateUser on disconnected connection")
	}
	if err := r.DeleteUser(uuid.NewString()); err == nil {
		t.Error("Expected error from DeleteUser on disconnected connection")
	}
	if _, err := r.GetLendings(); err == nil {
		t.Error("Expected error from GetLendings on disconnected connection")
	}
	if _, err := r.GetLendingByID(uuid.NewString()); err == nil {
		t.Error("Expected error from GetLendingByID on disconnected connection")
	}
	dummyLending := domain.Lending{
		ID:         uuid.NewString(),
		BookID:     uuid.NewString(),
		UserID:     uuid.NewString(),
		LendDate:   time.Now(),
		ReturnDate: time.Now(),
	}
	if _, err := r.CreateLending(dummyLending); err == nil {
		t.Error("Expected error from CreateLending on disconnected connection")
	}
	if _, err := r.UpdateLending(dummyLending); err == nil {
		t.Error("Expected error from UpdateLending on disconnected connection")
	}
	if err := r.DeleteLending(dummyLending.ID); err == nil {
		t.Error("Expected error from DeleteLending on disconnected connection")
	}
}
