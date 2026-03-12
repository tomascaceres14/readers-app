package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	resourcestatus "github.com/tomascaceres14/readers-app/backend-service/internal/resource_status"
	"github.com/tomascaceres14/readers-app/backend-service/utils"
)

type ResponseConfig struct {
	URL          string
	Exchange     string
	ExchangeType string
	Queue        string
	BindingKey   string
}

func NewResponseConfig() ResponseConfig {
	return ResponseConfig{
		URL:          utils.GetEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		Exchange:     utils.GetEnv("RABBITMQ_EXCHANGE", "scraping.exchange"),
		ExchangeType: utils.GetEnv("RABBITMQ_EXCHANGE_TYPE", "topic"),
		Queue:        utils.GetEnv("RABBITMQ_RESPONSE_QUEUE", "scraping.responses"),
		BindingKey:   utils.GetEnv("RABBITMQ_RESPONSE_BINDING_KEY", "scraping.response"),
	}
}

type ResponseMessage struct {
	ResourceID string `json:"resource_id"`
	Status     string `json:"status"`
	Title      string `json:"title,omitempty"`
	Excerpt    string `json:"excerpt,omitempty"`
	Language   string `json:"language,omitempty"`
}

type ResponseConsumer struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	config      ResponseConfig
	resourceSvc ResourceUpdater
}

type ResourceUpdater interface {
	UpdateAfterScrape(id uuid.UUID, statusName, title, excerpt, language string) error
	UpdateStatusFailed(id uuid.UUID, statusName string) error
}

func NewResponseConsumer(cfg ResponseConfig, resourceSvc ResourceUpdater) (*ResponseConsumer, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &ResponseConsumer{
		conn:        conn,
		channel:     ch,
		config:      cfg,
		resourceSvc: resourceSvc,
	}, nil
}

func (c *ResponseConsumer) Setup() error {
	err := c.channel.ExchangeDeclare(
		c.config.Exchange,
		c.config.ExchangeType,
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = c.channel.QueueDeclare(
		c.config.Queue,
		true, false, false, false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = c.channel.QueueBind(
		c.config.Queue,
		c.config.BindingKey,
		c.config.Exchange,
		false, nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return c.channel.Qos(1, 0, false)
}

func (c *ResponseConsumer) Listen() error {
	msgs, err := c.channel.Consume(
		c.config.Queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for msg := range msgs {
		var resp ResponseMessage
		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			log.Printf("Failed to unmarshal response message: %v", err)
			continue
		}

		fmt.Printf("Received response: %+v\n", resp)

		resourceID, err := uuid.Parse(resp.ResourceID)
		if err != nil {
			log.Printf("Failed to parse resource ID: %v", err)
			continue
		}

		if err := c.handleResponse(resourceID, resp); err != nil {
			log.Printf("Failed to handle response: %v", err)
		}
	}

	return nil
}

func (c *ResponseConsumer) handleResponse(resourceID uuid.UUID, resp ResponseMessage) error {
	if resp.Status == resourcestatus.OK {
		return c.handleSuccess(resourceID, resp)
	}
	if resp.Status == resourcestatus.FAILED {
		return c.handleFailure(resourceID)
	}
	return nil
}

func (c *ResponseConsumer) handleSuccess(resourceID uuid.UUID, resp ResponseMessage) error {
	if err := c.resourceSvc.UpdateAfterScrape(resourceID, resourcestatus.OK, resp.Title, resp.Excerpt, resp.Language); err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	return nil
}

func (c *ResponseConsumer) handleFailure(resourceID uuid.UUID) error {
	if err := c.resourceSvc.UpdateStatusFailed(resourceID, resourcestatus.FAILED); err != nil {
		return fmt.Errorf("failed to update resource status to FAILED: %w", err)
	}

	return nil
}

func (c *ResponseConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
