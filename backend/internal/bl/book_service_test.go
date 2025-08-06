package bl

import (
	"errors"
	"testing"

	"book_boy/backend/internal/models"
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
	//TODO IMPLEMENT
	return nil, nil
}

func (m *mockBookRepo) FilterBooks(filter models.BookFilter) ([]models.Book, error) {
	//TODO IMPLEMENT
	return nil, nil
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

func TestBookService_GetBookByID(t *testing.T) {
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
