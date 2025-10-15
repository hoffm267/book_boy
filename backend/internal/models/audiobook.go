package models

type Audiobook struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	TotalLength *CustomDuration `json:"total_length"`
}
