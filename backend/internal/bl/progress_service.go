package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
	"fmt"
)

type ProgressService interface {
	GetAll() ([]models.Progress, error)
	GetByID(id int) (*models.Progress, error)
	Create(progress *models.Progress) (int, error)
	Update(progress *models.Progress) error
	Delete(id int) error
	UpdateProgressPage(id, bookPage int) error
	UpdateProgressTime(id int, audiobookTime *models.CustomDuration) error
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

// Extensions
func (s *progressService) UpdateProgressPage(id, bookPage int) error {
	progress, totalPages, totalLength, err := s.repo.GetByIDWithTotals(id)
	if err != nil {
		return err
	}
	if progress == nil || progress.BookID == nil || totalPages <= 0 {
		return fmt.Errorf("book info missing for conversion")
	}

	if bookPage < 1 {
		bookPage = 1
	}
	if bookPage > totalPages {
		bookPage = totalPages
	}
	progress.BookPage = &bookPage

	// if audiobook total length available, compute & set AudiobookTime
	if totalLength != nil && totalLength.Duration > 0 {
		ts, _ := pageToTimestamp(totalPages, bookPage, totalLength.Duration)
		cd := models.CustomDuration{Duration: ts}
		progress.AudiobookTime = &cd
	}

	return s.repo.Update(progress)
}

func (s *progressService) UpdateProgressTime(progressID int, audiobookTime *models.CustomDuration) error {
	pr, totalPages, totalLength, err := s.repo.GetByIDWithTotals(progressID)
	if err != nil {
		return err
	}
	if pr == nil {
		return fmt.Errorf("book info missing for conversion")
	}

	pr.AudiobookTime = audiobookTime

	// if we have book total pages and audio total length, compute page
	if pr.BookID != nil && totalPages > 0 && totalLength != nil && totalLength.Duration > 0 {
		page, _ := timestampToPage(totalPages, audiobookTime.Duration, totalLength.Duration)
		pr.BookPage = &page
	}

	return s.repo.Update(pr)
}
