package service

import (
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
