package jssquish

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Simple in-memory repository to use for testing
type MemRepository struct {
	files   map[string]string
	checked []string
	opened  []string
}

func NewMemRepository(files map[string]string) *MemRepository {
	memRepo := &MemRepository{
		files: make(map[string]string),
	}
	for path, src := range files {
		memRepo.files[filepath.Clean(path)] = src
	}
	return memRepo
}

func (mr *MemRepository) IsFile(path string) bool {
	mr.checked = append(mr.checked, path)
	_, ok := mr.files[path]
	return ok
}

func (mr MemRepository) Open(path string) (io.ReadCloser, error) {
	mr.opened = append(mr.opened, path)
	src, ok := mr.files[path]
	if !ok {
		return nil, fmt.Errorf("Path does not exist: %s", path)
	}
	return ioutil.NopCloser(strings.NewReader(src)), nil
}

func (mr *MemRepository) Close() error {
	return nil
}
