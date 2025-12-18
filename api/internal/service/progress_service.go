package service

import (
	"book_boy/api/internal/domain"
	"book_boy/api/internal/repository"
	"fmt"
)

type ProgressService interface {
	GetAll() ([]domain.Progress, error)
	GetByID(id int) (*domain.Progress, error)
	GetByIDWithCompletion(id int) (*domain.Progress, error)
	Create(progress *domain.Progress) (int, error)
	Update(progress *domain.Progress) error
	Delete(id int) error
	UpdateProgressPage(id int, bookPage int) error
	UpdateProgressTime(id int, audiobookTime *domain.CustomDuration) error
	SetBook(id int, bookID int) error
	SetAudiobook(id int, audiobookID int) error
	FilterProgress(filter repository.ProgressFilter) ([]domain.Progress, error)
}

type progressService struct {
	repo repository.ProgressRepo
}

func NewProgressService(repo repository.ProgressRepo) ProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) GetAll() ([]domain.Progress, error) {
	return s.repo.GetAll()
}

func (s *progressService) GetByID(id int) (*domain.Progress, error) {
	return s.repo.GetByID(id)
}

func (s *progressService) GetByIDWithCompletion(id int) (*domain.Progress, error) {
	progress, totalPages, totalLength, err := s.repo.GetByIDWithTotals(id)
	if err != nil {
		return nil, err
	}
	progress.CompletionPercent = calculateCompletionPercent(progress, totalPages, totalLength)
	return progress, nil
}

func (s *progressService) Create(progress *domain.Progress) (int, error) {
	if err := progress.Validate(); err != nil {
		return 0, err
	}
	return s.repo.Create(progress)
}

func (s *progressService) Update(progress *domain.Progress) error {
	if err := progress.Validate(); err != nil {
		return err
	}
	return s.repo.Update(progress)
}

func (s *progressService) Delete(id int) error {
	return s.repo.Delete(id)
}

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

	if totalLength != nil && totalLength.Duration > 0 {
		ts, _ := pageToTimestamp(totalPages, bookPage, totalLength.Duration)
		cd := domain.CustomDuration{Duration: ts}
		progress.AudiobookTime = &cd
	}

	return s.repo.Update(progress)
}

func (s *progressService) UpdateProgressTime(progressID int, audiobookTime *domain.CustomDuration) error {
	pr, totalPages, totalLength, err := s.repo.GetByIDWithTotals(progressID)
	if err != nil {
		return err
	}
	if pr == nil {
		return fmt.Errorf("book info missing for conversion")
	}

	pr.AudiobookTime = audiobookTime

	if pr.BookID != nil && totalPages > 0 && totalLength != nil && totalLength.Duration > 0 {
		page, _ := timestampToPage(totalPages, audiobookTime.Duration, totalLength.Duration)
		pr.BookPage = &page
	}

	return s.repo.Update(pr)
}

func (s *progressService) SetBook(id int, bookID int) error {
	progress, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if progress == nil {
		return fmt.Errorf("progress not found")
	}

	progress.BookID = &bookID

	if err := s.repo.Update(progress); err != nil {
		return err
	}

	progress, totalPages, totalLength, err := s.repo.GetByIDWithTotals(id)
	if err != nil {
		return err
	}

	if progress.AudiobookTime != nil && totalPages > 0 && totalLength != nil && totalLength.Duration > 0 {
		page, err := timestampToPage(totalPages, progress.AudiobookTime.Duration, totalLength.Duration)
		if err == nil {
			progress.BookPage = &page
			return s.repo.Update(progress)
		}
	}

	return nil
}

func (s *progressService) SetAudiobook(id int, audiobookID int) error {
	progress, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if progress == nil {
		return fmt.Errorf("progress not found")
	}

	progress.AudiobookID = &audiobookID

	if err := s.repo.Update(progress); err != nil {
		return err
	}

	progress, totalPages, totalLength, err := s.repo.GetByIDWithTotals(id)
	if err != nil {
		return err
	}

	if progress.BookPage != nil && totalPages > 0 && totalLength != nil && totalLength.Duration > 0 {
		ts, err := pageToTimestamp(totalPages, *progress.BookPage, totalLength.Duration)
		if err == nil {
			cd := domain.CustomDuration{Duration: ts}
			progress.AudiobookTime = &cd
			return s.repo.Update(progress)
		}
	}

	return nil
}

func (s *progressService) FilterProgress(filter repository.ProgressFilter) ([]domain.Progress, error) {
	return s.repo.FilterProgress(filter)
}
