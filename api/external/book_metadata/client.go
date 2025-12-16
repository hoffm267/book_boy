package book_metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type BookMetadata struct {
	ISBN        string  `json:"isbn"`
	Title       *string `json:"title"`
	Author      *string `json:"author"`
	TotalPages  *int    `json:"total_pages"`
	CoverURL    *string `json:"cover_url"`
	Publisher   *string `json:"publisher"`
	PublishDate *string `json:"publish_date"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetBookMetadataByISBN(isbn string) (*BookMetadata, error) {
	url := fmt.Sprintf("%s/books/isbn/%s", c.baseURL, isbn)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch book: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("book not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("metadata service returned %d: %s", resp.StatusCode, string(body))
	}

	var book BookMetadata
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &book, nil
}
