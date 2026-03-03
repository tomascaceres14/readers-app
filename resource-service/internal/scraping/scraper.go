package scraping

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Scraper struct {
	Collector *colly.Collector
}

func NewScraper() *Scraper {
	c := colly.NewCollector()
	c.AllowedDomains = []string{}
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Referer", "https://www.google.com/")
	})

	return &Scraper{
		Collector: c,
	}
}

func (s *Scraper) Scrape(url string) {
	c := s.Collector
	var f *os.File
	defer f.Close()

	c.OnHTML("head > title", func(h *colly.HTMLElement) {
		// Sanitizar el título para usarlo como nombre de fichero
		title := strings.TrimSpace(h.Text)
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, " ", "_")

		if title == "" {
			title = url
		}

		var err error
		f, err = os.Create(fmt.Sprintf("files/%s.md", title))
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(f, "# %s\n*Source: %s*\n\n", h.Text, url)
	})

	c.OnHTML("article", func(h *colly.HTMLElement) {
		if f == nil {
			var err error
			f, err = os.Create(fmt.Sprintf("files/%s.md", url))
			if err != nil {
				panic(err)
			}
		}

		dom := h.DOM.Clone()
		dom.Find("nav, button, svg, header, footer, .banner, .ad, figure, h1").Remove()
		dom.Find("div").Each(func(i int, sel *goquery.Selection) {
			text := strings.TrimSpace(sel.Text())
			if text != "" {
				fmt.Fprintf(f, "%s\n", text)
			}
		})
	})

	c.OnHTML("main", func(h *colly.HTMLElement) {
		if f == nil {
			var err error
			f, err = os.Create(fmt.Sprintf("files/%s.md", url))
			if err != nil {
				panic(err)
			}
		}

		dom := h.DOM.Clone()
		dom.Find("nav, button, svg, header, footer, .banner, .ad, figure, h1").Remove()
		dom.Find("div").Each(func(i int, sel *goquery.Selection) {
			text := strings.TrimSpace(sel.Text())
			if text != "" {
				fmt.Fprintf(f, "%s\n", text)
			}
		})
	})

	if err := c.Visit(url); err != nil {
		log.Fatal("Error:", err)
	}

}
