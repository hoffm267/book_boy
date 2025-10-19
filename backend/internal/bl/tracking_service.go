package bl

import (
	"book_boy/backend/internal/dl"
	"book_boy/backend/internal/models"
	"fmt"
	"time"
)

type TrackingService interface {
	StartTracking(userID int, req *models.StartTrackingRequest) (*models.Progress, error)
	GetCurrentTracking(userID int) ([]models.CurrentTrackingResponse, error)
}

type trackingService struct {
	bookRepo      dl.BookRepo
	audiobookRepo dl.AudiobookRepo
	progressRepo  dl.ProgressRepo
}

func NewTrackingService(bookRepo dl.BookRepo, audiobookRepo dl.AudiobookRepo, progressRepo dl.ProgressRepo) TrackingService {
	return &trackingService{
		bookRepo:      bookRepo,
		audiobookRepo: audiobookRepo,
		progressRepo:  progressRepo,
	}
}

func (s *trackingService) StartTracking(userID int, req *models.StartTrackingRequest) (*models.Progress, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	progress := &models.Progress{
		UserID: userID,
	}

	if req.Format == "book" {
		book := &models.Book{
			Title:      req.Title,
			TotalPages: req.TotalPages,
			ISBN:       req.ISBN,
		}
		if err := book.Validate(); err != nil {
			return nil, err
		}

		bookID, err := s.bookRepo.Create(book)
		if err != nil {
			return nil, fmt.Errorf("failed to create book: %w", err)
		}

		progress.BookID = &bookID

		currentPage := 1
		if req.CurrentPage > 0 {
			currentPage = req.CurrentPage
		}
		progress.BookPage = &currentPage
	}

	if req.Format == "audiobook" {
		duration := &models.CustomDuration{}
		parsedDuration, err := time.ParseDuration(req.TotalLength)
		if err != nil {
			hms, parseErr := parseHMS(req.TotalLength)
			if parseErr != nil {
				return nil, fmt.Errorf("invalid total_length format: use HH:MM:SS or duration string")
			}
			parsedDuration = hms
		}
		duration.Duration = parsedDuration

		audiobook := &models.Audiobook{
			Title:       req.Title,
			TotalLength: duration,
		}
		if err := audiobook.Validate(); err != nil {
			return nil, err
		}

		audiobookID, err := s.audiobookRepo.Create(audiobook)
		if err != nil {
			return nil, fmt.Errorf("failed to create audiobook: %w", err)
		}

		progress.AudiobookID = &audiobookID

		if req.CurrentTime != "" {
			currentDuration := &models.CustomDuration{}
			parsedTime, err := time.ParseDuration(req.CurrentTime)
			if err != nil {
				hms, parseErr := parseHMS(req.CurrentTime)
				if parseErr != nil {
					return nil, fmt.Errorf("invalid current_time format: use HH:MM:SS or duration string")
				}
				parsedTime = hms
			}
			currentDuration.Duration = parsedTime
			progress.AudiobookTime = currentDuration
		} else {
			zeroDuration := &models.CustomDuration{Duration: 0}
			progress.AudiobookTime = zeroDuration
		}
	}

	if err := progress.Validate(); err != nil {
		return nil, err
	}

	progressID, err := s.progressRepo.Create(progress)
	if err != nil {
		return nil, fmt.Errorf("failed to create progress: %w", err)
	}

	progress.ID = progressID
	return progress, nil
}

func (s *trackingService) GetCurrentTracking(userID int) ([]models.CurrentTrackingResponse, error) {
	enriched, err := s.progressRepo.GetAllEnrichedByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current tracking: %w", err)
	}

	responses := make([]models.CurrentTrackingResponse, 0, len(enriched))
	for _, e := range enriched {
		responses = append(responses, models.CurrentTrackingResponse{
			ProgressID:        e.Progress.ID,
			Book:              e.Book,
			Audiobook:         e.Audiobook,
			CurrentPage:       e.Progress.BookPage,
			CurrentTime:       e.Progress.AudiobookTime,
			CompletionPercent: e.CompletionPercent,
			CreatedAt:         e.Progress.CreatedAt,
			UpdatedAt:         e.Progress.UpdatedAt,
		})
	}

	return responses, nil
}

func parseHMS(s string) (time.Duration, error) {
	var hours, minutes, seconds int
	_, err := fmt.Sscanf(s, "%d:%d:%d", &hours, &minutes, &seconds)
	if err != nil {
		return 0, err
	}
	return time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}
