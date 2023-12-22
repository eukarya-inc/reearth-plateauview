package cmsintegrationv3

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"strings"

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

const relatedConvStatus = "conv_status"

func handleRelatedDataset(ctx context.Context, s *Services, w *cmswebhook.Payload) error {
	// if event type is "item.create" and payload is metadata, skip it
	if w.Type == cmswebhook.EventItemCreate && w.ItemData.Item.OriginalItemID != nil ||
		w.ItemData == nil || w.ItemData.Item == nil || w.ItemData.Model == nil ||
		w.ItemData.Item.FieldByKey(relatedConvStatus) == nil {
		return nil
	}

	if w.ItemData.Model.Key != modelPrefix+relatedModel {
		log.Debugfc(ctx, "cmsintegrationv3: not related dataset")
		return nil
	}

	mainItem, err := s.GetMainItemWithMetadata(ctx, w.ItemData.Item)
	if err != nil {
		return err
	}

	item := RelatedItemFrom(mainItem)

	if err := convertRelatedDataset(ctx, s, w, item); err != nil {
		return err
	}

	if err := packRelatedDataset(ctx, s, w, item); err != nil {
		return err
	}

	return nil
}

func convertRelatedDataset(ctx context.Context, s *Services, w *cmswebhook.Payload, item *RelatedItem) (err error) {
	if tagIsNot(item.ConvertStatus, ConvertionStatusNotStarted) {
		log.Debugfc(ctx, "cmsintegrationv3: already converted")
		return nil
	}

	if !lo.SomeBy(relatedDataConvertionTargets, func(t string) bool {
		return len(item.Assets[t]) > 0
	}) {
		log.Debugfc(ctx, "cmsintegrationv3: no assets")
	}

	log.Infofc(ctx, "cmsintegrationv3: convertRelatedDataset")

	// update status
	if _, err := s.CMS.UpdateItem(ctx, item.ID, nil, (&RelatedItem{
		ConvertStatus: tagFrom(ConvertionStatusRunning),
	}).CMSItem().MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	defer func() {
		if err == nil {
			return
		}

		if _, err := s.CMS.UpdateItem(ctx, item.ID, nil, (&RelatedItem{
			MergeStatus: tagFrom(ConvertionStatusError),
		}).CMSItem().MetadataFields); err != nil {
			log.Warnf("cmsintegrationv3: failed to update item: %w", err)
		}

		// comment to the item
		if err := s.CMS.CommentToItem(ctx, item.ID, "G空間情報センター公開用zipファイルの作成に失敗しました。"); err != nil {
			log.Warnf("cmsintegrationv3: failed to add comment: %w", err)
		}
	}()

	// comment to the item
	if err := s.CMS.CommentToItem(ctx, item.ID, "変換を開始しました。"); err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	conv := map[string][]string{}
	for _, target := range relatedDataConvertionTargets {
		if len(item.Assets[target]) == 0 {
			continue
		}

		for _, t := range item.Assets[target] {
			asset, err := s.CMS.Asset(ctx, t)
			if err != nil {
				return fmt.Errorf("failed to get asset (%s): %w", target, err)
			}

			id := strings.TrimSuffix(path.Base(asset.URL), path.Ext(asset.URL))
			log.Debugf("cmsintegrationv3: convert %s (%s)", target, id)

			// download asset
			data, err := s.GETAsBytes(ctx, asset.URL)
			if err != nil {
				return fmt.Errorf("failed to download asset: %w", err)
			}

			fc, err := geojson.UnmarshalFeatureCollection(data)
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

			conv[target] = append(conv[target], assetID)
		}
	}

	// update item
	ritem := (&RelatedItem{
		ConvertedAssets: conv,
		ConvertStatus:   tagFrom(ConvertionStatusSuccess),
	}).CMSItem()
	if _, err := s.CMS.UpdateItem(ctx, item.ID, ritem.Fields, ritem.MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	// comment to the item
	if err := s.CMS.CommentToItem(ctx, item.ID, "変換に成功しました。"); err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}

func packRelatedDataset(ctx context.Context, s *Services, w *cmswebhook.Payload, item *RelatedItem) (err error) {
	if item.City == "" {
		log.Debugfc(ctx, "cmsintegrationv3: no city")
		return nil
	}

	if tagIsNot(item.MergeStatus, ConvertionStatusNotStarted) {
		log.Debugfc(ctx, "cmsintegrationv3: already merged")
		return nil
	}

	if missingTypes := lo.Filter(relatedDataTypes, func(t string, _ int) bool {
		return len(item.Assets[t]) == 0
	}); len(missingTypes) > 0 {
		log.Debugfc(ctx, "cmsintegrationv3: there are some missing assets: %v", missingTypes)
		return nil
	}

	log.Infofc(ctx, "cmsintegrationv3: packageRelatedDatasetForGeospatialjp")

	// update status
	if _, err := s.CMS.UpdateItem(ctx, item.ID, nil, (&RelatedItem{
		MergeStatus: tagFrom(ConvertionStatusRunning),
	}).CMSItem().MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	defer func() {
		if err == nil {
			return
		}

		if _, err := s.CMS.UpdateItem(ctx, item.ID, nil, (&RelatedItem{
			MergeStatus: tagFrom(ConvertionStatusError),
		}).CMSItem().MetadataFields); err != nil {
			log.Warnf("cmsintegrationv3: failed to update item: %w", err)
		}

		// comment to the item
		if err := s.CMS.CommentToItem(ctx, item.ID, fmt.Sprintf("G空間情報センター公開用zipファイルの作成に失敗しました。%s", err)); err != nil {
			log.Warnf("cmsintegrationv3: failed to add comment: %w", err)
		}
	}()

	// comment to the item
	if err := s.CMS.CommentToItem(ctx, item.ID, "G空間情報センター公開用zipファイルの作成を開始しました。"); err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	// get city
	cityItemRaw, err := s.CMS.GetItem(ctx, item.City, false)
	if err != nil {
		return fmt.Errorf("failed to get city: %w", err)
	}

	cityItem := CityItemFrom(cityItemRaw)
	zipName := fmt.Sprintf("%s_%s_related.zip", cityItem.CityCode, cityItem.CityNameEn)

	zipbuf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(zipbuf)

	assetPreset := false
	for _, target := range relatedDataTypes {
		name := fmt.Sprintf("%s_%s_%s", cityItem.CityCode, cityItem.CityNameEn, target)
		var features []*geojson.Feature
		noNeedToWriteAssets := len(item.Assets[target]) == 1

		for _, t := range item.Assets[target] {
			// get asset
			asset, err := s.CMS.Asset(ctx, t)
			if err != nil {

				return fmt.Errorf("failed to get asset (%s): %w", target, err)
			}

			// download asset
			data, err := s.GETAsBytes(ctx, asset.URL)
			if err != nil {
				return fmt.Errorf("failed to download asset (%s): %w", target, err)
			}

			// add to zip
			if !noNeedToWriteAssets {
				f, err := zw.Create(path.Base(asset.URL))
				if err != nil {
					return fmt.Errorf("failed to create zip file (%s): %w", target, err)
				}

				if _, err := f.Write(data); err != nil {
					return fmt.Errorf("failed to write zip file (%s): %w", target, err)
				}
			}

			fc := geojson.NewFeatureCollection()
			if err := json.Unmarshal(data, fc); err != nil {
				return fmt.Errorf("failed to decode asset (%s): %w", target, err)
			}

			features = append(features, fc.Features...)
			assetPreset = true
		}

		// merge multiple assets
		if len(features) > 0 {
			f, err := zw.Create(fmt.Sprintf("%s.geojson", target))
			if err != nil {
				return fmt.Errorf("failed to create zip file (%s): %w", target, err)
			}

			fc := map[string]any{
				"type":     "FeatureCollection",
				"name":     name,
				"features": features,
			}
			data, err := json.Marshal(fc)
			if err != nil {
				return fmt.Errorf("failed to marshal (%s): %w", target, err)
			}

			if _, err := f.Write(data); err != nil {
				return fmt.Errorf("failed to write zip file (%s): %w", target, err)
			}
		}
	}

	if !assetPreset {
		log.Debugfc(ctx, "cmsintegrationv3: no assets")
		return nil
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("failed to close zip: %w", err)
	}

	// upload zip
	assetID, err := s.CMS.UploadAssetDirectly(ctx, w.ProjectID(), zipName, zipbuf)
	if err != nil {
		return fmt.Errorf("failed to upload zip: %w", err)
	}

	// update item
	ritem := (&RelatedItem{
		Merged:      assetID,
		MergeStatus: tagFrom(ConvertionStatusSuccess),
	}).CMSItem()

	if _, err := s.CMS.UpdateItem(ctx, item.ID, ritem.Fields, ritem.MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	// comment to the item
	if err := s.CMS.CommentToItem(ctx, item.ID, "G空間情報センター公開用zipファイルの作成が完了しました。"); err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}
