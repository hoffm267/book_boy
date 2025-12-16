package service

import (
	"book_boy/api/external/book_metadata"
	"book_boy/api/internal/infra"
	"context"
	"fmt"
	"time"
)

type BookMetadataService struct {
	client *book_metadata.Client
	cache  *infra.Cache
}

func NewBookMetadataService(client *book_metadata.Client, cache *infra.Cache) *BookMetadataService {
	return &BookMetadataService{
		client: client,
		cache:  cache,
	}
}

func (s *BookMetadataService) GetBookByISBN(isbn string) (*book_metadata.BookMetadata, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("isbn:%s", isbn)

	var result book_metadata.BookMetadata
	err := s.cache.Get(ctx, cacheKey, &result)
	if err == nil {
		return &result, nil
	}

	result2, err := s.client.GetBookMetadataByISBN(isbn)
	if err != nil {
		return nil, err
	}

	s.cache.Set(ctx, cacheKey, result2, 24*time.Hour)
	return result2, nil
}
