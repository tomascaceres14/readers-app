package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-shiori/go-readability"
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
	url := "https://www.postgresql.org/docs/current/textsearch-intro.html"
	article, e := readability.FromURL(url, time.Minute)
	if e != nil {
		log.Fatal(e)
	}

	dstTxtFile, _ := os.Create("text-test.txt")
	defer dstTxtFile.Close()
	dstTxtFile.WriteString(article.TextContent)

	dstHTMLFile, _ := os.Create("html-temp.html")
	defer dstHTMLFile.Close()
	dstHTMLFile.WriteString(article.Content)

	fmt.Printf("URL     : %s\n", url)
	fmt.Printf("Title   : %s\n", article.Title)
	fmt.Printf("Author  : %s\n", article.Byline)
	fmt.Printf("Length  : %d\n", article.Length)
	fmt.Printf("Excerpt : %s\n", article.Excerpt)
	fmt.Printf("SiteName: %s\n", article.SiteName)
	fmt.Printf("Image   : %s\n", article.Image)
	fmt.Printf("Favicon : %s\n", article.Favicon)
	fmt.Println()

	os.Exit(0)
	consumer.Listen()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	consumer.Close()
}
