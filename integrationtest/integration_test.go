package integrationtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"libary-service/internal/domain"
)

const (
	directServiceURL   = "http://localhost:8082"
	injectedServiceURL = "http://localhost:8080"
)

func clearDB(t *testing.T) {
	dbURL := "postgres://postgres:password@localhost:5432/libraryDB"
	conn, err := pgx.Connect(context.Background(), dbURL)
	require.NoError(t, err, "Failed to connect to database")
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "DELETE FROM lendings")
	require.NoError(t, err, "Failed to clean lendings table")
	_, err = conn.Exec(context.Background(), "DELETE FROM users")
	require.NoError(t, err, "Failed to clean users table")
	_, err = conn.Exec(context.Background(), "DELETE FROM books")
	require.NoError(t, err, "Failed to clean books table")
}

func makeRequest(t *testing.T, method, url string, body []byte) *http.Response {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	require.NoError(t, err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "Failed to execute request")

	return resp
}

func makeJsonRequest(t *testing.T, method, url string, filePath string) *http.Response {
	jsonBytes, err := ioutil.ReadFile(filePath)
	require.NoError(t, err, "Failed to read JSON file")
	return makeRequest(t, method, url, jsonBytes)
}

func decodeResponse(t *testing.T, resp *http.Response, target interface{}) {
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	if len(bodyBytes) == 0 || (bodyBytes[0] != '{' && bodyBytes[0] != '[') {
		t.Fatalf("Expected JSON response but got: %s", string(bodyBytes))
	}

	err = json.Unmarshal(bodyBytes, target)
	require.NoError(t, err, "Failed to decode JSON response")
}

func TestServices(t *testing.T) {
	t.Run("Testing Direct-Service", func(t *testing.T) { serviceTest(t, directServiceURL) })
	t.Run("Testing Injected-Service", func(t *testing.T) { serviceTest(t, injectedServiceURL) })
}

type reqLending struct {
	ID         *string   `json:"id,omitempty" db:"id"`
	BookID     string    `json:"book_id" db:"book_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	LendDate   time.Time `json:"lend_date" db:"lend_date"`
	ReturnDate time.Time `json:"return_date,omitempty" db:"return_date"`
}

func serviceTest(t *testing.T, baseURL string) {
	clearDB(t)

	t.Run("Testing Books", func(t *testing.T) { testBooks(t, baseURL) })
	t.Run("Testing Users", func(t *testing.T) { testUsers(t, baseURL) })
	t.Run("Testing Lendings", func(t *testing.T) { testLendings(t, baseURL) })
}

func testBooks(t *testing.T, baseURL string) {
	book := domain.Book{
		Title:  "The Fellowship of the Ring",
		Author: "J. R. R. Tolkien",
	}

	resp := makeJsonRequest(t, http.MethodPost, baseURL+"/books", "book_create.json")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdBook domain.Book
	decodeResponse(t, resp, &createdBook)

	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, book.Title, createdBook.Title)
	assert.Equal(t, book.Author, createdBook.Author)

	resp = makeRequest(t, http.MethodGet, fmt.Sprintf("%s/books/%s", baseURL, createdBook.ID), nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedBook domain.Book
	decodeResponse(t, resp, &retrievedBook)
	assert.Equal(t, createdBook, retrievedBook)

	resp = makeJsonRequest(t, http.MethodPut, fmt.Sprintf("%s/books/%s", baseURL, createdBook.ID), "book_update.json")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedBook domain.Book
	decodeResponse(t, resp, &updatedBook)
	assert.Equal(t, "The Two Towers", updatedBook.Title)

	resp = makeRequest(t, http.MethodGet, baseURL+"/books", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var books []domain.Book
	decodeResponse(t, resp, &books)
	assert.Equal(t, len(books), 1)
	assert.Equal(t, books[0], updatedBook)

	resp = makeRequest(t, http.MethodDelete, fmt.Sprintf("%s/books/%s", baseURL, createdBook.ID), nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	resp = makeRequest(t, http.MethodGet, fmt.Sprintf("%s/books/%s", baseURL, createdBook.ID), nil)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func testUsers(t *testing.T, baseURL string) {
	resp := makeJsonRequest(t, http.MethodPost, baseURL+"/users", "user_create.json")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdUser domain.User
	decodeResponse(t, resp, &createdUser)

	assert.NotEmpty(t, createdUser.ID)
	assert.Equal(t, "Max Mustermann", createdUser.Name)
	assert.Equal(t, "max@mustermann.de", createdUser.Email)

	resp = makeRequest(t, http.MethodGet, fmt.Sprintf("%s/users/%s", baseURL, createdUser.ID), nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedUser domain.User
	decodeResponse(t, resp, &retrievedUser)
	assert.Equal(t, createdUser.ID, retrievedUser.ID)

	resp = makeJsonRequest(t, http.MethodPut, fmt.Sprintf("%s/users/%s", baseURL, retrievedUser.ID), "user_update.json")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedUser domain.User
	decodeResponse(t, resp, &updatedUser)
	assert.Equal(t, "Erika Mustermann", updatedUser.Name)
	assert.Equal(t, "erika@mustermann.de", updatedUser.Email)

	resp = makeRequest(t, http.MethodGet, baseURL+"/users", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var users []domain.User
	decodeResponse(t, resp, &users)
	assert.Equal(t, len(users), 1)
	assert.Equal(t, updatedUser, users[0])
}

func testLendings(t *testing.T, baseURL string) {
	var createdBook domain.Book
	var createdUser domain.User

	resp := makeJsonRequest(t, http.MethodPost, baseURL+"/books", "book_create.json")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	decodeResponse(t, resp, &createdBook)

	resp = makeJsonRequest(t, http.MethodPost, baseURL+"/users", "user_create.json")
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	decodeResponse(t, resp, &createdUser)

	lending := reqLending{
		BookID:   createdBook.ID,
		UserID:   createdUser.ID,
		LendDate: time.Now(),
	}

	lendingBytes, err := json.Marshal(lending)
	assert.NoError(t, err)
	resp = makeRequest(t, http.MethodPost, baseURL+"/lendings", lendingBytes)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var createdLending domain.Lending
	decodeResponse(t, resp, &createdLending)

	assert.NotEmpty(t, createdLending.ID)
	assert.Equal(t, lending.BookID, createdLending.BookID)
	assert.Equal(t, lending.UserID, createdLending.UserID)

	resp = makeRequest(t, http.MethodGet, fmt.Sprintf("%s/lendings/%s", baseURL, createdLending.ID), nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedLending domain.Lending
	decodeResponse(t, resp, &retrievedLending)
	assert.Equal(t, createdLending.ID, retrievedLending.ID)

	lending.ReturnDate = time.Now().Add(7 * 24 * time.Hour)
	lendingBytes, err = json.Marshal(lending)
	assert.NoError(t, err)
	resp = makeRequest(t, http.MethodPut, fmt.Sprintf("%s/lendings/%s", baseURL, retrievedLending.ID), lendingBytes)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedLending domain.Lending
	decodeResponse(t, resp, &updatedLending)
	assert.NotZero(t, updatedLending.ReturnDate)

	resp = makeRequest(t, http.MethodGet, baseURL+"/lendings", nil)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var lendings []domain.Lending
	decodeResponse(t, resp, &lendings)
	assert.GreaterOrEqual(t, len(lendings), 1)

	resp = makeRequest(t, http.MethodDelete, fmt.Sprintf("%s/lendings/%s", baseURL, createdLending.ID), nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	resp = makeRequest(t, http.MethodDelete, fmt.Sprintf("%s/books/%s", baseURL, createdBook.ID), nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	resp = makeRequest(t, http.MethodDelete, fmt.Sprintf("%s/users/%s", baseURL, createdUser.ID), nil)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
