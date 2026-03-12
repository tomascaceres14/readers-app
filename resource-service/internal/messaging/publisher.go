package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tomascaceres14/readers-app/resource-service/utils"
)

const (
	StatusOK     = "OK"
	StatusFAILED = "FAILED"
)

type PublisherConfig struct {
	URL          string
	Exchange     string
	ExchangeType string
}

func NewPublisherConfig() PublisherConfig {
	return PublisherConfig{
		URL:          utils.GetEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		Exchange:     utils.GetEnv("RABBITMQ_EXCHANGE", "scraping.exchange"),
		ExchangeType: utils.GetEnv("RABBITMQ_EXCHANGE_TYPE", "topic"),
	}
}

type ResponseMessage struct {
	ResourceID string `json:"resource_id"`
	Status     string `json:"status"`
	Title      string `json:"title,omitempty"`
	Excerpt    string `json:"excerpt,omitempty"`
	Language   string `json:"language,omitempty"`
}

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  PublisherConfig
}

func NewPublisher(cfg PublisherConfig) (*Publisher, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Publisher{conn: conn, channel: ch, config: cfg}, nil
}

func (p *Publisher) Setup() error {
	err := p.channel.ExchangeDeclare(
		p.config.Exchange,
		p.config.ExchangeType,
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = p.channel.QueueDeclare(
		"scraping.responses",
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare responses queue: %w", err)
	}

	err = p.channel.QueueBind(
		"scraping.responses",
		"scraping.response",
		p.config.Exchange,
		false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind responses queue: %w", err)
	}

	return nil
}

func (p *Publisher) PublishResponse(msg ResponseMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(ctx,
		p.config.Exchange,
		"scraping.response",
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
