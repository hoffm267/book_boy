package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type BookService interface {
	GetAllBooks() ([]models.Book, error)
}

type bookService struct {
	repo dl.BookRepo
}

func NewBookService(repo dl.BookRepo) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) GetAllBooks() ([]models.Book, error) {
	return s.repo.GetAll()
}
