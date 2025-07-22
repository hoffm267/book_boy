package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type ProgressService interface {
	GetAll() ([]models.Progress, error)
}

type progressService struct {
	repo dl.ProgressRepo
}

func NewProgressService(repo dl.ProgressRepo) ProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) GetAll() ([]models.Progress, error) {
	return s.repo.GetAll()
}
