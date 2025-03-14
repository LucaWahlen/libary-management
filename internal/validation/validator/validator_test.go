package validator

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"libary-service/generated/mocks"
	"libary-service/internal/domain"
)

func TestCheckBook(t *testing.T) {
	v := New(nil)

	testCases := []struct {
		name           string
		book           domain.Book
		expectedErrors []string
	}{
		{
			name: "valid book",
			book: domain.Book{
				Title:  "The Fellowship of the Ring",
				Author: "J.R.R. Tolkien",
				ID:     "",
			},
			expectedErrors: nil,
		},
		{
			name: "missing title",
			book: domain.Book{
				Title:  "",
				Author: "J.R.R. Tolkien",
				ID:     "",
			},
			expectedErrors: []string{"title is required"},
		},
		{
			name: "missing author",
			book: domain.Book{
				Title:  "The Two Towers",
				Author: "",
				ID:     "",
			},
			expectedErrors: []string{"author is required"},
		},
		{
			name: "non-empty id",
			book: domain.Book{
				Title:  "The Return of the King",
				Author: "J.R.R. Tolkien",
				ID:     uuid.New().String(),
			},
			expectedErrors: []string{"id should be empty"},
		},
		{
			name: "multiple errors",
			book: domain.Book{
				Title:  "",
				Author: "",
				ID:     uuid.New().String(),
			},
			expectedErrors: []string{"title is required", "author is required", "id should be empty"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.CheckBook(tc.book)
			if len(tc.expectedErrors) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				for _, substr := range tc.expectedErrors {
					assert.Contains(t, err.Error(), substr)
				}
			}
		})
	}
}

func TestCheckUser(t *testing.T) {
	v := New(nil)

	testCases := []struct {
		name           string
		user           domain.User
		expectedErrors []string
	}{
		{
			name: "valid user",
			user: domain.User{
				Name:  "Max Mustermann",
				Email: "max@mustermann.de",
				ID:    "",
			},
			expectedErrors: nil,
		},
		{
			name: "missing name",
			user: domain.User{
				Name:  "",
				Email: "max@mustermann.de",
				ID:    "",
			},
			expectedErrors: []string{"name is required"},
		},
		{
			name: "missing email",
			user: domain.User{
				Name:  "Erika Mustermann",
				Email: "",
				ID:    "",
			},
			expectedErrors: []string{"email is required"},
		},
		{
			name: "non-empty id",
			user: domain.User{
				Name:  "Max Mustermann",
				Email: "max@mustermann.de",
				ID:    uuid.New().String(),
			},
			expectedErrors: []string{"id should be empty"},
		},
		{
			name: "multiple errors",
			user: domain.User{
				Name:  "",
				Email: "",
				ID:    uuid.New().String(),
			},
			expectedErrors: []string{"name is required", "email is required", "id should be empty"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.CheckUser(tc.user)
			if len(tc.expectedErrors) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				for _, substr := range tc.expectedErrors {
					assert.Contains(t, err.Error(), substr)
				}
			}
		})
	}
}

func TestCheckLending(t *testing.T) {
	lendDate := time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)
	returnDate := lendDate.Add(24 * time.Hour)

	validBook := domain.Book{
		ID:     uuid.New().String(),
		Title:  "The Fellowship of the Ring",
		Author: "J.R.R. Tolkien",
	}
	validUser := domain.User{
		ID:    uuid.New().String(),
		Name:  "Max Mustermann",
		Email: "max@mustermann.de",
	}

	testCases := []struct {
		name           string
		lending        domain.Lending
		bookErr        error
		userErr        error
		expectedErrors []string
	}{
		{
			name: "valid lending without return date",
			lending: domain.Lending{
				ID:         "",
				BookID:     validBook.ID,
				UserID:     validUser.ID,
				LendDate:   lendDate,
				ReturnDate: time.Time{},
			},
			bookErr:        nil,
			userErr:        nil,
			expectedErrors: nil,
		},
		{
			name: "valid lending with valid return date",
			lending: domain.Lending{
				ID:         "",
				BookID:     validBook.ID,
				UserID:     validUser.ID,
				LendDate:   lendDate,
				ReturnDate: returnDate,
			},
			bookErr:        nil,
			userErr:        nil,
			expectedErrors: nil,
		},
		{
			name: "non-empty id",
			lending: domain.Lending{
				ID:         uuid.New().String(),
				BookID:     validBook.ID,
				UserID:     validUser.ID,
				LendDate:   lendDate,
				ReturnDate: returnDate,
			},
			bookErr:        nil,
			userErr:        nil,
			expectedErrors: []string{"id should be empty"},
		},
		{
			name: "missing book",
			lending: domain.Lending{
				ID:         "",
				BookID:     uuid.New().String(),
				UserID:     validUser.ID,
				LendDate:   lendDate,
				ReturnDate: returnDate,
			},
			bookErr:        errors.New("not found"),
			userErr:        nil,
			expectedErrors: []string{"book not found"},
		},
		{
			name: "missing user",
			lending: domain.Lending{
				ID:         "",
				BookID:     validBook.ID,
				UserID:     uuid.New().String(),
				LendDate:   lendDate,
				ReturnDate: returnDate,
			},
			bookErr:        nil,
			userErr:        errors.New("not found"),
			expectedErrors: []string{"user not found"},
		},
		{
			name: "missing lend_date",
			lending: domain.Lending{
				ID:         "",
				BookID:     validBook.ID,
				UserID:     validUser.ID,
				LendDate:   time.Time{},
				ReturnDate: time.Time{},
			},
			bookErr:        nil,
			userErr:        nil,
			expectedErrors: []string{"lend_date is required"},
		},
		{
			name: "invalid return date (lend_date not before return_date)",
			lending: domain.Lending{
				ID:         "",
				BookID:     validBook.ID,
				UserID:     validUser.ID,
				LendDate:   lendDate,
				ReturnDate: lendDate,
			},
			bookErr:        nil,
			userErr:        nil,
			expectedErrors: []string{"lend_date is less than return_date"},
		},
		{
			name: "multiple errors in lending",
			lending: domain.Lending{
				ID:         uuid.New().String(),
				BookID:     uuid.New().String(),
				UserID:     uuid.New().String(),
				LendDate:   lendDate,
				ReturnDate: lendDate,
			},
			bookErr:        errors.New("not found"),
			userErr:        errors.New("not found"),
			expectedErrors: []string{"id should be empty", "book not found", "user not found", "lend_date is less than return_date"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.Repository)
			mockRepo.On("GetBookByID", tc.lending.BookID).Return(validBook, tc.bookErr)
			mockRepo.On("GetUserByID", tc.lending.UserID).Return(validUser, tc.userErr)

			val := New(mockRepo)
			err := val.CheckLending(tc.lending)
			if len(tc.expectedErrors) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				for _, substr := range tc.expectedErrors {
					assert.Contains(t, err.Error(), substr)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
