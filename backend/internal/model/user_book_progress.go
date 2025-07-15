package model

import "time"

type User_book_progress struct {
	ID             int           `json:"id"`
	User_ID        int           `json:"user_id"`
	Book_ID        int           `json:"book_id"`
	Audiobook_ID   int           `json:"audiobook_id"`
	Book_Page      int           `json:"book_page"`
	Audiobook_Time time.Duration `json:"audiobook_time"`
	Created_At     time.Time     `json:"created_at"`
	Updated_At     time.Time     `json:"updated_at"`
}
