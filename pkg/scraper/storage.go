package scraper

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileSystemWriter struct {
	StorageRoot string
}

func (w *FileSystemWriter) WriteForumIndex(f *ForumIndex) error {
	if err := os.MkdirAll(w.StorageRoot, 0755); err != nil {
		return err
	}

	data, err := f.ToJSON()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(w.StorageRoot, "index.json"), data, 0755); err != nil {
		return err
	}

	return nil
}

func (w *FileSystemWriter) WriteThread(t *Thread) error {
	threadPath := filepath.Join(w.StorageRoot, fmt.Sprintf("thread-%d", t.ID))
	if err := os.MkdirAll(threadPath, 0755); err != nil {
		return err
	}

	data, err := t.ToJSON()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(threadPath, "thread.json"), data, 0755); err != nil {
		return err
	}

	return nil
}

func (w *FileSystemWriter) WritePost(p *Post) error {
	threadPath := filepath.Join(w.StorageRoot, fmt.Sprintf("thread-%d", p.ThreadID))
	if err := os.MkdirAll(threadPath, 0755); err != nil {
		return err
	}

	data, err := p.ToJSON()
	if err != nil {
		return err
	}

	postPath := filepath.Join(threadPath, fmt.Sprintf("%d.json", p.ID))
	if err := ioutil.WriteFile(postPath, data, 0755); err != nil {
		return err
	}

	return nil
}
