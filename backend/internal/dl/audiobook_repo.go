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
	GetSimilarTitles(title string) ([]models.Audiobook, error)
}

type audiobookRepo struct {
	db *sql.DB
}

func NewAudiobookRepo(db *sql.DB) AudiobookRepo {
	return &audiobookRepo{db: db}
}

func (r *audiobookRepo) GetAll() ([]models.Audiobook, error) {
	rows, err := r.db.Query("SELECT id, title, total_length FROM audiobooks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audiobooks []models.Audiobook
	for rows.Next() {
		var audiobook models.Audiobook
		if err := rows.Scan(&audiobook.ID, &audiobook.Title, &audiobook.TotalLength); err != nil {
			return nil, err
		}
		audiobooks = append(audiobooks, audiobook)
	}
	return audiobooks, nil
}

func (r *audiobookRepo) GetByID(id int) (*models.Audiobook, error) {
	row := r.db.QueryRow("SELECT id, title, total_length FROM audiobooks WHERE id = $1", id)

	var audiobook models.Audiobook
	if err := row.Scan(&audiobook.ID, &audiobook.Title, &audiobook.TotalLength); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &audiobook, nil
}

func (r *audiobookRepo) Create(ab *models.Audiobook) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO audiobooks (title, length) VALUES ($1, $2) RETURNING id",
		ab.Title, ab.TotalLength,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *audiobookRepo) Update(audiobook *models.Audiobook) error {
	_, err := r.db.Exec(
		"UPDATE audiobooks SET title = $1, SET length = $2 WHERE id = $3",
		audiobook.Title, audiobook.TotalLength, audiobook.ID,
	)
	return err
}

func (r *audiobookRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM audiobooks WHERE id = $1", id)
	return err
}

func (r *audiobookRepo) GetSimilarTitles(title string) ([]models.Audiobook, error) {
	rows, err := r.db.Query("SELECT id, title, total_length FROM audiobooks WHERE title % $1 ORDER BY similarity(title, $1) DESC", title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audiobooks []models.Audiobook
	for rows.Next() {
		var audiobook models.Audiobook
		if err := rows.Scan(&audiobook.ID, &audiobook.Title, &audiobook.TotalLength); err != nil {
			return nil, err
		}
		audiobooks = append(audiobooks, audiobook)
	}
	return audiobooks, nil
}
