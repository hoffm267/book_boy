package service

import (
	"book_boy/api/internal/domain"
	"math"
	"testing"
	"time"
)

func TestPageToTimestamp(t *testing.T) {
	tests := []struct {
		name        string
		totalPages  int
		bookPage    int
		totalLength time.Duration
		want        time.Duration
		wantErr     bool
	}{
		{
			name:        "first page maps to start",
			totalPages:  300,
			bookPage:    1,
			totalLength: 10 * time.Hour,
			want:        0,
			wantErr:     false,
		},
		{
			name:        "last page maps to end",
			totalPages:  300,
			bookPage:    300,
			totalLength: 10 * time.Hour,
			want:        10 * time.Hour,
			wantErr:     false,
		},
		{
			name:        "middle page",
			totalPages:  300,
			bookPage:    150,
			totalLength: 10 * time.Hour,
			want:        4*time.Hour + 59*time.Minute, // (149/299) * 10 hours ≈ 4.983 hours
			wantErr:     false,
		},
		{
			name:        "single page book",
			totalPages:  1,
			bookPage:    1,
			totalLength: 2 * time.Hour,
			want:        2 * time.Hour,
			wantErr:     false,
		},
		{
			name:        "page clamped to min",
			totalPages:  300,
			bookPage:    0,
			totalLength: 10 * time.Hour,
			want:        0,
			wantErr:     false,
		},
		{
			name:        "page clamped to max",
			totalPages:  300,
			bookPage:    500,
			totalLength: 10 * time.Hour,
			want:        10 * time.Hour,
			wantErr:     false,
		},
		{
			name:        "invalid total pages",
			totalPages:  0,
			bookPage:    1,
			totalLength: 10 * time.Hour,
			want:        0,
			wantErr:     true,
		},
		{
			name:        "invalid total length",
			totalPages:  300,
			bookPage:    1,
			totalLength: 0,
			want:        0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pageToTimestamp(tt.totalPages, tt.bookPage, tt.totalLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("pageToTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Allow 1 second tolerance for rounding
				diff := math.Abs(float64(got - tt.want))
				if diff > float64(time.Second) {
					t.Errorf("pageToTimestamp() = %v, want %v (diff: %v)", got, tt.want, time.Duration(diff))
				}
			}
		})
	}
}

func TestTimestampToPage(t *testing.T) {
	tests := []struct {
		name          string
		totalPages    int
		audiobookTime time.Duration
		totalLength   time.Duration
		want          int
		wantErr       bool
	}{
		{
			name:          "start maps to first page",
			totalPages:    300,
			audiobookTime: 0,
			totalLength:   10 * time.Hour,
			want:          1,
			wantErr:       false,
		},
		{
			name:          "end maps to last page",
			totalPages:    300,
			audiobookTime: 10 * time.Hour,
			totalLength:   10 * time.Hour,
			want:          300,
			wantErr:       false,
		},
		{
			name:          "middle timestamp",
			totalPages:    300,
			audiobookTime: 5 * time.Hour,
			totalLength:   10 * time.Hour,
			want:          151, // (0.5 * 299) + 1 = 150.5 -> rounds to 151
			wantErr:       false,
		},
		{
			name:          "single page book",
			totalPages:    1,
			audiobookTime: 1 * time.Hour,
			totalLength:   2 * time.Hour,
			want:          1,
			wantErr:       false,
		},
		{
			name:          "time clamped to min",
			totalPages:    300,
			audiobookTime: -5 * time.Minute,
			totalLength:   10 * time.Hour,
			want:          1,
			wantErr:       false,
		},
		{
			name:          "time clamped to max",
			totalPages:    300,
			audiobookTime: 15 * time.Hour,
			totalLength:   10 * time.Hour,
			want:          300,
			wantErr:       false,
		},
		{
			name:          "invalid total pages",
			totalPages:    0,
			audiobookTime: 1 * time.Hour,
			totalLength:   10 * time.Hour,
			want:          0,
			wantErr:       true,
		},
		{
			name:          "invalid total length",
			totalPages:    300,
			audiobookTime: 1 * time.Hour,
			totalLength:   0,
			want:          0,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timestampToPage(tt.totalPages, tt.audiobookTime, tt.totalLength)
			if (err != nil) != tt.wantErr {
				t.Errorf("timestampToPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("timestampToPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	// Test that converting page -> time -> page gives us back the original page
	totalPages := 300
	totalLength := 10 * time.Hour

	testPages := []int{1, 50, 150, 250, 300}

	for _, page := range testPages {
		t.Run("round trip page "+string(rune(page)), func(t *testing.T) {
			// Convert page to timestamp
			timestamp, err := pageToTimestamp(totalPages, page, totalLength)
			if err != nil {
				t.Fatalf("pageToTimestamp failed: %v", err)
			}

			// Convert timestamp back to page
			gotPage, err := timestampToPage(totalPages, timestamp, totalLength)
			if err != nil {
				t.Fatalf("timestampToPage failed: %v", err)
			}

			// Should get back the original page (or very close due to rounding)
			if math.Abs(float64(gotPage-page)) > 1 {
				t.Errorf("Round trip failed: page %d -> %v -> page %d", page, timestamp, gotPage)
			}
		})
	}
}

func TestSymmetry(t *testing.T) {
	// Test that the algorithms are symmetric
	totalPages := 300
	totalLength := 10 * time.Hour

	// Test key points
	tests := []struct {
		page int
		time time.Duration
	}{
		{1, 0},
		{300, 10 * time.Hour},
	}

	for _, tt := range tests {
		// Page to Time
		gotTime, err := pageToTimestamp(totalPages, tt.page, totalLength)
		if err != nil {
			t.Fatalf("pageToTimestamp failed: %v", err)
		}
		if math.Abs(float64(gotTime-tt.time)) > float64(time.Second) {
			t.Errorf("Page %d should map to %v, got %v", tt.page, tt.time, gotTime)
		}

		// Time to Page
		gotPage, err := timestampToPage(totalPages, tt.time, totalLength)
		if err != nil {
			t.Fatalf("timestampToPage failed: %v", err)
		}
		if gotPage != tt.page {
			t.Errorf("Time %v should map to page %d, got %d", tt.time, tt.page, gotPage)
		}
	}
}

func TestCalculateCompletionPercent(t *testing.T) {
	tests := []struct {
		name        string
		progress    *domain.Progress
		totalPages  int
		totalLength *domain.CustomDuration
		want        int
	}{
		{
			name:        "book only 50%",
			progress:    &domain.Progress{BookPage: ptrInt(50)},
			totalPages:  100,
			totalLength: nil,
			want:        50,
		},
		{
			name:        "book only 100%",
			progress:    &domain.Progress{BookPage: ptrInt(100)},
			totalPages:  100,
			totalLength: nil,
			want:        100,
		},
		{
			name:        "audiobook only 50%",
			progress:    &domain.Progress{AudiobookTime: &domain.CustomDuration{Duration: 30 * time.Minute}},
			totalPages:  0,
			totalLength: &domain.CustomDuration{Duration: 60 * time.Minute},
			want:        50,
		},
		{
			name:        "audiobook only 100%",
			progress:    &domain.Progress{AudiobookTime: &domain.CustomDuration{Duration: 60 * time.Minute}},
			totalPages:  0,
			totalLength: &domain.CustomDuration{Duration: 60 * time.Minute},
			want:        100,
		},
		{
			name: "both formats averaged",
			progress: &domain.Progress{
				BookPage:      ptrInt(50),
				AudiobookTime: &domain.CustomDuration{Duration: 30 * time.Minute},
			},
			totalPages:  100,
			totalLength: &domain.CustomDuration{Duration: 60 * time.Minute},
			want:        50,
		},
		{
			name: "both formats different percentages",
			progress: &domain.Progress{
				BookPage:      ptrInt(80),
				AudiobookTime: &domain.CustomDuration{Duration: 20 * time.Minute},
			},
			totalPages:  100,
			totalLength: &domain.CustomDuration{Duration: 60 * time.Minute},
			want:        57, // (80 + 33.33) / 2 ≈ 56.67 -> rounds to 57
		},
		{
			name:        "no data returns 0",
			progress:    &domain.Progress{},
			totalPages:  0,
			totalLength: nil,
			want:        0,
		},
		{
			name:        "book page nil returns 0",
			progress:    &domain.Progress{BookPage: nil},
			totalPages:  100,
			totalLength: nil,
			want:        0,
		},
		{
			name:        "book over 100% clamped",
			progress:    &domain.Progress{BookPage: ptrInt(150)},
			totalPages:  100,
			totalLength: nil,
			want:        100,
		},
		{
			name:        "audiobook over 100% clamped",
			progress:    &domain.Progress{AudiobookTime: &domain.CustomDuration{Duration: 90 * time.Minute}},
			totalPages:  0,
			totalLength: &domain.CustomDuration{Duration: 60 * time.Minute},
			want:        100,
		},
		{
			name:        "zero total pages with book page",
			progress:    &domain.Progress{BookPage: ptrInt(50)},
			totalPages:  0,
			totalLength: nil,
			want:        0,
		},
		{
			name:        "zero total length with audiobook time",
			progress:    &domain.Progress{AudiobookTime: &domain.CustomDuration{Duration: 30 * time.Minute}},
			totalPages:  0,
			totalLength: &domain.CustomDuration{Duration: 0},
			want:        0,
		},
		{
			name:        "nil total length with audiobook time",
			progress:    &domain.Progress{AudiobookTime: &domain.CustomDuration{Duration: 30 * time.Minute}},
			totalPages:  0,
			totalLength: nil,
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateCompletionPercent(tt.progress, tt.totalPages, tt.totalLength)
			if got != tt.want {
				t.Errorf("calculateCompletionPercent() = %d, want %d", got, tt.want)
			}
		})
	}
}
