package bl

import (
	"book_boy/backend/internal/models"
	"fmt"
	"math"
	"time"
)

// pageToTimestamp converts a 1-based book page to a timestamp inside the audiobook.
// - page is clamped to [1, bookTotalPages]
// - maps page 1 -> 0s (start), last page -> end of audiobook
// - uses (totalPages-1) to ensure last page maps to 100% completion
func pageToTimestamp(totalPages int, bookPage int, totalLength time.Duration) (time.Duration, error) {
	if totalPages <= 0 {
		return 0, fmt.Errorf("bookTotalPages must be > 0")
	}
	if totalLength <= 0 {
		return 0, fmt.Errorf("audioLen must be > 0")
	}
	if bookPage < 1 {
		bookPage = 1
	}
	if bookPage > totalPages {
		bookPage = totalPages
	}
	if totalPages == 1 {
		return totalLength, nil
	}
	prop := float64(bookPage-1) / float64(totalPages-1)
	secs := prop * totalLength.Seconds()
	return time.Duration(math.Round(secs)) * time.Second, nil
}

// timestampToPage converts an audiobook timestamp to a 1-based book page.
// - ts is clamped to [0, audioLen].
// - result is clipped to [1, bookTotalPages].
// - uses (totalPages-1) to ensure end of audiobook maps to last page
func timestampToPage(totalPages int, audiobookTime, totalLength time.Duration) (int, error) {
	if totalPages <= 0 {
		return 0, fmt.Errorf("bookTotalPages must be > 0")
	}
	if totalLength <= 0 {
		return 0, fmt.Errorf("audioLen must be > 0")
	}
	if audiobookTime < 0 {
		audiobookTime = 0
	}
	if audiobookTime > totalLength {
		audiobookTime = totalLength
	}
	if totalPages == 1 {
		return 1, nil
	}
	pageF := audiobookTime.Seconds()/totalLength.Seconds()*float64(totalPages-1) + 1
	page := int(math.Round(pageF))
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}
	return page, nil
}

func calculateCompletionPercent(progress *models.Progress, totalPages int, totalLength *models.CustomDuration) int {
	var bookPercent, audioPercent float64
	hasBook := progress.BookPage != nil && totalPages > 0
	hasAudio := progress.AudiobookTime != nil && totalLength != nil && totalLength.Duration > 0

	if hasBook {
		bookPercent = float64(*progress.BookPage) / float64(totalPages) * 100
		if bookPercent > 100 {
			bookPercent = 100
		}
	}

	if hasAudio {
		audioPercent = progress.AudiobookTime.Duration.Seconds() / totalLength.Duration.Seconds() * 100
		if audioPercent > 100 {
			audioPercent = 100
		}
	}

	if hasBook && hasAudio {
		return int(math.Round((bookPercent + audioPercent) / 2))
	}
	if hasBook {
		return int(math.Round(bookPercent))
	}
	if hasAudio {
		return int(math.Round(audioPercent))
	}
	return 0
}
