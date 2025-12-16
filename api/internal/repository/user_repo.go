package repository

import (
	"database/sql"

	"book_boy/api/internal/domain"
)

type UserRepo interface {
	GetAll() ([]domain.User, error)
	GetByID(id int) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Create(user *domain.User) (int, error)
	Update(user *domain.User) error
	Delete(id int) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetAll() ([]domain.User, error) {
	rows, err := r.db.Query("SELECT id, username, email, password_hash, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) GetByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE id = $1", id).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(user *domain.User) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Email, user.PasswordHash,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *userRepo) Update(user *domain.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET username = $1, email = $2 WHERE id = $3",
		user.Username, user.Email, user.ID,
	)
	return err
}

func (r *userRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
