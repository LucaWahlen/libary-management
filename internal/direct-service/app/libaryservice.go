package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"libary-service/internal/direct-service/repository"
	"libary-service/internal/domain"
	"net/http"
	"strings"
)

// extractID parses an ID from the request URL path.
// The basePath should be formated like this: "/books/"
func extractID(r *http.Request, basePath string) (string, error) {
	path := r.URL.Path
	if !strings.HasPrefix(path, basePath) {
		return "", errors.New("invalid path")
	}
	idStr := strings.TrimPrefix(path, basePath)
	idStr = strings.Trim(idStr, "/")
	return idStr, nil
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repository.GetBooks()
	if err != nil {
		http.Error(w, "Error retrieving books", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/books/")
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := repository.GetBookByID(id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book domain.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var errs []error
	if book.Title == "" {
		errs = append(errs, fmt.Errorf("title is required"))
	}
	if book.Author == "" {
		errs = append(errs, fmt.Errorf("author is required"))
	}
	if book.ID != "" {
		errs = append(errs, fmt.Errorf("id should be empty"))
	}
	if err := errors.Join(errs...); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = uuid.New().String()

	createdBook, err := repository.CreateBook(book)
	if err != nil {
		http.Error(w, "Error creating book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBook)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/books/")
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book domain.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var errs []error
	if book.Title == "" {
		errs = append(errs, fmt.Errorf("title is required"))
	}
	if book.Author == "" {
		errs = append(errs, fmt.Errorf("author is required"))
	}
	if book.ID != "" {
		errs = append(errs, fmt.Errorf("id should be empty"))
	}
	if err := errors.Join(errs...); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = id

	updatedBook, err := repository.UpdateBook(book)
	if err != nil {
		http.Error(w, "Error updating book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/books/")
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteBook(id); err != nil {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.GetUsers()
	if err != nil {
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/users/")
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := repository.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var errs []error
	if user.Name == "" {
		errs = append(errs, fmt.Errorf("name is required"))
	}
	if user.Email == "" {
		errs = append(errs, fmt.Errorf("email is required"))
	}
	if user.ID != "" {
		errs = append(errs, fmt.Errorf("id should be empty"))
	}
	if err := errors.Join(errs...); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()

	createdUser, err := repository.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/users/")
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var errs []error
	if user.Name == "" {
		errs = append(errs, fmt.Errorf("name is required"))
	}
	if user.Email == "" {
		errs = append(errs, fmt.Errorf("email is required"))
	}
	if user.ID != "" {
		errs = append(errs, fmt.Errorf("id should be empty"))
	}
	if err := errors.Join(errs...); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = id

	updatedUser, err := repository.UpdateUser(user)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/users/")
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteUser(id); err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetLendings(w http.ResponseWriter, r *http.Request) {
	lendings, err := repository.GetLendings()
	if err != nil {
		http.Error(w, "Error retrieving lendings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lendings)
}

func GetLendingByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/lendings/")
	if err != nil {
		http.Error(w, "Invalid lending ID", http.StatusBadRequest)
		return
	}

	lending, err := repository.GetLendingByID(id)
	if err != nil {
		http.Error(w, "Lending not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lending)
}

func CreateLending(w http.ResponseWriter, r *http.Request) {
	var lending domain.Lending
	if err := json.NewDecoder(r.Body).Decode(&lending); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var errs []error
	if lending.ID != "" {
		errs = append(errs, fmt.Errorf("id should be empty"))
	}

	_, bookMissing := repository.GetBookByID(lending.BookID)
	if bookMissing != nil {
		errs = append(errs, fmt.Errorf("book not found"))
	}

	_, userMissing := repository.GetUserByID(lending.UserID)
	if userMissing != nil {
		errs = append(errs, fmt.Errorf("user not found"))
	}

	if lending.LendDate.IsZero() {
		errs = append(errs, fmt.Errorf("lend_date is required"))
	}

	if !lending.ReturnDate.IsZero() {
		if !lending.LendDate.Before(lending.ReturnDate) {
			errs = append(errs, fmt.Errorf("lend_date is less than return_date"))
		}
	}
	if err := errors.Join(errs...); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lending.ID = uuid.New().String()

	createdLending, err := repository.CreateLending(lending)
	if err != nil {
		http.Error(w, "Error creating lending", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdLending)
}

func UpdateLending(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/lendings/")
	if err != nil {
		http.Error(w, "Invalid lending ID", http.StatusBadRequest)
		return
	}

	var lending domain.Lending
	if err := json.NewDecoder(r.Body).Decode(&lending); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var errs []error
	if lending.ID != "" {
		errs = append(errs, fmt.Errorf("id should be empty"))
	}

	_, bookMissing := repository.GetBookByID(lending.BookID)
	if bookMissing != nil {
		errs = append(errs, fmt.Errorf("book not found"))
	}

	_, userMissing := repository.GetUserByID(lending.UserID)
	if userMissing != nil {
		errs = append(errs, fmt.Errorf("user not found"))
	}

	if lending.LendDate.IsZero() {
		errs = append(errs, fmt.Errorf("lend_date is required"))
	}

	if !lending.ReturnDate.IsZero() {
		if !lending.LendDate.Before(lending.ReturnDate) {
			errs = append(errs, fmt.Errorf("lend_date is less than return_date"))
		}
	}
	if err := errors.Join(errs...); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lending.ID = id

	updatedLending, err := repository.UpdateLending(lending)
	if err != nil {
		http.Error(w, "Error updating lending", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedLending)
}

func DeleteLending(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/lendings/")
	if err != nil {
		http.Error(w, "Invalid lending ID", http.StatusBadRequest)
		return
	}

	if err := repository.DeleteLending(id); err != nil {
		http.Error(w, "Error deleting lending", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
