package book_metadata

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type SearchResult struct {
	Books []BookMetadata `json:"books"`
	Total int            `json:"total"`
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

func (c *Client) SearchBooks(query string, limit int) (*SearchResult, error) {
	endpoint := fmt.Sprintf("%s/books/search", c.baseURL)

	params := url.Values{}
	params.Add("q", query)
	params.Add("limit", fmt.Sprintf("%d", limit))

	url := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to search books: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("metadata service returned %d: %s", resp.StatusCode, string(body))
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
