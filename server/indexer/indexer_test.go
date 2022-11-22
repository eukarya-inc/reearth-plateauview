package indexer

import "testing"

var (
	_TEST_TILESET_JSON = "testdata/tileset.json"
	_TEST_CONFIG_JSON = "testdata/config.json"
	_TEST_RESUILT_PATH = "testdata/result"
)

func TestIndexer(t *testing.T) {
	if err := Indexer(_TEST_TILESET_JSON, _TEST_CONFIG_JSON, _TEST_RESUILT_PATH); err != nil {
		t.Errorf("failed to generate indexes: %v", err)
	}
}
