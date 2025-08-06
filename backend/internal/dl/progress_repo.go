package dl

import (
	"database/sql"

	"book_boy/backend/internal/models"
)

type ProgressRepo interface {
	GetAll() ([]models.Progress, error)
	GetByID(id int) (*models.Progress, error)
	Create(progress *models.Progress) (int, error)
	Update(progress *models.Progress) error
	Delete(id int) error
	GetByIDWithTotals(id int) (*models.Progress, int, *models.CustomDuration, error)
}

type progressRepo struct {
	db *sql.DB
}

func NewProgressRepo(db *sql.DB) ProgressRepo {
	return &progressRepo{db: db}
}

func (r *progressRepo) GetAll() ([]models.Progress, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, book_id, audiobook_id, book_page, audiobook_time, created_at, updated_at
		FROM progress
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var progresses []models.Progress
	for rows.Next() {
		var progress models.Progress
		err := rows.Scan(
			&progress.ID, &progress.UserID, &progress.BookID, &progress.AudiobookID,
			&progress.BookPage, &progress.AudiobookTime, &progress.CreatedAt, &progress.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		progresses = append(progresses, progress)
	}

	return progresses, nil
}

func (r *progressRepo) GetByID(id int) (*models.Progress, error) {
	row := r.db.QueryRow(`
		SELECT id, user_id, book_id, audiobook_id, book_page, audiobook_time, created_at, updated_at
		FROM progress WHERE id = $1
	`, id)

	var p models.Progress
	err := row.Scan(
		&p.ID, &p.UserID, &p.BookID, &p.AudiobookID,
		&p.BookPage, &p.AudiobookTime, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *progressRepo) Create(progress *models.Progress) (int, error) {
	var id int
	err := r.db.QueryRow(`
		INSERT INTO progress (user_id, book_id, audiobook_id, book_page, audiobook_time)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, progress.UserID, progress.BookID, progress.AudiobookID, progress.BookPage, progress.AudiobookTime).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *progressRepo) Update(progress *models.Progress) error {
	_, err := r.db.Exec(`
		UPDATE progress
		SET user_id = $1, book_id = $2, audiobook_id = $3, book_page = $4, audiobook_time = $5, updated_at = NOW()
		WHERE id = $6
	`, progress.UserID, progress.BookID, progress.AudiobookID, progress.BookPage, progress.AudiobookTime, progress.ID)
	return err
}

func (r *progressRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM progress WHERE id = $1", id)
	return err
}

// Extensions
func (r *progressRepo) GetByIDWithTotals(id int) (*models.Progress, int, *models.CustomDuration, error) {
	query := `
    SELECT
    	p.id, p.user_id, p.book_id, p.audiobook_id, p.book_page, p.audiobook_time, p.created_at, p.updated_at,
    	COALESCE(b.total_pages, 0),
    	a.total_length
    FROM progress p
    LEFT JOIN books b ON p.book_id = b.id
    LEFT JOIN audiobooks a ON p.audiobook_id = a.id
    WHERE p.id = $1
    `
	var pr models.Progress
	var totalPages int
	var totalLength *models.CustomDuration

	err := r.db.QueryRow(query, id).Scan(
		&pr.ID, &pr.UserID, &pr.BookID, &pr.AudiobookID,
		&pr.BookPage, &pr.AudiobookTime, &pr.CreatedAt, &pr.UpdatedAt,
		&totalPages,
		&totalLength,
	)
	if err != nil {
		return nil, 0, nil, err
	}
	return &pr, totalPages, totalLength, nil
}
