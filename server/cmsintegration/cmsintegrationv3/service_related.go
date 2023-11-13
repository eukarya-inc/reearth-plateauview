package cmsintegrationv3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/dataconv"
	geojson "github.com/paulmach/go.geojson"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
)

var relatedDataConvertionTargets = []string{
	"border",
	"landmark",
	"station",
}

func convertRelatedDataset(ctx context.Context, s *Services, w *cmswebhook.Payload) error {
	if w.ItemData.Model.Key != modelPrefix+relatedModel {
		log.Debugfc(ctx, "cmsintegrationv3: not related dataset")
		return nil
	}

	item := RelatedItemFrom(w.ItemData.Item)

	if item.City == "" {
		log.Debugfc(ctx, "cmsintegrationv3: no city")
		return nil
	}

	targets := lo.Filter(relatedDataConvertionTargets, func(t string, _ int) bool {
		return updatedField(w, t) != nil
	})
	if len(targets) == 0 {
		log.Debugfc(ctx, "cmsintegrationv3: no changes")
		return nil
	}

	log.Infofc(ctx, "cmsintegrationv3: convertRelatedDataset")

	// get city item
	cityItemRaw, err := s.CMS.GetItem(ctx, item.City, false)
	if err != nil {
		return fmt.Errorf("failed to get city item: %w", err)
	}

	cityItem := CityItemFrom(cityItemRaw)
	if cityItem.CityName == "" || cityItem.CityCode == "" {
		return fmt.Errorf("city item is not valid: %v", cityItem)
	}

	for _, target := range targets {
		if item.Assets[target] == "" || item.ConvertedAssets[target] != "" ||
			item.ConvertStatus[target] != "" && item.ConvertStatus[target] != ConvertionStatusNotStarted {
			continue
		}

		id := fmt.Sprintf("%s_%s_%s", cityItem.CityName, cityItem.CityCode, target)
		log.Debugf("cmsintegrationv3: convert %s (%s)", target, id)

		// download asset
		asset, err := s.DownloadAssetAsBytes(ctx, item.Assets[target])
		if err != nil {
			return fmt.Errorf("failed to download asset: %w", err)
		}

		fc, err := geojson.UnmarshalFeatureCollection(asset)
		if err != nil {
			return fmt.Errorf("failed to unmarshal asset: %w", err)
		}

		// conv
		var res any
		if target == "border" {
			res, err = dataconv.ConvertBorder(fc, id)
		} else if target == "landmark" || target == "station" {
			res, err = dataconv.ConvertLandmark(fc, id)
		}

		if err != nil || res == nil {
			return fmt.Errorf("failed to convert: %w", err)
		}

		uploadBody, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("failed to marshal: %w", err)
		}

		// upload
		assetID, err := s.CMS.UploadAssetDirectly(ctx, w.ProjectID(), id+".czml", bytes.NewReader(uploadBody))
		if err != nil {
			return fmt.Errorf("failed to upload asset: %w", err)
		}

		// update item
		ritem := (&RelatedItem{
			ConvertedAssets: map[string]string{
				target: assetID,
			},
			ConvertStatus: map[string]ConvertionStatus{
				target: ConvertionStatusSuccess,
			},
		}).CMSItem()
		_, err = s.CMS.UpdateItem(ctx, item.ID, ritem.Fields, ritem.MetadataFields)
		if err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}
	}

	// comment to the item
	if err := s.CMS.CommentToItem(ctx, item.ID, "変換に成功しました。"); err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}

func packageRelatedDatasetForGeospatialjp(ctx context.Context, s *Services, w *cmswebhook.Payload) error {
	log.Infofc(ctx, "cmsintegrationv3: packageRelatedDatasetForGeospatialjp")
	// TODO

	return nil
}
