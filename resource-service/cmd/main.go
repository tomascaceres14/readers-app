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
	consumerCfg := messaging.NewConsumerConfig()
	consumer, err := messaging.NewScrapingConsumer(consumerCfg, scraper)
	if err != nil {
		log.Fatal(err)
	}

	if err := consumer.Setup(); err != nil {
		log.Fatal(err)
	}
	url := "https://nleiva.medium.com/learn-go-45d4b9c177c7"

	if err := scraper.Scrape(url); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
	consumer.Listen()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	consumer.Close()
}
