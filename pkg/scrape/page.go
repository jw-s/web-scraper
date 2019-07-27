package scrape

// Page describes a HTML document with it's url and href children page links.
type Page struct {
	URL      string
	Children []string
}
