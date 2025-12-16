package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"book_boy/api/internal/domain"
)

type ProgressRepo interface {
	GetAll() ([]domain.Progress, error)
	GetByID(id int) (*domain.Progress, error)
	Create(progress *domain.Progress) (int, error)
	Update(progress *domain.Progress) error
	Delete(id int) error
	GetByIDWithTotals(id int) (*domain.Progress, int, *domain.CustomDuration, error)
	FilterProgress(filter ProgressFilter) ([]domain.Progress, error)
	GetAllEnrichedByUser(userID int) ([]domain.EnrichedProgress, error)
}

type progressRepo struct {
	db *sql.DB
}

func NewProgressRepo(db *sql.DB) ProgressRepo {
	return &progressRepo{db: db}
}

func (r *progressRepo) GetAll() ([]domain.Progress, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, book_id, audiobook_id, book_page, audiobook_time, created_at, updated_at
		FROM progress
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var progresses []domain.Progress
	for rows.Next() {
		var progress domain.Progress
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

func (r *progressRepo) GetByID(id int) (*domain.Progress, error) {
	row := r.db.QueryRow(`
		SELECT id, user_id, book_id, audiobook_id, book_page, audiobook_time, created_at, updated_at
		FROM progress WHERE id = $1
	`, id)

	var p domain.Progress
	err := row.Scan(
		&p.ID, &p.UserID, &p.BookID, &p.AudiobookID,
		&p.BookPage, &p.AudiobookTime, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *progressRepo) Create(progress *domain.Progress) (int, error) {
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

func (r *progressRepo) Update(progress *domain.Progress) error {
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

func (r *progressRepo) GetByIDWithTotals(id int) (*domain.Progress, int, *domain.CustomDuration, error) {
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
	var pr domain.Progress
	var totalPages int
	var totalLength *domain.CustomDuration

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

func (r *progressRepo) FilterProgress(filter ProgressFilter) ([]domain.Progress, error) {
	query := "SELECT id, user_id, book_id, audiobook_id, book_page, audiobook_time, created_at, updated_at FROM progress"
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.ID != nil {
		conditions = append(conditions, fmt.Sprintf("id = $%d", argIndex))
		args = append(args, *filter.ID)
		argIndex++
	}
	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, *filter.UserID)
		argIndex++
	}
	if filter.BookID != nil {
		conditions = append(conditions, fmt.Sprintf("book_id = $%d", argIndex))
		args = append(args, *filter.BookID)
		argIndex++
	}
	if filter.AudiobookID != nil {
		conditions = append(conditions, fmt.Sprintf("audiobook_id = $%d", argIndex))
		args = append(args, *filter.AudiobookID)
		argIndex++
	}
	if filter.Status != nil {
		if *filter.Status == domain.ProgressStatusInProgress {
			conditions = append(conditions, "(book_id IS NOT NULL OR audiobook_id IS NOT NULL)")
		} else if *filter.Status == domain.ProgressStatusCompleted {
			conditions = append(conditions, "FALSE")
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progresses []domain.Progress
	for rows.Next() {
		var progress domain.Progress
		if err := rows.Scan(&progress.ID, &progress.UserID, &progress.BookID, &progress.AudiobookID, &progress.BookPage, &progress.AudiobookTime, &progress.CreatedAt, &progress.UpdatedAt); err != nil {
			return nil, err
		}
		progresses = append(progresses, progress)
	}
	return progresses, nil
}

func (r *progressRepo) GetAllEnrichedByUser(userID int) ([]domain.EnrichedProgress, error) {
	query := `
		SELECT
			p.id, p.user_id, p.book_id, p.audiobook_id, p.book_page, p.audiobook_time, p.created_at, p.updated_at,
			b.id, b.isbn, b.title, b.total_pages,
			a.id, a.title, a.total_length
		FROM progress p
		LEFT JOIN books b ON p.book_id = b.id
		LEFT JOIN audiobooks a ON p.audiobook_id = a.id
		WHERE p.user_id = $1
		ORDER BY p.updated_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.EnrichedProgress
	for rows.Next() {
		var e domain.EnrichedProgress
		var bookID, bookTotalPages *int
		var bookISBN, bookTitle *string
		var audiobookID *int
		var audiobookTitle *string
		var audiobookTotalLength *domain.CustomDuration

		err := rows.Scan(
			&e.Progress.ID, &e.Progress.UserID, &e.Progress.BookID, &e.Progress.AudiobookID,
			&e.Progress.BookPage, &e.Progress.AudiobookTime, &e.Progress.CreatedAt, &e.Progress.UpdatedAt,
			&bookID, &bookISBN, &bookTitle, &bookTotalPages,
			&audiobookID, &audiobookTitle, &audiobookTotalLength,
		)
		if err != nil {
			return nil, err
		}

		if bookID != nil {
			isbn := ""
			if bookISBN != nil {
				isbn = *bookISBN
			}
			title := ""
			if bookTitle != nil {
				title = *bookTitle
			}
			totalPages := 0
			if bookTotalPages != nil {
				totalPages = *bookTotalPages
			}
			e.Book = &domain.Book{
				ID:         *bookID,
				ISBN:       isbn,
				Title:      title,
				TotalPages: totalPages,
			}
			e.TotalPages = totalPages
		}

		if audiobookID != nil {
			title := ""
			if audiobookTitle != nil {
				title = *audiobookTitle
			}
			e.Audiobook = &domain.Audiobook{
				ID:          *audiobookID,
				Title:       title,
				TotalLength: audiobookTotalLength,
			}
			e.TotalLength = audiobookTotalLength
		}

		if e.Progress.BookPage != nil && e.TotalPages > 0 {
			percent := float64(*e.Progress.BookPage) / float64(e.TotalPages) * 100
			if percent > 100 {
				percent = 100
			}
			e.CompletionPercent = int(percent)
		} else if e.Progress.AudiobookTime != nil && e.TotalLength != nil && e.TotalLength.Duration > 0 {
			percent := e.Progress.AudiobookTime.Duration.Seconds() / e.TotalLength.Duration.Seconds() * 100
			if percent > 100 {
				percent = 100
			}
			e.CompletionPercent = int(percent)
		}

		results = append(results, e)
	}

	return results, nil
}
