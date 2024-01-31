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

	log.Debugfc(ctx, "geospatialjpv3: seed: %s", ppp.Sprint(dataItem))
	return seed, nil
}

func valueToAsset(v *cms.Value) string {
	vv := map[string]any{}
	if err := v.JSON(&vv); err != nil {
		return ""
	}
	if url, ok := vv["url"].(string); ok {
		return url
	}
	return ""
}
