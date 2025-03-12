package mockvalidator

import "libary-service/internal/domain"

type MockValidator struct{}

func (m MockValidator) CheckCreateBook(book domain.Book) error {
	return nil
}

func (m MockValidator) CheckUpdateBook(book domain.Book) error {
	return nil
}

func (m MockValidator) CheckCreateUser(user domain.User) error {
	return nil
}

func (m MockValidator) CheckUpdateUser(user domain.User) error {
	return nil
}

func (m MockValidator) CheckCreateLending(lending domain.Lending) error {
	return nil
}

func (m MockValidator) CheckUpdateLending(lending domain.Lending) error {
	return nil
}
