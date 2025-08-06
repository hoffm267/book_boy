package models

type Book struct {
	ID         int    `json:"id"`
	ISBN       string `json:"isbn"`
	Title      string `json:"title"`
	TotalPages int    `json:"total_pages"`
}

type BookFilter struct {
	ID         *int
	ISBN       *string
	Title      *string
	TotalPages *int
}
