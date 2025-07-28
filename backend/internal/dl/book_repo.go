package dl

import (
	"database/sql"

	"book_boy/backend/internal/models"
)

type BookRepo interface {
	GetAll() ([]models.Book, error)
	GetByID(bookID int) (*models.Book, error)
	Create(book *models.Book) (int, error)
	Update(book *models.Book) error
	Delete(bookID int) error
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

func (r *bookRepo) Create(book *models.Book) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO books (isbn, title) VALUES ($1, $2) RETURNING id",
		book.ISBN, book.Title,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *bookRepo) Update(book *models.Book) error {
	_, err := r.db.Exec(
		"UPDATE books SET isbn = $1, title = $2 WHERE id = $3",
		book.ISBN, book.Title, book.ID,
	)
	return err
}

func (r *bookRepo) Delete(bookID int) error {
	_, err := r.db.Exec("DELETE FROM books WHERE id = $1", bookID)
	return err
}
