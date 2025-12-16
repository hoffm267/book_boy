package repository

import (
	"database/sql"

	"book_boy/api/internal/domain"
)

type AudiobookRepo interface {
	GetAll() ([]domain.Audiobook, error)
	GetByID(id int) (*domain.Audiobook, error)
	Create(audiobook *domain.Audiobook) (int, error)
	Update(audiobook *domain.Audiobook) error
	Delete(id int) error
	GetSimilarTitles(title string) ([]domain.Audiobook, error)
}

type audiobookRepo struct {
	db *sql.DB
}

func NewAudiobookRepo(db *sql.DB) AudiobookRepo {
	return &audiobookRepo{db: db}
}

func (r *audiobookRepo) GetAll() ([]domain.Audiobook, error) {
	rows, err := r.db.Query("SELECT id, user_id, title, total_length FROM audiobooks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audiobooks []domain.Audiobook
	for rows.Next() {
		var audiobook domain.Audiobook
		if err := rows.Scan(&audiobook.ID, &audiobook.UserID, &audiobook.Title, &audiobook.TotalLength); err != nil {
			return nil, err
		}
		audiobooks = append(audiobooks, audiobook)
	}
	return audiobooks, nil
}

func (r *audiobookRepo) GetByID(id int) (*domain.Audiobook, error) {
	row := r.db.QueryRow("SELECT id, user_id, title, total_length FROM audiobooks WHERE id = $1", id)

	var audiobook domain.Audiobook
	if err := row.Scan(&audiobook.ID, &audiobook.UserID, &audiobook.Title, &audiobook.TotalLength); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &audiobook, nil
}

func (r *audiobookRepo) Create(audiobook *domain.Audiobook) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO audiobooks (user_id, title, total_length) VALUES ($1, $2, $3) RETURNING id",
		audiobook.UserID, audiobook.Title, audiobook.TotalLength,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *audiobookRepo) Update(audiobook *domain.Audiobook) error {
	_, err := r.db.Exec(
		"UPDATE audiobooks SET title = $1, total_length = $2 WHERE id = $3 AND user_id = $4",
		audiobook.Title, audiobook.TotalLength, audiobook.ID, audiobook.UserID,
	)
	return err
}

func (r *audiobookRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM audiobooks WHERE id = $1", id)
	return err
}

func (r *audiobookRepo) GetSimilarTitles(title string) ([]domain.Audiobook, error) {
	rows, err := r.db.Query("SELECT id, user_id, title, total_length FROM audiobooks WHERE title % $1 ORDER BY similarity(title, $1) DESC", title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audiobooks []domain.Audiobook
	for rows.Next() {
		var audiobook domain.Audiobook
		if err := rows.Scan(&audiobook.ID, &audiobook.UserID, &audiobook.Title, &audiobook.TotalLength); err != nil {
			return nil, err
		}
		audiobooks = append(audiobooks, audiobook)
	}
	return audiobooks, nil
}
