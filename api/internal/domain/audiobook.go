package domain

import "book_boy/api/internal/errors"

type Audiobook struct {
	ID          int             `json:"id"`
	UserID      int             `json:"user_id"`
	Title       string          `json:"title" binding:"required,min=1,max=500"`
	TotalLength *CustomDuration `json:"total_length" binding:"required"`
}

func (a *Audiobook) Validate() error {
	if a.Title == "" {
		return errors.ErrInvalidInput("title cannot be empty")
	}
	if len(a.Title) > 500 {
		return errors.ErrInvalidInput("title cannot exceed 500 characters")
	}
	if a.TotalLength == nil || a.TotalLength.Duration <= 0 {
		return errors.ErrInvalidInput("total_length must be greater than 0")
	}
	return nil
}
