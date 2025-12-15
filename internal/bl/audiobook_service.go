package bl

import (
	"book_boy/internal/dl"
	"book_boy/internal/infra"
	"book_boy/internal/models"
	"context"
	"fmt"
	"time"
)

type AudiobookService interface {
	GetAll() ([]models.Audiobook, error)
	GetByID(id int) (*models.Audiobook, error)
	Create(audiobook *models.Audiobook) (int, error)
	Update(audiobook *models.Audiobook) error
	GetSimilarTitles(title string) ([]models.Audiobook, error)
	Delete(id int) error
}

type audiobookService struct {
	repo  dl.AudiobookRepo
	cache *infra.Cache
}

func NewAudiobookService(repo dl.AudiobookRepo, cache *infra.Cache) AudiobookService {
	return &audiobookService{repo: repo, cache: cache}
}

func (s *audiobookService) GetAll() ([]models.Audiobook, error) {
	return s.repo.GetAll()
}

func (s *audiobookService) GetByID(id int) (*models.Audiobook, error) {
	if s.cache != nil {
		ctx := context.Background()
		cacheKey := fmt.Sprintf("audiobook:%d", id)

		var audiobook models.Audiobook
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

func (s *audiobookService) Create(audiobook *models.Audiobook) (int, error) {
	if err := audiobook.Validate(); err != nil {
		return 0, err
	}
	return s.repo.Create(audiobook)
}

func (s *audiobookService) Update(audiobook *models.Audiobook) error {
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

func (s *audiobookService) GetSimilarTitles(title string) ([]models.Audiobook, error) {
	return s.repo.GetSimilarTitles(title)
}
