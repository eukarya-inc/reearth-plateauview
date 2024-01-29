package geospatialjpv3

import (
	"context"
	"fmt"

	"github.com/k0kubun/pp/v3"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

type GeospatialItem struct {
	CityGML string `json:"citygml,omitempty"`
	Plateau string `json:"plateau,omitempty"`
	Related string `json:"related,omitempty"`
}

func getGeospatialItems(ctx context.Context, c cms.Interface, cityItem *CityItem) (GeospatialItem, error) {
	res := GeospatialItem{}

	item, err := c.GetItem(ctx, cityItem.GeospatialjpData, true)
	if err != nil {
		return res, fmt.Errorf("failed to get item: %w", err)
	}

	{
		pp := pp.New()
		pp.SetColoringEnabled(false)
		s := pp.Sprint(item)
		log.Debugfc(ctx, "geospatialjpv3: geoItem: %s", s)
	}

	for _, field := range item.Fields {
		if field.Key == "citygml" {
			if field.Value != nil {
				res.CityGML = field.Value.(string)
			} else {
				res.CityGML = ""
			}
		}

		if field.Key == "plateau" {
			if field.Value != nil {
				res.Plateau = field.Value.(string)
			} else {
				res.Plateau = ""
			}
		}

		if field.Key == "related" {
			if field.Value != nil {
				res.Related = field.Value.(string)
			} else {
				res.Related = ""
			}
		}
	}

	return res, nil
}
