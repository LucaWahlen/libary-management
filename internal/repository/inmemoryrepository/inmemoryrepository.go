package inmemoryrepository

import (
	"errors"
	"libary-service/internal/domain"
	"sync"
)

type InMemoryRepository struct {
	mu       sync.Mutex
	books    map[string]domain.Book
	users    map[string]domain.User
	lendings map[string]domain.Lending
}

func New() *InMemoryRepository {
	return &InMemoryRepository{
		books:    make(map[string]domain.Book),
		users:    make(map[string]domain.User),
		lendings: make(map[string]domain.Lending),
	}
}

func (repo *InMemoryRepository) Connect() error {
	return nil
}

func (repo *InMemoryRepository) Disconnect() error {
	return nil
}

func (repo *InMemoryRepository) GetBooks() ([]domain.Book, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	books := make([]domain.Book, 0, len(repo.books))
	for _, b := range repo.books {
		books = append(books, b)
	}
	return books, nil
}

func (repo *InMemoryRepository) GetBookByID(id string) (domain.Book, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if book, ok := repo.books[id]; ok {
		return book, nil
	}
	return domain.Book{}, errors.New("book not found")
}

func (repo *InMemoryRepository) CreateBook(book domain.Book) (domain.Book, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.books[book.ID] = book
	return book, nil
}

func (repo *InMemoryRepository) UpdateBook(updated domain.Book) (domain.Book, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.books[updated.ID]; !ok {
		return domain.Book{}, errors.New("book not found")
	}
	repo.books[updated.ID] = updated
	return updated, nil
}

func (repo *InMemoryRepository) DeleteBook(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.books[id]; !ok {
		return errors.New("book not found")
	}
	delete(repo.books, id)
	return nil
}

func (repo *InMemoryRepository) GetUsers() ([]domain.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	users := make([]domain.User, 0, len(repo.users))
	for _, u := range repo.users {
		users = append(users, u)
	}
	return users, nil
}

func (repo *InMemoryRepository) GetUserByID(id string) (domain.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if user, ok := repo.users[id]; ok {
		return user, nil
	}
	return domain.User{}, errors.New("user not found")
}

func (repo *InMemoryRepository) CreateUser(user domain.User) (domain.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.users[user.ID] = user
	return user, nil
}

func (repo *InMemoryRepository) UpdateUser(updated domain.User) (domain.User, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.users[updated.ID]; !ok {
		return domain.User{}, errors.New("user not found")
	}
	repo.users[updated.ID] = updated
	return updated, nil
}

func (repo *InMemoryRepository) DeleteUser(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(repo.users, id)
	return nil
}

func (repo *InMemoryRepository) GetLendings() ([]domain.Lending, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	lendings := make([]domain.Lending, 0, len(repo.lendings))
	for _, l := range repo.lendings {
		lendings = append(lendings, l)
	}
	return lendings, nil
}

func (repo *InMemoryRepository) GetLendingByID(id string) (domain.Lending, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if lending, ok := repo.lendings[id]; ok {
		return lending, nil
	}
	return domain.Lending{}, errors.New("lending not found")
}

func (repo *InMemoryRepository) CreateLending(lending domain.Lending) (domain.Lending, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.lendings[lending.ID] = lending
	return lending, nil
}

func (repo *InMemoryRepository) UpdateLending(updated domain.Lending) (domain.Lending, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.lendings[updated.ID]; !ok {
		return domain.Lending{}, errors.New("lending not found")
	}
	repo.lendings[updated.ID] = updated
	return updated, nil
}

func (repo *InMemoryRepository) DeleteLending(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, ok := repo.lendings[id]; !ok {
		return errors.New("lending not found")
	}
	delete(repo.lendings, id)
	return nil
}
