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
	// if event type is "item.create" and payload is metadata, skip it
	if w.Type != cmswebhook.EventItemCreate && w.ItemData.Item.OriginalItemID != nil {
		return nil
	}

	featureType := strings.TrimPrefix(w.ItemData.Model.Key, modelPrefix)
	mainItem, err := s.GetMainItemWithMetadata(ctx, w.ItemData.Item)
	if err != nil {
		return err
	}

	item := FeatureItemFrom(mainItem)
	if item.ConvertionStatus != "" && item.ConvertionStatus != ConvertionStatusNotStarted {
		log.Debugfc(ctx, "cmsintegrationv3: already converted")
		return nil
	}

	if item.SkipQC && item.SkipConvert || item.CityGML == "" || item.City == "" {
		log.Debugfc(ctx, "cmsintegrationv3: skip convert")
		return nil
	}

	log.Infofc(ctx, "cmsintegrationv3: sendRequestToFME")

	ty := fmeTypeQcConv
	if item.SkipQC {
		ty = fmeTypeConv
	} else if item.SkipConvert {
		ty = fmeTypeQC
	}

	// update convertion status
	err = s.UpdateFeatureItemStatus(ctx, mainItem.ID, ty, ConvertionStatusRunning)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	// get CityGML asset
	cityGMLAsset, err := s.CMS.Asset(ctx, item.CityGML)
	if err != nil {
		_ = failToConvert(ctx, s, mainItem.ID, ty, "CityGMLが見つかりません。")
		return fmt.Errorf("failed to get citygml asset: %w", err)
	}

	// get city item
	cityItemRaw, err := s.CMS.GetItem(ctx, item.City, false)
	if err != nil {
		_ = failToConvert(ctx, s, mainItem.ID, ty, "都市アイテムが見つかりません。")
		return fmt.Errorf("failed to get city item: %w", err)
	}

	cityItem := CityItemFrom(cityItemRaw)
	if cityItem.CodeLists == "" {
		_ = failToConvert(ctx, s, mainItem.ID, ty, "コードリストが都市アイテムに登録されていないため品質検査・変換を開始できません。")
		return fmt.Errorf("city item has no codelist")
	}

	// get codelist asset
	codelistAsset, err := s.CMS.Asset(ctx, cityItem.CodeLists)
	if err != nil {
		_ = failToConvert(ctx, s, mainItem.ID, ty, "コードリストが見つかりません。")
		return fmt.Errorf("failed to get codelist asset: %w", err)
	}

	// request to fme
	err = s.FME.Request(ctx, fmeRequest{
		ID: fmeID{
			ItemID:      mainItem.ID,
			ProjectID:   w.ProjectID(),
			FeatureType: featureType,
			Type:        string(ty),
		}.String(conf.Secret),
		Target:    cityGMLAsset.URL,
		PRCS:      cityItem.PRCS.ESPGCode(),
		Codelists: codelistAsset.URL,
		ResultURL: resultURL(conf),
		Type:      ty,
		// Config: ,
	})
	if err != nil {
		_ = failToConvert(ctx, s, mainItem.ID, ty, "FMEへのリクエストに失敗しました。%v", err)
		return fmt.Errorf("failed to request to fme: %w", err)
	}

	// post a comment to the item
	err = s.CMS.CommentToItem(ctx, mainItem.ID, fmt.Sprintf("%sを開始しました。", ty.Title()))
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

	log.Infofc(ctx, "cmsintegrationv3: receiveResultFromFME")

	logmsg := ""
	if f.LogURL != "" {
		logmsg = "ログ： " + f.LogURL
	}

	// handle error
	if f.Status == "error" {
		log.Warnfc(ctx, "cmsintegrationv3: failed to convert: %v", f.LogURL)
		_ = failToConvert(ctx, s, id.ItemID, fmeRequestType(id.Type), "%sに失敗しました。%s", fmeRequestType(id.Type).Title(), logmsg)
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
				log.Errorfc(ctx, "cmsintegrationv3: failed to upload asset (%s): %v", url, err)
				return nil
			}
			dataAssets = append(dataAssets, aid)
		}
	}

	// upload dic
	var dicAssetID string
	if assets.Dic != "" {
		log.Debugfc(ctx, "cmsintegrationv3: upload dic: %s", assets.Dic)
		var err error
		dicAssetID, err = s.CMS.UploadAsset(ctx, id.ProjectID, assets.Dic)
		if err != nil {
			return fmt.Errorf("failed to upload dic: %w", err)
		}
	}

	// upload maxlod
	var maxlodAssetID string
	if assets.MaxLOD != "" {
		log.Debugfc(ctx, "cmsintegrationv3: upload maxlod: %s", assets.MaxLOD)
		var err error
		maxlodAssetID, err = s.CMS.UploadAsset(ctx, id.ProjectID, assets.MaxLOD)
		if err != nil {
			return fmt.Errorf("failed to upload maxlod: %w", err)
		}
	}

	// upload qc result
	var qcResult string
	if f.Status != "error" && f.LogURL != "" {
		log.Debugfc(ctx, "cmsintegrationv3: upload qc result: %s", f.LogURL)
		var err error
		qcResult, err = s.CMS.UploadAsset(ctx, id.ProjectID, f.LogURL)
		if err != nil {
			return fmt.Errorf("failed to upload qc result: %w", err)
		}
	}

	// update item
	convStatus := ConvertionStatus("")
	qcStatus := ConvertionStatus("")

	if f.Type == "conv" {
		convStatus = ConvertionStatusSuccess
		if id.Type == string(fmeTypeQcConv) {
			qcStatus = ConvertionStatusSuccess
		}
	} else if f.Type == "qc" {
		convStatus = ConvertionStatusSuccess
	}

	item := (&FeatureItem{
		Data:             dataAssets,
		Dic:              dicAssetID,
		MaxLOD:           maxlodAssetID,
		ConvertionStatus: convStatus,
		QCStatus:         qcStatus,
		QCResult:         qcResult,
	}).CMSItem()

	_, err := s.CMS.UpdateItem(ctx, id.ItemID, item.Fields, item.MetadataFields)
	if err != nil {
		j1, _ := json.Marshal(item.Fields)
		j2, _ := json.Marshal(item.MetadataFields)
		log.Debugfc(ctx, "cmsintegrationv3: item update for %s: %s, %s", id.ItemID, j1, j2)
		return fmt.Errorf("failed to update item: %w", err)
	}

	// comment to the item
	err = s.CMS.CommentToItem(ctx, id.ItemID, fmt.Sprintf("%sが完了しました。%s", fmeRequestType(id.Type).Title(), logmsg))
	if err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}

func failToConvert(ctx context.Context, s *Services, itemID string, convType fmeRequestType, message string, args ...any) error {
	if err := s.UpdateFeatureItemStatus(ctx, itemID, convType, ConvertionStatusError); err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	if err := s.CMS.CommentToItem(ctx, itemID, fmt.Sprintf(message, args...)); err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	return nil
}
