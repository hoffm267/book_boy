package service

import (
	"errors"
	"testing"
	"time"

	"book_boy/api/internal/domain"
	"book_boy/api/internal/repository"
)

type mockProgressRepo struct {
	Data map[int]domain.Progress
	Err  error
}

func (m *mockProgressRepo) GetAll() ([]domain.Progress, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var all []domain.Progress
	for _, v := range m.Data {
		all = append(all, v)
	}
	return all, nil
}

func (m *mockProgressRepo) GetByID(id int) (*domain.Progress, error) {
	if val, ok := m.Data[id]; ok {
		return &val, nil
	}
	return nil, nil
}

func (m *mockProgressRepo) Create(progress *domain.Progress) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	id := len(m.Data) + 1
	progress.ID = id
	m.Data[progress.ID] = *progress
	return id, nil
}

func (m *mockProgressRepo) Update(progress *domain.Progress) error {
	if m.Err != nil {
		return m.Err
	}
	if _, ok := m.Data[progress.ID]; !ok {
		return errors.New("not found")
	}
	m.Data[progress.ID] = *progress
	return nil
}

func (m *mockProgressRepo) Delete(id int) error {
	if m.Err != nil {
		return m.Err
	}
	if _, ok := m.Data[id]; !ok {
		return errors.New("not found")
	}
	delete(m.Data, id)
	return nil
}

func (m *mockProgressRepo) GetByIDWithTotals(id int) (progress *domain.Progress, totalPages int, totalLength *domain.CustomDuration, err error) {
	if m.Err != nil {
		return nil, 0, nil, m.Err
	}
	if prog, ok := m.Data[id]; ok {
		return &prog, 500, prog.AudiobookTime, nil
	}
	return nil, 0, nil, nil
}

func (m *mockProgressRepo) SetBook(id int, bookId int) error {
	if m.Err != nil {
		return m.Err
	}
	if prog, ok := m.Data[id]; ok {
		prog.BookID = &bookId
		m.Data[id] = prog
		return nil
	}
	return errors.New("not found")
}

func (m *mockProgressRepo) SetAudiobook(id int, audiobookId int) error {
	if m.Err != nil {
		return m.Err
	}
	if prog, ok := m.Data[id]; ok {
		prog.AudiobookID = &audiobookId
		m.Data[id] = prog
		return nil
	}
	return errors.New("not found")
}

func (m *mockProgressRepo) FilterProgress(filter repository.ProgressFilter) ([]domain.Progress, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var results []domain.Progress
	for _, prog := range m.Data {
		match := true
		if filter.ID != nil && prog.ID != *filter.ID {
			match = false
		}
		if filter.UserID != nil && prog.UserID != *filter.UserID {
			match = false
		}
		if filter.BookID != nil {
			if prog.BookID == nil || *prog.BookID != *filter.BookID {
				match = false
			}
		}
		if filter.AudiobookID != nil {
			if prog.AudiobookID == nil || *prog.AudiobookID != *filter.AudiobookID {
				match = false
			}
		}
		if match {
			results = append(results, prog)
		}
	}
	return results, nil
}

func (m *mockProgressRepo) GetAllEnrichedByUser(userID int) ([]domain.EnrichedProgress, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var results []domain.EnrichedProgress
	for _, prog := range m.Data {
		if prog.UserID == userID {
			results = append(results, domain.EnrichedProgress{
				Progress:          prog,
				CompletionPercent: 0,
			})
		}
	}
	return results, nil
}

func TestProgressService(t *testing.T) {
	mockData := map[int]domain.Progress{
		1: {
			ID:        1,
			UserID:    1,
			BookID:    ptrInt(1),
			BookPage:  ptrInt(10),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo := &mockProgressRepo{
		Data: mockData,
		Err:  nil,
	}

	svc := NewProgressService(mockRepo)

	t.Run("GetAll", func(t *testing.T) {
		res, err := svc.GetAll()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res) != len(mockData) {
			t.Fatalf("expected %d progress records, got %d", len(mockData), len(res))
		}
	})

	t.Run("GetByID found", func(t *testing.T) {
		res, err := svc.GetByID(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.ID != 1 {
			t.Fatalf("expected ID 1, got %+v", res)
		}
	})

	t.Run("GetByID not found", func(t *testing.T) {
		res, err := svc.GetByID(999)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res != nil {
			t.Fatalf("expected nil, got %+v", res)
		}
	})

	t.Run("Create", func(t *testing.T) {
		progress := domain.Progress{
			UserID:    2,
			BookID:    ptrInt(2),
			BookPage:  ptrInt(20),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := svc.Create(&progress)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if progress.ID == 0 {
			t.Fatal("expected ID to be set")
		}
	})

	t.Run("Update found", func(t *testing.T) {
		bookID := 1
		progress := domain.Progress{
			ID:       1,
			UserID:   1,
			BookID:   &bookID,
			BookPage: ptrInt(15),
		}
		err := svc.Update(&progress)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Update not found", func(t *testing.T) {
		progress := domain.Progress{ID: 999}
		err := svc.Update(&progress)
		if err == nil {
			t.Fatal("expected error for non-existent record")
		}
	})

	t.Run("Delete found", func(t *testing.T) {
		err := svc.Delete(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Delete not found", func(t *testing.T) {
		err := svc.Delete(999)
		if err == nil {
			t.Fatal("expected error for non-existent record")
		}
	})
}

func TestProgressService_UpdateProgressPage(t *testing.T) {
	bookID := 1
	bookPage := 50
	mockRepo := &mockProgressRepo{
		Data: map[int]domain.Progress{
			1: {
				ID:       1,
				UserID:   1,
				BookID:   &bookID,
				BookPage: &bookPage,
			},
		},
	}
	svc := NewProgressService(mockRepo)

	err := svc.UpdateProgressPage(1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	updated := mockRepo.Data[1]
	if updated.BookPage == nil || *updated.BookPage != 100 {
		t.Fatalf("expected BookPage to be 100, got %v", updated.BookPage)
	}
}

func TestProgressService_UpdateProgressTime(t *testing.T) {
	bookID := 1
	mockRepo := &mockProgressRepo{
		Data: map[int]domain.Progress{
			1: {
				ID:     1,
				UserID: 1,
				BookID: &bookID,
			},
		},
	}
	svc := NewProgressService(mockRepo)

	newTime := &domain.CustomDuration{Duration: 30 * time.Minute}
	err := svc.UpdateProgressTime(1, newTime)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	updated := mockRepo.Data[1]
	if updated.AudiobookTime == nil || updated.AudiobookTime.Duration != 30*time.Minute {
		t.Fatalf("expected AudiobookTime to be 30 minutes, got %v", updated.AudiobookTime)
	}
}

func TestProgressService_SetBook(t *testing.T) {
	mockRepo := &mockProgressRepo{
		Data: map[int]domain.Progress{
			1: {ID: 1, UserID: 1},
		},
	}
	svc := NewProgressService(mockRepo)

	err := svc.SetBook(1, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProgressService_SetAudiobook(t *testing.T) {
	mockRepo := &mockProgressRepo{
		Data: map[int]domain.Progress{
			1: {ID: 1, UserID: 1},
		},
	}
	svc := NewProgressService(mockRepo)

	err := svc.SetAudiobook(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProgressService_FilterProgress(t *testing.T) {
	bookID1 := 1
	bookID2 := 2
	mockRepo := &mockProgressRepo{
		Data: map[int]domain.Progress{
			1: {ID: 1, UserID: 1, BookID: &bookID1},
			2: {ID: 2, UserID: 1, BookID: &bookID2},
			3: {ID: 3, UserID: 2, BookID: &bookID1},
		},
	}
	svc := NewProgressService(mockRepo)

	userID := 1
	filter := repository.ProgressFilter{UserID: &userID}
	results, err := svc.FilterProgress(filter)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 progress records for user 1, got %d", len(results))
	}
	for _, prog := range results {
		if prog.UserID != 1 {
			t.Errorf("expected UserID 1, got %d", prog.UserID)
		}
	}
}

func ptrInt(i int) *int { return &i }
