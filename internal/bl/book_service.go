package bl

import (
	"book_boy/internal/dl"
	"book_boy/internal/infra"
	"book_boy/internal/models"
	"context"
	"fmt"
	"time"
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
	repo  dl.BookRepo
	cache *infra.Cache
	queue *infra.Queue
}

func NewBookService(repo dl.BookRepo, cache *infra.Cache, queue *infra.Queue) BookService {
	return &bookService{repo: repo, cache: cache, queue: queue}
}

func (s *bookService) GetAll() ([]models.Book, error) {
	return s.repo.GetAll()
}

func (s *bookService) GetByID(id int) (*models.Book, error) {
	if s.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("book:%d", id)

		var book models.Book
		if err := s.cache.Get(ctx, cacheKey, &book); err == nil {
			return &book, nil
		}
	}

	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		s.cache.Set(context.Background(), fmt.Sprintf("book:%d", id), result, 10*time.Minute)
	}
	return result, nil
}

func (s *bookService) Create(book *models.Book) (int, error) {
	if err := book.Validate(); err != nil {
		return 0, err
	}

	bookID, err := s.repo.Create(book)
	if err != nil {
		return 0, err
	}

	if s.queue != nil && book.ISBN != "" {
		job := infra.MetadataFetchJob{
			BookID: bookID,
			ISBN:   book.ISBN,
		}
		if err := s.queue.Publish("metadata_fetch", job); err != nil {
			fmt.Printf("Warning: failed to queue metadata fetch: %v\n", err)
		}
	}

	return bookID, nil
}

func (s *bookService) Update(book *models.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	if err := s.repo.Update(book); err != nil {
		return err
	}
	if s.cache != nil {
		if err := s.cache.Delete(context.Background(), fmt.Sprintf("book:%d", book.ID)); err != nil {
			return err
		}
	}
	return nil
}

func (s *bookService) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	if s.cache != nil {
		s.cache.Delete(context.Background(), fmt.Sprintf("book:%d", id))
	}
	return nil
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
