package dl

import (
	"book_boy/backend/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ProgressRepo interface {
	GetAll() ([]models.Progress, error)
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
