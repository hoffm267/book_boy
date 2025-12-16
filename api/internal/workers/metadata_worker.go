package workers

import (
	"book_boy/api/external/book_metadata"
	"book_boy/api/internal/service"
	"book_boy/api/internal/infra"
	"encoding/json"
	"fmt"
	"log"
)

type MetadataWorker struct {
	queue          *infra.Queue
	bookService    service.BookService
	metadataClient *book_metadata.Client
}

func NewMetadataWorker(queue *infra.Queue, bookService service.BookService, metadataClient *book_metadata.Client) *MetadataWorker {
	return &MetadataWorker{
		queue:          queue,
		bookService:    bookService,
		metadataClient: metadataClient,
	}
}

func (w *MetadataWorker) Start() error {
	log.Println("Starting metadata worker...")

	handler := func(message []byte) error {
		var job infra.MetadataFetchJob
		if err := json.Unmarshal(message, &job); err != nil {
			return fmt.Errorf("failed to unmarshal job: %w", err)
		}

		log.Printf("Processing metadata fetch for book ID %d, ISBN %s", job.BookID, job.ISBN)

		metadata, err := w.metadataClient.GetBookMetadataByISBN(job.ISBN)
		if err != nil {
			log.Printf("Warning: failed to fetch metadata for ISBN %s: %v", job.ISBN, err)
			return nil
		}

		book, err := w.bookService.GetByID(job.BookID)
		if err != nil {
			return fmt.Errorf("failed to get book: %w", err)
		}

		if metadata.Title != nil {
			book.Title = *metadata.Title
		}
		if metadata.TotalPages != nil {
			book.TotalPages = *metadata.TotalPages
		}

		if err := w.bookService.Update(book); err != nil {
			return fmt.Errorf("failed to update book: %w", err)
		}

		log.Printf("Successfully updated book ID %d with metadata", job.BookID)
		return nil
	}

	return w.queue.Subscribe("metadata_fetch", handler)
}
