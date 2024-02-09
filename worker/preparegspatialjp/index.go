package preparegspatialjp

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/dustin/go-humanize"
)

type IndexSeed struct {
	CityName       string
	Year           int
	V              int
	CityGMLZipPath string
	PlateuaZipPath string
	RelatedZipPath string
	// name: path
	Generic map[string]string
}

type IndexItem struct {
	Name     string
	Children []*IndexItem
}

func GenerateIndex(ctx context.Context, seed *IndexSeed) (string, error) {
	citygmlFS, citygmlFSCloser, err := openZip(seed.CityGMLZipPath)
	if err != nil {
		return "", fmt.Errorf("failed to open citygml zip: %w", err)
	}

	plateauFS, plateauFSCloser, err := openZip(seed.PlateuaZipPath)
	if err != nil {
		return "", fmt.Errorf("failed to open plateau zip: %w", err)
	}

	relatedFS, relatedFSCloser, err := openZip(seed.RelatedZipPath)
	if err != nil {
		return "", fmt.Errorf("failed to open related zip: %w", err)
	}

	defer func() {
		if citygmlFSCloser != nil {
			_ = citygmlFSCloser()
		}
		if plateauFSCloser != nil {
			_ = plateauFSCloser()
		}
		if relatedFSCloser != nil {
			_ = relatedFSCloser()
		}
	}()

	citygml, err := generateCityGMLIndexItem(ctx, seed, citygmlFS)
	if err != nil {
		return "", fmt.Errorf("failed to generate index items: %w", err)
	}

	plateau, err := generatePlateauIndexItem(ctx, seed, plateauFS)
	if err != nil {
		return "", fmt.Errorf("failed to generate index items: %w", err)
	}

	related, err := generateRelatedIndexItem(ctx, seed, relatedFS)
	if err != nil {
		return "", fmt.Errorf("failed to generate index items: %w", err)
	}

	generics, err := generateGenericdIndexItems(ctx, seed, seed.Generic)
	if err != nil {
		return "", fmt.Errorf("failed to generate index items: %w", err)
	}

	items := append([]*IndexItem{citygml, plateau, related}, generics...)

	leading := fmt.Sprintf("%sの%d年度版データを標準製品仕様書V%dに基づいて作成した提供データ目録です。\n\n", seed.CityName, seed.Year, seed.V)
	return leading + renderIndexItems(items, 0), nil
}

func renderIndexItems(t []*IndexItem, depth int) (res string) {
	for _, c := range t {
		res += renderIndexItem(c, depth)
	}
	return
}

func renderIndexItem(t *IndexItem, depth int) (res string) {
	if t == nil {
		return ""
	}
	res = fmt.Sprintf("%s- %s\n", strings.Repeat("  ", depth*2), t.Name)
	for _, c := range t.Children {
		res += renderIndexItem(c, depth+1)
	}
	return
}

func generateCityGMLIndexItem(ctx context.Context, seed *IndexSeed, f fs.FS) (*IndexItem, error) {
	return walk(f, "", "/", func(path string, d fs.DirEntry, err error) (*IndexItem, error) {
		panic("not implemented") // TODO: Implement
	})
}

func generatePlateauIndexItem(ctx context.Context, seed *IndexSeed, f fs.FS) (*IndexItem, error) {
	return walk(f, "", "/", func(path string, d fs.DirEntry, err error) (*IndexItem, error) {
		panic("not implemented") // TODO: Implement
	})
}

func generateRelatedIndexItem(ctx context.Context, seed *IndexSeed, f fs.FS) (*IndexItem, error) {
	return walk(f, "", "/", func(path string, d fs.DirEntry, err error) (*IndexItem, error) {
		panic("not implemented") // TODO: Implement
	})
}

func generateGenericdIndexItems(ctx context.Context, seed *IndexSeed, urls map[string]string) (res []*IndexItem, err error) {
	for name, url := range urls {
		size, err := httpSize(url)
		if err != nil {
			return nil, fmt.Errorf("%s からファイルをダウンロードできませんでした。: %w", url, err)
		}

		res = append(res, &IndexItem{
			Name: fmt.Sprintf("**%s**：ユースケースデータ (%s)", fileNameFromURL(url), humanize.Bytes(uint64(size))),
			Children: []*IndexItem{
				{
					Name: name,
				},
			},
		})
	}
	return
}
