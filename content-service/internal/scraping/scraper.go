package scraping

import "github.com/gocolly/colly"

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
