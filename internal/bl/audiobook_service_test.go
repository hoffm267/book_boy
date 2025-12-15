package bl

import (
	"errors"
	"testing"
	"time"

	"book_boy/internal/models"
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
	for _, audiobook := range m.Audiobooks {
		if audiobook.ID == id {
			return &audiobook, nil
		}
	}
	return nil, nil
}

func (m *mockAudiobookRepo) Create(audiobook *models.Audiobook) (int, error) {
	m.LastCreated = audiobook
	if m.Err != nil {
		return 0, m.Err
	}
	return 123, nil
}

func (m *mockAudiobookRepo) Update(audiobook *models.Audiobook) error {
	m.LastUpdated = audiobook
	return m.Err
}

func (m *mockAudiobookRepo) Delete(id int) error {
	m.LastDeleted = id
	return m.Err
}

func (m *mockAudiobookRepo) GetSimilarTitles(title string) ([]models.Audiobook, error) {
	//TODO IMPLEMENT
	return nil, nil
}

// ---- TESTS ----

func TestAudiobookService_GetAll(t *testing.T) {
	d1 := time.Hour + 45*time.Minute + 30*time.Second
	d2 := 90 * time.Minute

	mockData := []models.Audiobook{
		{ID: 1, Title: "Test Book A", TotalLength: &models.CustomDuration{Duration: d1}},
		{ID: 2, Title: "Test Book B", TotalLength: &models.CustomDuration{Duration: d2}},
	}

	mockRepo := &mockAudiobookRepo{Audiobooks: mockData}
	svc := NewAudiobookService(mockRepo, nil)

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
	d1 := time.Hour + 45*time.Minute + 30*time.Second
	mockRepo := &mockAudiobookRepo{
		Audiobooks: []models.Audiobook{{ID: 1, Title: "One", TotalLength: &models.CustomDuration{Duration: d1}}},
	}
	svc := NewAudiobookService(mockRepo, nil)

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
	svc := NewAudiobookService(mockRepo, nil)

	duration := &models.CustomDuration{}
	duration.Duration = 3600000000000
	audiobook := &models.Audiobook{Title: "New Book", TotalLength: duration}
	id, err := svc.Create(audiobook)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 123 {
		t.Errorf("expected ID 123, got %d", id)
	}
	if mockRepo.LastCreated != audiobook {
		t.Errorf("Create was not called with correct data")
	}
}

func TestAudiobookService_Update(t *testing.T) {
	mockRepo := &mockAudiobookRepo{}
	svc := NewAudiobookService(mockRepo, nil)

	d1 := time.Hour + 45*time.Minute + 30*time.Second
	audiobook := &models.Audiobook{ID: 1, Title: "Updated Book", TotalLength: &models.CustomDuration{Duration: d1}}
	err := svc.Update(audiobook)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mockRepo.LastUpdated != audiobook {
		t.Errorf("Update was not called with correct data")
	}
}

func TestAudiobookService_Delete(t *testing.T) {
	mockRepo := &mockAudiobookRepo{}
	svc := NewAudiobookService(mockRepo, nil)

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
	svc := NewAudiobookService(mockRepo, nil)

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

func TestAudiobookService_GetSimilarTitles(t *testing.T) {
	d1 := time.Hour + 30*time.Minute
	mockRepo := &mockAudiobookRepo{
		Audiobooks: []models.Audiobook{
			{ID: 1, Title: "The Great Gatsby", TotalLength: &models.CustomDuration{Duration: d1}},
			{ID: 2, Title: "The Catcher in the Rye", TotalLength: &models.CustomDuration{Duration: d1}},
		},
	}
	svc := NewAudiobookService(mockRepo, nil)

	audiobooks, err := svc.GetSimilarTitles("Great")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if audiobooks != nil {
		t.Logf("GetSimilarTitles returned %d audiobooks (implementation pending)", len(audiobooks))
	}
}
