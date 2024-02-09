package preparegspatialjp

import (
	"fmt"
	"io/fs"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
)

// TODO: get dic data for fld
func generatePlateauIndexItem(seed *IndexSeed, name string, size uint64, f fs.FS) (*IndexItem, error) {
	data := map[string]plateauItemSeed{}

	if err := fs.WalkDir(f, "", func(p string, d fs.DirEntry, err error) error {
		base := path.Base(p)
		featureType := extractFeatureType(base)
		if featureType == "" {
			return nil
		}

		lod := extractLOD(base)

		if _, ok := data[featureType]; !ok {
			data[featureType] = plateauItemSeed{
				Type:  featureType,
				Name:  base,
				Title: featureTypees[featureType],
				LOD:   nil,
			}
		}
		if lod > -1 {
			d := data[featureType]
			d.LOD = append(d.LOD, lod)
			data[featureType] = d
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to walk plateau zip: %w", err)
	}

	items := plateauItems(data)
	children := []*IndexItem{}
	for _, d := range items {
		children = append(children, d.Item())
	}

	return &IndexItem{
		Name:     fmt.Sprintf("**%s**：3D Tiles, MVT（v%d）(%s)", name, seed.V, humanize.Bytes(size)),
		Children: children,
	}, nil
}

func extractFeatureType(name string) string {
	for _, f := range featureTypes {
		if strings.Contains("_"+name+"_", f) {
			return f
		}
	}

	return ""
}

func extractDataFormat(name string) string {
	if strings.Contains(name, "_3dtiles") {
		return "3D Tiles"
	}
	if strings.Contains(name, "_mvt") {
		return "MVT"
	}
	return ""
}

var lod = []int{0, 1, 2, 3, 4}

func extractLOD(name string) int {
	for _, l := range lod {
		if strings.Contains(name, fmt.Sprintf("_lod%d", l)) {
			return l
		}
	}
	return -1
}

type plateauItemSeed struct {
	Type  string
	Name  string
	Title string
	LOD   []int
}

func (p plateauItemSeed) Item() *IndexItem {
	format := extractDataFormat(p.Name)

	if len(p.LOD) == 0 {
		return &IndexItem{
			Name: fmt.Sprintf("**%s**：%s（%s）", p.Name, p.Title, format),
		}
	}
	if len(p.LOD) == 1 {
		return &IndexItem{
			Name: fmt.Sprintf("**%s**：%s（LOD%d, %s）", p.Name, p.Title, p.LOD[0], format),
		}
	}

	children := make([]*IndexItem, len(p.LOD))
	for i, l := range p.LOD {
		children[i] = &IndexItem{
			Name: fmt.Sprintf("%s（LOD%d, %s）", p.Title, l, format),
		}
	}

	return &IndexItem{
		Name:     fmt.Sprintf("**%s**：%s（%s）", p.Name, p.Title, format),
		Children: children,
	}
}

func plateauItems(m map[string]plateauItemSeed) []plateauItemSeed {
	items := make([]plateauItemSeed, 0, len(m))
	for _, v := range m {
		items = append(items, v)
	}

	// sort by key. follow fratureTypes order
	sort.Slice(items, func(i, j int) bool {
		index1 := slices.Index(featureTypes, items[i].Type)
		index2 := slices.Index(featureTypes, items[j].Type)
		return index1 < index2
	})

	return items
}
