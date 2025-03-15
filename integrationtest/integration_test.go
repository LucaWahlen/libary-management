package integration_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"libary-service/internal/domain"
)

const (
	directServiceURL   = "http://localhost:8081"
	injectedServiceURL = "http://localhost:8080"
)

func setupTestDatabase(t *testing.T) {
	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/libraryDB"
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	require.NoError(t, err, "Failed to connect to database")
	defer conn.Close(context.Background())

	// Clean database tables
	_, err = conn.Exec(context.Background(), "DELETE FROM lendings")
	require.NoError(t, err, "Failed to clean lendings table")

	_, err = conn.Exec(context.Background(), "DELETE FROM users")
	require.NoError(t, err, "Failed to clean users table")

	_, err = conn.Exec(context.Background(), "DELETE FROM books")
	require.NoError(t, err, "Failed to clean books table")
}

func makeRequest(t *testing.T, method, url string, body interface{}) *http.Response {
	var req *http.Request
	var err error

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		require.NoError(t, err, "Failed to marshal request body")
		req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to execute request")

	return resp
}

func decodeResponse(t *testing.T, resp *http.Response, target interface{}) {
	defer resp.Body.Close()
	err := json.NewDecoder(resp.Body).Decode(target)
	require.NoError(t, err, "Failed to decode response")
}

func TestServices(t *testing.T) {

	serviceURLs := map[string]string{
		"Direct":   directServiceURL,
		"Injected": injectedServiceURL,
	}

	for name, baseURL := range serviceURLs {
		t.Run(name, func(t *testing.T) {
			t.Run(fmt.Sprintf("Testing %s-Service", name), func(t *testing.T) { serviceTest(t, baseURL, name) })
		})
	}
}

func serviceTest(t *testing.T, baseURL, serviceName string) {
	setupTestDatabase(t)

	book := domain.Book{
		Title:  "Test Book Direct",
		Author: "Test Author Direct",
	}

	resp := makeRequest(t, http.MethodPost, baseURL+"/books", book)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdBook domain.Book
	decodeResponse(t, resp, &createdBook)

	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, book.Title, createdBook.Title)
	assert.Equal(t, book.Author, createdBook.Author)

	// Get book by ID
	resp = makeRequest(t, http.MethodGet, fmt.Sprintf("%s/books/%s", baseURL, createdBook.ID), nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedBook domain.Book
	decodeResponse(t, resp, &retrievedBook)
	assert.Equal(t, createdBook.ID, retrievedBook.ID)

	// Update book
	retrievedBook.Title = "Updated Book Direct"
	resp = makeRequest(t, http.MethodPut, fmt.Sprintf("%s/books/%s", baseURL, retrievedBook.ID), retrievedBook)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedBook domain.Book
	decodeResponse(t, resp, &updatedBook)
	assert.Equal(t, "Updated Book Direct", updatedBook.Title)

	// Get all books
	resp = makeRequest(t, http.MethodGet, baseURL+"/books", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var books []domain.Book
	decodeResponse(t, resp, &books)
	assert.GreaterOrEqual(t, len(books), 1)

	// Delete book
	resp = makeRequest(t, http.MethodDelete, fmt.Sprintf("%s/books/%s", baseURL, createdBook.ID), nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify deletion
	resp = makeRequest(t, http.MethodGet, fmt.Sprintf("%s/books/%s", baseURL, createdBook.ID), nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Create user
	user := domain.User{
		Name:  "Test User Direct",
		Email: "test.direct@example.com",
	}

	resp = makeRequest(t, http.MethodPost, baseURL+"/users", user)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdUser domain.User
	decodeResponse(t, resp, &createdUser)

	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, user.Name, createdUser.Name)
	assert.Equal(t, user.Email, createdUser.Email)

	// Get user by ID
	resp = makeRequest(t, http.MethodGet, fmt.Sprintf("%s/users/%s", baseURL, createdUser.ID), nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedUser domain.User
	decodeResponse(t, resp, &retrievedUser)
	assert.Equal(t, createdUser.ID, retrievedUser.ID)

	// Update user
	retrievedUser.Name = "Updated User Direct"
	resp = makeRequest(t, http.MethodPut, fmt.Sprintf("%s/users/%s", baseURL, retrievedUser.ID), retrievedUser)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedUser domain.User
	decodeResponse(t, resp, &updatedUser)
	assert.Equal(t, "Updated User Direct", updatedUser.Name)

	// Get all users
	resp = makeRequest(t, http.MethodGet, baseURL+"/users", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var users []domain.User
	decodeResponse(t, resp, &users)
	assert.GreaterOrEqual(t, len(users), 1)

	// Note: We don't delete the user here because we'll use it for the lending test

	book = domain.Book{
		Title:  "Lending Test Book Direct",
		Author: "Lending Test Author Direct",
	}

	resp = makeRequest(t, http.MethodPost, baseURL+"/books", book)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	decodeResponse(t, resp, &createdBook)

	user = domain.User{
		Name:  "Lending Test User Direct",
		Email: "lending.test.direct@example.com",
	}

	resp = makeRequest(t, http.MethodPost, baseURL+"/users", user)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	decodeResponse(t, resp, &createdUser)

	// Create lending
	lending := domain.Lending{
		BookID:   createdBook.ID,
		UserID:   createdUser.ID,
		LendDate: time.Now(),
	}

	resp = makeRequest(t, http.MethodPost, baseURL+"/lendings", lending)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdLending domain.Lending
	decodeResponse(t, resp, &createdLending)

	assert.NotEmpty(t, createdLending.ID)
	assert.Equal(t, lending.BookID, createdLending.BookID)
	assert.Equal(t, lending.UserID, createdLending.UserID)

	// Get lending by ID
	resp = makeRequest(t, http.MethodGet, fmt.Sprintf("%s/lendings/%s", baseURL, createdLending.ID), nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedLending domain.Lending
	decodeResponse(t, resp, &retrievedLending)
	assert.Equal(t, createdLending.ID, retrievedLending.ID)

	// Update lending (mark as returned)
	returnTime := time.Now().Add(7 * 24 * time.Hour) // Return after a week
	retrievedLending.ReturnDate = returnTime
	resp = makeRequest(t, http.MethodPut, fmt.Sprintf("%s/lendings/%s", baseURL, retrievedLending.ID), retrievedLending)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedLending domain.Lending
	decodeResponse(t, resp, &updatedLending)
	assert.NotZero(t, updatedLending.ReturnDate)

	// Get all lendings
	resp = makeRequest(t, http.MethodGet, baseURL+"/lendings", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var lendings []domain.Lending
	decodeResponse(t, resp, &lendings)
	assert.GreaterOrEqual(t, len(lendings), 1)

	// Delete lending
	resp = makeRequest(t, http.MethodDelete, fmt.Sprintf("%s/lendings/%s", baseURL, createdLending.ID), nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Clean up created book and user
	resp = makeRequest(t, http.MethodDelete, fmt.Sprintf("%s/books/%s", baseURL, createdBook.ID), nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	resp = makeRequest(t, http.MethodDelete, fmt.Sprintf("%s/users/%s", baseURL, createdUser.ID), nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
