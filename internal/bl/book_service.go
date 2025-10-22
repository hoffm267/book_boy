package bl

import (
	"book_boy/internal/dl"
	"book_boy/internal/models"
)

type BookService interface {
	GetAll() ([]models.Book, error)
	GetByID(id int) (*models.Book, error)
	Create(book *models.Book) (int, error)
	Update(book *models.Book) error
	Delete(id int) error
	GetByTitle(title string) (*models.Book, error)
	GetSimilarTitles(title string) ([]models.Book, error)
	FilterBooks(filter models.BookFilter) ([]models.Book, error)
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

func (s *bookService) GetByID(id int) (*models.Book, error) {
	return s.repo.GetByID(id)
}

func (s *bookService) Create(book *models.Book) (int, error) {
	if err := book.Validate(); err != nil {
		return 0, err
	}
	return s.repo.Create(book)
}

func (s *bookService) Update(book *models.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	return s.repo.Update(book)
}

func (s *bookService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *bookService) GetByTitle(title string) (*models.Book, error) {
	return s.repo.GetByTitle(title)
}

func (s *bookService) GetSimilarTitles(title string) ([]models.Book, error) {
	return s.repo.GetSimilarTitles(title)
}

func (s *bookService) FilterBooks(filter models.BookFilter) ([]models.Book, error) {
	return s.repo.FilterBooks(filter)
}
