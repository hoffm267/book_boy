package bl

import (
	"fmt"
	"math"
	"time"
)

// pageToTimestamp converts a 1-based book page to a timestamp inside the audiobook.
// - page is clamped to [1, bookTotalPages]
// - maps page 1 -> 0s (start of book). Change to center-of-page by using (page-0.5) if desired.
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
	prop := float64(bookPage-1) / float64(totalPages)
	secs := prop * totalLength.Seconds()
	return time.Duration(math.Round(secs)) * time.Second, nil
}

// timestampToPage converts an audiobook timestamp to a 1-based book page.
// - ts is clamped to [0, audioLen].
// - result is clipped to [1, bookTotalPages].
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
	pageF := audiobookTime.Seconds()/totalLength.Seconds()*float64(totalPages) + 1
	page := int(math.Round(pageF))
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}
	return page, nil
}
