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

		if err := f.CheckQualityAndConvertAll(c.Request().Context(), fme.Request{
			ID: ID{
				// TODO: get these values from body
				ItemID:  "",
				AssetID: "",
			}.String(secret),
			Target: "",
			PRCS:   "6669", // TODO2: accept prcs code from webhook
		}); err != nil {
			log.Errorf("webhook: failed to request fme: %w", err)
			return nil
		}

		return nil
	}
}
