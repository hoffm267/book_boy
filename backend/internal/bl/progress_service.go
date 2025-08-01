package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type ProgressService interface {
	GetAll() ([]models.Progress, error)
	GetByID(id int) (*models.Progress, error)
	Create(progress *models.Progress) (int, error)
	Update(progress *models.Progress) error
	Delete(id int) error
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

func (s *progressService) GetByID(id int) (*models.Progress, error) {
	return s.repo.GetByID(id)
}

func (s *progressService) Create(progress *models.Progress) (int, error) {
	return s.repo.Create(progress)
}

func (s *progressService) Update(progress *models.Progress) error {
	return s.repo.Update(progress)
}

func (s *progressService) Delete(id int) error {
	return s.repo.Delete(id)
}
