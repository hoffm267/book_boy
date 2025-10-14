package bl

import (
	"errors"
	"testing"
	"time"

	"book_boy/backend/internal/models"
)

type mockProgressRepo struct {
	Data map[int]models.Progress
	Err  error
}

func (m *mockProgressRepo) GetAll() ([]models.Progress, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var all []models.Progress
	for _, v := range m.Data {
		all = append(all, v)
	}
	return all, nil
}

func (m *mockProgressRepo) GetByID(id int) (*models.Progress, error) {
	if val, ok := m.Data[id]; ok {
		return &val, nil
	}
	return nil, nil
}

func (m *mockProgressRepo) Create(progress *models.Progress) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	id := len(m.Data) + 1
	progress.ID = id
	m.Data[progress.ID] = *progress
	return id, nil
}

func (m *mockProgressRepo) Update(progress *models.Progress) error {
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

func (m *mockProgressRepo) GetByIDWithTotals(id int) (progress *models.Progress, bookPage int, audiobookTime *models.CustomDuration, err error) {
	if m.Err != nil {
		return nil, 0, nil, m.Err
	}
	if prog, ok := m.Data[id]; ok {
		page := 0
		if prog.BookPage != nil {
			page = *prog.BookPage
		}
		return &prog, page, prog.AudiobookTime, nil
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

func (m *mockProgressRepo) FilterProgress(filter models.ProgressFilter) ([]models.Progress, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var results []models.Progress
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

func TestProgressService(t *testing.T) {
	mockData := map[int]models.Progress{
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
		progress := models.Progress{
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
		progress := models.Progress{
			ID:       1,
			UserID:   1,
			BookPage: ptrInt(15),
		}
		err := svc.Update(&progress)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Update not found", func(t *testing.T) {
		progress := models.Progress{ID: 999}
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

func ptrInt(i int) *int { return &i }
