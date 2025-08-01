package models

type Book struct {
	ID         int    `json:"id"`
	ISBN       string `json:"isbn"`
	Title      string `json:"title"`
	TotalPages int    `json:"total_pages"`
}
