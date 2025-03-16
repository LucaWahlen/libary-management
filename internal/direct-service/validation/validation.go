package validation

import (
	"errors"
	"fmt"
	"libary-service/internal/direct-service/repository"
	"libary-service/internal/domain"
)

func CheckBook(book domain.Book) error {
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

func CheckUser(user domain.User) error {
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

func CheckLending(lending domain.Lending) error {
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

	return errors.Join(errs...)
}
