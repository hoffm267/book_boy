package dl

import (
	"database/sql"

	"book_boy/backend/internal/models"
)

type BookRepo interface {
	GetAll() ([]models.Book, error)
	GetByID(id int) (*models.Book, error)
	Create(book *models.Book) (int, error)
	Update(book *models.Book) error
	Delete(id int) error
	GetByTitle(title string) (*models.Book, error)
}

type bookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) BookRepo {
	return &bookRepo{db: db}
}

// CRUD
func (r *bookRepo) GetAll() ([]models.Book, error) {
	//TODO find out if this is worth doing
	rows, err := r.db.Query("SELECT id, isbn, title, total_pages FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.ISBN, &book.Title, &book.TotalPages); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *bookRepo) GetByID(id int) (*models.Book, error) {
	row := r.db.QueryRow("SELECT id, isbn, title, total_pages FROM books WHERE id = $1", id)

	var book models.Book
	err := row.Scan(&book.ID, &book.ISBN, &book.Title, &book.TotalPages)
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

func (r *bookRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM books WHERE id = $1", id)
	return err
}

// Extension
func (r *bookRepo) GetByTitle(title string) (*models.Book, error) {
	row := r.db.QueryRow("SELECT id, isbn, title, total_pages FROM books WHERE title = $1", title)

	var book models.Book
	err := row.Scan(&book.ID, &book.ISBN, &book.Title, &book.TotalPages)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 404 not found?
		}
		return nil, err
	}

	return &book, nil
}
