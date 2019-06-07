package main

import (
	"github.com/gocolly/colly"
	"github.com/plombardi89/ab-scrape/pkg/types"
	"log"
	"net/url"
	"path/filepath"
	"time"
)

var (
	outputRoot = "dump/"
	forumURL   = "http://archboston.org/community"
)

func main() {
	forumURL, err := url.Parse(forumURL)
	if err != nil {
		log.Fatalln(err)
	}

	c := colly.NewCollector()
	index := &types.ForumIndex{URL: forumURL.String()}
	scrapeErrs := index.Scrape(c)
	if len(scrapeErrs) > 0 {
		for _, err := range scrapeErrs {
			log.Printf("error: %v\n", err)
		}

		log.Fatalln("Errors!")
	}

	writer := types.FileSystemWriter{
		StorageRoot: filepath.Join(outputRoot, forumURL.Host, time.Now().UTC().Format("20060102150405")),
	}

	writerErr := writer.WriteForumIndex(index)
	if writerErr != nil {
		log.Fatalln(writerErr)
	}

	for _, f := range index.Forums {
		log.Printf("INFO: scraping forum [%d]: %s\n", f.ID, f.Name)
		f.Scrape(colly.NewCollector(), writer)
	}
}
