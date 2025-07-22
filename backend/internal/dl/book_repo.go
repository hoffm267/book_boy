package dl

import (
	"book_boy/backend/internal/models"
	"database/sql"
)

type BookRepo interface {
	GetAll() ([]models.Book, error)
}

type bookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) BookRepo {
	return &bookRepo{db: db}
}

func (r *bookRepo) GetAll() ([]models.Book, error) {
	rows, err := r.db.Query("SELECT id, isbn, title FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.ISBN, &b.Title); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}
