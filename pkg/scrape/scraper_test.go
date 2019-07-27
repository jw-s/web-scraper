package scrape

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakePageGetter map[string][]string

func (f *fakePageGetter) Get(url string) *Page {
	if childLinks, ok := (*f)[url]; ok {
		return &Page{
			URL:      url,
			Children: childLinks,
		}
	}
	return &Page{
		URL: url,
	}
}

func TestScrape(t *testing.T) {
	tests := []struct {
		pages          []string
		fakePageGetter fakePageGetter
		url            string
	}{
		{
			url: "https://google.com",
			pages: []string{
				"https://google.com",
				"https://google.com/search",
				"https://google.com/help",
				"https://google.com/about",
				"https://google.com/faq",
				"https://google.com/about/contact",
				"https://google.com/about/tc",
			},
			fakePageGetter: fakePageGetter{
				"https://google.com": []string{
					"https://google.com/search",
					"https://google.com/help",
					"https://google.com/about",
				},
				"https://google.com/search": []string{
					"https://google.com/help",
					"https://google.com/about",
				},
				"https://google.com/about": []string{
					"https://google.com/about/tc",
					"https://google.com/about/contact",
					"https://google.com",
				},
				"https://google.com/help": []string{
					"https://google.com/faq",
				},
			},
		},
	}

	for _, test := range tests {
		scraper := NewScraper(&test.fakePageGetter)
		pages := scraper.Scrape(test.url)
		assert.ElementsMatch(t, pages, test.pages)
	}

}
