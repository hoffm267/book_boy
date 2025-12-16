package domain

import (
	"book_boy/api/internal/errors"
	"time"
)

type ProgressStatus string

const (
	ProgressStatusInProgress ProgressStatus = "in_progress"
	ProgressStatusCompleted  ProgressStatus = "completed"
)

type Progress struct {
	ID                int             `json:"id"`
	UserID            int             `json:"user_id"`
	BookID            *int            `json:"book_id,omitempty"`
	AudiobookID       *int            `json:"audiobook_id,omitempty"`
	BookPage          *int            `json:"book_page,omitempty" binding:"omitempty,min=1"`
	AudiobookTime     *CustomDuration `json:"audiobook_time,omitempty"`
	CompletionPercent int             `json:"completion_percent,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

func (p *Progress) Validate() error {
	if p.BookID == nil && p.AudiobookID == nil {
		return errors.ErrInvalidInput("progress must have at least a book_id or audiobook_id")
	}
	if p.BookPage != nil && *p.BookPage <= 0 {
		return errors.ErrInvalidInput("book_page must be greater than 0")
	}
	if p.AudiobookTime != nil && p.AudiobookTime.Duration < 0 {
		return errors.ErrInvalidInput("audiobook_time cannot be negative")
	}
	return nil
}

type EnrichedProgress struct {
	Progress          Progress
	Book              *Book
	Audiobook         *Audiobook
	TotalPages        int
	TotalLength       *CustomDuration
	CompletionPercent int
}
