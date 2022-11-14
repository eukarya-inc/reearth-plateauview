package cmsintegration

import (
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cms"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/fme"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/webhook"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
)

func WebhookHandler(f fme.Interface, cms cms.Interface, modelID, citygmlFieldKey, bldgFieldKey, secret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		w := webhook.GetPayload(c.Request().Context())
		if w == nil {
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}

		if err := c.JSON(http.StatusOK, "ok"); err != nil {
			return err
		}

		if w.Type != "item.update" && w.Type != "item.create" {
			log.Infof("webhook: invalid event type: %s", w.Type)
			return nil
		}

		if w.Data.Item.ModelID != modelID {
			log.Infof("webhook: invalid model id: %s", w.Data.Item.ModelID)
			return nil
		}

		assetFieldID := w.Data.Schema.FieldIDByKey(citygmlFieldKey)
		assetField := w.Data.Item.Field(assetFieldID)
		if assetField == nil {
			log.Infof("webhook: field not found: fieldId=%s", assetFieldID)
			return nil
		}
		asset := assetField.Value
		if asset == nil || asset.ID == "" || asset.URL == "" {
			log.Infof("webhook: invalid citygml field value: %+v", assetField)
			return nil
		}

		if err := cms.Comment(c.Request().Context(), asset.ID, "品質検査及び3D Tilesへの変換を開始しました。"); err != nil {
			log.Errorf("notify: failed to comment: %w", err)
			return nil
		}

		if err := f.CheckQualityAndConvertAll(c.Request().Context(), fme.Request{
			ID: ID{
				ItemID:       w.Data.Item.ID,
				AssetID:      asset.ID,
				TilesFieldID: bldgFieldKey,
			}.String(secret),
			Target: asset.URL,
			PRCS:   "6669", // TODO2: accept prcs code from webhook
		}); err != nil {
			log.Errorf("webhook: failed to request fme: %w", err)
			return nil
		}

		return nil
	}
}
