package bl

import (
	"testing"

	"book_boy/backend/internal/models"
)

type mockUserRepo struct {
	Users []models.User
	Err   error
}

func (m *mockUserRepo) GetAll() ([]models.User, error) {
	return m.Users, m.Err
}

func TestUserService_GetAll(t *testing.T) {
	mockData := []models.User{
		{ID: 1, Username: "alice"},
		{ID: 2, Username: "bob"},
	}

	mockRepo := &mockUserRepo{
		Users: mockData,
		Err:   nil,
	}

	svc := NewUserService(mockRepo)
	result, err := svc.GetAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != len(mockData) {
		t.Errorf("expected %d users, got %d", len(mockData), len(result))
	}

	for i := range result {
		if result[i] != mockData[i] {
			t.Errorf("mismatch at index %d: expected %+v, got %+v", i, mockData[i], result[i])
		}
	}
}
