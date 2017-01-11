package jssquish

import (
	"io"
)

func Main(repo Repository, entrypoint string, out io.Writer) error {
	var (
		resolver = NewResolver(repo)
		writer   = NewWriter(out)
	)

	fs := &FileSet{
		repo:     repo,
		resolver: resolver,
		writer:   writer,
		entries:  make(map[string]*srcEntry),
	}

	return fs.Create(entrypoint)
}
