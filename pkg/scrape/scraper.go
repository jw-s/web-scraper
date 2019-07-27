package scrape

import "github.com/sirupsen/logrus"

// Scraper crawls a website, starting from the root.
type Scraper struct {
	pageScrapper PageScraper
}

// NewScraper creates a Scraper.
func NewScraper(pageScrapper PageScraper) *Scraper {
	return &Scraper{
		pageScrapper: pageScrapper,
	}
}

// Scrape takes a root url and crawls through all the root and children links.
func (s *Scraper) Scrape(url string) []string {
	collector := make(chan *Page) //TODO make channel buffer size configurable
	go s.scrape(collector, url)
	return s.collect(collector)
}

func (s *Scraper) collect(pageCh chan *Page) []string {
	var pages []string
	outstandingScrapes := 1 //seed job
	logrus.Info("Starting page collection")
loop:
	for {
		select {
		case page := <-pageCh:
			logrus.Infof("Received page %v", *page)
			outstandingScrapes--
			for _, p := range pages {
				if p == page.URL {
					logrus.Debugf("Page URL has already been crawled: %s. skipping!", p)
					continue loop
				}
			}

			pages = append(pages, page.URL)

		children:
			for _, url := range page.Children {
				for _, p := range pages {
					if p == url {
						logrus.Debugf("Child Page URL has already been crawled: %s. skipping!", p)
						continue children
					}
				}
				go s.scrape(pageCh, url)
				outstandingScrapes++
			}
		default:
			if outstandingScrapes == 0 {
				logrus.Infof("page count: %v", len(pages))
				return pages
			}
		}
	}
}

func (s *Scraper) scrape(pageCh chan<- *Page, url string) {
	page := s.pageScrapper.Get(url)
	logrus.Infof("Got page %v, sending page to collector", *page)
	pageCh <- page
}
