package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type BookService interface {
	GetAll() ([]models.Book, error)
	GetByID(id int) (*models.Book, error)
	Create(book *models.Book) (int, error)
	Update(book *models.Book) error
	Delete(id int) error
}

type bookService struct {
	repo dl.BookRepo
}

func NewBookService(repo dl.BookRepo) BookService {
	return &bookService{repo: repo}
}

// CRUD
func (s *bookService) GetAll() ([]models.Book, error) {
	return s.repo.GetAll()
}

func (s *bookService) GetByID(id int) (*models.Book, error) {
	return s.repo.GetByID(id)
}

func (s *bookService) Create(book *models.Book) (int, error) {
	return s.repo.Create(book)
}

func (s *bookService) Update(book *models.Book) error {
	return s.repo.Update(book)
}

func (s *bookService) Delete(id int) error {
	return s.repo.Delete(id)
}

// Extensions
func (s *bookService) GetByTitle(title string) (*models.Book, error) {
	return s.repo.GetByTitle(title)
}
