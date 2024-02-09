package preparegspatialjp

import (
	"fmt"
	"io/fs"

	"github.com/dustin/go-humanize"
)

func generateGenericdIndexItems(seed *IndexSeed, data []GspatialjpIndexItemGroup) (res []*IndexItem, err error) {
	for _, d := range data {
		u := d.AssetURL()
		if u == "" {
			continue
		}
		if d.Type == "" {
			d.Type = "ユースケースデータ"
		}

		size, err := httpSize(u)
		if err != nil {
			return nil, fmt.Errorf("%s からファイルをダウンロードできませんでした。: %w", u, err)
		}

		res = append(res, &IndexItem{
			Name: fmt.Sprintf(
				"**%s**：%s (%s)",
				fileNameFromURL(u),
				d.Type,
				humanize.Bytes(size),
			),
			Children: []*IndexItem{
				{
					Name: d.Name,
				},
			},
		})
	}
	err = fs.SkipDir
	return
}
