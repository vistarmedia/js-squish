package jssquish

import (
	"bytes"
	pth "path"
)

type srcEntry struct {
	id   int
	deps map[string]*srcEntry
}

// A `FileSet` maintains a unique set of fully-qualified source files. For each
// file added to the set, its require statements are parsed, and any transitive
// dependencies added to the `FileSet`
type FileSet struct {
	nextId   int
	repo     Repository
	resolver *Resolver
	writer   *Writer
	entries  map[string]*srcEntry
}

// Creates a `FileSet` with the given starting point to walk files. The value
// for `impt` will be resolved through the full node module resolution
// algorithm, so paths, directories with an index, and directories with a
// package.json are all valid.
func (fs *FileSet) Create(impt string) error {
	if err := fs.writer.Open(); err != nil {
		return err
	}

	if _, err := fs.add(impt, "."); err != nil {
		return err
	}

	return fs.writer.Close()
}

// Internally adds a import to this `FileSet` from the perspective of the
// directory `from`.
func (fs *FileSet) add(impt, from string) (*srcEntry, error) {
	// Determine an id for this import
	id := fs.nextId
	fs.nextId++

	// Resolve the import
	path, err := fs.resolver.Resolve(impt, from)
	if err != nil {
		return nil, err
	}

	// If it's already known, we're done!
	if entry, ok := fs.entries[path]; ok {
		return entry, nil
	}

	// Resolve all imports
	src, imports, err := fs.read(path)
	if err != nil {
		return nil, err
	}

	// Ensure each dependency is fully resolved, and add it as a dependency to the
	// current import
	pwd := pth.Dir(path)
	deps := make(map[string]*srcEntry)
	for _, impt := range imports {
		if entry, err := fs.add(impt, pwd); err != nil {
			return nil, err
		} else {
			deps[impt] = entry
		}
	}

	entry := &srcEntry{
		id:   id,
		deps: deps,
	}
	fs.entries[path] = entry

	// Write the contents out to the file
	if err := fs.writer.Write(src, id, deps); err != nil {
		return nil, err
	}

	return entry, nil
}

func (fs *FileSet) read(path string) (*bytes.Buffer, []string, error) {
	r, err := fs.repo.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer r.Close()

	src := &bytes.Buffer{}
	if _, err := src.ReadFrom(r); err != nil {
		return nil, nil, err
	}

	imports, err := ParseRequires(src, path)
	if err != nil {
		return nil, nil, err
	}

	return src, imports, nil
}
