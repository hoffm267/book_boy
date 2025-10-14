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
	for _, user := range m.Users {
		result = append(result, user)
	}
	return result, m.Err
}

func (m *mockUserRepo) GetByID(id int) (*models.User, error) {
	if user, ok := m.Users[id]; ok {
		return &user, nil
	}
	return nil, nil
}

func (m *mockUserRepo) GetByEmail(email string) (*models.User, error) {
	for _, user := range m.Users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) Create(user *models.User) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	id := len(m.Users) + 1
	user.ID = id
	m.Users[user.ID] = *user
	return id, nil
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
	_, err := service.Create(user)
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
