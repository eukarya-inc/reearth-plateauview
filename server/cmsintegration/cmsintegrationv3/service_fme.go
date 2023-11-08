package cmsintegrationv3

import (
	"context"
	"encoding/json"
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
	if item.ConvertionStatus != "" && item.ConvertionStatus != ConvertionStatusNotStarted {
		log.Debugfc(ctx, "cmsintegrationv2: already converted")
		return nil
	}

	if item.SkipQC && item.SkipConvert || item.CityGML == "" || item.City == "" {
		log.Debugfc(ctx, "cmsintegrationv2: skip convert")
		return nil
	}

	log.Infofc(ctx, "cmsintegrationv2: sendRequestToFME")

	// get CityGML asset
	cityGMLAsset, err := s.CMS.Asset(ctx, item.CityGML)
	if err != nil {
		_ = failToConvert(ctx, s, itemID, "CityGMLが見つかりません。")
		return fmt.Errorf("failed to get citygml asset: %w", err)
	}

	// get city item
	cityItemRaw, err := s.CMS.GetItem(ctx, item.City, false)
	if err != nil {
		_ = failToConvert(ctx, s, itemID, "都市アイテムが見つかりません。")
		return fmt.Errorf("failed to get city item: %w", err)
	}

	cityItem := CityItemFrom(cityItemRaw)
	if cityItem.CodeLists == "" {
		_ = failToConvert(ctx, s, itemID, "コードリストが都市アイテムに登録されていないため品質検査・変換を開始できません。")
		return fmt.Errorf("city item has no codelist")
	}

	// get codelist asset
	codelistAsset, err := s.CMS.Asset(ctx, cityItem.CodeLists)
	if err != nil {
		_ = failToConvert(ctx, s, itemID, "コードリストが見つかりません。")
		return fmt.Errorf("failed to get codelist asset: %w", err)
	}

	// update convertion status
	err = s.UpdateFeatureItemStatus(ctx, itemID, ConvertionStatusRunning)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	ty := fmeTypeQcConv
	if item.SkipQC {
		ty = fmeTypeConv
	} else if item.SkipConvert {
		ty = fmeTypeQC
	}

	// request to fme
	err = s.FME.Request(ctx, fmeRequest{
		ID: fmeID{
			ItemID:      itemID,
			ProjectID:   w.ProjectID(),
			FeatureType: featureType,
		}.String(conf.Secret),
		Target:    cityGMLAsset.URL,
		PRCS:      cityItem.PRCS.ESPGCode(),
		Codelists: codelistAsset.URL,
		ResultURL: resultURL(conf),
		Type:      ty,
		// Config: ,
	})
	if err != nil {
		_ = failToConvert(ctx, s, itemID, "FMEへのリクエストに失敗しました。")
		return fmt.Errorf("failed to request to fme: %w", err)
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

	if f.Status == "" {
		return fmt.Errorf("invalid status")
	}

	log.Infofc(ctx, "cmsintegrationv2: receiveResultFromFME")

	logmsg := ""
	if f.LogURL != "" {
		logmsg = "ログ： " + f.LogURL
	}

	// handle error
	if f.Status == "error" {
		log.Warnfc(ctx, "cmsintegrationv2: failed to convert: %v", f.LogURL)
		_ = failToConvert(ctx, s, id.ItemID, "品質検査・変換に失敗しました。"+logmsg)
		return nil
	}

	// get url from the result
	assets := f.GetResultURLs(id.FeatureType)

	// upload assets
	var dataAssets []string
	if len(assets.Data) > 0 {
		dataAssets = make([]string, 0, len(assets.Data))
		for _, url := range assets.Data {
			aid, err := s.CMS.UploadAsset(ctx, id.ProjectID, url)
			if err != nil {
				log.Errorfc(ctx, "cmsintegrationv2: failed to upload asset (%s): %v", url, err)
				return nil
			}
			dataAssets = append(dataAssets, aid)
		}
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
	convStatus := ConvertionStatus("")
	if f.Type == "conv" {
		convStatus = ConvertionStatusSuccess
	}
	item := (&FeatureItem{
		Data:             dataAssets,
		Dic:              dicAssetID,
		MaxLOD:           maxlodAssetID,
		ConvertionStatus: convStatus,
	}).CMSItem()

	_, err := s.CMS.UpdateItem(ctx, id.ItemID, item.Fields, item.MetadataFields)
	if err != nil {
		j1, _ := json.Marshal(item.Fields)
		j2, _ := json.Marshal(item.MetadataFields)
		log.Debugfc(ctx, "cmsintegrationv3: item update for %s: %s, %s", id.ItemID, j1, j2)
		return fmt.Errorf("failed to update item: %w", err)
	}

	// comment to the item
	err = s.CMS.CommentToItem(ctx, id.ItemID, "品質検査・変換が完了しました。"+logmsg)
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}

func failToConvert(ctx context.Context, s *Services, itemID, message string) error {
	if err := s.UpdateFeatureItemStatus(ctx, itemID, ConvertionStatusError); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	if err := s.CMS.CommentToItem(ctx, itemID, message); err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}
