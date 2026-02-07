package service

import (
	"errors"
	"testing"

	"book_boy/api/internal/domain"
	"book_boy/api/internal/repository"
)

type mockBookRepo struct {
	Books       map[int]domain.Book
	Err         error
	LastCreated *domain.Book
	LastUpdated *domain.Book
	LastDeleted int
}

func (m *mockBookRepo) GetAll() ([]domain.Book, error) {
	books := make([]domain.Book, 0, len(m.Books))
	for _, book := range m.Books {
		books = append(books, book)
	}
	return books, m.Err
}

func (m *mockBookRepo) GetByID(id int) (*domain.Book, error) {
	book, ok := m.Books[id]
	if !ok {
		return nil, nil
	}
	return &book, nil
}

func (m *mockBookRepo) Create(book *domain.Book) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	id := len(m.Books) + 1
	book.ID = id
	m.Books[id] = *book
	m.LastCreated = book
	return id, nil
}

func (m *mockBookRepo) Update(book *domain.Book) error {
	if m.Err != nil {
		return m.Err
	}
	if _, exists := m.Books[book.ID]; !exists {
		return errors.New("book not found")
	}
	m.Books[book.ID] = *book
	m.LastUpdated = book
	return nil
}

func (m *mockBookRepo) Delete(id int) error {
	if m.Err != nil {
		return m.Err
	}
	delete(m.Books, id)
	m.LastDeleted = id
	return nil
}

func (m *mockBookRepo) GetByTitle(title string) (*domain.Book, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	for _, book := range m.Books {
		if book.Title == title {
			return &book, nil
		}
	}
	return nil, nil
}

func (m *mockBookRepo) GetSimilarTitles(title string) ([]domain.Book, error) {
	//TODO IMPLEMENT
	return nil, nil
}

func (m *mockBookRepo) FilterBooks(filter repository.BookFilter) ([]domain.Book, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var results []domain.Book
	for _, book := range m.Books {
		match := true
		if filter.ID != nil && book.ID != *filter.ID {
			match = false
		}
		if filter.ISBN != nil && book.ISBN != *filter.ISBN {
			match = false
		}
		if filter.Title != nil && book.Title != *filter.Title {
			match = false
		}
		if filter.TotalPages != nil && book.TotalPages != *filter.TotalPages {
			match = false
		}
		if match {
			results = append(results, book)
		}
	}
	return results, nil
}

func TestBookService_GetAll(t *testing.T) {
	mockRepo := &mockBookRepo{
		Books: map[int]domain.Book{
			1: {ID: 1, ISBN: "1111", Title: "Test Book A", TotalPages: 500},
			2: {ID: 2, ISBN: "2222", Title: "Test Book B", TotalPages: 500},
		},
	}
	svc := NewBookService(mockRepo, nil, nil)

	result, err := svc.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(mockRepo.Books) {
		t.Fatalf("expected %d books, got %d", len(mockRepo.Books), len(result))
	}

	for _, book := range result {
		expected, ok := mockRepo.Books[book.ID]
		if !ok {
			t.Errorf("unexpected book ID %d in result", book.ID)
			continue
		}
		if book != expected {
			t.Errorf("for book ID %d: expected %+v, got %+v", book.ID, expected, book)
		}
	}
}

func TestBookService_GetByID(t *testing.T) {
	mockRepo := &mockBookRepo{
		Books: map[int]domain.Book{
			1: {ID: 1, ISBN: "1111", Title: "Test Book A", TotalPages: 500},
			2: {ID: 2, ISBN: "2222", Title: "Test Book B", TotalPages: 500},
		},
		Err: nil,
	}

	svc := NewBookService(mockRepo, nil, nil)

	t.Run("found", func(t *testing.T) {
		book, err := svc.GetByID(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if book == nil || book.ID != 1 {
			t.Fatalf("expected book ID 1, got %+v", book)
		}
	})

	t.Run("not found", func(t *testing.T) {
		book, err := svc.GetByID(-1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if book != nil {
			t.Fatalf("expected nil, got %+v", book)
		}
	})
}

func TestBookService_Create(t *testing.T) {
	mockRepo := &mockBookRepo{Books: make(map[int]domain.Book)}
	svc := NewBookService(mockRepo, nil, nil)

	book := &domain.Book{ISBN: "3333", Title: "New Book", TotalPages: 123}
	id, err := svc.Create(book)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == 0 || mockRepo.Books[id].Title != "New Book" {
		t.Fatalf("book not created properly")
	}
}

func TestBookService_Update(t *testing.T) {
	mockRepo := &mockBookRepo{
		Books: map[int]domain.Book{
			1: {ID: 1, ISBN: "1111", Title: "Old Title", TotalPages: 322},
		},
	}
	svc := NewBookService(mockRepo, nil, nil)

	book := &domain.Book{ID: 1, ISBN: "1111", Title: "Updated Title", TotalPages: 700}
	err := svc.Update(book)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mockRepo.Books[1].Title != "Updated Title" {
		t.Fatalf("book title not updated")
	}
}

func TestBookService_Delete(t *testing.T) {
	mockRepo := &mockBookRepo{
		Books: map[int]domain.Book{
			1: {ID: 1, ISBN: "1111", Title: "Delete Me", TotalPages: 100},
		},
	}
	svc := NewBookService(mockRepo, nil, nil)

	err := svc.Delete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, exists := mockRepo.Books[1]; exists {
		t.Fatalf("book not deleted")
	}
}

func TestBookService_GetByTitle(t *testing.T) {
	mockRepo := &mockBookRepo{
		Books: map[int]domain.Book{
			1: {ID: 1, ISBN: "1111", Title: "Unique Title", TotalPages: 200},
			2: {ID: 2, ISBN: "2222", Title: "Another Book", TotalPages: 300},
		},
	}
	svc := NewBookService(mockRepo, nil, nil)

	book, err := svc.GetByTitle("Unique Title")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if book == nil || book.Title != "Unique Title" {
		t.Fatalf("expected book with title 'Unique Title', got %+v", book)
	}

	notFound, err := svc.GetByTitle("Nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if notFound != nil {
		t.Fatalf("expected nil for nonexistent title, got %+v", notFound)
	}
}

func TestBookService_GetSimilarTitles(t *testing.T) {
	mockRepo := &mockBookRepo{
		Books: map[int]domain.Book{
			1: {ID: 1, ISBN: "1111", Title: "Harry Potter", TotalPages: 300},
			2: {ID: 2, ISBN: "2222", Title: "Lord of the Rings", TotalPages: 400},
		},
	}
	svc := NewBookService(mockRepo, nil, nil)

	books, err := svc.GetSimilarTitles("Potter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if books != nil {
		t.Logf("GetSimilarTitles returned %d books (implementation pending)", len(books))
	}
}

func TestBookService_Create_ValidationError(t *testing.T) {
	mockRepo := &mockBookRepo{Books: make(map[int]domain.Book)}
	svc := NewBookService(mockRepo, nil, nil)

	book := &domain.Book{ISBN: "1234", Title: "Valid", TotalPages: -5}
	_, err := svc.Create(book)
	if err == nil {
		t.Fatal("expected validation error for negative total_pages")
	}
}

func TestBookService_Errors(t *testing.T) {
	mockRepo := &mockBookRepo{Books: make(map[int]domain.Book), Err: errors.New("db error")}
	svc := NewBookService(mockRepo, nil, nil)

	if _, err := svc.GetAll(); err == nil {
		t.Error("expected GetAll to return error")
	}
	if _, err := svc.Create(&domain.Book{ISBN: "1234", Title: "Test", TotalPages: 100}); err == nil {
		t.Error("expected Create to return error")
	}
	if err := svc.Update(&domain.Book{ID: 1, ISBN: "1234", Title: "Test", TotalPages: 100}); err == nil {
		t.Error("expected Update to return error")
	}
	if err := svc.Delete(1); err == nil {
		t.Error("expected Delete to return error")
	}
	if _, err := svc.GetByTitle("test"); err == nil {
		t.Error("expected GetByTitle to return error")
	}
	if _, err := svc.FilterBooks(repository.BookFilter{}); err == nil {
		t.Error("expected FilterBooks to return error")
	}
}

func TestBookService_FilterBooks(t *testing.T) {
	mockRepo := &mockBookRepo{
		Books: map[int]domain.Book{
			1: {ID: 1, ISBN: "1111", Title: "Book A", TotalPages: 200},
			2: {ID: 2, ISBN: "2222", Title: "Book B", TotalPages: 300},
			3: {ID: 3, ISBN: "3333", Title: "Book C", TotalPages: 200},
		},
	}
	svc := NewBookService(mockRepo, nil, nil)

	pages := 200
	filter := repository.BookFilter{TotalPages: &pages}
	books, err := svc.FilterBooks(filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 books with 200 pages, got %d", len(books))
	}
	for _, book := range books {
		if book.TotalPages != 200 {
			t.Errorf("expected TotalPages 200, got %d", book.TotalPages)
		}
	}
}
