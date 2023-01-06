package cmsintegration

import (
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/cms/cmswebhook"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/fme"
	"github.com/reearth/reearthx/log"
)

const (
	modelKey        = "plateau"
	cityGMLFieldKey = "citygml"
	bldgFieldKey    = "bldg"
)

func WebhookHandler(c Config) (cmswebhook.Handler, error) {
	s, err := NewServices(c)
	if err != nil {
		return nil, err
	}

	return func(req *http.Request, w *cmswebhook.Payload) error {
		if !w.Operator.IsUser() {
			log.Infof("cmsintegration webhook: invalid event operator: %+v", w.Operator)
			return nil
		}

		ctx := req.Context()

		if w.Type != "item.update" && w.Type != "item.create" {
			log.Infof("cmsintegration webhook: invalid event type: %s", w.Type)
			return nil
		}

		if w.Data.Model.Key != modelKey {
			log.Infof("cmsintegration webhook: invalid model id: %s, key: %s", w.Data.Item.ModelID, w.Data.Model.Key)
			return nil
		}

		assetField := w.Data.FieldByKey(cityGMLFieldKey)
		if assetField == nil || assetField.Value == nil {
			log.Infof("cmsintegration webhook: asset field not found")
			return nil
		}
		if v := assetField.ValueString(); v == nil || *v == "" {
			log.Infof("cmsintegration webhook: asset field empty")
			return nil
		}

		bldgField := w.Data.FieldByKey(bldgFieldKey)
		if bldgField != nil && bldgField.Value != nil {
			if s := bldgField.ValueStrings(); len(s) > 0 {
				log.Infof("cmsintegration webhook: 3dtiles already converted: field=%+v", bldgField)
				return nil
			}
		}

		assetID := assetField.ValueString()
		if assetID == nil || *assetID == "" {
			log.Infof("cmsintegration webhook: invalid field value: %+v", assetField)
			return nil
		}

		asset, err := s.CMS.Asset(ctx, *assetID)
		if err != nil || asset == nil || asset.ID == "" {
			log.Infof("cmsintegration webhook: cannot fetch asset: %w", err)
			return nil
		}

		fmeReq := fme.Request{
			ID: ID{
				ItemID:      w.Data.Item.ID,
				AssetID:     asset.ID,
				ProjectID:   w.Data.Schema.ProjectID,
				BldgFieldID: bldgField.ID,
			}.String(c.Secret),
			Target: asset.URL,
			PRCS:   "6669", // TODO2: accept prcs code from webhook
		}

		if s.FME == nil {
			log.Infof("webhook: fme mocked: %+v", fmeReq)
		} else if err := s.FME.CheckQualityAndConvertAll(ctx, fmeReq); err != nil {
			log.Errorf("cmsintegration webhook: failed to request fme: %w", err)
			return nil
		}

		if err := s.CMS.Comment(ctx, asset.ID, "CityGMLの品質検査及び3D Tilesへの変換を開始しました。"); err != nil {
			log.Errorf("cmsintegration webhook: failed to comment: %w", err)
			return nil
		}

		log.Infof("cmsintegration webhook: done")

		return nil
	}, nil
}
