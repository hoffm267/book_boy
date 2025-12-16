package domain

import "book_boy/api/internal/errors"

type Book struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	ISBN       string `json:"isbn" binding:"required"`
	Title      string `json:"title" binding:"omitempty,min=1,max=500"`
	TotalPages int    `json:"total_pages" binding:"omitempty,min=1"`
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return errors.ErrInvalidInput("title cannot be empty")
	}
	if len(b.Title) > 500 {
		return errors.ErrInvalidInput("title cannot exceed 500 characters")
	}
	if b.TotalPages <= 0 {
		return errors.ErrInvalidInput("total_pages must be greater than 0")
	}
	return nil
}
