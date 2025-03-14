package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"libary-service/generated/mocks"
	"libary-service/internal/domain"
)

func TestExtractID(t *testing.T) {
	id := uuid.NewString()
	testCases := []struct {
		name          string
		path          string
		basePath      string
		expectedID    string
		expectedError bool
	}{
		{"valid path", "/books/" + id, "/books/", id, false},
		{"valid path with trailing slash", "/books/" + id + "/", "/books/", id, false},
		{"invalid path", "/invalid/" + id, "/books/", "", true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", tc.path, nil)
			got, err := extractID(req, tc.basePath)
			if tc.expectedError {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, got)
			}
		})
	}
}

func TestGetBooks(t *testing.T) {
	id1 := uuid.NewString()
	id2 := uuid.NewString()
	testCases := []struct {
		name           string
		books          []domain.Book
		repositoryErr  error
		expectedStatus int
	}{
		{"success", []domain.Book{
			{ID: id1, Title: "The Fellowship of the Ring", Author: "J.R.R. Tolkien"},
			{ID: id2, Title: "The Two Towers", Author: "J.R.R. Tolkien"},
		}, nil, http.StatusOK},
		{"repository error", nil, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetBooks").Return(tc.books, tc.repositoryErr)
			mockValidation := new(mocks.Validation)
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("GET", "/books", nil)
			rr := httptest.NewRecorder()
			service.GetBooks(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusOK {
				var responseBooks []domain.Book
				err := json.Unmarshal(rr.Body.Bytes(), &responseBooks)
				assert.NoError(t, err)
				assert.Equal(t, tc.books, responseBooks)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetBookByID(t *testing.T) {
	bookID := uuid.NewString()
	book := domain.Book{ID: bookID, Title: "The Fellowship of the Ring", Author: "J.R.R. Tolkien"}
	testCases := []struct {
		name           string
		path           string
		book           domain.Book
		repositoryErr  error
		expectedStatus int
	}{
		{"success", "/books/" + bookID, book, nil, http.StatusOK},
		{"invalid id", "/invalid/" + bookID, domain.Book{}, nil, http.StatusBadRequest},
		{"book not found", "/books/" + bookID, domain.Book{}, errors.New("book not found"), http.StatusNotFound},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetBookByID", bookID).Return(tc.book, tc.repositoryErr).Maybe()
			mockValidation := new(mocks.Validation)
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("GET", tc.path, nil)
			rr := httptest.NewRecorder()
			service.GetBookByID(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusOK {
				var responseBook domain.Book
				err := json.Unmarshal(rr.Body.Bytes(), &responseBook)
				assert.NoError(t, err)
				assert.Equal(t, tc.book, responseBook)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateBook(t *testing.T) {
	validBook := domain.Book{Title: "The Fellowship of the Ring", Author: "J.R.R. Tolkien"}
	newID := uuid.NewString()
	testCases := []struct {
		name           string
		requestBody    interface{}
		validationErr  error
		createdBook    domain.Book
		repositoryErr  error
		expectedStatus int
	}{
		{"success", validBook, nil, domain.Book{ID: newID, Title: "The Fellowship of the Ring", Author: "J.R.R. Tolkien"}, nil, http.StatusCreated},
		{"invalid request body", "invalid json", nil, domain.Book{}, nil, http.StatusBadRequest},
		{"validation error", validBook, errors.New("validation error"), domain.Book{}, nil, http.StatusBadRequest},
		{"repository error", validBook, nil, domain.Book{}, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockValidation := new(mocks.Validation)
			var requestBytes []byte
			var err error
			if str, ok := tc.requestBody.(string); ok {
				requestBytes = []byte(str)
			} else {
				requestBytes, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}
			mockValidation.On("CheckBook", mock.AnythingOfType("domain.Book")).Return(tc.validationErr).Maybe()
			mockRepo.On("CreateBook", mock.AnythingOfType("domain.Book")).Return(tc.createdBook, tc.repositoryErr).Maybe()
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(requestBytes))
			rr := httptest.NewRecorder()
			service.CreateBook(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusCreated {
				var responseBook domain.Book
				err := json.Unmarshal(rr.Body.Bytes(), &responseBook)
				assert.NoError(t, err)
				assert.NotEmpty(t, responseBook.ID)
				assert.Equal(t, tc.createdBook.Title, responseBook.Title)
				assert.Equal(t, tc.createdBook.Author, responseBook.Author)
			}
			mockRepo.AssertExpectations(t)
			mockValidation.AssertExpectations(t)
		})
	}
}

func TestUpdateBook(t *testing.T) {
	validBook := domain.Book{Title: "The Two Towers", Author: "J.R.R. Tolkien"}
	bookID := uuid.NewString()
	testCases := []struct {
		name           string
		path           string
		requestBody    interface{}
		validationErr  error
		updatedBook    domain.Book
		repositoryErr  error
		expectedStatus int
	}{
		{"success", "/books/" + bookID, validBook, nil, domain.Book{ID: bookID, Title: "The Two Towers", Author: "J.R.R. Tolkien"}, nil, http.StatusOK},
		{"invalid path", "/invalid/" + bookID, validBook, nil, domain.Book{}, nil, http.StatusBadRequest},
		{"invalid request body", "/books/" + bookID, "invalid json", nil, domain.Book{}, nil, http.StatusBadRequest},
		{"validation error", "/books/" + bookID, validBook, errors.New("validation error"), domain.Book{}, nil, http.StatusBadRequest},
		{"repository error", "/books/" + bookID, validBook, nil, domain.Book{}, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockValidation := new(mocks.Validation)
			var requestBytes []byte
			var err error
			if str, ok := tc.requestBody.(string); ok {
				requestBytes = []byte(str)
			} else {
				requestBytes, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}
			mockValidation.On("CheckBook", mock.AnythingOfType("domain.Book")).Return(tc.validationErr).Maybe()
			mockRepo.On("UpdateBook", mock.AnythingOfType("domain.Book")).Return(tc.updatedBook, tc.repositoryErr).Maybe()
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("PUT", tc.path, bytes.NewBuffer(requestBytes))
			rr := httptest.NewRecorder()
			service.UpdateBook(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusOK {
				var responseBook domain.Book
				err := json.Unmarshal(rr.Body.Bytes(), &responseBook)
				assert.NoError(t, err)
				assert.Equal(t, tc.updatedBook, responseBook)
			}
			mockRepo.AssertExpectations(t)
			mockValidation.AssertExpectations(t)
		})
	}
}

func TestDeleteBook(t *testing.T) {
	bookID := uuid.NewString()
	testCases := []struct {
		name           string
		path           string
		repositoryErr  error
		expectedStatus int
	}{
		{"success", "/books/" + bookID, nil, http.StatusNoContent},
		{"invalid path", "/invalid/" + bookID, nil, http.StatusBadRequest},
		{"repository error", "/books/" + bookID, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockValidation := new(mocks.Validation)
			mockRepo.On("DeleteBook", bookID).Return(tc.repositoryErr).Maybe()
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("DELETE", tc.path, nil)
			rr := httptest.NewRecorder()
			service.DeleteBook(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUsers(t *testing.T) {
	userID1 := uuid.NewString()
	userID2 := uuid.NewString()
	testCases := []struct {
		name           string
		users          []domain.User
		repositoryErr  error
		expectedStatus int
	}{
		{"success", []domain.User{
			{ID: userID1, Name: "Max Mustermann", Email: "max@mustermann.de"},
			{ID: userID2, Name: "Erika Mustermann", Email: "erika@mustermann.de"},
		}, nil, http.StatusOK},
		{"repository error", nil, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetUsers").Return(tc.users, tc.repositoryErr)
			mockValidation := new(mocks.Validation)
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("GET", "/users", nil)
			rr := httptest.NewRecorder()
			service.GetUsers(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusOK {
				var responseUsers []domain.User
				err := json.Unmarshal(rr.Body.Bytes(), &responseUsers)
				assert.NoError(t, err)
				assert.Equal(t, tc.users, responseUsers)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	userID := uuid.NewString()
	user := domain.User{ID: userID, Name: "Max Mustermann", Email: "max@mustermann.de"}
	testCases := []struct {
		name           string
		path           string
		user           domain.User
		repositoryErr  error
		expectedStatus int
	}{
		{"success", "/users/" + userID, user, nil, http.StatusOK},
		{"invalid id", "/invalid/" + userID, domain.User{}, nil, http.StatusBadRequest},
		{"user not found", "/users/" + userID, domain.User{}, errors.New("user not found"), http.StatusNotFound},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetUserByID", userID).Return(tc.user, tc.repositoryErr).Maybe()
			mockValidation := new(mocks.Validation)
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("GET", tc.path, nil)
			rr := httptest.NewRecorder()
			service.GetUserByID(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusOK {
				var responseUser domain.User
				err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
				assert.NoError(t, err)
				assert.Equal(t, tc.user, responseUser)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateUser(t *testing.T) {
	validUser := domain.User{Name: "Max Mustermann", Email: "max@mustermann.de"}
	newUserID := uuid.NewString()
	testCases := []struct {
		name           string
		requestBody    interface{}
		validationErr  error
		createdUser    domain.User
		repositoryErr  error
		expectedStatus int
	}{
		{"success", validUser, nil, domain.User{ID: newUserID, Name: "Max Mustermann", Email: "max@mustermann.de"}, nil, http.StatusCreated},
		{"invalid request body", "invalid json", nil, domain.User{}, nil, http.StatusBadRequest},
		{"validation error", validUser, errors.New("validation error"), domain.User{}, nil, http.StatusBadRequest},
		{"repository error", validUser, nil, domain.User{}, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockValidation := new(mocks.Validation)
			var requestBytes []byte
			var err error
			if str, ok := tc.requestBody.(string); ok {
				requestBytes = []byte(str)
			} else {
				requestBytes, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}
			mockValidation.On("CheckUser", mock.AnythingOfType("domain.User")).Return(tc.validationErr).Maybe()
			mockRepo.On("CreateUser", mock.AnythingOfType("domain.User")).Return(tc.createdUser, tc.repositoryErr).Maybe()
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBytes))
			rr := httptest.NewRecorder()
			service.CreateUser(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusCreated {
				var responseUser domain.User
				err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
				assert.NoError(t, err)
				assert.NotEmpty(t, responseUser.ID)
				assert.Equal(t, tc.createdUser.Name, responseUser.Name)
				assert.Equal(t, tc.createdUser.Email, responseUser.Email)
			}
			mockRepo.AssertExpectations(t)
			mockValidation.AssertExpectations(t)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	validUser := domain.User{Name: "Erika Mustermann", Email: "erika@mustermann.de"}
	userID := uuid.NewString()
	testCases := []struct {
		name           string
		path           string
		requestBody    interface{}
		validationErr  error
		updatedUser    domain.User
		repositoryErr  error
		expectedStatus int
	}{
		{"success", "/users/" + userID, validUser, nil, domain.User{ID: userID, Name: "Erika Mustermann", Email: "erika@mustermann.de"}, nil, http.StatusOK},
		{"invalid path", "/invalid/" + userID, validUser, nil, domain.User{}, nil, http.StatusBadRequest},
		{"invalid request body", "/users/" + userID, "invalid json", nil, domain.User{}, nil, http.StatusBadRequest},
		{"validation error", "/users/" + userID, validUser, errors.New("validation error"), domain.User{}, nil, http.StatusBadRequest},
		{"repository error", "/users/" + userID, validUser, nil, domain.User{}, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockValidation := new(mocks.Validation)
			var requestBytes []byte
			var err error
			if str, ok := tc.requestBody.(string); ok {
				requestBytes = []byte(str)
			} else {
				requestBytes, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}
			mockValidation.On("CheckUser", mock.AnythingOfType("domain.User")).Return(tc.validationErr).Maybe()
			mockRepo.On("UpdateUser", mock.AnythingOfType("domain.User")).Return(tc.updatedUser, tc.repositoryErr).Maybe()
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("PUT", tc.path, bytes.NewBuffer(requestBytes))
			rr := httptest.NewRecorder()
			service.UpdateUser(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusOK {
				var responseUser domain.User
				err := json.Unmarshal(rr.Body.Bytes(), &responseUser)
				assert.NoError(t, err)
				assert.Equal(t, tc.updatedUser, responseUser)
			}
			mockRepo.AssertExpectations(t)
			mockValidation.AssertExpectations(t)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	userID := uuid.NewString()
	testCases := []struct {
		name           string
		path           string
		repositoryErr  error
		expectedStatus int
	}{
		{"success", "/users/" + userID, nil, http.StatusNoContent},
		{"invalid path", "/invalid/" + userID, nil, http.StatusBadRequest},
		{"repository error", "/users/" + userID, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockValidation := new(mocks.Validation)
			mockRepo.On("DeleteUser", userID).Return(tc.repositoryErr).Maybe()
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("DELETE", tc.path, nil)
			rr := httptest.NewRecorder()
			service.DeleteUser(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetLendings(t *testing.T) {
	lendDate1 := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	returnDate1 := lendDate1.Add(7 * 24 * time.Hour)
	lendDate2 := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	returnDate2 := lendDate2.Add(14 * 24 * time.Hour)
	id1 := uuid.NewString()
	id2 := uuid.NewString()
	bookID1 := uuid.NewString()
	bookID2 := uuid.NewString()
	userID1 := uuid.NewString()
	userID2 := uuid.NewString()
	testCases := []struct {
		name           string
		lendings       []domain.Lending
		repositoryErr  error
		expectedStatus int
	}{
		{"success", []domain.Lending{
			{ID: id1, BookID: bookID1, UserID: userID1, LendDate: lendDate1, ReturnDate: returnDate1},
			{ID: id2, BookID: bookID2, UserID: userID2, LendDate: lendDate2, ReturnDate: returnDate2},
		}, nil, http.StatusOK},
		{"repository error", nil, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetLendings").Return(tc.lendings, tc.repositoryErr)
			mockValidation := new(mocks.Validation)
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("GET", "/lendings", nil)
			rr := httptest.NewRecorder()
			service.GetLendings(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusOK {
				var responseLendings []domain.Lending
				err := json.Unmarshal(rr.Body.Bytes(), &responseLendings)
				assert.NoError(t, err)
				assert.Equal(t, len(tc.lendings), len(responseLendings))
				for i, expectedLending := range tc.lendings {
					actualLending := responseLendings[i]
					assert.Equal(t, expectedLending.ID, actualLending.ID)
					assert.Equal(t, expectedLending.BookID, actualLending.BookID)
					assert.Equal(t, expectedLending.UserID, actualLending.UserID)
					assert.Equal(t, expectedLending.LendDate, actualLending.LendDate)
					assert.Equal(t, expectedLending.ReturnDate, actualLending.ReturnDate)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetLendingByID(t *testing.T) {
	lendDate := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	returnDate := lendDate.Add(7 * 24 * time.Hour)
	lendingID := uuid.NewString()
	bookID := uuid.NewString()
	userID := uuid.NewString()
	lending := domain.Lending{ID: lendingID, BookID: bookID, UserID: userID, LendDate: lendDate, ReturnDate: returnDate}
	testCases := []struct {
		name           string
		path           string
		lending        domain.Lending
		repositoryErr  error
		expectedStatus int
	}{
		{"success", "/lendings/" + lendingID, lending, nil, http.StatusOK},
		{"invalid id", "/invalid/" + lendingID, domain.Lending{}, nil, http.StatusBadRequest},
		{"lending not found", "/lendings/" + lendingID, domain.Lending{}, errors.New("lending not found"), http.StatusNotFound},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetLendingByID", lendingID).Return(tc.lending, tc.repositoryErr).Maybe()
			mockValidation := new(mocks.Validation)
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("GET", tc.path, nil)
			rr := httptest.NewRecorder()
			service.GetLendingByID(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusOK {
				var responseLending domain.Lending
				err := json.Unmarshal(rr.Body.Bytes(), &responseLending)
				assert.NoError(t, err)
				assert.Equal(t, tc.lending.ID, responseLending.ID)
				assert.Equal(t, tc.lending.BookID, responseLending.BookID)
				assert.Equal(t, tc.lending.UserID, responseLending.UserID)
				assert.Equal(t, tc.lending.LendDate, responseLending.LendDate)
				assert.Equal(t, tc.lending.ReturnDate, responseLending.ReturnDate)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateLending(t *testing.T) {
	lendDate := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	returnDate := lendDate.Add(7 * 24 * time.Hour)
	bookID := uuid.NewString()
	userID := uuid.NewString()
	validLending := domain.Lending{BookID: bookID, UserID: userID, LendDate: lendDate, ReturnDate: returnDate}
	newLendingID := uuid.NewString()
	createdLending := domain.Lending{ID: newLendingID, BookID: bookID, UserID: userID, LendDate: lendDate, ReturnDate: returnDate}
	testCases := []struct {
		name           string
		requestBody    interface{}
		validationErr  error
		createdLending domain.Lending
		repositoryErr  error
		expectedStatus int
	}{
		{"success", validLending, nil, createdLending, nil, http.StatusCreated},
		{"invalid request body", "invalid json", nil, domain.Lending{}, nil, http.StatusBadRequest},
		{"validation error", validLending, errors.New("validation error"), domain.Lending{}, nil, http.StatusBadRequest},
		{"repository error", validLending, nil, domain.Lending{}, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockValidation := new(mocks.Validation)
			var requestBytes []byte
			var err error
			if str, ok := tc.requestBody.(string); ok {
				requestBytes = []byte(str)
			} else {
				requestBytes, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}
			mockValidation.On("CheckLending", mock.AnythingOfType("domain.Lending")).Return(tc.validationErr).Maybe()
			mockRepo.On("CreateLending", mock.AnythingOfType("domain.Lending")).Return(tc.createdLending, tc.repositoryErr).Maybe()
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("POST", "/lendings", bytes.NewBuffer(requestBytes))
			rr := httptest.NewRecorder()
			service.CreateLending(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusCreated {
				var responseLending domain.Lending
				err := json.Unmarshal(rr.Body.Bytes(), &responseLending)
				assert.NoError(t, err)
				assert.NotEmpty(t, responseLending.ID)
				assert.Equal(t, tc.createdLending.BookID, responseLending.BookID)
				assert.Equal(t, tc.createdLending.UserID, responseLending.UserID)
				assert.Equal(t, tc.createdLending.LendDate, responseLending.LendDate)
				assert.Equal(t, tc.createdLending.ReturnDate, responseLending.ReturnDate)
			}
			mockRepo.AssertExpectations(t)
			mockValidation.AssertExpectations(t)
		})
	}
}

func TestUpdateLending(t *testing.T) {
	lendDate := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	returnDate := lendDate.Add(14 * 24 * time.Hour)
	bookID := uuid.NewString()
	userID := uuid.NewString()
	validLending := domain.Lending{BookID: bookID, UserID: userID, LendDate: lendDate, ReturnDate: returnDate}
	lendingID := uuid.NewString()
	updatedLending := domain.Lending{ID: lendingID, BookID: bookID, UserID: userID, LendDate: lendDate, ReturnDate: returnDate}
	testCases := []struct {
		name           string
		path           string
		requestBody    interface{}
		validationErr  error
		updatedLending domain.Lending
		repositoryErr  error
		expectedStatus int
	}{
		{"success", "/lendings/" + lendingID, validLending, nil, updatedLending, nil, http.StatusOK},
		{"invalid path", "/invalid/" + lendingID, validLending, nil, domain.Lending{}, nil, http.StatusBadRequest},
		{"invalid request body", "/lendings/" + lendingID, "invalid json", nil, domain.Lending{}, nil, http.StatusBadRequest},
		{"validation error", "/lendings/" + lendingID, validLending, errors.New("validation error"), domain.Lending{}, nil, http.StatusBadRequest},
		{"repository error", "/lendings/" + lendingID, validLending, nil, domain.Lending{}, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockValidation := new(mocks.Validation)
			var requestBytes []byte
			var err error
			if str, ok := tc.requestBody.(string); ok {
				requestBytes = []byte(str)
			} else {
				requestBytes, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}
			mockValidation.On("CheckLending", mock.AnythingOfType("domain.Lending")).Return(tc.validationErr).Maybe()
			mockRepo.On("UpdateLending", mock.AnythingOfType("domain.Lending")).Return(tc.updatedLending, tc.repositoryErr).Maybe()
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("PUT", tc.path, bytes.NewBuffer(requestBytes))
			rr := httptest.NewRecorder()
			service.UpdateLending(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedStatus == http.StatusOK {
				var responseLending domain.Lending
				err := json.Unmarshal(rr.Body.Bytes(), &responseLending)
				assert.NoError(t, err)
				assert.Equal(t, tc.updatedLending.ID, responseLending.ID)
				assert.Equal(t, tc.updatedLending.BookID, responseLending.BookID)
				assert.Equal(t, tc.updatedLending.UserID, responseLending.UserID)
				assert.Equal(t, tc.updatedLending.LendDate, responseLending.LendDate)
				assert.Equal(t, tc.updatedLending.ReturnDate, responseLending.ReturnDate)
			}
			mockRepo.AssertExpectations(t)
			mockValidation.AssertExpectations(t)
		})
	}
}

func TestDeleteLending(t *testing.T) {
	lendingID := uuid.NewString()
	testCases := []struct {
		name           string
		path           string
		repositoryErr  error
		expectedStatus int
	}{
		{"success", "/lendings/" + lendingID, nil, http.StatusNoContent},
		{"invalid path", "/invalid/" + lendingID, nil, http.StatusBadRequest},
		{"repository error", "/lendings/" + lendingID, errors.New("database error"), http.StatusInternalServerError},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockValidation := new(mocks.Validation)
			mockRepo.On("DeleteLending", lendingID).Return(tc.repositoryErr).Maybe()
			service := NewLibaryService(mockRepo, mockValidation)
			req, _ := http.NewRequest("DELETE", tc.path, nil)
			rr := httptest.NewRecorder()
			service.DeleteLending(rr, req)
			assert.Equal(t, tc.expectedStatus, rr.Code)
			mockRepo.AssertExpectations(t)
		})
	}
}
