package preparegspatialjp

import (
	"context"

	cms "github.com/reearth/reearth-cms-api/go"
)

func GetRelatedZipAssetIDAndURL(ctx context.Context, cms *cms.CMS, cityItem *CityItem) (string, string, error) {
	if cityItem.RelatedDataset == "" {
		return "", "", nil
	}

	item, err := cms.GetItem(ctx, cityItem.RelatedDataset, true)
	if err != nil {
		return "", "", err
	}

	v := item.FieldByKey("merged").Value
	v2, ok := v.(map[string]any)
	if !ok {
		return "", "", nil
	}

	id, ok := v2["id"].(string)
	if !ok {
		return "", "", nil
	}

	url, ok := v2["url"].(string)
	if !ok {
		return "", "", nil
	}

	return id, url, nil
}
