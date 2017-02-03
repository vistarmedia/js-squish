package jssquish

import (
	"io"
)

func Main(
	repo Repository,
	entrypoint string,
	environment *string,
	out io.Writer) error {
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

	if environment == nil {
		return fs.Create(entrypoint)
	} else {
		return fs.CreateWithNodeEnv(entrypoint, environment)
	}
}
