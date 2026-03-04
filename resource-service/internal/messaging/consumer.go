package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tomascaceres14/readers-app/resource-service/internal/scraping"
	"github.com/tomascaceres14/readers-app/resource-service/utils"
)

type ConsumerConfig struct {
	URL          string
	Exchange     string
	ExchangeType string
	BindingKey   string
	Queue        string
}

func NewConsumerConfig() ConsumerConfig {
	return ConsumerConfig{
		URL:          utils.GetEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		Exchange:     utils.GetEnv("RABBITMQ_EXCHANGE", "scraping.exchange"),
		ExchangeType: utils.GetEnv("RABBITMQ_EXCHANGE_TYPE", "topic"),
		BindingKey:   utils.GetEnv("RABBITMQ_CONSUMER_BINDING_KEY", "scraping.request"),
		Queue:        utils.GetEnv("RABBITMQ_CONSUMER_QUEUE", "scraping.requests"),
	}
}

type Message struct {
	UserID     string `json:"user_id"`
	ResourceID string `json:"resource_id"`
	URL        string `json:"url"`
}

type ScrapingConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	scraper *scraping.Scraper
	config  ConsumerConfig
}

func NewScrapingConsumer(cfg ConsumerConfig, scraper *scraping.Scraper) (*ScrapingConsumer, error) {

	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &ScrapingConsumer{conn: conn, channel: ch, config: cfg, scraper: scraper}, nil
}

func (c *ScrapingConsumer) Setup() error {
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

func (c *ScrapingConsumer) Listen() error {
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

	fmt.Println("Consumer started, waiting for messages...")

	for msg := range msgs {
		var message Message
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		fmt.Printf("Received message: %+v\n", message)
		if err := c.scraper.Scrape(message.URL); err != nil {
			fmt.Println("Error scraping: ", err)
		}
	}

	return nil
}

func (c *ScrapingConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
