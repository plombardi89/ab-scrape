package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type Thread struct {
	// ID is the system identifier for the Thread
	ID int

	// ForumID is the system identifier for the Forum.
	ForumID int

	// URL is the location where the Thread can be retrieved.
	URL string

	// Title is the name or Title of the Thread.
	Title string
}

func ThreadFromSelection(f *Forum, s *goquery.Selection) (*Thread, error) {
	t := &Thread{
		ForumID: f.ID,
	}

	threadTitleHTML := s.Find("a[id*='thread_title_']")
	idAttr, found := threadTitleHTML.Attr("id")
	if !found {
		return nil, errors.New("no thread ID attribute found in the HTML structure")
	}

	rawThreadID := strings.Replace(idAttr, "thread_title_", "", 1)
	threadID, stringToIntErr := strconv.Atoi(rawThreadID)
	if stringToIntErr != nil {
		return nil, stringToIntErr
	}

	t.ID = threadID
	t.Title = threadTitleHTML.Text()
	t.URL = fmt.Sprintf("%s/showthread.php?t=%d", f.ForumIndexURL, t.ID)

	return t, nil
}

func (t *Thread) ToJSON() ([]byte, error) {
	bytes, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
