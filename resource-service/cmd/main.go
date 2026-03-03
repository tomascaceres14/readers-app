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
	scrapeQ, err := messaging.NewScrapingConsumer(consumerCfg, scraper)
	if err != nil {
		log.Fatal(err)
	}

	if err := scrapeQ.Setup(); err != nil {
		log.Fatal(err)
	}
	scrapeQ.Listen()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	scrapeQ.Close()
}
