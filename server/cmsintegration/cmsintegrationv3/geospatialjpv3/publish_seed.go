package geospatialjpv3

import (
	"context"
	"fmt"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

type Seed struct {
	CityGML            string
	Plateau            string
	Related            string
	Desc               string
	CityGMLDescription string
	PlateauDescription string
	RelatedDescription string
	Version            int
}

const defaultVersion = 3

func getSeed(ctx context.Context, c cms.Interface, cityItem *CityItem) (Seed, error) {
	seed := Seed{
		Version: defaultVersion,
	}

	rawDataItem, err := c.GetItem(ctx, cityItem.GeospatialjpData, true)
	if err != nil {
		return seed, fmt.Errorf("failed to get data item: %w", err)
	}

	rawIndexItem, err := c.GetItem(ctx, cityItem.GeospatialjpIndex, true)
	if err != nil {
		return seed, fmt.Errorf("failed to get index item: %w", err)
	}

	var dataItem CMSDataItem
	rawDataItem.Unmarshal(&dataItem)

	var indexItem CMSIndexItem
	rawIndexItem.Unmarshal(&indexItem)

	log.Debugfc(ctx, "geospatialjpv3: rawDataItem: %s", ppp.Sprint(rawDataItem))
	log.Debugfc(ctx, "geospatialjpv3: rawIndexItem: %s", ppp.Sprint(rawIndexItem))
	log.Debugfc(ctx, "geospatialjpv3: dataItem: %s", ppp.Sprint(dataItem))
	log.Debugfc(ctx, "geospatialjpv3: indexItem: %s", ppp.Sprint(indexItem))

	if dataItem.CityGML != nil {
		seed.CityGML = valueToAsset(dataItem.CityGML)
	}
	if dataItem.Plateau != nil {
		seed.Plateau = valueToAsset(dataItem.Plateau)
	}
	if dataItem.Related != nil {
		seed.Related = valueToAsset(dataItem.Related)
	}
	if indexItem.Desc != "" {
		seed.Desc = indexItem.Desc
	}

	seed.CityGMLDescription = indexItem.DescCityGML
	seed.PlateauDescription = indexItem.DescPlateau
	seed.RelatedDescription = indexItem.DescRelated

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
