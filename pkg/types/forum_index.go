package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"log"
	"net/url"
	"sort"
	"strconv"
)

type ForumIndex struct {
	// URL is the location where the Forum index page can be retrieved.
	URL string

	// Forums is a slice of Forum structs discovered during a ForumIndex scrape.
	Forums []Forum
}

func (f *ForumIndex) Scrape(c *colly.Collector) []error {
	var scrapeErrors = make([]error, 0)

	// Use an XPath query to select the 4th table which contains the forum index then we build up the ForumIndex
	// struct based on the discovered contents
	c.OnXML("/html/body/div/div/div/table[3]", func(e *colly.XMLElement) {
		node, ok := e.DOM.(*html.Node)
		if !ok {
			log.Fatalln("failed to cast to *html.Node")
		}

		doc := goquery.NewDocumentFromNode(node)

		// parse the <tbody> elements in the table. Even count <tbody> contain the forum category information while
		// the odd count <tbody> contain the actual forum name, description and hyperlinks.
		categories := make([]string, 0)
		forumsInCategory := make(map[string][]Forum)

		doc.ChildrenFiltered("tbody").Each(func(i int, s *goquery.Selection) {
			if i%2 == 0 {
				categories = append(categories, s.Find("tr > td > a").Text())
			} else {
				category := categories[i-len(categories)]
				forumsInCategory[category] = make([]Forum, 0)

				s.Find("td.alt1Active").Each(func(i int, s *goquery.Selection) {
					forumNameAndPath := s.Find("div:nth-child(1) > a")

					forumNameAndPath.Attr("href")
					href, found := forumNameAndPath.Attr("href")
					if !found {
						scrapeErrors = append(scrapeErrors, errors.New("missing attribute 'href' on forum row"))
					}

					hrefURL, hrefURLErr := url.Parse(href)
					if hrefURLErr != nil {
						scrapeErrors = append(scrapeErrors, hrefURLErr)
					}

					// forumID is stored as an integer value on the "f" query parameter
					forumID, err := strconv.Atoi(hrefURL.Query().Get("f"))
					if err != nil {
						scrapeErrors = append(scrapeErrors, err)
					}

					f := Forum{
						ID:            forumID,
						URL:           fmt.Sprintf("%s/forumdisplay.php?f=%d", f.URL, forumID),
						ForumIndexURL: f.URL,
						Category:      category,
						Name:          forumNameAndPath.Text(),
					}

					forumsInCategory[category] = append(forumsInCategory[category], f)
				})
			}
		})

		for _, v := range forumsInCategory {
			f.Forums = append(f.Forums, v...)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("INFO: visiting %s\n", r.URL)
	})

	_ = c.Visit(f.URL)

	sort.Slice(f.Forums, func(i, j int) bool {
		return f.Forums[i].ID < f.Forums[j].ID
	})

	return scrapeErrors
}

func (f *ForumIndex) ToJSON() ([]byte, error) {
	bytes, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
