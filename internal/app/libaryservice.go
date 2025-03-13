package app

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"libary-service/internal/domain"
	"libary-service/internal/repository"
	"libary-service/internal/validation"
	"net/http"
	"strings"
)

type LibaryService struct {
	repository repository.Repository
	validation validation.Validation
}

func NewLibaryService(repository repository.Repository, validation validation.Validation) *LibaryService {
	return &LibaryService{repository: repository, validation: validation}
}

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

func (s *LibaryService) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := s.repository.GetBooks()
	if err != nil {
		http.Error(w, "Error retrieving books", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (s *LibaryService) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/books/")
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := s.repository.GetBookByID(id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (s *LibaryService) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book domain.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := s.validation.CheckBook(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = uuid.New().String()

	createdBook, err := s.repository.CreateBook(book)
	if err != nil {
		http.Error(w, "Error creating book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdBook)
}

func (s *LibaryService) UpdateBook(w http.ResponseWriter, r *http.Request) {
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

	if err := s.validation.CheckBook(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	book.ID = id

	updatedBook, err := s.repository.UpdateBook(book)
	if err != nil {
		http.Error(w, "Error updating book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedBook)
}

func (s *LibaryService) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/books/")
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if err := s.repository.DeleteBook(id); err != nil {
		http.Error(w, "Error deleting book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *LibaryService) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.repository.GetUsers()
	if err != nil {
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (s *LibaryService) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/users/")
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := s.repository.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *LibaryService) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := s.validation.CheckUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()

	createdUser, err := s.repository.CreateUser(user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (s *LibaryService) UpdateUser(w http.ResponseWriter, r *http.Request) {
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

	if err := s.validation.CheckUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = id

	updatedUser, err := s.repository.UpdateUser(user)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

func (s *LibaryService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/users/")
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := s.repository.DeleteUser(id); err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *LibaryService) GetLendings(w http.ResponseWriter, r *http.Request) {
	lendings, err := s.repository.GetLendings()
	if err != nil {
		http.Error(w, "Error retrieving lendings", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lendings)
}

func (s *LibaryService) GetLendingByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/lendings/")
	if err != nil {
		http.Error(w, "Invalid lending ID", http.StatusBadRequest)
		return
	}

	lending, err := s.repository.GetLendingByID(id)
	if err != nil {
		http.Error(w, "Lending not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lending)
}

func (s *LibaryService) CreateLending(w http.ResponseWriter, r *http.Request) {
	var lending domain.Lending
	if err := json.NewDecoder(r.Body).Decode(&lending); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := s.validation.CheckLending(lending); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lending.ID = uuid.New().String()

	createdLending, err := s.repository.CreateLending(lending)
	if err != nil {
		http.Error(w, "Error creating lending", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdLending)
}

func (s *LibaryService) UpdateLending(w http.ResponseWriter, r *http.Request) {
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

	if err := s.validation.CheckLending(lending); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lending.ID = id

	updatedLending, err := s.repository.UpdateLending(lending)
	if err != nil {
		http.Error(w, "Error updating lending", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedLending)
}

func (s *LibaryService) DeleteLending(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r, "/lendings/")
	if err != nil {
		http.Error(w, "Invalid lending ID", http.StatusBadRequest)
		return
	}

	if err := s.repository.DeleteLending(id); err != nil {
		http.Error(w, "Error deleting lending", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
