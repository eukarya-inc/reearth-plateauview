package preparegspatialjp

import (
	"context"
	"fmt"

	cms "github.com/reearth/reearth-cms-api/go"
)

type FeatureItem struct {
	ID      string `json:"id,omitempty" cms:"id"`
	CityGML string `json:"citygml,omitempty" cms:"-"`
	Data    string `json:"data,omitempty" cms:"-"`
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
	fi := FeatureItem{}
	item.Unmarshal(&fi)

	field := item.FieldByKey("citygml")
	if field != nil {
		asset, ok := field.Value.(map[string]any)
		if ok {
			url, _ := asset["url"].(string)
			fi.CityGML = url
		}
	}

	field = item.FieldByKey("data")
	if field != nil {
		asset, ok := field.Value.(map[string]any)
		if ok {
			url, _ := asset["url"].(string)
			fi.Data = url
		}
	}

	return fi
}
