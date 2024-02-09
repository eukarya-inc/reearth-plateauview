package preparegspatialjp

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/dustin/go-humanize"
)

func generateRelatedIndexItem(seed *IndexSeed, name string, size uint64, f fs.FS) (*IndexItem, error) {
	items := []string{}

	if err := fs.WalkDir(f, "", func(p string, d fs.DirEntry, err error) error {
		if p == "" {
			return nil
		}

		if t := detectRelatedDataType(p); t != "" {
			items = append(items, t)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk related zip: %w", err)
	}

	children := []*IndexItem{}
	for _, d := range items {
		children = append(children, &IndexItem{
			Name: fmt.Sprintf("**%s**", relatedDataTypeMap[d]),
		})
	}

	return &IndexItem{
		Name:     fmt.Sprintf("**%s**：関連データセット（v%d）(%s)", name, seed.V, humanize.Bytes(size)),
		Children: children,
	}, nil
}

func detectRelatedDataType(name string) string {
	for _, t := range relatedDataTypes {
		if t == name || strings.Contains(name, "_"+t) {
			return t
		}
	}
	return ""
}
