package scrape

import (
	"github.com/gocolly/colly"
	"github.com/plombardi89/ab-scrape/pkg/types"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

type Scraper struct {
	ForumIndexLocation url.URL
	OutputDirectory    string
}

func (s *Scraper) Scrape() []error {
	errors := make([]error, 0)
	c := colly.NewCollector()

	index := &types.ForumIndex{
		URL: s.ForumIndexLocation.String(),
	}

	errs := index.Scrape(c)
	if len(errs) > 0 {
		errors = append(errors, errs...)
		return errors
	}

	if err := s.ensureStorageExists(); err != nil {
		errors = append(errors, err)
		return errors
	}

	forumIndexJSON, err := index.ToJSON()
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	fileWriteErr := ioutil.WriteFile(filepath.Join(s.OutputDirectory, ""), forumIndexJSON, 0755)
	if fileWriteErr != nil {
		errors = append(errors, fileWriteErr)
		return errors
	}

	return nil
}

func (s *Scraper) ensureStorageExists() error {
	return os.MkdirAll(s.OutputDirectory, 0755)
}

func (s *Scraper) storeForumIndex() error {

}
