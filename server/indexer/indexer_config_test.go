package indexer

import (
	"testing"
)

func TestParseIndexerConfig(t *testing.T) {
	_, err := ParseIndexerConfigFile("testdata/config.json")

	if err != nil {
		t.Errorf("failed to parse the indexerconfig json")
	}
}
