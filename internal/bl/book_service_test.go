package bl

import (
	"errors"
	"testing"

	"book_boy/internal/models"
)

type mockBookRepo struct {
	Books       map[int]models.Book
	Err         error
	LastCreated *models.Book
	LastUpdated *models.Book
	LastDeleted int
}

func (m *mockBookRepo) GetAll() ([]models.Book, error) {
	books := make([]models.Book, 0, len(m.Books))
	for _, book := range m.Books {
		books = append(books, book)
	}
	return books, m.Err
}

func (m *mockBookRepo) GetByID(id int) (*models.Book, error) {
	book, ok := m.Books[id]
	if !ok {
		return nil, nil
	}
	return &book, nil
}

func (m *mockBookRepo) Create(book *models.Book) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	id := len(m.Books) + 1
	book.ID = id
	m.Books[id] = *book
	m.LastCreated = book
	return id, nil
}

func (m *mockBookRepo) Update(book *models.Book) error {
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

func (m *mockBookRepo) GetByTitle(title string) (*models.Book, error) {
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

func (m *mockBookRepo) GetSimilarTitles(title string) ([]models.Book, error) {
	//TODO IMPLEMENT
	return nil, nil
}

func (m *mockBookRepo) FilterBooks(filter models.BookFilter) ([]models.Book, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var results []models.Book
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
		Books: map[int]models.Book{
			1: {ID: 1, ISBN: "1111", Title: "Test Book A", TotalPages: 500},
			2: {ID: 2, ISBN: "2222", Title: "Test Book B", TotalPages: 500},
		},
	}
	svc := NewBookService(mockRepo)

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
		Books: map[int]models.Book{
			1: {ID: 1, ISBN: "1111", Title: "Test Book A", TotalPages: 500},
			2: {ID: 2, ISBN: "2222", Title: "Test Book B", TotalPages: 500},
		},
		Err: nil,
	}

	svc := NewBookService(mockRepo)

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
	mockRepo := &mockBookRepo{Books: make(map[int]models.Book)}
	svc := NewBookService(mockRepo)

	book := &models.Book{ISBN: "3333", Title: "New Book", TotalPages: 123}
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
		Books: map[int]models.Book{
			1: {ID: 1, ISBN: "1111", Title: "Old Title", TotalPages: 322},
		},
	}
	svc := NewBookService(mockRepo)

	book := &models.Book{ID: 1, ISBN: "1111", Title: "Updated Title", TotalPages: 700}
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
		Books: map[int]models.Book{
			1: {ID: 1, ISBN: "1111", Title: "Delete Me", TotalPages: 100},
		},
	}
	svc := NewBookService(mockRepo)

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
		Books: map[int]models.Book{
			1: {ID: 1, ISBN: "1111", Title: "Unique Title", TotalPages: 200},
			2: {ID: 2, ISBN: "2222", Title: "Another Book", TotalPages: 300},
		},
	}
	svc := NewBookService(mockRepo)

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
		Books: map[int]models.Book{
			1: {ID: 1, ISBN: "1111", Title: "Harry Potter", TotalPages: 300},
			2: {ID: 2, ISBN: "2222", Title: "Lord of the Rings", TotalPages: 400},
		},
	}
	svc := NewBookService(mockRepo)

	books, err := svc.GetSimilarTitles("Potter")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if books != nil {
		t.Logf("GetSimilarTitles returned %d books (implementation pending)", len(books))
	}
}

func TestBookService_FilterBooks(t *testing.T) {
	mockRepo := &mockBookRepo{
		Books: map[int]models.Book{
			1: {ID: 1, ISBN: "1111", Title: "Book A", TotalPages: 200},
			2: {ID: 2, ISBN: "2222", Title: "Book B", TotalPages: 300},
			3: {ID: 3, ISBN: "3333", Title: "Book C", TotalPages: 200},
		},
	}
	svc := NewBookService(mockRepo)

	pages := 200
	filter := models.BookFilter{TotalPages: &pages}
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
