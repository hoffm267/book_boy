package service

import (
	"book_boy/api/internal/domain"
	"book_boy/api/internal/repository"
	"book_boy/api/internal/infra"
	"context"
	"fmt"
	"time"
)

type AudiobookService interface {
	GetAll() ([]domain.Audiobook, error)
	GetByID(id int) (*domain.Audiobook, error)
	Create(audiobook *domain.Audiobook) (int, error)
	Update(audiobook *domain.Audiobook) error
	GetSimilarTitles(title string) ([]domain.Audiobook, error)
	Delete(id int) error
}

type audiobookService struct {
	repo  repository.AudiobookRepo
	cache *infra.Cache
}

func NewAudiobookService(repo repository.AudiobookRepo, cache *infra.Cache) AudiobookService {
	return &audiobookService{repo: repo, cache: cache}
}

func (s *audiobookService) GetAll() ([]domain.Audiobook, error) {
	return s.repo.GetAll()
}

func (s *audiobookService) GetByID(id int) (*domain.Audiobook, error) {
	if s.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("audiobook:%d", id)

		var audiobook domain.Audiobook
		if err := s.cache.Get(ctx, cacheKey, &audiobook); err == nil {
			return &audiobook, nil
		}
	}

	result, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		s.cache.Set(context.Background(), fmt.Sprintf("audiobook:%d", id), result, 10*time.Minute)
	}
	return result, nil
}

func (s *audiobookService) Create(audiobook *domain.Audiobook) (int, error) {
	if err := audiobook.Validate(); err != nil {
		return 0, err
	}
	return s.repo.Create(audiobook)
}

func (s *audiobookService) Update(audiobook *domain.Audiobook) error {
	if err := audiobook.Validate(); err != nil {
		return err
	}
	if err := s.repo.Update(audiobook); err != nil {
		return err
	}
	if s.cache != nil {
		s.cache.Delete(context.Background(), fmt.Sprintf("audiobook:%d", audiobook.ID))
	}
	return nil
}

func (s *audiobookService) Delete(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	if s.cache != nil {
		s.cache.Delete(context.Background(), fmt.Sprintf("audiobook:%d", id))
	}
	return nil
}

func (s *audiobookService) GetSimilarTitles(title string) ([]domain.Audiobook, error) {
	return s.repo.GetSimilarTitles(title)
}
