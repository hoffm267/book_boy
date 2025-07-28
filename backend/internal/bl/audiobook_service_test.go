package bl

import (
	"errors"
	"testing"

	"book_boy/backend/internal/models"
)

type mockAudiobookRepo struct {
	Audiobooks   []models.Audiobook
	Err          error
	LastCreated  *models.Audiobook
	LastUpdated  *models.Audiobook
	LastDeleted  int
	GetByIDInput int
}

func (m *mockAudiobookRepo) GetAll() ([]models.Audiobook, error) {
	return m.Audiobooks, m.Err
}

func (m *mockAudiobookRepo) GetByID(id int) (*models.Audiobook, error) {
	m.GetByIDInput = id
	if m.Err != nil {
		return nil, m.Err
	}
	for _, ab := range m.Audiobooks {
		if ab.ID == id {
			return &ab, nil
		}
	}
	return nil, nil
}

func (m *mockAudiobookRepo) Create(ab *models.Audiobook) (int, error) {
	m.LastCreated = ab
	if m.Err != nil {
		return 0, m.Err
	}
	return 123, nil
}

func (m *mockAudiobookRepo) Update(ab *models.Audiobook) error {
	m.LastUpdated = ab
	return m.Err
}

func (m *mockAudiobookRepo) Delete(id int) error {
	m.LastDeleted = id
	return m.Err
}

// ---- TESTS ----

func TestAudiobookService_GetAll(t *testing.T) {
	mockData := []models.Audiobook{
		{ID: 1, ISBN: "1111", Title: "Test Book A"},
		{ID: 2, ISBN: "2222", Title: "Test Book B"},
	}

	mockRepo := &mockAudiobookRepo{Audiobooks: mockData}
	svc := NewAudiobookService(mockRepo)

	result, err := svc.GetAll()
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

func TestAudiobookService_GetByID(t *testing.T) {
	mockRepo := &mockAudiobookRepo{
		Audiobooks: []models.Audiobook{{ID: 1, ISBN: "1111", Title: "One"}},
	}
	svc := NewAudiobookService(mockRepo)

	result, err := svc.GetByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.ID != 1 {
		t.Fatalf("expected ID 1, got %+v", result)
	}
}

func TestAudiobookService_Create(t *testing.T) {
	mockRepo := &mockAudiobookRepo{}
	svc := NewAudiobookService(mockRepo)

	ab := &models.Audiobook{ISBN: "3333", Title: "New Book"}
	id, err := svc.Create(ab)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 123 {
		t.Errorf("expected ID 123, got %d", id)
	}
	if mockRepo.LastCreated != ab {
		t.Errorf("Create was not called with correct data")
	}
}

func TestAudiobookService_Update(t *testing.T) {
	mockRepo := &mockAudiobookRepo{}
	svc := NewAudiobookService(mockRepo)

	ab := &models.Audiobook{ID: 1, ISBN: "4444", Title: "Updated Book"}
	err := svc.Update(ab)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mockRepo.LastUpdated != ab {
		t.Errorf("Update was not called with correct data")
	}
}

func TestAudiobookService_Delete(t *testing.T) {
	mockRepo := &mockAudiobookRepo{}
	svc := NewAudiobookService(mockRepo)

	err := svc.Delete(99)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mockRepo.LastDeleted != 99 {
		t.Errorf("expected LastDeleted to be 99, got %d", mockRepo.LastDeleted)
	}
}

func TestAudiobookService_Errors(t *testing.T) {
	mockRepo := &mockAudiobookRepo{Err: errors.New("db error")}
	svc := NewAudiobookService(mockRepo)

	if _, err := svc.GetAll(); err == nil {
		t.Error("expected GetAll to return error")
	}
	if _, err := svc.GetByID(1); err == nil {
		t.Error("expected GetByID to return error")
	}
	if _, err := svc.Create(&models.Audiobook{}); err == nil {
		t.Error("expected Create to return error")
	}
	if err := svc.Update(&models.Audiobook{}); err == nil {
		t.Error("expected Update to return error")
	}
	if err := svc.Delete(1); err == nil {
		t.Error("expected Delete to return error")
	}
}
