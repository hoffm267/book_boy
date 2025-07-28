package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type BookService interface {
	GetAll() ([]models.Book, error)
	GetByID(bookID int) (*models.Book, error)
	Create(book *models.Book) (int, error)
	Update(book *models.Book) error
	Delete(bookID int) error
}

type bookService struct {
	repo dl.BookRepo
}

func NewBookService(repo dl.BookRepo) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) GetAll() ([]models.Book, error) {
	return s.repo.GetAll()
}

func (s *bookService) GetByID(bookID int) (*models.Book, error) {
	return s.repo.GetByID(bookID)
}

func (s *bookService) Create(book *models.Book) (int, error) {
	return s.repo.Create(book)
}

func (s *bookService) Update(book *models.Book) error {
	return s.repo.Update(book)
}

func (s *bookService) Delete(bookID int) error {
	return s.repo.Delete(bookID)
}
