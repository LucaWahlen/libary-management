package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"libary-service/internal/domain"
	"log"
	"os"
	"time"
)

var db *pgx.Conn

func Connect() error {
	dbURL := os.Getenv("DATABASE_URL")
	var err error
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return err
	}
	db = conn
	log.Printf("Successfully connected to database")
	return nil
}

func Disconnect() error {
	if db != nil {
		db.Close(context.Background())
		log.Printf("Successfully disconnected from database")
	}
	return errors.New("error closing database connection")
}

func GetBooks() ([]domain.Book, error) {
	rows, err := db.Query(context.Background(), "SELECT id, title, author FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		var b domain.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func GetBookByID(id string) (domain.Book, error) {
	var b domain.Book
	err := db.QueryRow(context.Background(), "SELECT id, title, author FROM books WHERE id = $1", id).
		Scan(&b.ID, &b.Title, &b.Author)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Book{}, errors.New("book not found")
	} else if err != nil {
		return domain.Book{}, err
	}
	return b, nil
}

func CreateBook(book domain.Book) (domain.Book, error) {
	_, err := db.Exec(context.Background(), "INSERT INTO books (id, title, author) VALUES ($1, $2, $3)",
		book.ID, book.Title, book.Author)
	if err != nil {
		return domain.Book{}, err
	}
	return book, nil
}

func UpdateBook(book domain.Book) (domain.Book, error) {
	result, err := db.Exec(context.Background(), "UPDATE books SET title = $2, author = $3 WHERE id = $1",
		book.ID, book.Title, book.Author)
	if err != nil {
		return domain.Book{}, err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.Book{}, errors.New("book not found")
	}
	return book, nil
}

func DeleteBook(id string) error {
	result, err := db.Exec(context.Background(), "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("book not found")
	}
	return nil
}

func GetUsers() ([]domain.User, error) {
	rows, err := db.Query(context.Background(), "SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserByID(id string) (domain.User, error) {
	var u domain.User
	err := db.QueryRow(context.Background(), "SELECT id, name, email FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.User{}, errors.New("user not found")
	} else if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func CreateUser(user domain.User) (domain.User, error) {
	_, err := db.Exec(context.Background(), "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)",
		user.ID, user.Name, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func UpdateUser(user domain.User) (domain.User, error) {
	result, err := db.Exec(context.Background(), "UPDATE users SET name = $2, email = $3 WHERE id = $1",
		user.ID, user.Name, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.User{}, errors.New("user not found")
	}
	return user, nil
}

func DeleteUser(id string) error {
	result, err := db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func GetLendings() ([]domain.Lending, error) {
	rows, err := db.Query(context.Background(), "SELECT id, book_id, user_id, lend_date, return_date FROM lendings")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lendings []domain.Lending
	for rows.Next() {
		var l domain.Lending
		var returnDate sql.NullTime
		if err := rows.Scan(&l.ID, &l.BookID, &l.UserID, &l.LendDate, &returnDate); err != nil {
			return nil, err
		}
		if returnDate.Valid {
			l.ReturnDate = returnDate.Time
		} else {
			l.ReturnDate = time.Time{}
		}
		lendings = append(lendings, l)
	}
	return lendings, nil
}

func GetLendingByID(id string) (domain.Lending, error) {
	var l domain.Lending
	var returnDate sql.NullTime
	err := db.QueryRow(context.Background(), "SELECT id, book_id, user_id, lend_date, return_date FROM lendings WHERE id = $1", id).
		Scan(&l.ID, &l.BookID, &l.UserID, &l.LendDate, &returnDate)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Lending{}, errors.New("lending not found")
	} else if err != nil {
		return domain.Lending{}, err
	}
	if returnDate.Valid {
		l.ReturnDate = returnDate.Time
	} else {
		l.ReturnDate = time.Time{}
	}
	return l, nil
}

func CreateLending(lending domain.Lending) (domain.Lending, error) {
	var returnDate interface{}
	if lending.ReturnDate.IsZero() {
		returnDate = nil
	} else {
		returnDate = lending.ReturnDate
	}
	_, err := db.Exec(context.Background(),
		"INSERT INTO lendings (id, book_id, user_id, lend_date, return_date) VALUES ($1, $2, $3, $4, $5)",
		lending.ID, lending.BookID, lending.UserID, lending.LendDate, returnDate,
	)
	if err != nil {
		return domain.Lending{}, err
	}
	return lending, nil
}

func UpdateLending(lending domain.Lending) (domain.Lending, error) {
	var returnDate interface{}
	if lending.ReturnDate.IsZero() {
		returnDate = nil
	} else {
		returnDate = lending.ReturnDate
	}
	result, err := db.Exec(context.Background(),
		"UPDATE lendings SET book_id = $2, user_id = $3, lend_date = $4, return_date = $5 WHERE id = $1",
		lending.ID, lending.BookID, lending.UserID, lending.LendDate, returnDate,
	)
	if err != nil {
		return domain.Lending{}, err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.Lending{}, errors.New("lending not found")
	}
	return lending, nil
}

func DeleteLending(id string) error {
	result, err := db.Exec(context.Background(), "DELETE FROM lendings WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("lending not found")
	}
	return nil
}
