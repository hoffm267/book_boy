package service

import (
	"book_boy/api/internal/domain"
	"errors"
	"testing"
)

type mockUserRepo struct {
	Users map[int]domain.User
	Err   error
}

func (m *mockUserRepo) GetAll() ([]domain.User, error) {
	var result []domain.User
	for _, user := range m.Users {
		result = append(result, user)
	}
	return result, m.Err
}

func (m *mockUserRepo) GetByID(id int) (*domain.User, error) {
	if user, ok := m.Users[id]; ok {
		return &user, nil
	}
	return nil, nil
}

func (m *mockUserRepo) GetByEmail(email string) (*domain.User, error) {
	for _, user := range m.Users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) Create(user *domain.User) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	id := len(m.Users) + 1
	user.ID = id
	m.Users[user.ID] = *user
	return id, nil
}

func (m *mockUserRepo) Update(user *domain.User) error {
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

func TestUserService_GetAll_Error(t *testing.T) {
	repo := &mockUserRepo{Users: make(map[int]domain.User), Err: errors.New("db error")}
	svc := NewUserService(repo)

	_, err := svc.GetAll()
	if err == nil {
		t.Fatal("expected error from repo")
	}
}

func TestUserService_GetByID_NotFound(t *testing.T) {
	repo := &mockUserRepo{Users: make(map[int]domain.User)}
	svc := NewUserService(repo)

	user, err := svc.GetByID(999)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil for non-existent user, got %+v", user)
	}
}

func TestUserService_CRUD(t *testing.T) {
	repo := &mockUserRepo{Users: make(map[int]domain.User)}
	service := NewUserService(repo)

	user := &domain.User{Username: "test_user"}
	_, err := service.Create(user)
	if err != nil || user.ID == 0 {
		t.Fatalf("Create failed: %v", err)
	}

	fetched, _ := service.GetByID(user.ID)
	if fetched == nil || fetched.Username != "test_user" {
		t.Fatalf("GetByID failed")
	}

	user.Username = "updated_user"
	if err := service.Update(user); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	updated, _ := service.GetByID(user.ID)
	if updated.Username != "updated_user" {
		t.Fatalf("Update did not persist")
	}

	if err := service.Delete(user.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	deleted, _ := service.GetByID(user.ID)
	if deleted != nil {
		t.Fatalf("Delete did not persist")
	}
}
