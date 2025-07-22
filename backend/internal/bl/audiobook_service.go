package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type AudiobookService interface {
	GetAll() ([]models.Audiobook, error)
}

type audiobookService struct {
	repo dl.AudiobookRepo
}

func NewAudiobookService(repo dl.AudiobookRepo) AudiobookService {
	return &audiobookService{repo: repo}
}

func (s *audiobookService) GetAll() ([]models.Audiobook, error) {
	return s.repo.GetAll()
}
