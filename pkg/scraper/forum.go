package scraper

import (
	"bytes"
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
		f.LastPage = extractLastPageNumber(doc)

		for i := 2; i <= f.LastPage; i++ {
			_ = c.Visit(fmt.Sprintf("%s/forumdisplay.php?f=%d&order=desc&page=%d", f.ForumIndexURL, f.ID, i))
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

			postCollector := colly.NewCollector()
			postCollector.OnRequest(func(r *colly.Request) {
				log.Printf("INFO: visiting thread %s\n", r.URL)
			})

			postCollector.OnResponse(func(response *colly.Response) {
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
				if err != nil {
					log.Fatalln(err)
				}

				threadPages := extractLastPageNumber(doc)
				doc.Find("table[id*='post']").Each(func(i int, selection *goquery.Selection) {
					post, err := NewPostFromSelection(selection)
					if err != nil {
						log.Fatalln(err)
					}
					post.ThreadID = thread.ID

					err = writer.WritePost(post)
					if err != nil {
						log.Printf("ERROR: unable to write post: %v", err)
						return
					}
				})

				for i := 2; i <= threadPages; i++ {
					_ = postCollector.Visit(fmt.Sprintf("%s/showthread.php?t=%d&page=%d", f.ForumIndexURL, thread.ID, i))
				}
			})

			_ = postCollector.Visit(thread.URL)
			return true
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Printf("INFO: visiting %s\n", r.URL)
	})

	_ = c.Visit(f.URL)
}

func extractLastPageNumber(doc *goquery.Document) int {
	var result int

	doc.Find("div.pagenav > table > tbody > tr > td.vbmenu_control").EachWithBreak(func(i int, s *goquery.Selection) bool {
		pageAOfB := strings.Trim(strings.Replace(strings.Replace(s.Text(), "Page", "", 1), "of", "", 1), " ")
		minMaxPages := strings.Fields(pageAOfB)
		lastPage, err := strconv.Atoi(minMaxPages[1])
		if err != nil {
			log.Fatalln(err)
		}

		result = lastPage
		return false
	})

	return result
}
