package scrape

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	defaultScrapeTimeout = 15 * time.Second
)

// PageScraper provides the interaction of fetching a HTML page.
type PageScraper interface {
	Get(string) *Page
}

// PageGetter provides a way of fetching HTML pages for a given domain.
type PageGetter struct {
	domain string
	client *http.Client
}

// NewPageGetter creates a new PageGetter.
func NewPageGetter(domain string, client *http.Client) *PageGetter {
	if client == nil {
		client = &http.Client{
			Timeout: defaultScrapeTimeout,
		}
	}
	return &PageGetter{
		domain: domain,
		client: client,
	}
}

// Get returns a page for the given URL.
func (s *PageGetter) Get(url string) (page *Page) {
	page = &Page{
		URL: url,
	}

	p, err := s.client.Get(url)
	if err != nil {
		logrus.Error(err)
		return
	}

	defer p.Body.Close()

	if p.StatusCode >= 200 && p.StatusCode <= 399 {
		page.Children = s.getLinks(p.Body)
	}
	return
}
