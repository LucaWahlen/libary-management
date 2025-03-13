package validator

import (
	"errors"
	"fmt"
	"libary-service/internal/domain"
	"libary-service/internal/repository"
)

type Validator struct {
	repository repository.Repository
}

func New(repository repository.Repository) *Validator {
	return &Validator{repository}
}

func (v Validator) CheckBook(book domain.Book) error {
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
	return errors.Join(errs...)
}

func (v Validator) CheckUser(user domain.User) error {
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
	return errors.Join(errs...)
}

func (v Validator) CheckLending(lending domain.Lending) error {
	var errs []error

	if lending.ID != "" {
		errs = append(errs, fmt.Errorf("id should be empty"))
	}

	_, bookMissing := v.repository.GetBookByID(lending.BookID)
	if bookMissing != nil {
		errs = append(errs, fmt.Errorf("book not found"))
	}

	_, userMissing := v.repository.GetUserByID(lending.UserID)
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

	return errors.Join(errs...)
}
