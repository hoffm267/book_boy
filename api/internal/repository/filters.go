package repository

import "book_boy/api/internal/domain"

type BookFilter struct {
	ID         *int
	UserID     *int
	ISBN       *string
	Title      *string
	TotalPages *int
}

type ProgressFilter struct {
	ID          *int
	UserID      *int
	BookID      *int
	AudiobookID *int
	Status      *domain.ProgressStatus
}
