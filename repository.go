package jssquish

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	pth "path"
	"path/filepath"
)

// A `Repository` is a collection of files. The node resolution algorithm will
// make many checks to see if various files exist, so it's assume the `IsFile`
// implementation will be fast.
type Repository interface {
	IsFile(path string) bool
	Open(path string) (io.ReadCloser, error)
	Close() error
}

// `Repository` implementation which loads the contents of all files into memory
// on construction. Because consumers will also likely copy the contents into
// memory, this will not work for repositories which have a deflated size
// larger than a gig or two.
type MemoryRepository map[string]*bytes.Buffer

// Creates a new `MemoryRepository` from a `js_tar`. It will load the entire
// decompressed contents of the repository into memory.
func NewMemoryJsTarRepository(f *os.File) (MemoryRepository, error) {
	repo := make(MemoryRepository)

	err := iterateTar(f, func(hdr *tar.Header, fi os.FileInfo,
		r io.Reader) error {

		buf := &bytes.Buffer{}
		if _, err := buf.ReadFrom(r); err != nil {
			return err
		}
		repo[filepath.Clean(hdr.Name)] = buf
		return nil
	})

	if err != nil {
		return nil, err
	}

	return repo, nil
}

// Checks the cache to see if the requested path has been seen
func (repo MemoryRepository) IsFile(path string) bool {
	_, ok := repo[path]
	return ok
}

// Returns the contents of the in-memory `bytes.Buffer` for the requested file,
// or errors. Note that it's safe to read from this buffer multiple times, which
// is not the case for general `Repository` implementations.
func (repo MemoryRepository) Open(path string) (io.ReadCloser, error) {
	if buf, ok := repo[path]; ok {
		return ioutil.NopCloser(buf), nil
	} else {
		return nil, fmt.Errorf("Could not open path: %s", path)
	}
}

func (mr MemoryRepository) Close() error {
	mr = make(MemoryRepository)
	return nil
}

// Repository which keeps all files on disk, and uses it as a metadata store.
// Note that it's important to call `Close` on this `Repository`, as it may
// otherwise leak files. Similarly, it's important to not point this right to a
// directory of important files, as it will remove them when closed.
// It will not preserve ownership, mode, or other metadata when expanding a
// file.
// When expanding, it will create a cache of filenames to keep `IsFile` snappy.
type DiskRepository struct {
	root  string
	cache map[string]bool
}

// Creates a new `DiskRepository` by expanding the contents of a `js_tar` to a
// temporary directory.
func NewDiskJsTarRepository(f *os.File) (*DiskRepository, error) {
	dir, err := ioutil.TempDir(os.Getenv("TMPDIR"), "js-squish")
	if err != nil {
		return nil, err
	}

	files := make(map[string]bool)
	err = iterateTar(f, func(hdr *tar.Header, fi os.FileInfo, r io.Reader) error {

		dstPath := pth.Join(dir, hdr.Name)
		if dstPath == dir {
			return nil
		}
		if fi.IsDir() {
			return nil
		}
		if os.MkdirAll(pth.Dir(dstPath), 0); err != nil {
			return err
		} else {
			files[filepath.Clean(hdr.Name)] = true
		}

		if dst, err := os.Create(dstPath); err != nil {
			return err
		} else if _, err := io.Copy(dst, r); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	repo := &DiskRepository{
		root:  dir,
		cache: files,
	}
	return repo, nil
}

// Consults the cache to see if this requested file has been extracted
func (dr *DiskRepository) IsFile(path string) bool {
	_, ok := dr.cache[path]
	return ok
}

// Opens the file directly from disk. This will return the actual `os.File`
// instance, so its up the caller to close it in order to not leak handles.
// Also note that this won't prevent walking up and out of a directory, and is
// generally not safe to consume untrusted inputs.
func (dr *DiskRepository) Open(path string) (io.ReadCloser, error) {
	if !dr.IsFile(path) {
		return nil, fmt.Errorf("Could not open path: %s in %s", path, dr.root)
	}
	absolute := pth.Join(dr.root, path)
	return os.Open(absolute)
}

// Removes the temporary directory this implementation used to store the
// expanded files. This is blocking and may not return quickly.
func (dr *DiskRepository) Close() error {
	return os.RemoveAll(dr.root)
}

type tarFlowControl uint8

func (tarFlowControl) Error() string {
	return "Stop Iteration"
}

const stopIteration = tarFlowControl(0)

// Iterate over a tar file, invoking each function as a file is encountered. To
// abort iteration, the function can return the above defined `stopIteration`
// error value.
// Because the tar reading implementation is serial and stateful, It is not safe
// to call this function with the same file concurrently.
func iterateTar(f *os.File, each func(*tar.Header, os.FileInfo,
	io.Reader) error) error {

	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	gzf, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzf.Close()

	tr := tar.NewReader(gzf)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		err = each(hdr, hdr.FileInfo(), tr)
		if err == stopIteration {
			return nil
		}
		if err != nil {
			return err
		}
	}
	return nil
}
