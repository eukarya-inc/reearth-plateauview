package geospatialjpv3

import (
	"context"
	"fmt"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

type Seed struct {
	CityGML string
	Plateau string
	Related string
	Version int
}

const defaultVersion = 3

func getSeed(ctx context.Context, c cms.Interface, cityItem *CityItem) (Seed, error) {
	seed := Seed{
		Version: defaultVersion,
	}

	rawDataItem, err := c.GetItem(ctx, cityItem.GeospatialjpData, true)
	if err != nil {
		return seed, fmt.Errorf("failed to get item: %w", err)
	}

	log.Debugfc(ctx, "geospatialjpv3: rawDataItem: %s", ppp.Sprint(rawDataItem))

	var dataItem GspatialjpItem
	rawDataItem.Unmarshal(&dataItem)

	if dataItem.CityGML != nil {
		seed.CityGML = valueToAsset(dataItem.CityGML)
	}
	if dataItem.Plateau != nil {
		seed.Plateau = valueToAsset(dataItem.Plateau)
	}
	if dataItem.Related != nil {
		seed.Related = valueToAsset(dataItem.Related)
	}

	return seed, nil
}

func valueToAsset(v map[string]any) string {
	if v == nil {
		return ""
	}
	if url, ok := v["url"].(string); ok {
		return url
	}
	return ""
}
