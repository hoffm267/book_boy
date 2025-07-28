package bl

import (
	"book_boy/backend/internal/models"
	"errors"
	"testing"
)

type mockUserRepo struct {
	Users map[int]models.User
	Err   error
}

func (m *mockUserRepo) GetAll() ([]models.User, error) {
	var result []models.User
	for _, u := range m.Users {
		result = append(result, u)
	}
	return result, m.Err
}

func (m *mockUserRepo) GetByID(id int) (*models.User, error) {
	if u, ok := m.Users[id]; ok {
		return &u, nil
	}
	return nil, nil
}

func (m *mockUserRepo) Create(user *models.User) error {
	if m.Err != nil {
		return m.Err
	}
	user.ID = len(m.Users) + 1
	m.Users[user.ID] = *user
	return nil
}

func (m *mockUserRepo) Update(user *models.User) error {
	if _, exists := m.Users[user.ID]; !exists {
		return errors.New("user not found")
	}
	m.Users[user.ID] = *user
	return nil
}

func (m *mockUserRepo) Delete(id int) error {
	delete(m.Users, id)
	return nil
}

func TestUserService_CRUD(t *testing.T) {
	repo := &mockUserRepo{Users: make(map[int]models.User)}
	service := NewUserService(repo)

	// Create
	user := &models.User{Username: "test_user"}
	err := service.Create(user)
	if err != nil || user.ID == 0 {
		t.Fatalf("Create failed: %v", err)
	}

	// Read
	fetched, _ := service.GetByID(user.ID)
	if fetched == nil || fetched.Username != "test_user" {
		t.Fatalf("GetByID failed")
	}

	// Update
	user.Username = "updated_user"
	if err := service.Update(user); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Confirm update
	updated, _ := service.GetByID(user.ID)
	if updated.Username != "updated_user" {
		t.Fatalf("Update did not persist")
	}

	// Delete
	if err := service.Delete(user.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Confirm deletion
	deleted, _ := service.GetByID(user.ID)
	if deleted != nil {
		t.Fatalf("Delete did not persist")
	}
}
