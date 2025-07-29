package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
)

type ProgressService interface {
	GetAll() ([]models.Progress, error)
	GetByID(id int) (*models.Progress, error)
	Create(p *models.Progress) (int, error)
	Update(p *models.Progress) error
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

func (s *progressService) Create(p *models.Progress) (int, error) {
	return s.repo.Create(p)
}

func (s *progressService) Update(p *models.Progress) error {
	return s.repo.Update(p)
}

func (s *progressService) Delete(id int) error {
	return s.repo.Delete(id)
}
