package dl

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"book_boy/backend/internal/models"
)

type ProgressRepo interface {
	GetAll() ([]models.Progress, error)
	GetByID(id int) (*models.Progress, error)
	Create(progress *models.Progress) error
	Update(progress *models.Progress) error
	Delete(id int) error
}

type progressRepo struct {
	db *sql.DB
}

func NewProgressRepo(db *sql.DB) ProgressRepo {
	return &progressRepo{db: db}
}

// TODO change JSON return (probably on controller) to be HH:MM:SS format from nanoseconds
func parsePgInterval(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("unexpected interval format: %s", s)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	seconds, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, err
	}

	duration := time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds*float64(time.Second))

	return duration, nil
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

	var progress []models.Progress
	for rows.Next() {
		var p models.Progress
		var audiobookTime sql.NullString
		err := rows.Scan(
			&p.ID, &p.UserID, &p.BookID, &p.AudiobookID,
			&p.BookPage, &audiobookTime, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if audiobookTime.Valid {
			dur, err := parsePgInterval(audiobookTime.String)
			if err != nil {
				return nil, err
			}
			p.AudiobookTime = &dur
		} else {
			p.AudiobookTime = nil
		}

		progress = append(progress, p)
	}

	return progress, nil
}

func (r *progressRepo) GetByID(id int) (*models.Progress, error) {
	row := r.db.QueryRow(`
		SELECT id, user_id, book_id, audiobook_id, book_page, audiobook_time, created_at, updated_at
		FROM progress WHERE id = $1
	`, id)

	var p models.Progress
	var audiobookTime sql.NullString
	err := row.Scan(
		&p.ID, &p.UserID, &p.BookID, &p.AudiobookID,
		&p.BookPage, &audiobookTime, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if audiobookTime.Valid {
		dur, err := parsePgInterval(audiobookTime.String)
		if err != nil {
			return nil, err
		}
		p.AudiobookTime = &dur
	} else {
		p.AudiobookTime = nil
	}
	return &p, nil
}

func (r *progressRepo) Create(p *models.Progress) error {
	_, err := r.db.Exec(`
		INSERT INTO progress (user_id, book_id, audiobook_id, book_page, audiobook_time)
		VALUES ($1, $2, $3, $4, $5)
	`, p.UserID, p.BookID, p.AudiobookID, p.BookPage, p.AudiobookTime)
	return err
}

func (r *progressRepo) Update(p *models.Progress) error {
	_, err := r.db.Exec(`
		UPDATE progress
		SET user_id = $1, book_id = $2, audiobook_id = $3, book_page = $4, audiobook_time = $5, updated_at = NOW()
		WHERE id = $6
	`, p.UserID, p.BookID, p.AudiobookID, p.BookPage, p.AudiobookTime, p.ID)
	return err
}

func (r *progressRepo) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM progress WHERE id = $1", id)
	return err
}
