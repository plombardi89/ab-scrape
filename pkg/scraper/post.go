package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	// ID is the system identifier for the thread.
	ID int

	// ThreadID is the system identifier for the thread the post belongs to.
	ThreadID int

	// Author is the name of the authoring user for the post.
	Author string

	// CreatedAt is when the post was written.
	CreatedAt time.Time

	// Content is the unstructured content of the post.
	Content string

	// URL is the location where the post can be retrieved.
	URL string
}

func NewPostFromSelection(s *goquery.Selection) (*Post, error) {
	p := &Post{}

	idAttr, found := s.Attr("id")
	if !found {
		return nil, errors.New("unable to find id attribute for post")
	}

	rawPostID := strings.Replace(idAttr, "post", "", 1)
	postID, err := strconv.Atoi(rawPostID)
	if err != nil {
		return nil, err
	}

	postMessageSelection := s.Find(fmt.Sprintf("div#post_message_%d", postID))
	p.ID = postID
	p.Content = postMessageSelection.Text()
	p.Author = s.Find("a.bigusername").Text()

	dateData := s.Find(fmt.Sprintf("#post%d > tbody:nth-child(1) > tr:nth-child(1) > td:nth-child(1)", p.ID)).Text()

	// Posts are not always marked with absolute dates ... sometimes it's relative, for example, "Yesterday, 06:24 AM".
	// In this case we need to convert that to something useful so we can process it later.
	now := time.Now()
	dateData = strings.TrimSpace(dateData)
	if strings.HasPrefix(dateData, "Yesterday") {
		yesterday := time.Now().Add(-(time.Hour * 24))
		dateData = strings.Replace(dateData, "Yesterday", yesterday.Format("1-2-2006"), 1)
	} else if strings.HasPrefix(dateData, "Today") {
		dateData = strings.Replace(dateData, "Today", now.Format("1-2-2006"), 1)
	}

	postCreatedAt, err := time.Parse("1-2-2006, 03:04 PM", dateData)
	if err != nil {
		log.Fatalln(err)
	}
	p.CreatedAt = postCreatedAt

	return p, nil
}

func (p *Post) ToJSON() ([]byte, error) {
	bytes, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
