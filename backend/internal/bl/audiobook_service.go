package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
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
	repo dl.AudiobookRepo
}

func NewAudiobookService(repo dl.AudiobookRepo) AudiobookService {
	return &audiobookService{repo: repo}
}

func (s *audiobookService) GetAll() ([]models.Audiobook, error) {
	return s.repo.GetAll()
}

func (s *audiobookService) GetByID(id int) (*models.Audiobook, error) {
	return s.repo.GetByID(id)
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
	return s.repo.Update(audiobook)
}

func (s *audiobookService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *audiobookService) GetSimilarTitles(title string) ([]models.Audiobook, error) {
	return s.repo.GetSimilarTitles(title)
}
