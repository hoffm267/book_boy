package models

import "time"

type Progress struct {
	ID            int             `json:"id"`
	UserID        int             `json:"user_id"`
	BookID        *int            `json:"book_id"`
	AudiobookID   *int            `json:"audiobook_id"`
	BookPage      *int            `json:"book_page"`
	AudiobookTime *CustomDuration `json:"audiobook_time"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type ProgressFilter struct {
	ID          *int
	UserID      *int
	BookID      *int
	AudiobookID *int
}
