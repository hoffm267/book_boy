package dl

import (
	"book_boy/backend/internal/models"
	"database/sql"
)

type AudiobookRepo interface {
	GetAll() ([]models.Audiobook, error)
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
