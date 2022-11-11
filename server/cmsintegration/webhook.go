package cmsintegration

import (
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/fme"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/webhook"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
)

func WebhookHandler(f fme.Interface, secret string) echo.HandlerFunc {
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

		// TODO: get these values from body
		id := ID{
			ItemID:  "",
			AssetID: "",
		}.String(secret)
		var target, prcs string

		if err := f.CheckQualityAndConvertAll(c.Request().Context(), fme.Request{
			ID:     id,
			Target: target,
			PRCS:   prcs,
		}); err != nil {
			log.Errorf("webhook: failed to request fme: %w", err)
			return nil
		}

		return nil
	}
}
