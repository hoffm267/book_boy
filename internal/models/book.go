package models

type Book struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	ISBN       string `json:"isbn" binding:"required"`
	Title      string `json:"title" binding:"omitempty,min=1,max=500"`
	TotalPages int    `json:"total_pages" binding:"omitempty,min=1"`
}

func (b *Book) Validate() error {
	if b.Title == "" {
		return ErrInvalidInput("title cannot be empty")
	}
	if len(b.Title) > 500 {
		return ErrInvalidInput("title cannot exceed 500 characters")
	}
	if b.TotalPages <= 0 {
		return ErrInvalidInput("total_pages must be greater than 0")
	}
	return nil
}

type BookFilter struct {
	ID         *int
	UserID     *int
	ISBN       *string
	Title      *string
	TotalPages *int
}
