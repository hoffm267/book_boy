package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type BookService interface {
	GetAll() ([]models.Book, error)
	GetByID(bookID int) (*models.Book, error)
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
