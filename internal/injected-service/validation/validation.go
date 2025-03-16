//go:generate mockery --name=Validation --output=../../../generated/mocks --case=underscore
package validation

import (
	"libary-service/internal/domain"
)

type Validation interface {
	CheckBook(book domain.Book) error
	CheckUser(user domain.User) error
	CheckLending(lending domain.Lending) error
}
