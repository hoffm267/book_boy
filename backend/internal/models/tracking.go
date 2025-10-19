package models

import "time"

type StartTrackingRequest struct {
	Format      string `json:"format" binding:"required,oneof=book audiobook"`
	Title       string `json:"title" binding:"required,min=1,max=500"`
	Author      string `json:"author"`
	TotalPages  int    `json:"total_pages"`
	TotalLength string `json:"total_length"`
	ISBN        string `json:"isbn"`
	CurrentPage int    `json:"current_page"`
	CurrentTime string `json:"current_time"`
}

func (r *StartTrackingRequest) Validate() error {
	if r.Format == "book" && r.TotalPages <= 0 {
		return ErrInvalidInput("total_pages is required and must be > 0 for books")
	}
	if r.Format == "audiobook" && r.TotalLength == "" {
		return ErrInvalidInput("total_length is required for audiobooks")
	}
	return nil
}

type CurrentTrackingResponse struct {
	ProgressID        int             `json:"progress_id"`
	Book              *Book           `json:"book"`
	Audiobook         *Audiobook      `json:"audiobook"`
	CurrentPage       *int            `json:"current_page"`
	CurrentTime       *CustomDuration `json:"current_time"`
	CompletionPercent int             `json:"completion_percent"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
