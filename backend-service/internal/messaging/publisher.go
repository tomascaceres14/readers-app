package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tomascaceres14/readers-app/backend-service/utils"
)

type RabbitConfig struct {
	URL          string
	Exchange     string
	RoutingKey   string
	ExchangeType string
}

func NewRabbitConfig() RabbitConfig {
	return RabbitConfig{
		URL:          utils.GetEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		Exchange:     utils.GetEnv("RABBITMQ_EXCHANGE", "scraping.exchange"),
		RoutingKey:   utils.GetEnv("RABBITMQ_ROUTING_KEY", "scrape"),
		ExchangeType: utils.GetEnv("RABBITMQ_EXCHANGE_TYPE", amqp.ExchangeTopic),
	}
}

type Message struct {
	UserID     string `json:"user_id"`
	ResourceID string `json:"resource_id"`
	URL        string `json:"url"`
}

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  RabbitConfig
}

func NewPublisher() (*Publisher, error) {
	config := NewRabbitConfig()

	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Publisher{conn: conn, channel: ch, config: config}, nil
}

func (p *Publisher) Setup() error {
	err := p.channel.ExchangeDeclare(
		p.config.Exchange,
		p.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	return nil
}

func (p *Publisher) PublishScrapingTask(msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = p.channel.PublishWithContext(ctx,
		p.config.Exchange,
		p.config.RoutingKey,
		false,
		false,
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
