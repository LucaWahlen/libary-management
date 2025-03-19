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

	// hier vermischst du fuer mein Gefuehl Eingangsdatenvalidierung mit BusinessLogig
	// d.h. die referentielle Integritaet wuerde ich nicht hier pruefen
	// Tatsaechlich ist diese Integritaet ja auch in der DB modelliert und du bekommst dort einen
	// Fehler zurueck - allerdings dann nicht aufgeschluesselt ob das Buch oder der User fehlt.
	// Trotzdem wuerde ich diese Logik dann eher in der "app" Schicht sehen (als Businesslogik)
	// oder gekapselt in "CreateLending" (d.h. dort verschiedene Fehler zurueckgeben, die in der App-Schicht)
	// dann ausgewertet werden koennen.
	//
	// Aber ich verstehe auch dass es fuer deinen dependency injection Anwendungsfall schoen ist, wenn auch die Validation
	// Schicht eine Dependency hat ;-)
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

	// fuer mich besser lesbar so
	if !lending.ReturnDate.IsZero() && !lending.ReturnDate.After(lending.ReturnDate) {
		// if !lending.ReturnDate.IsZero() {
		// 	if !lending.LendDate.Before(lending.ReturnDate) {
		errs = append(errs, fmt.Errorf("lend_date is less than return_date"))
	}

	return errors.Join(errs...)
}
