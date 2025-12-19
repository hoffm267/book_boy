package domain

type BookCreatedEvent struct {
	BookID    int    `json:"book_id"`
	ISBN      string `json:"isbn"`
	CreatedAt string `json:"created_at"`
}

type BookMetadataFetchedEvent struct {
	BookID     int    `json:"book_id"`
	ISBN       string `json:"isbn"`
	Title      string `json:"title"`
	TotalPages int    `json:"total_pages"`
	Author     string `json:"author,omitempty"`
	Publisher  string `json:"publisher,omitempty"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}
