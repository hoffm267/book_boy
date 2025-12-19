package service

import (
	"book_boy/api/internal/domain"
	"book_boy/api/internal/repository"
	"book_boy/api/internal/infra"
	"context"
	"fmt"
	"time"
)

type BookService interface {
	GetAll() ([]domain.Book, error)
	GetByID(id int) (*domain.Book, error)
	Create(book *domain.Book) (int, error)
	Update(book *domain.Book) error
	Delete(id int) error
	GetByTitle(title string) (*domain.Book, error)
	GetSimilarTitles(title string) ([]domain.Book, error)
	FilterBooks(filter repository.BookFilter) ([]domain.Book, error)
}

type bookService struct {
	repo      repository.BookRepo
	cache     *infra.Cache
	publisher *infra.EventPublisher
}

func NewBookService(repo repository.BookRepo, cache *infra.Cache, publisher *infra.EventPublisher) BookService {
	return &bookService{repo: repo, cache: cache, publisher: publisher}
}

func (s *bookService) GetAll() ([]domain.Book, error) {
	return s.repo.GetAll()
}

func (s *bookService) GetByID(id int) (*domain.Book, error) {
	if s.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("book:%d", id)

		var book domain.Book
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

func (s *bookService) Create(book *domain.Book) (int, error) {
	if err := book.Validate(); err != nil {
		return 0, err
	}

	bookID, err := s.repo.Create(book)
	if err != nil {
		return 0, err
	}

	if s.publisher != nil && book.ISBN != "" {
		event := domain.BookCreatedEvent{
			BookID:    bookID,
			ISBN:      book.ISBN,
			CreatedAt: time.Now().Format(time.RFC3339),
		}
		if err := s.publisher.Publish("book.created", event); err != nil {
			fmt.Printf("Warning: failed to publish book.created event: %v\n", err)
		}
	}

	return bookID, nil
}

func (s *bookService) Update(book *domain.Book) error {
	if err := book.Validate(); err != nil {
		return err
	}
	if err := s.repo.Update(book); err != nil {
		return err
	}
	if s.cache != nil {
		s.cache.Delete(context.Background(), fmt.Sprintf("book:%d", book.ID))
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

func (s *bookService) GetByTitle(title string) (*domain.Book, error) {
	return s.repo.GetByTitle(title)
}

func (s *bookService) GetSimilarTitles(title string) ([]domain.Book, error) {
	return s.repo.GetSimilarTitles(title)
}

func (s *bookService) FilterBooks(filter repository.BookFilter) ([]domain.Book, error) {
	return s.repo.FilterBooks(filter)
}
