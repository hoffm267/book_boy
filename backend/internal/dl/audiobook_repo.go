package dl

import (
	"database/sql"

	"book_boy/backend/internal/models"
)

type AudiobookRepo interface {
	GetAll() ([]models.Audiobook, error)
	GetByID(id int) (*models.Audiobook, error)
	Create(ab *models.Audiobook) (int, error)
	Update(ab *models.Audiobook) error
	Delete(id int) error
}

type audiobookRepo struct {
	db *sql.DB
}

func NewAudiobookRepo(db *sql.DB) AudiobookRepo {
	return &audiobookRepo{db: db}
}

func (r *audiobookRepo) GetAll() ([]models.Audiobook, error) {
	rows, err := r.db.Query("SELECT id, isbn, title FROM audiobooks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audiobooks []models.Audiobook
	for rows.Next() {
		var ab models.Audiobook
		if err := rows.Scan(&ab.ID, &ab.ISBN, &ab.Title); err != nil {
			return nil, err
		}
		audiobooks = append(audiobooks, ab)
	}
	return audiobooks, nil
}

func (r *audiobookRepo) GetByID(id int) (*models.Audiobook, error) {
	row := r.db.QueryRow("SELECT id, isbn, title FROM audiobooks WHERE id = $1", id)

	var ab models.Audiobook
	if err := row.Scan(&ab.ID, &ab.ISBN, &ab.Title); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &ab, nil
}

func (r *audiobookRepo) Create(ab *models.Audiobook) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO audiobooks (isbn, title) VALUES ($1, $2) RETURNING id",
		ab.ISBN, ab.Title,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *audiobookRepo) Update(ab *models.Audiobook) error {
	_, err := r.db.Exec(
		"UPDATE audiobooks SET isbn = $1, title = $2 WHERE id = $3",
		ab.ISBN, ab.Title, ab.ID,
	)
	return err
}

func (r *audiobookRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM audiobooks WHERE id = $1", id)
	return err
}
