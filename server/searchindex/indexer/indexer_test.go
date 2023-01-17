package indexer

import (
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

	input := NewFSFS(os.DirFS("testdata"))
	output := NewOSOutputFS(filepath.Join("testdata", "result"))
	indexer := NewIndexer(config, input, output)

	err = indexer.BuildAndWrite()
	assert.NoError(t, err)
}
