package types

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"log"
	"strconv"
	"strings"
)

type Forum struct {
	// System ID of the forum
	ID int

	// URL is the location where the Forum thread list can be retrieved.
	URL string

	// ForumIndexURL is the location where the Forum Index can be retrieved.
	ForumIndexURL string

	// Category is the named category where the Forum is located in the Forum Index.
	Category string

	// Name is the name of the forum.
	Name string

	// LastPage is the last numbered page for the Forum.
	LastPage int
}

func (f *Forum) Scrape(c *colly.Collector, writer FileSystemWriter) {
	// Grab the <table> with the page navigation in it.
	c.OnXML("/html/body/div/div/div/form/table[1]", func(e *colly.XMLElement) {
		node, ok := e.DOM.(*html.Node)
		if !ok {
			log.Fatalln("failed to cast to *html.Node")
		}

		doc := goquery.NewDocumentFromNode(node)
		f.extractLastPageNumber(doc)

		for i := 2; i <= f.LastPage; i++ {
			_ = c.Visit(fmt.Sprintf("%sf=%d&order=desc&page=%d", f.ForumIndexURL, f.ID, i))
		}
	})

	threadBitsID := fmt.Sprintf("tbody#threadbits_forum_%d", f.ID)
	c.OnHTML(threadBitsID, func(e *colly.HTMLElement) {
		e.DOM.Find("tr").EachWithBreak(func(i int, s *goquery.Selection) bool {
			thread, err := ThreadFromSelection(f, s)
			if err != nil {
				log.Fatalln(err)
			}

			thread.ForumID = f.ID
			if writerErr := writer.WriteThread(thread); writerErr != nil {
				return false
			}

			_ = c.Visit(thread.URL)
			return true
		})
	})

	//c.OnHTML("table[id*='post']", func(e *colly.HTMLElement) {
	//	post, err := NewPostFromSelection(e.DOM)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	err = writer.WritePost(post)
	//	if err != nil {
	//		log.Printf("ERROR: unable to write post: %v", err)
	//		return
	//	}
	//
	//	//log.Printf("INFO: wrote post %s\n", postPath)
	//})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("INFO: visiting %s\n", r.URL)
	})

	_ = c.Visit(f.URL)
}

func (f *Forum) extractLastPageNumber(doc *goquery.Document) {
	doc.Find("div.pagenav > table > tbody > tr > td.vbmenu_control").Each(func(i int, s *goquery.Selection) {
		pageAOfB := strings.Trim(strings.Replace(strings.Replace(s.Text(), "Page", "", 1), "of", "", 1), " ")
		minMaxPages := strings.Fields(pageAOfB)
		lastPage, err := strconv.Atoi(minMaxPages[1])
		if err != nil {
			log.Fatalln(err)
		}

		f.LastPage = lastPage
	})
}
