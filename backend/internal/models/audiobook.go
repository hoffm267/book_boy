package models

type Audiobook struct {
	ID          int             `json:"id"`
	Title       string          `json:"title" binding:"required,min=1,max=500"`
	TotalLength *CustomDuration `json:"total_length" binding:"required"`
}

func (a *Audiobook) Validate() error {
	if a.Title == "" {
		return ErrInvalidInput("title cannot be empty")
	}
	if len(a.Title) > 500 {
		return ErrInvalidInput("title cannot exceed 500 characters")
	}
	if a.TotalLength == nil || a.TotalLength.Duration <= 0 {
		return ErrInvalidInput("total_length must be greater than 0")
	}
	return nil
}
