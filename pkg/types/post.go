package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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

	return p, nil
}

func (p *Post) ToJSON() ([]byte, error) {
	bytes, err := json.MarshalIndent(p, "", "    ")
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

//type FilePostWriter struct {
//	outputRoot string
//}
//
//func (w* FilePostWriter) Write(p *Post) (string, error) {
//	threadPath := fmt.Sprintf("%s/threads/%d", main.outputRoot, p.ThreadID)
//	if err := os.MkdirAll(threadPath, 0755); err != nil {
//		return "", err
//	}
//
//	jsonBytes, marshallErr := json.MarshalIndent(p, "", "  ")
//	if marshallErr != nil {
//		return "", marshallErr
//	}
//
//	postPath := fmt.Sprintf("%s/post-%d.json", threadPath, p.ID)
//	if err := ioutil.WriteFile(postPath, jsonBytes, 0755); err != nil {
//		return "", err
//	}
//
//	return postPath, nil
//}
