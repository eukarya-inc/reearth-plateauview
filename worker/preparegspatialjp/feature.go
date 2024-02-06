package preparegspatialjp

import (
	"context"
	"fmt"

	cms "github.com/reearth/reearth-cms-api/go"
)

type FeatureItem struct {
	ID      string   `json:"id,omitempty"`
	CityGML string   `json:"citygml,omitempty"`
	Data    []string `json:"data,omitempty"`
}

func getAllFeatureItems(ctx context.Context, c *cms.CMS, cityItem *CityItem) (map[string]FeatureItem, error) {
	items := map[string]FeatureItem{}

	for key, ref := range cityItem.References {
		if ref == "" {
			continue
		}
		item, err := c.GetItem(ctx, ref, true)
		if err != nil {
			return nil, fmt.Errorf("failed to get item: %w", err)
		}

		fi := FeatureItemFrom(item)
		items[key] = fi
	}

	return items, nil
}

func FeatureItemFrom(item *cms.Item) FeatureItem {
	res := FeatureItem{}

	type internalGroup struct {
		Data []any `cms:"data"`
	}

	type internalItem struct {
		ID      string          `cms:"id"`
		CityGML any             `cms:"citygml"`
		Data    []any           `cms:"data"`
		Items   []internalGroup `cms:"items,group"`
	}

	fi := internalItem{}
	item.Unmarshal(&fi)

	res.ID = fi.ID

	if fi.CityGML != nil {
		asset, ok := fi.CityGML.(map[string]any)
		if ok {
			url, _ := asset["url"].(string)
			res.CityGML = url
		}
	}

	if fi.Data != nil {
		for _, d := range fi.Data {
			asset, ok := d.(map[string]any)
			if ok {
				url, _ := asset["url"].(string)
				res.Data = append(res.Data, url)
			}
		}
	}

	if fi.Items != nil {
		for _, item := range fi.Items {
			for _, d := range item.Data {
				asset, ok := d.(map[string]any)
				if ok {
					url, _ := asset["url"].(string)
					res.Data = append(res.Data, url)
				}
			}
		}
	}

	return res
}
