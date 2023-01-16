package indexer

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
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

type NopCloser struct {
	w io.Writer
}

func NewNopCloser(w io.Writer) *NopCloser {
	if w == nil {
		return nil
	}
	return &NopCloser{w: w}
}

func (n *NopCloser) Write(p []byte) (int, error) {
	return n.w.Write(p)
}

func (n *NopCloser) Close() error {
	return nil
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

type ZipOutputFS struct {
	base string
	w    *zip.Writer
}

func NewZIpOutputFS(base string, w io.Writer) *ZipOutputFS {
	return &ZipOutputFS{base: base, w: zip.NewWriter(w)}
}

func (f *ZipOutputFS) Open(name string) (WriteCloser, error) {
	w, err := f.w.Create(path.Join(f.base, name))
	return NewNopCloser(w), err
}

type HTTPFS struct {
	c    *http.Client
	base string
}

func NewHTTPFS(c *http.Client, base string) *HTTPFS {
	if c == nil {
		c = http.DefaultClient
	}
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
