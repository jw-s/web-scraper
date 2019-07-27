package main

import (
	"flag"
	"github.com/jw-s/web-scraper/pkg/scrape"
	"github.com/sirupsen/logrus"
)

var (
	baseURL = flag.String("base-url", "https://google.com", "Base url for scraping")
	debug   = flag.Bool("debug", false, "Debug log level")
)

func init() {
	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {

	var (
		pageScrapper = scrape.NewPageGetter(*baseURL, nil)
		scraper      = scrape.NewScraper(pageScrapper)
	)
	logrus.Info(scraper.Scrape(*baseURL))
}
