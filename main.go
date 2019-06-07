package main

import (
	"github.com/gocolly/colly"
	"github.com/plombardi89/ab-scrape/pkg/types"
	"log"
	"net/url"
	"path/filepath"
	"time"
)

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/PuerkitoBio/goquery"
//	"github.com/gocolly/colly"
//	"golang.org/x/net/html"
//	"io/ioutil"
//	"log"
//	"net/url"
//	"os"
//	"sort"
//	"strconv"
//	"strings"
//	"time"
//)
//
var (
	outputRoot = "dump/"
	forumURL   = "http://archboston.org/community"
)

//
//var outputRoot = "dump/"
//
//type ForumIndex struct {
//	URL       string             `json:""`
//	Forums    map[string][]Forum `json:""`
//}
//
//type Forum struct {
//	ID          int    `json:""`
//	Category    string `json:""`
//	Name        string `json:""`
//	Description string `json:""`
//	URL         string `json:""`
//	lastPage    int
//}
//
//type Thread struct {
//	ID       int
//	ForumID  int
//	Title    string
//	URL      string
//	lastPage int
//}
//
//func printSelection(s *goquery.Selection) {
//	fmt.Println("===================")
//	fmt.Println(goquery.OuterHtml(s))
//	fmt.Println("===================")
//}
//
//func ThreadFromSelection(s *goquery.Selection) (*Thread, error) {
//	t := &Thread{}
//
//	threadTitleHTML := s.Find("a[id*='thread_title_']")
//	idAttr, found := threadTitleHTML.Attr("id")
//	if !found {
//		log.Fatalln("no thread id attribute")
//	}
//
//	rawThreadID := strings.Replace(idAttr, "thread_title_", "", 1)
//	threadID, convErr := strconv.Atoi(rawThreadID)
//	if convErr != nil {
//		log.Fatalln(convErr)
//	}
//
//	t.ID = threadID
//	t.Title = threadTitleHTML.Text()
//	t.URL = fmt.Sprintf("%s/showthread.php?t=%d", forumURL, t.ID)
//
//
//	//fmt.Println(threadTItleHTML.Attr("id"))
//
//	//printSelection(s.Find("td[id*='td_threadtitle_']"))
//	//printSelection(s.Children())
//
//	return t, nil
//}
//
//func (t *Thread) Scrape() {
//
//}
//
//func (f *Forum) extractLastPageNumber(doc *goquery.Document) {
//	doc.Find("div.pagenav > table > tbody > tr > td.vbmenu_control").Each(func(i int, s *goquery.Selection) {
//		pageAOfB := strings.Trim(strings.Replace(strings.Replace(s.Text(), "Page", "", 1), "of", "", 1), " ")
//		minMaxPages := strings.Fields(pageAOfB)
//		lastPage, err := strconv.Atoi(minMaxPages[1])
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		f.lastPage = lastPage
//	})
//}
//
//func (f *Forum) Scrape(c *colly.Collector) {
//
//	// Grab the <table> with the page navigation in it.
//	c.OnXML("/html/body/div/div/div/form/table[1]", func(e *colly.XMLElement) {
//		node, ok := e.DOM.(*html.Node)
//		if !ok {
//			log.Fatalln("failed to cast to *html.Node")
//		}
//
//		doc := goquery.NewDocumentFromNode(node)
//		f.extractLastPageNumber(doc)
//
//		//threadBitsID := fmt.Sprintf("tbody#threadbits_forum_%d", f.ID)
//		////fmt.Println(threadBitsID)
//		//doc.Find(threadBitsID).Each(func(i int, selection *goquery.Selection) {
//		//	fmt.Println(goquery.OuterHtml(selection))
//		//})
//	})
//
//	threadBitsID := fmt.Sprintf("tbody#threadbits_forum_%d", f.ID)
//	c.OnHTML(threadBitsID, func(e *colly.HTMLElement) {
//		e.DOM.Find("tr").Each(func(i int, s *goquery.Selection) {
//			thread, err := ThreadFromSelection(s)
//			if err != nil {
//				log.Fatalln(err)
//			}
//
//			thread.ForumID = f.ID
//			//fmt.Println(thread)
//
//			_ = c.Visit(thread.URL)
//		})
//	})
//
//	postWriter := FilePostWriter{}
//	c.OnHTML("table[id*='post']", func(e *colly.HTMLElement) {
//		post, err := NewPostFromSelection(e.DOM)
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		postPath, err := postWriter.Write(post)
//		if err != nil {
//			return
//		}
//
//		log.Printf("INFO: wrote post %s\n", postPath)
//	})
//
//	c.OnRequest(func(r *colly.Request) {
//		log.Printf("INFO: visiting %s\n", r.URL)
//	})
//
//	_ = c.Visit(f.URL)
//}
//
//func (f *ForumIndex) Scrape(c *colly.Collector) {
//	// Use an XPath query to select the 4th table which contains the forum index then we build up the ForumIndex
//	// struct based on the discovered contents
//	c.OnXML("/html/body/div/div/div/table[3]", func(e *colly.XMLElement) {
//		node, ok := e.DOM.(*html.Node)
//		if !ok {
//			log.Fatalln("failed to cast to *html.Node")
//		}
//
//		doc := goquery.NewDocumentFromNode(node)
//
//		// parse the <tbody> elements in the table. Even count <tbody> contain the forum category information while
//		// the odd count <tbody> contain the actual forum name, description and hyperlinks.
//		categories := make([]string, 0)
//		forumsInCategory := make(map[string][]Forum)
//		doc.ChildrenFiltered("tbody").Each(func(i int, s *goquery.Selection) {
//			if i%2 == 0 {
//				categories = append(categories, s.Find("tr > td > a").Text())
//			} else {
//				category := categories[i-len(categories)]
//				forumsInCategory[category] = make([]Forum, 0)
//
//				s.Find("td.alt1Active").Each(func(i int, s *goquery.Selection) {
//					forumNameAndPath := s.Find("div:nth-child(1) > a")
//
//					forumNameAndPath.Attr("href")
//					href, found := forumNameAndPath.Attr("href")
//					if !found {
//						log.Fatalln("missing attribute 'href' on forum row")
//					}
//
//					hrefURL, hrefURLErr := url.Parse(href)
//					if hrefURLErr != nil {
//						log.Fatalln(hrefURLErr)
//					}
//pathlib
//					// forumID is stored as an integer value on the "f" query parameter
//					forumID, err := strconv.Atoi(hrefURL.Query().Get("f"))
//					if err != nil {
//						log.Fatalln(err)
//					}
//
//					f := Forum{
//						ID:          forumID,
//						URL:         fmt.Sprintf("%s/forumdisplay.php?f=%d", f.URL, forumID),
//						Category:    category,
//						Name:        forumNameAndPath.Text(),
//						Description: "", // TODO: grab the description
//					}
//
//					forumsInCategory[category] = append(forumsInCategory[category], f)
//				})
//			}
//		})
//
//		f.Forums = forumsInCategory
//	})
//
//	c.OnRequest(func(r *colly.Request) {
//		log.Printf("INFO: visiting %s\n", r.URL)
//	})
//
//	_ = c.Visit(f.URL)
//}

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

	//indexJSON, marshallErr := json.MarshalIndent(index, "", "    ")
	//if marshallErr != nil {
	//	log.Fatalln(marshallErr)
	//}
	//
	//if indexWriteErr := ioutil.WriteFile(fmt.Sprintf("%s/index.json", scrapedDataDir), indexJSON, 0755); indexWriteErr != nil {
	//	log.Fatalln(indexWriteErr)
	//}
	//
	//forums := make([]Forum, 0)
	//for _, f := range index.Forums {
	//	forums = append(forums, f...)
	//}
	//
	//sort.Slice(forums, func(i, j int) bool {
	//	return forums[i].URL < forums[j].URL
	//})
	//
	//forums[0].Scrape(colly.NewCollector())

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
