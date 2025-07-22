package bl

import (
	"testing"

	"book_boy/backend/internal/models"
)

type mockAudiobookRepo struct {
	Audiobooks []models.Audiobook
	Err        error
}

func (m *mockAudiobookRepo) GetAll() ([]models.Audiobook, error) {
	return m.Audiobooks, m.Err
}

func TestAudibookService_GetAll(t *testing.T) {
	mockData := []models.Audiobook{
		{ID: 1, ISBN: "1111", Title: "Test Book A"},
		{ID: 2, ISBN: "2222", Title: "Test Book B"},
	}

	mockRepo := &mockAudiobookRepo{
		Audiobooks: mockData,
		Err:        nil,
	}

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
