package indexer

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type FS interface {
	Open(name string) (io.ReadCloser, error)
}

type OutputFS interface {
	Open(name string) (WriteCloser, error)
}

type WriteCloser interface {
	io.Writer
	io.Closer
}

type FSFS struct {
	fs fs.FS
}

func NewFSFS(f fs.FS) *FSFS {
	return &FSFS{fs: f}
}

func (f *FSFS) Open(name string) (io.ReadCloser, error) {
	file, err := f.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

type OSOutputFS struct {
	base string
}

func NewOSOutputFS(base string) *OSOutputFS {
	return &OSOutputFS{base: base}
}

func (f *OSOutputFS) Open(name string) (w WriteCloser, err error) {
	return os.Create(filepath.Join(f.base, name))
}

type HTTPFS struct {
	c    *http.Client
	base string
}

func NewHTTPFS(c *http.Client, base string) *HTTPFS {
	return &HTTPFS{c: c, base: base}
}

func (f *HTTPFS) Open(name string) (io.ReadCloser, error) {
	u, err := url.JoinPath(f.base, name)
	if err != nil {
		return nil, err
	}

	res, err := f.c.Get(u)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("status code is %d", res.StatusCode)
	}

	return res.Body, nil
}
