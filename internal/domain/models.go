package domain

import "time"

type Book struct {
	ID     string `json:"id,omitempty" db:"id"`
	Title  string `json:"title" db:"title"`
	Author string `json:"author" db:"author"`
}

type User struct {
	ID    string `json:"id,omitempty" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
}

type Lending struct {
	ID         string    `json:"id" db:"id"`
	BookID     string    `json:"book_id" db:"book_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	LendDate   time.Time `json:"lend_date" db:"lend_date"`
	ReturnDate time.Time `json:"return_date,omitempty" db:"return_date"`
}
