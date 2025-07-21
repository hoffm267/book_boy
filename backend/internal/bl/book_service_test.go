package bl

import (
	"testing"

	"book_boy/backend/internal/models"
)

type mockBookRepo struct {
	Books []models.Book
	Err   error
}

func (m *mockBookRepo) GetAll() ([]models.Book, error) {
	return m.Books, m.Err
}

func TestBookService_GetAllBooks(t *testing.T) {
	mockData := []models.Book{
		{ID: 1, ISBN: "1111", Title: "Test Book A"},
		{ID: 2, ISBN: "2222", Title: "Test Book B"},
	}

	mockRepo := &mockBookRepo{
		Books: mockData,
		Err:   nil,
	}

	svc := NewBookService(mockRepo)
	result, err := svc.GetAllBooks()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(mockData) {
		t.Fatalf("expected %d books, got %d", len(mockData), len(result))
	}

	for i := range result {
		if result[i] != mockData[i] {
			t.Errorf("mismatch at index %d: expected %+v, got %+v", i, mockData[i], result[i])
		}
	}
}
