package service

import (
	"book_boy/api/internal/domain"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	if os.Getenv("JWT_SECRET") == "" {
		os.Setenv("JWT_SECRET", "test-secret-key")
	}
	os.Exit(m.Run())
}

type mockAuthUserRepo struct {
	Users        map[int]domain.User
	UsersByEmail map[string]domain.User
	Err          error
	NextID       int
}

func (m *mockAuthUserRepo) GetAll() ([]domain.User, error) {
	var result []domain.User
	for _, user := range m.Users {
		result = append(result, user)
	}
	return result, m.Err
}

func (m *mockAuthUserRepo) GetByID(id int) (*domain.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if user, ok := m.Users[id]; ok {
		return &user, nil
	}
	return nil, nil
}

func (m *mockAuthUserRepo) GetByEmail(email string) (*domain.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if user, ok := m.UsersByEmail[email]; ok {
		return &user, nil
	}
	return nil, nil
}

func (m *mockAuthUserRepo) Create(user *domain.User) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	m.NextID++
	user.ID = m.NextID
	m.Users[user.ID] = *user
	m.UsersByEmail[user.Email] = *user
	return user.ID, nil
}

func (m *mockAuthUserRepo) Update(user *domain.User) error {
	if m.Err != nil {
		return m.Err
	}
	if _, exists := m.Users[user.ID]; !exists {
		return nil
	}
	m.Users[user.ID] = *user
	return nil
}

func (m *mockAuthUserRepo) Delete(id int) error {
	if m.Err != nil {
		return m.Err
	}
	if user, exists := m.Users[id]; exists {
		delete(m.UsersByEmail, user.Email)
		delete(m.Users, id)
	}
	return nil
}

func setupAuthTest() *mockAuthUserRepo {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	repo := &mockAuthUserRepo{
		Users:        make(map[int]domain.User),
		UsersByEmail: make(map[string]domain.User),
		NextID:       0,
	}

	existingUser := domain.User{
		ID:           1,
		Username:     "existing_user",
		Email:        "existing@example.com",
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	repo.Users[1] = existingUser
	repo.UsersByEmail["existing@example.com"] = existingUser
	repo.NextID = 1

	return repo
}

func TestAuthService_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		repo := setupAuthTest()
		svc := NewAuthService(repo)

		req := &domain.RegisterRequest{
			Username: "newuser",
			Email:    "newuser@example.com",
			Password: "password123",
		}

		user, err := svc.Register(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user == nil {
			t.Fatal("expected user to be returned")
		}
		if user.Username != "newuser" {
			t.Errorf("expected username 'newuser', got %s", user.Username)
		}
		if user.Email != "newuser@example.com" {
			t.Errorf("expected email 'newuser@example.com', got %s", user.Email)
		}
		if user.ID == 0 {
			t.Error("expected user ID to be set")
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123"))
		if err != nil {
			t.Error("password was not hashed correctly")
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		repo := setupAuthTest()
		svc := NewAuthService(repo)

		req := &domain.RegisterRequest{
			Username: "duplicate",
			Email:    "existing@example.com",
			Password: "password123",
		}

		user, err := svc.Register(req)
		if err == nil {
			t.Fatal("expected error for duplicate email")
		}
		if user != nil {
			t.Error("expected nil user on error")
		}
		if err.Error() != "user with this email already exists" {
			t.Errorf("expected duplicate error message, got: %v", err)
		}
	})

	t.Run("repository GetByEmail error", func(t *testing.T) {
		repo := setupAuthTest()
		repo.Err = jwt.ErrInvalidKey
		svc := NewAuthService(repo)

		req := &domain.RegisterRequest{
			Username: "erroruser",
			Email:    "error@example.com",
			Password: "password123",
		}

		user, err := svc.Register(req)
		if err == nil {
			t.Fatal("expected error from repository")
		}
		if user != nil {
			t.Error("expected nil user on error")
		}
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		repo := setupAuthTest()
		svc := NewAuthService(repo)

		req := &domain.LoginRequest{
			Email:    "existing@example.com",
			Password: "password123",
		}

		token, user, err := svc.Login(req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected token to be returned")
		}
		if user == nil {
			t.Fatal("expected user to be returned")
		}
		if user.Email != "existing@example.com" {
			t.Errorf("expected email 'existing@example.com', got %s", user.Email)
		}
		if user.Username != "existing_user" {
			t.Errorf("expected username 'existing_user', got %s", user.Username)
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "your-secret-key-change-this-in-production"
		}
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !parsedToken.Valid {
			t.Errorf("token is not valid: %v", err)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		repo := setupAuthTest()
		svc := NewAuthService(repo)

		req := &domain.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		token, user, err := svc.Login(req)
		if err == nil {
			t.Fatal("expected error for invalid email")
		}
		if token != "" {
			t.Error("expected empty token on error")
		}
		if user != nil {
			t.Error("expected nil user on error")
		}
		if err.Error() != "invalid email or password" {
			t.Errorf("expected 'invalid email or password', got: %v", err)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		repo := setupAuthTest()
		svc := NewAuthService(repo)

		req := &domain.LoginRequest{
			Email:    "existing@example.com",
			Password: "wrongpassword",
		}

		token, user, err := svc.Login(req)
		if err == nil {
			t.Fatal("expected error for invalid password")
		}
		if token != "" {
			t.Error("expected empty token on error")
		}
		if user != nil {
			t.Error("expected nil user on error")
		}
		if err.Error() != "invalid email or password" {
			t.Errorf("expected 'invalid email or password', got: %v", err)
		}
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	repo := setupAuthTest()
	svc := NewAuthService(repo)

	t.Run("valid token", func(t *testing.T) {
		req := &domain.LoginRequest{
			Email:    "existing@example.com",
			Password: "password123",
		}
		token, _, err := svc.Login(req)
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}

		parsedToken, err := svc.ValidateToken(token)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if parsedToken == nil || !parsedToken.Valid {
			t.Error("expected valid token")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		_, err := svc.ValidateToken("invalid.token.here")
		if err == nil {
			t.Error("expected error for invalid token")
		}
	})

	t.Run("expired token", func(t *testing.T) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "your-secret-key-change-this-in-production"
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  1,
			"email":    "existing@example.com",
			"username": "existing_user",
			"exp":      time.Now().Add(-1 * time.Hour).Unix(),
			"iat":      time.Now().Add(-2 * time.Hour).Unix(),
		})

		tokenString, _ := token.SignedString([]byte(secret))

		parsedToken, err := svc.ValidateToken(tokenString)
		if err == nil {
			t.Error("expected error for expired token")
		}
		if parsedToken != nil && parsedToken.Valid {
			t.Error("expired token should not be valid")
		}
	})

	t.Run("malformed token", func(t *testing.T) {
		_, err := svc.ValidateToken("not.a.valid.jwt.token.at.all")
		if err == nil {
			t.Error("expected error for malformed token")
		}
	})
}

func TestAuthService_GetUserFromToken(t *testing.T) {
	repo := setupAuthTest()
	svc := NewAuthService(repo)

	t.Run("valid token returns user", func(t *testing.T) {
		req := &domain.LoginRequest{
			Email:    "existing@example.com",
			Password: "password123",
		}
		token, _, err := svc.Login(req)
		if err != nil {
			t.Fatalf("login failed: %v", err)
		}

		user, err := svc.GetUserFromToken(token)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user == nil {
			t.Fatal("expected user to be returned")
		}
		if user.ID != 1 {
			t.Errorf("expected user ID 1, got %d", user.ID)
		}
		if user.Email != "existing@example.com" {
			t.Errorf("expected email 'existing@example.com', got %s", user.Email)
		}
		if user.Username != "existing_user" {
			t.Errorf("expected username 'existing_user', got %s", user.Username)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		user, err := svc.GetUserFromToken("invalid.token.here")
		if err == nil {
			t.Error("expected error for invalid token")
		}
		if user != nil {
			t.Error("expected nil user on error")
		}
	})

	t.Run("user not found in database", func(t *testing.T) {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "your-secret-key-change-this-in-production"
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  999,
			"email":    "nonexistent@example.com",
			"username": "nonexistent",
			"exp":      time.Now().Add(24 * time.Hour).Unix(),
			"iat":      time.Now().Unix(),
		})

		tokenString, _ := token.SignedString([]byte(secret))

		user, err := svc.GetUserFromToken(tokenString)
		if err == nil {
			t.Error("expected error when user not found")
		}
		if user != nil {
			t.Error("expected nil user when not found")
		}
		if err.Error() != "user not found" {
			t.Errorf("expected 'user not found', got: %v", err)
		}
	})
}
