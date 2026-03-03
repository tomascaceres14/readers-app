package scraping

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Scraper struct {
	Collector *colly.Collector
}

func NewScraper() *Scraper {

	scraper := &Scraper{
		Collector: colly.NewCollector(),
	}

	scraper.Collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Referer", "https://www.google.com/")
	})

	return scraper
}

func (s *Scraper) Scrape(url string) {
	f, err := os.Create("files/sample.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	c := s.Collector

	c.OnHTML("head > title", func(h *colly.HTMLElement) {
		fmt.Fprintf(f, "# %s\n*Source: %s*\n\n", h.Text, url)
	})

	c.OnHTML("main", func(h *colly.HTMLElement) {
		dom := h.DOM.Clone()
		dom.Find("nav, button, svg, header, footer, .banner, .ad").Remove()

		dom.Find("div").Each(func(i int, s *goquery.Selection) {
			text := strings.TrimSpace(s.Text())
			if text != "" {
				f.WriteString(text)
				f.WriteString("\n")
			}
		})
	})

	c.AllowedDomains = []string{}

	if err := c.Visit(url); err != nil {
		fmt.Println("Error:", err)
	}
}
