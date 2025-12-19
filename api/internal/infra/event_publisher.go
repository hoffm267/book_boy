package infra

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher struct {
	channel  *amqp.Channel
	exchange string
}

func NewEventPublisher(conn *amqp.Connection, exchange string) (*EventPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &EventPublisher{
		channel:  ch,
		exchange: exchange,
	}, nil
}

func (p *EventPublisher) Publish(routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.channel.Publish(
		p.exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	fmt.Printf("Published event: %s with payload: %s\n", routingKey, string(body))
	return nil
}

func (p *EventPublisher) Close() error {
	if p.channel != nil {
		return p.channel.Close()
	}
	return nil
}
