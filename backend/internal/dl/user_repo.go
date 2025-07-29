package dl

import (
	"database/sql"

	"book_boy/backend/internal/models"
)

type UserRepo interface {
	GetAll() ([]models.User, error)
	GetByID(id int) (*models.User, error)
	Create(user *models.User) (int, error)
	Update(user *models.User) error
	Delete(id int) error
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) GetAll() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, username FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepo) GetByID(id int) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow("SELECT id, username FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) Create(user *models.User) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO users (username) VALUES ($1) RETURNING id",
		user.Username,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *userRepo) Update(user *models.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET username = $1 WHERE id = $2",
		user.Username, user.ID,
	)
	return err
}

func (r *userRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}
