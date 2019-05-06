package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"sort"
	"strconv"
	"time"
)

var (
	forumURL = "http://archboston.org/community"
)

var outputRoot = "dump/"

type ForumIndex struct {
	URL       string             `json:""`
	Forums    map[string][]Forum `json:""`
	collector *colly.Collector
}

type Forum struct {
	ID          int    `json:""`
	Category    string `json:""`
	Name        string `json:""`
	Description string `json:""`
	URL         string `json:""`
}

type Thread struct {
	ID      int
	ForumID int
	Title   string
}

type Post struct {
	ID        int
	ThreadID  int
	Author    string
	Timestamp time.Time
	Content   string
	Signature string
	URL       url.URL
}

func (t *Thread) Scrape() {

}

func (f *Forum) getLastPage(doc *goquery.Document) int {
	fmt.Println(doc)
	fmt.Println(doc.Find("div.pagenav").Text())

	return 0
}

func (f *Forum) Scrape() {
	c := colly.NewCollector()

	// Grab the <table> with the page navigation in it.
	c.OnXML("/html/body/div/div/div/table[1]", func(e *colly.XMLElement) {
		node, ok := e.DOM.(*html.Node)
		if !ok {
			log.Fatalln("failed to cast to *html.Node")
		}

		doc := goquery.NewDocumentFromNode(node)
		f.getLastPage(doc)
	})

	// Grab the <table> with the actual threads in it.


	//c.OnXML("", func(e *colly.XMLElement) {
	//	node, ok := e.DOM.(*html.Node)
	//	if !ok {
	//		log.Fatalln("failed to cast to *html.Node")
	//	}
	//
	//	doc := goquery.NewDocumentFromNode(node)
	//	lastPage := f.getLastPage(doc)
	//
	//	log.Printf("Last Page: %d\n", lastPage)
	//})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("INFO: visiting %s\n", r.URL)
	})

	_ = c.Visit(f.URL)
}

func (f *ForumIndex) Scrape() {
	// Use an XPath query to select the 4th table which contains the forum index then we build up the ForumIndex
	// struct based on the discovered contents
	f.collector.OnXML("/html/body/div/div/div/table[3]", func(e *colly.XMLElement) {
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
						log.Fatalln("missing attribute 'href' on forum row")
					}

					hrefURL, hrefURLErr := url.Parse(href)
					if hrefURLErr != nil {
						log.Fatalln(hrefURLErr)
					}

					// forumID is stored as an integer value on the "f" query parameter
					forumID, err := strconv.Atoi(hrefURL.Query().Get("f"))
					if err != nil {
						log.Fatalln(err)
					}

					f := Forum{
						ID:          forumID,
						URL:         fmt.Sprintf("%s/forumdisplay.php?f=%d", f.URL, forumID),
						Category:    category,
						Name:        forumNameAndPath.Text(),
						Description: "", // TODO: grab the description
					}

					forumsInCategory[category] = append(forumsInCategory[category], f)
				})
			}
		})

		f.Forums = forumsInCategory
	})

	f.collector.OnRequest(func(r *colly.Request) {
		log.Printf("INFO: visiting %s\n", r.URL)
	})

	_ = f.collector.Visit(f.URL)
}

func main() {
	forumURL, err := url.Parse(forumURL)
	if err != nil {
		log.Fatalln(err)
	}

	scrapedDataDir := fmt.Sprintf("%s/%s-%s", outputRoot, forumURL.Host, time.Now().UTC().Format("20060102150405"))

	if err := os.MkdirAll(scrapedDataDir, 0755); err != nil {
		log.Fatalln(err)
	}

	c := colly.NewCollector()
	index := &ForumIndex{URL: forumURL.String(), collector: c}
	index.Scrape()

	indexJSON, marshallErr := json.MarshalIndent(index, "", "    ")
	if marshallErr != nil {
		log.Fatalln(marshallErr)
	}

	if indexWriteErr := ioutil.WriteFile(fmt.Sprintf("%s/index.json", scrapedDataDir), indexJSON, 0755); indexWriteErr != nil {
		log.Fatalln(indexWriteErr)
	}

	forums := make([]Forum, 0)
	for _, f := range index.Forums {
		forums = append(forums, f...)
	}

	sort.Slice(forums, func(i, j int) bool {
		return forums[i].URL < forums[j].URL
	})

	forums[0].Scrape()

	//for _, forum := range forums {
	//	forum.Scrape()
	//}

	//fmt.Println(string(indexJSON))

	//c := colly.NewCollector()

	// Find and visit all links

	//c.OnHTML("tbody", func(e *colly.HTMLElement) {
	//	e.DOM.Find("td.tcat").Each(func(i int, s *goquery.Selection) {
	//		category := s.Find("a")
	//		fmt.Println(category.Text())
	//	})
	//})
	//
	//c.OnRequest(func(r *colly.Request) {
	//	fmt.Println("Visiting", r.URL)
	//})
	//
	//_ = c.Visit(forumURL.String())
}
