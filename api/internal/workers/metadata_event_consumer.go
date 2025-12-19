package workers

import (
	"book_boy/api/internal/domain"
	"book_boy/api/internal/infra"
	"book_boy/api/internal/service"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MetadataEventConsumer struct {
	conn       *amqp.Connection
	service    service.BookService
	sseManager *infra.SSEManager
}

func NewMetadataEventConsumer(conn *amqp.Connection, svc service.BookService, sseMgr *infra.SSEManager) *MetadataEventConsumer {
	return &MetadataEventConsumer{
		conn:       conn,
		service:    svc,
		sseManager: sseMgr,
	}
}

func (c *MetadataEventConsumer) Start() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		"book_events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(
		"api_metadata_results",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		queue.Name,
		"book.metadata_fetched",
		"book_events",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Metadata event consumer started, waiting for book.metadata_fetched events...")

	go func() {
		for msg := range msgs {
			if err := c.handleMetadataFetched(msg.Body); err != nil {
				log.Printf("Error handling metadata event: %v\n", err)
				msg.Nack(false, true)
			} else {
				msg.Ack(false)
			}
		}
	}()

	return nil
}

func (c *MetadataEventConsumer) handleMetadataFetched(body []byte) error {
	var event domain.BookMetadataFetchedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("Received metadata for book %d: %s (%d pages)\n", event.BookID, event.Title, event.TotalPages)

	if !event.Success {
		log.Printf("Metadata fetch failed for book %d: %s\n", event.BookID, event.Error)
		return nil
	}

	book, err := c.service.GetByID(event.BookID)
	if err != nil {
		return fmt.Errorf("failed to get book: %w", err)
	}

	book.Title = event.Title
	book.TotalPages = event.TotalPages

	if err := c.service.Update(book); err != nil {
		return fmt.Errorf("failed to update book: %w", err)
	}

	log.Printf("Successfully updated book %d with metadata\n", event.BookID)

	c.sseManager.Broadcast("book.metadata_fetched", book)

	return nil
}
