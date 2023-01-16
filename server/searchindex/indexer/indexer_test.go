package indexer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestIndexer(t *testing.T) {
	dir, err := os.Stat("testdata")
	if err != nil || !dir.IsDir() {
		t.Skip()
	}

	input := NewFSFS(os.DirFS("testdata"))
	output := NewOSOutputFS(filepath.Join("testdata", "result"))
	var config Config
	assert.NoError(t, json.Unmarshal(lo.Must(os.ReadFile("testdata/config.json")), &config))
	indexer := NewIndexer(&config, input, output)

	err = indexer.BuildAndWrite()
	assert.NoError(t, err)
}
