package dl

import (
	"database/sql"

	"book_boy/backend/internal/models"
)

type BookRepo interface {
	GetAll() ([]models.Book, error)
	GetByID(bookID int) (*models.Book, error)
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

func (r *bookRepo) GetByID(bookID int) (*models.Book, error) {
	row := r.db.QueryRow("SELECT id, isbn, title FROM books WHERE id = $1", bookID)

	var book models.Book
	err := row.Scan(&book.ID, &book.ISBN, &book.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 404 not found?
		}
		return nil, err
	}

	return &book, nil
}
