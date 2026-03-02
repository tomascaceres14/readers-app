package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/tomascaceres14/readers-app/content-service/utils"
)

type RabbitConfig struct {
	URL          string
	Exchange     string
	Queue        string
	RoutingKey   string
	ExchangeType string
}

func NewRabbitConfig() RabbitConfig {
	return RabbitConfig{
		URL:          utils.GetEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		Exchange:     utils.GetEnv("RABBITMQ_EXCHANGE", "scraping.exchange"),
		Queue:        utils.GetEnv("RABBITMQ_QUEUE", "scraping.queue"),
		RoutingKey:   utils.GetEnv("RABBITMQ_ROUTING_KEY", "scrape"),
		ExchangeType: utils.GetEnv("RABBITMQ_EXCHANGE_TYPE", amqp.ExchangeTopic),
	}
}

type Message struct {
	UserID     string `json:"user_id"`
	ResourceID string `json:"resource_id"`
	URL        string `json:"url"`
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  RabbitConfig
}

func NewConsumer() (*Consumer, error) {
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

	return &Consumer{conn: conn, channel: ch, config: config}, nil
}

func (c *Consumer) Setup() error {
	err := c.channel.ExchangeDeclare(
		c.config.Exchange,
		c.config.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	_, err = c.channel.QueueDeclare(
		c.config.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = c.channel.QueueBind(
		c.config.Queue,
		c.config.RoutingKey,
		c.config.Exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

func (c *Consumer) Listen() error {
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

	log.Println("Consumer started, waiting for messages...")

	for msg := range msgs {
		var message Message
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			continue
		}

		fmt.Printf("Received message: %+v\n", message)
		fmt.Printf("  UserID: %s\n", message.UserID)
		fmt.Printf("  ResourceID: %s\n", message.ResourceID)
		fmt.Printf("  URL: %s\n", message.URL)
	}

	return nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func main() {
	consumer, err := NewConsumer()
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	if err := consumer.Setup(); err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := consumer.Listen(); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
