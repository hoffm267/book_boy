package models

// Keeping seperate from book because more fields will be added later, i.e. duration or narrator
type Audiobook struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	TotalLength *CustomDuration `json:"total_length"`
}
