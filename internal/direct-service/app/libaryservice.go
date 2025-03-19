package app

import (
	"encoding/json"
	"errors"
	"libary-service/internal/direct-service/repository"
	"libary-service/internal/direct-service/validation"
	"libary-service/internal/domain"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// die Funktion "extractID" brauchst du ja nur, weil du zwar std konforme Handler Funktionen nutzt (und damit kein Zugriff auf die gin-Funktionalitaet hast),
// aber nicht den std ServeMux, sondern gin (ansonten koenntest du seit go 1.22 r.PathValue nehmen, um die id zu bekommen)
// -> Ich wuerde heute eher den std ServeMux nehmen, da er eben seit go 1.22 diese Moeglichkeiten hat (und du nutzt von gin ja nix?) 

// extractID parses an ID from the request URL path.
// The basePath should be formated like this: "/books/"
func extractID(r *http.Request, basePath string) (string, error) {
	r.PathValue(name string)
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
		// das stimmt hier nicht ganz, koennte ja auch ein anderer Fehler sein
		// http.Error(w, "Book not found", http.StatusNotFound)
		if errors.Is(err, repository.ErrBookNotFound) {
			http.Error(w, "Book not found", http.StatusNotFound)
			return
		} else {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
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

	if err := validation.CheckBook(book); err != nil {
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

	if err := validation.CheckBook(book); err != nil {
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

	if err := validation.CheckUser(user); err != nil {
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

	if err := validation.CheckUser(user); err != nil {
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

	if err := validation.CheckLending(lending); err != nil {
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

	if err := validation.CheckLending(lending); err != nil {
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
