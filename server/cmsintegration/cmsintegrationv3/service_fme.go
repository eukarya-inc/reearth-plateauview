package cmsintegrationv3

import (
	"context"
	"fmt"
	"strings"

	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
)

func sendRequestToFME(ctx context.Context, s *Services, conf *Config, w *cmswebhook.Payload) error {
	featureType := strings.TrimPrefix(w.ItemData.Model.Key, modelPrefix)
	itemID := w.ItemData.Item.ID
	if w.ItemData.Item.MetadataItemID == nil {
		log.Debugfc(ctx, "cmsintegrationv2: no metadata item id")
		return nil
	}

	item := FeatureItemFrom(w.ItemData.Item)
	if item.ConvertStatus != "" && item.ConvertStatus != ConvertionStatusNotStarted {
		log.Debugfc(ctx, "cmsintegrationv2: already converted")
		return nil
	}

	log.Infofc(ctx, "cmsintegrationv2: sendRequestToFME")

	// get city item
	cityItemRaw, err := s.CMS.GetItem(ctx, item.City, false)
	if err != nil {
		return fmt.Errorf("failed to get city item: %w", err)
	}

	cityItem := CityItemFrom(cityItemRaw)

	// get CityGML asset
	cityGMLAsset, err := s.CMS.Asset(ctx, item.CityGML)
	if err != nil {
		return fmt.Errorf("failed to get citygml asset: %w", err)
	}

	// get codelist asset
	codelistAsset, err := s.CMS.Asset(ctx, cityItem.CodeList)
	if err != nil {
		return fmt.Errorf("failed to get codelist asset: %w", err)
	}

	// request to fme
	err = s.FME.Request(ctx, fmeRequest{
		ID: fmeID{
			ItemID:      itemID,
			ProjectID:   w.ProjectID(),
			FeatureType: featureType,
		}.String(conf.Secret),
		Target:    cityGMLAsset.URL,
		PRCS:      cityItem.PRCS,
		Codelists: codelistAsset.URL,
		ResultURL: resultURL(conf),
		// Config: ,
	})
	if err != nil {
		return fmt.Errorf("failed to request to fme: %w", err)
	}

	// update convertion status
	err = s.UpdateFeatureItemStatus(ctx, itemID, ConvertionStatusRunning)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	// post a comment to the item
	err = s.CMS.CommentToItem(ctx, itemID, "品質検査・変換を開始しました。")
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}

func receiveResultFromFME(ctx context.Context, s *Services, conf *Config, f fmeResult) error {
	id := f.ParseID(conf.Secret)
	if id.ItemID == "" {
		return fmt.Errorf("invalid id: %s", f.ID)
	}

	log.Infofc(ctx, "cmsintegrationv2: receiveResultFromFME")

	logmsg := ""
	if f.LogURL != "" {
		logmsg = "ログ： " + f.LogURL
	}

	// handle error
	if f.Status == "error" {
		log.Warnfc(ctx, "cmsintegrationv2: failed to convert: %v", f.LogURL)

		// update item status
		err := s.UpdateFeatureItemStatus(ctx, id.ItemID, ConvertionStatusError)
		if err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}

		// comment to the item
		err = s.CMS.CommentToItem(ctx, id.ItemID, "品質検査・変換に失敗しました。"+logmsg)
		if err != nil {
			return fmt.Errorf("failed to add comment: %w", err)
		}

		return nil
	}

	// get url from the result
	assets := f.GetResultURLs(id.FeatureType)

	// upload assets
	dataAssets := make([]string, 0, len(assets.Data))
	for _, url := range assets.Data {
		aid, err := s.CMS.UploadAsset(ctx, id.ProjectID, url)
		if err != nil {
			log.Errorfc(ctx, "cmsintegrationv2: failed to upload asset (%s): %v", url, err)
			return nil
		}
		dataAssets = append(dataAssets, aid)
	}

	// upload dic
	var dicAssetID string
	if assets.Dic != "" {
		log.Debugfc(ctx, "cmsintegrationv2: upload dic: %s", assets.Dic)
		var err error
		dicAssetID, err = s.CMS.UploadAsset(ctx, id.ProjectID, assets.Dic)
		if err != nil {
			return fmt.Errorf("failed to upload dic: %w", err)
		}
	}

	// upload maxlod
	var maxlodAssetID string
	if assets.MaxLOD != "" {
		log.Debugfc(ctx, "cmsintegrationv2: upload maxlod: %s", assets.MaxLOD)
		var err error
		maxlodAssetID, err = s.CMS.UploadAsset(ctx, id.ProjectID, assets.MaxLOD)
		if err != nil {
			return fmt.Errorf("failed to upload maxlod: %w", err)
		}
	}

	// update item
	err := s.CompleteFeatureItemConvertion(ctx, id.ItemID, dataAssets, dicAssetID, maxlodAssetID)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	// comment to the item
	err = s.CMS.CommentToItem(ctx, id.ItemID, "品質検査・変換が完了しました。"+logmsg)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}
