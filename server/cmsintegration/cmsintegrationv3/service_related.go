package cmsintegrationv3

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/dataconv"
	geojson "github.com/paulmach/go.geojson"
	cms "github.com/reearth/reearth-cms-api/go"
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
	log.Debugfc(ctx, "cmsintegrationv3: related dataset: %#v", item)

	if err := convertRelatedDataset(ctx, s, w, item); err != nil {
		return err
	}

	if err := packRelatedDataset(ctx, s, w, item); err != nil {
		return err
	}

	return nil
}

func convertRelatedDataset(ctx context.Context, s *Services, w *cmswebhook.Payload, item *RelatedItem) (err error) {
	project := w.ProjectID()
	convTargets := make([]string, 0, len(relatedDataConvertionTargets))
	newStatus := map[string]*cms.Tag{}
	newItems := map[string]RelatedItemDatum{}

	for _, target := range relatedDataConvertionTargets {
		if tagIsNot(item.ConvertStatus[target], ConvertionStatusNotStarted) {
			log.Debugfc(ctx, "cmsintegrationv3: already converted")
			continue
		}

		if len(item.Items[target].Asset) == 0 {
			continue
		}

		convTargets = append(convTargets, target)
		newStatus[target] = tagFrom(ConvertionStatusRunning)
	}

	if len(convTargets) == 0 {
		log.Debugfc(ctx, "cmsintegrationv3: no conv targets")
		return nil
	}

	log.Infofc(ctx, "cmsintegrationv3: convertRelatedDataset: %v", convTargets)

	// update status
	if _, err := s.CMS.UpdateItem(ctx, item.ID, nil, (&RelatedItem{
		ConvertStatus: newStatus,
	}).CMSItem().MetadataFields); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	defer func() {
		for k, v := range newStatus {
			if tagIs(v, ConvertionStatusRunning) {
				newStatus[k] = tagFrom(ConvertionStatusNotStarted)
			}
		}

		// update item
		ritem := (&RelatedItem{
			Items:         newItems,
			ConvertStatus: newStatus,
		}).CMSItem()
		if _, err2 := s.CMS.UpdateItem(ctx, item.ID, ritem.Fields, ritem.MetadataFields); err2 != nil {
			err = fmt.Errorf("failed to update item: %w", err2)
		}

		// comment to the item
		succeeded := lo.EveryBy(convTargets, func(t string) bool {
			return tagIs(newStatus[t], ConvertionStatusSuccess)
		})

		var comment string
		if succeeded {
			comment = "変換が完了しました。"
		} else {
			comment = fmt.Sprintf("変換が完了しましたが、一部のデータの変換に失敗しました。\n%s", err)
		}

		if err2 := s.CMS.CommentToItem(ctx, item.ID, comment); err2 != nil {
			err = fmt.Errorf("failed to add comment: %w", err2)
		}
	}()

	// comment to the item
	if err := s.CMS.CommentToItem(ctx, item.ID, "変換を開始しました。"); err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	for _, target := range convTargets {
		log.Debugf("cmsintegrationv3: convert %s: %#v", target)

		d := item.Items[target]
		var converror bool
		var convassets []string
		for _, a := range d.Asset {
			newAsset, err2 := convRelatedType(ctx, project, target, a, s)
			if err2 != nil {
				err = errors.Join(err, fmt.Errorf("%s: %s", target, err2))
				converror = true
			} else {
				convassets = append(convassets, newAsset)
			}
		}

		if converror {
			newStatus[target] = tagFrom(ConvertionStatusError)
		} else {
			newStatus[target] = tagFrom(ConvertionStatusSuccess)
			newItems[target] = RelatedItemDatum{
				ID:        d.ID,
				Converted: convassets,
			}
		}
	}

	log.Infofc(ctx, "cmsintegrationv3: convertRelatedDataset: converted: %v", convTargets)
	return err
}

func convRelatedType(ctx context.Context, project, target, assetID string, s *Services) (_ string, err error) {
	asset, err := s.CMS.Asset(ctx, assetID)
	if err != nil {
		return "", fmt.Errorf("failed to get asset (%s): %w", target, err)
	}

	id := strings.TrimSuffix(path.Base(asset.URL), path.Ext(asset.URL))
	log.Debugf("cmsintegrationv3: convert %s (%s)", target, id)

	// download asset
	data, err := s.GETAsBytes(ctx, asset.URL)
	if err != nil {
		return "", fmt.Errorf("failed to download asset: %w", err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal asset: %w", err)
	}

	// conv
	var res any
	if target == "border" {
		res, err = dataconv.ConvertBorder(fc, id)
	} else if target == "landmark" || target == "station" {
		res, err = dataconv.ConvertLandmark(fc, id)
	}

	if err != nil || res == nil {
		return "", fmt.Errorf("failed to convert: %w", err)
	}

	uploadBody, err := json.Marshal(res)
	if err != nil {
		return "", fmt.Errorf("failed to marshal: %w", err)
	}

	// upload
	newAssetID, err := s.CMS.UploadAssetDirectly(ctx, project, id+".czml", bytes.NewReader(uploadBody))
	if err != nil {
		return "", fmt.Errorf("failed to upload asset: %w", err)
	}

	log.Debugf("cmsintegrationv3: converted %s (%s) to %s", target, id, newAssetID)
	return newAssetID, nil
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
		return len(item.Items[t].Asset) == 0
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

	var mergedAssetID string

	defer func() {
		var status ConvertionStatus
		if err == nil {
			status = ConvertionStatusSuccess
		} else {
			status = ConvertionStatusError
		}

		newItem := (&RelatedItem{
			Merged:      mergedAssetID,
			MergeStatus: tagFrom(status),
		}).CMSItem()
		if _, err := s.CMS.UpdateItem(ctx, item.ID, newItem.Fields, newItem.MetadataFields); err != nil {
			log.Errorfc(ctx, "cmsintegrationv3: failed to update item: %w", err)
		}

		// comment to the item
		var comment string
		if err == nil {
			comment = "G空間情報センター公開用zipファイルの作成が完了しました。"
		} else {
			comment = fmt.Sprintf("G空間情報センター公開用zipファイルの作成に失敗しました。\n%s", err)
		}

		if err := s.CMS.CommentToItem(ctx, item.ID, comment); err != nil {
			log.Errorfc(ctx, "cmsintegrationv3: failed to add comment: %w", err)
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

	zipbuf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(zipbuf)

	assetMerged := false
	var year int
	for _, target := range relatedDataTypes {
		d := item.Items[target]

		if len(d.Asset) == 0 {
			continue
		}

		y, err := packRelatedDatasetTarget(ctx, target, d.Asset, zw, cityItem.CityCode, cityItem.CityNameEn, s)
		if err != nil {
			return err
		}

		assetMerged = true
		year = y
	}

	if !assetMerged {
		log.Debugfc(ctx, "cmsintegrationv3: no assets")
		return nil
	}

	if err := zw.Close(); err != nil {
		return fmt.Errorf("failed to close zip: %w", err)
	}

	// upload zip
	zipName := fmt.Sprintf("%s_%s_%d_related.zip", cityItem.CityCode, cityItem.CityNameEn, year)
	mergedAssetID, err = s.CMS.UploadAssetDirectly(ctx, w.ProjectID(), zipName, zipbuf)
	if err != nil {
		return fmt.Errorf("failed to upload zip: %w", err)
	}

	log.Infofc(ctx, "cmsintegrationv3: packageRelatedDatasetForGeospatialjp: done")
	return nil
}

func packRelatedDatasetTarget(ctx context.Context, target string, assets []string, zw *zip.Writer, cityCode, cityNameEn string, s *Services) (int, error) {
	var mergedFeatures []*geojson.Feature
	noNeedToWriteAssets := len(assets) == 1

	var assetName *relatedAssetName

	for i, asset := range assets {
		asset, err := s.CMS.Asset(ctx, asset)
		if err != nil {
			return 0, fmt.Errorf("(%s/%d): アセットが見つかりません: %v", target, i+1, err)
		}

		an := parseRelatedAssetName(asset.URL)
		if an == nil {
			return 0, fmt.Errorf("(%s/%d/%s): ファイル名が命名規則に沿っていません。 \"[市区町村コード5桁]_[市区町村名英名]_[提供事業者名]_[整備年度4桁]_[landmark,shelterなど].geojson\" としてください。: %v", target, i+1, path.Base(asset.URL), err)
		}

		if !noNeedToWriteAssets && an.CityCode == cityCode {
			return 0, fmt.Errorf("(%s/%d/%s): アセット名の市区町村コードが全体の市区町村コードと同じです。区ごとに登録する場合はファイル名中のコードを各区のコードにしてください。", target, i+1, path.Base(asset.URL))
		}

		if assetName == nil {
			assetName = an
		}

		// download asset
		data, err := s.GETAsBytes(ctx, asset.URL)
		if err != nil {
			return 0, fmt.Errorf("failed to download asset (%s): %w", target, err)
		}

		// parse GeoJSON
		fc := geojson.NewFeatureCollection()
		if err := json.Unmarshal(data, fc); err != nil {
			return 0, fmt.Errorf("(%s/%d/%s): GeoJSONとして読み込むことができませんでした。正しいGeoJSONかどうかファイルの内容を確認してください。: %v", target, i+1, path.Base(asset.URL), err)
		}

		mergedFeatures = append(mergedFeatures, fc.Features...)

		// add to zip
		if !noNeedToWriteAssets {
			f, err := zw.Create(path.Base(asset.URL))
			if err != nil {
				return 0, fmt.Errorf("failed to create zip file (%s): %w", target, err)
			}

			if _, err := f.Write(data); err != nil {
				return 0, fmt.Errorf("failed to write zip file (%s): %w", target, err)
			}
		}
	}

	// merge multiple assets
	if len(mergedFeatures) > 0 {
		name := fmt.Sprintf("%s_%s_%s_%d_%s", cityCode, cityNameEn, assetName.Provider, assetName.Year, assetName.Type)
		f, err := zw.Create(fmt.Sprintf("%s.geojson", name))
		if err != nil {
			return 0, fmt.Errorf("failed to create zip file (%s): %w", target, err)
		}

		fc := map[string]any{
			"type":     "FeatureCollection",
			"name":     name,
			"features": mergedFeatures,
		}
		data, err := json.Marshal(fc)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal (%s): %w", target, err)
		}

		if _, err := f.Write(data); err != nil {
			return 0, fmt.Errorf("failed to write zip file (%s): %w", target, err)
		}
	}

	return assetName.Year, nil
}

type relatedAssetName struct {
	CityCode string
	CityName string
	Provider string
	Year     int
	Type     string
	Ext      string
}

var reRelatedAssetName = regexp.MustCompile(`^([0-9]{5})_([^_]+)_([^_]+)_([0-9]+)_(.+)\.(.+)$`)

func parseRelatedAssetName(name string) *relatedAssetName {
	name = path.Base(name)
	m := reRelatedAssetName.FindStringSubmatch(name)
	if m == nil {
		return nil
	}

	y, _ := strconv.Atoi(m[4])

	return &relatedAssetName{
		CityCode: m[1],
		CityName: m[2],
		Provider: m[3],
		Year:     y,
		Type:     m[5],
		Ext:      m[6],
	}
}
