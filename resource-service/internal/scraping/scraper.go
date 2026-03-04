package scraping

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	readability "codeberg.org/readeck/go-readability/v2"
	"github.com/gocolly/colly/v2"
)

type Scraper struct {
	Collector *colly.Collector
}

func NewScraper() *Scraper {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)
	c.AllowedDomains = []string{}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Sec-Ch-Ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Referer", "https://www.google.com/")
	})

	return &Scraper{
		Collector: c,
	}
}

func (s *Scraper) Scrape(fetchUrl string) error {
	c := colly.NewCollector()
	var htmlContent string

	c.OnResponse(func(r *colly.Response) {
		htmlContent = string(r.Body)
	})

	if err := c.Visit(fetchUrl); err != nil {
		return fmt.Errorf("failed to visit url: %w", err)
	}

	parsedUrl, err := url.Parse(fetchUrl)
	if err != nil {
		return fmt.Errorf("failed to parse url: %w", err)
	}

	article, err := readability.FromReader(strings.NewReader(htmlContent), parsedUrl)
	if err != nil {
		return fmt.Errorf("failed to parse content: %w", err)
	}

	filename := strings.ReplaceAll(article.Title(), " ", "_")
	filename = strings.ReplaceAll(filename, "/", "_")

	f, err := os.Create(fmt.Sprintf("files/%s.md", filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "# %s\n### %s\n*Source: %s*\n\n",
		article.Title(),
		article.Excerpt(),
		fetchUrl,
	)

	article.RenderText(f)

	return err
}
