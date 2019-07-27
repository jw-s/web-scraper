package scrape

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/html"
	"io"
	"net/url"
)

func (s *PageGetter) getLinks(r io.Reader) (links []string) {
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						url, sameHost := NormalizeURL(a.Val, s.domain)
						if sameHost {
							links = append(links, url.String())
							break
						}
					}
				}
			}
		}
	}
}

// NormalizeURL takes the href link and baseURL
// If the href is relative, the baseURL is added to make it absolute
// If the href is absolute then it's a no-op.
// returns absolute URL and if the URL is part of the baseURL domain
func NormalizeURL(href, base string) (*url.URL, bool) {
	uri, err := url.Parse(href)
	if err != nil {
		logrus.Error(err)
		return nil, false
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		logrus.Error(err)
		return nil, false
	}
	uri = baseURL.ResolveReference(uri)

	return uri, uri.Host == baseURL.Host
}
