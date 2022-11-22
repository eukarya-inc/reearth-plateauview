package indexer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type IndexConfig struct {
	Kind string `json:"kind"`
}

type IndexesConfig struct {
	IdProperty string                 `json:"idProperty"`
	Indexes    map[string]IndexConfig `json:"indexes"`
}

func IndexerConfigFromJson(data io.Reader) (*IndexesConfig, error) {
	var ic *IndexesConfig
	if err := json.NewDecoder(data).Decode(&ic); err != nil {
		return nil, fmt.Errorf("decode failed: %v", err)
	}
	return ic, nil
}

func ParseIndexerConfigFile(fileName string) (*IndexesConfig, error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open failed: %v", err)
	}
	res, err := IndexerConfigFromJson(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("json conversion failed: %v", err)
	}
	return res, nil
}
