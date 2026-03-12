package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tomascaceres14/readers-app/resource-service/internal/messaging"
	"github.com/tomascaceres14/readers-app/resource-service/internal/scraping"
)

func main() {

	scraper := scraping.NewScraper()

	publisherCfg := messaging.NewPublisherConfig()
	publisher, err := messaging.NewPublisher(publisherCfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := publisher.Setup(); err != nil {
		log.Fatal(err)
	}

	consumerCfg := messaging.NewConsumerConfig()
	consumer, err := messaging.NewScrapingConsumer(consumerCfg, scraper, publisher)
	if err != nil {
		log.Fatal(err)
	}

	if err := consumer.Setup(); err != nil {
		log.Fatal(err)
	}

	consumer.Listen()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	consumer.Close()
	publisher.Close()
}
