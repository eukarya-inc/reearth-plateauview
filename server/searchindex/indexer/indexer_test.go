package indexer

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var config = &Config{
	IdProperty: "gml_id",
	Indexes: map[string]Index{
		"名称":           {Kind: "enum"},
		"用途":           {Kind: "enum"},
		"住所":           {Kind: "enum"},
		"名建物利用現況_中分類称": {Kind: "enum"},
		"建物利用現況_小分類":   {Kind: "enum"},
	},
}

func TestIndexer(t *testing.T) {
	dir, err := os.Stat("testdata")
	if err != nil || !dir.IsDir() {
		t.Skip()
	}

	b := bytes.NewBuffer(nil)
	zw := zip.NewWriter(b)

	input := NewFSFS(os.DirFS("testdata"))
	output := NewZipOutputFS(zw, "")
	indexer := NewIndexer(config, input, output)

	err = indexer.BuildAndWrite()
	assert.NoError(t, err)

	br := bytes.NewReader(b.Bytes())
	zr, err := zip.NewReader(br, 0)
	assert.NoError(t, err)
	err = extractZip(zr, filepath.Join("testdata", "result"))
	assert.NoError(t, err)
}

func extractZip(zr *zip.Reader, base string) error {
	_ = os.MkdirAll(base, os.ModePerm)
	for _, f := range zr.File {
		f := func() error {
			r, err := f.Open()
			if err != nil {
				return err
			}
			defer r.Close()
			f, err := os.Create(filepath.Join(base, f.Name))
			if err != nil {
				return err
			}
			_, err = io.Copy(f, r)
			return err
		}
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
