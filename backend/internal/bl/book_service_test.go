package bl

import (
	"testing"

	"book_boy/backend/internal/models"
)

type mockBookRepo struct {
	Books map[int]models.Book
	Err   error
}

func (m *mockBookRepo) GetAll() ([]models.Book, error) {
	books := make([]models.Book, 0, len(m.Books))
	for _, b := range m.Books {
		books = append(books, b)
	}
	return books, m.Err
}

func (m *mockBookRepo) GetByID(bookID int) (*models.Book, error) {
	b, ok := m.Books[bookID]
	if !ok {
		return nil, nil
	}
	return &b, nil
}

func TestBookService_GetAll(t *testing.T) {
	mockData := map[int]models.Book{
		1: {ID: 1, ISBN: "1111", Title: "Test Book A"},
		2: {ID: 2, ISBN: "2222", Title: "Test Book B"},
	}

	mockRepo := &mockBookRepo{
		Books: mockData,
		Err:   nil,
	}

	svc := NewBookService(mockRepo)
	result, err := svc.GetAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(mockData) {
		t.Fatalf("expected %d books, got %d", len(mockData), len(result))
	}

	for _, book := range result {
		expected, ok := mockData[book.ID]
		if !ok {
			t.Errorf("unexpected book ID %d in result", book.ID)
			continue
		}
		if book != expected {
			t.Errorf("for book ID %d: expected %+v, got %+v", book.ID, expected, book)
		}
	}
}

func TestBookService_GetBookByID(t *testing.T) {
	mockData := map[int]models.Book{
		1: {ID: 1, ISBN: "1111", Title: "Test Book A"},
		2: {ID: 2, ISBN: "2222", Title: "Test Book B"},
	}

	mockRepo := &mockBookRepo{
		Books: mockData,
		Err:   nil,
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
