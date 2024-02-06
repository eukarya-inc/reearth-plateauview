package preparegspatialjp

import (
	"context"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/samber/lo"
)

func GetRelatedZipAssetID(ctx context.Context, cms *cms.CMS, cityItem *CityItem) (string, error) {
	if cityItem.RelatedDataset == "" {
		return "", nil
	}

	item, err := cms.GetItem(ctx, cityItem.RelatedDataset, false)
	if err != nil {
		return "", err
	}

	asset := item.FieldByKey("merged").GetValue().String()

	return lo.FromPtr(asset), nil
}
