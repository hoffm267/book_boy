package service

import (
	"errors"
	"testing"
	"time"

	"book_boy/api/internal/domain"
)

func TestTrackingService_StartTracking_Book(t *testing.T) {
	bookRepo := &mockBookRepo{Books: make(map[int]domain.Book)}
	audiobookRepo := &mockAudiobookRepo{}
	progressRepo := &mockProgressRepo{Data: make(map[int]domain.Progress)}
	svc := NewTrackingService(bookRepo, audiobookRepo, progressRepo)

	req := &domain.StartTrackingRequest{
		Format:     "book",
		Title:      "The Great Gatsby",
		TotalPages: 300,
		ISBN:       "9780743273565",
	}

	progress, err := svc.StartTracking(1, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progress == nil {
		t.Fatal("expected progress, got nil")
	}
	if progress.BookID == nil {
		t.Fatal("expected BookID to be set")
	}
	if progress.BookPage == nil || *progress.BookPage != 1 {
		t.Fatalf("expected BookPage to be 1, got %v", progress.BookPage)
	}
	if len(bookRepo.Books) != 1 {
		t.Fatalf("expected 1 book created, got %d", len(bookRepo.Books))
	}
}

func TestTrackingService_StartTracking_BookWithCurrentPage(t *testing.T) {
	bookRepo := &mockBookRepo{Books: make(map[int]domain.Book)}
	audiobookRepo := &mockAudiobookRepo{}
	progressRepo := &mockProgressRepo{Data: make(map[int]domain.Progress)}
	svc := NewTrackingService(bookRepo, audiobookRepo, progressRepo)

	req := &domain.StartTrackingRequest{
		Format:      "book",
		Title:       "The Great Gatsby",
		TotalPages:  300,
		ISBN:        "9780743273565",
		CurrentPage: 50,
	}

	progress, err := svc.StartTracking(1, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progress.BookPage == nil || *progress.BookPage != 50 {
		t.Fatalf("expected BookPage to be 50, got %v", progress.BookPage)
	}
}

func TestTrackingService_StartTracking_Audiobook(t *testing.T) {
	bookRepo := &mockBookRepo{Books: make(map[int]domain.Book)}
	audiobookRepo := &mockAudiobookRepo{}
	progressRepo := &mockProgressRepo{Data: make(map[int]domain.Progress)}
	svc := NewTrackingService(bookRepo, audiobookRepo, progressRepo)

	req := &domain.StartTrackingRequest{
		Format:      "audiobook",
		Title:       "The Great Gatsby",
		TotalLength: "10:30:00",
	}

	progress, err := svc.StartTracking(1, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progress == nil {
		t.Fatal("expected progress, got nil")
	}
	if progress.AudiobookID == nil {
		t.Fatal("expected AudiobookID to be set")
	}
	if progress.AudiobookTime == nil || progress.AudiobookTime.Duration != 0 {
		t.Fatalf("expected AudiobookTime to be 0, got %v", progress.AudiobookTime)
	}
}

func TestTrackingService_StartTracking_AudiobookWithCurrentTime(t *testing.T) {
	bookRepo := &mockBookRepo{Books: make(map[int]domain.Book)}
	audiobookRepo := &mockAudiobookRepo{}
	progressRepo := &mockProgressRepo{Data: make(map[int]domain.Progress)}
	svc := NewTrackingService(bookRepo, audiobookRepo, progressRepo)

	req := &domain.StartTrackingRequest{
		Format:      "audiobook",
		Title:       "The Great Gatsby",
		TotalLength: "10:30:00",
		CurrentTime: "2:15:00",
	}

	progress, err := svc.StartTracking(1, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := 2*time.Hour + 15*time.Minute
	if progress.AudiobookTime == nil || progress.AudiobookTime.Duration != expected {
		t.Fatalf("expected AudiobookTime %v, got %v", expected, progress.AudiobookTime)
	}
}

func TestTrackingService_StartTracking_ValidationErrors(t *testing.T) {
	bookRepo := &mockBookRepo{Books: make(map[int]domain.Book)}
	audiobookRepo := &mockAudiobookRepo{}
	progressRepo := &mockProgressRepo{Data: make(map[int]domain.Progress)}
	svc := NewTrackingService(bookRepo, audiobookRepo, progressRepo)

	t.Run("book missing total_pages", func(t *testing.T) {
		req := &domain.StartTrackingRequest{
			Format: "book",
			Title:  "Some Book",
		}
		_, err := svc.StartTracking(1, req)
		if err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("audiobook missing total_length", func(t *testing.T) {
		req := &domain.StartTrackingRequest{
			Format: "audiobook",
			Title:  "Some Audiobook",
		}
		_, err := svc.StartTracking(1, req)
		if err == nil {
			t.Fatal("expected validation error")
		}
	})

	t.Run("audiobook invalid total_length format", func(t *testing.T) {
		req := &domain.StartTrackingRequest{
			Format:      "audiobook",
			Title:       "Some Audiobook",
			TotalLength: "not-a-duration",
		}
		_, err := svc.StartTracking(1, req)
		if err == nil {
			t.Fatal("expected error for invalid duration format")
		}
	})

	t.Run("audiobook invalid current_time format", func(t *testing.T) {
		req := &domain.StartTrackingRequest{
			Format:      "audiobook",
			Title:       "Some Audiobook",
			TotalLength: "10:00:00",
			CurrentTime: "not-a-time",
		}
		_, err := svc.StartTracking(1, req)
		if err == nil {
			t.Fatal("expected error for invalid current_time format")
		}
	})
}

func TestTrackingService_StartTracking_RepoErrors(t *testing.T) {
	t.Run("book repo create error", func(t *testing.T) {
		bookRepo := &mockBookRepo{Books: make(map[int]domain.Book), Err: errors.New("db error")}
		audiobookRepo := &mockAudiobookRepo{}
		progressRepo := &mockProgressRepo{Data: make(map[int]domain.Progress)}
		svc := NewTrackingService(bookRepo, audiobookRepo, progressRepo)

		req := &domain.StartTrackingRequest{
			Format:     "book",
			Title:      "Test Book",
			TotalPages: 200,
		}
		_, err := svc.StartTracking(1, req)
		if err == nil {
			t.Fatal("expected error from book repo")
		}
	})

	t.Run("audiobook repo create error", func(t *testing.T) {
		bookRepo := &mockBookRepo{Books: make(map[int]domain.Book)}
		audiobookRepo := &mockAudiobookRepo{Err: errors.New("db error")}
		progressRepo := &mockProgressRepo{Data: make(map[int]domain.Progress)}
		svc := NewTrackingService(bookRepo, audiobookRepo, progressRepo)

		req := &domain.StartTrackingRequest{
			Format:      "audiobook",
			Title:       "Test Audiobook",
			TotalLength: "5:00:00",
		}
		_, err := svc.StartTracking(1, req)
		if err == nil {
			t.Fatal("expected error from audiobook repo")
		}
	})

	t.Run("progress repo create error", func(t *testing.T) {
		bookRepo := &mockBookRepo{Books: make(map[int]domain.Book)}
		audiobookRepo := &mockAudiobookRepo{}
		progressRepo := &mockProgressRepo{Data: make(map[int]domain.Progress), Err: errors.New("db error")}
		svc := NewTrackingService(bookRepo, audiobookRepo, progressRepo)

		req := &domain.StartTrackingRequest{
			Format:     "book",
			Title:      "Test Book",
			TotalPages: 200,
		}
		_, err := svc.StartTracking(1, req)
		if err == nil {
			t.Fatal("expected error from progress repo")
		}
	})
}

func TestTrackingService_GetCurrentTracking(t *testing.T) {
	bookID := 1
	page := 50
	progressRepo := &mockProgressRepo{
		Data: map[int]domain.Progress{
			1: {ID: 1, UserID: 1, BookID: &bookID, BookPage: &page},
			2: {ID: 2, UserID: 2, BookID: &bookID, BookPage: &page},
		},
	}
	svc := NewTrackingService(nil, nil, progressRepo)

	responses, err := svc.GetCurrentTracking(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 1 {
		t.Fatalf("expected 1 response for user 1, got %d", len(responses))
	}
	if responses[0].ProgressID != 1 {
		t.Errorf("expected ProgressID 1, got %d", responses[0].ProgressID)
	}
}

func TestTrackingService_GetCurrentTracking_Error(t *testing.T) {
	progressRepo := &mockProgressRepo{
		Data: make(map[int]domain.Progress),
		Err:  errors.New("db error"),
	}
	svc := NewTrackingService(nil, nil, progressRepo)

	_, err := svc.GetCurrentTracking(1)
	if err == nil {
		t.Fatal("expected error from progress repo")
	}
}

func TestTrackingService_GetCurrentTracking_Empty(t *testing.T) {
	progressRepo := &mockProgressRepo{Data: make(map[int]domain.Progress)}
	svc := NewTrackingService(nil, nil, progressRepo)

	responses, err := svc.GetCurrentTracking(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(responses) != 0 {
		t.Fatalf("expected 0 responses, got %d", len(responses))
	}
}
