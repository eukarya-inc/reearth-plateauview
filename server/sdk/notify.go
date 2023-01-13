package sdk

import (
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/eukarya-inc/reearth-plateauview/server/fme"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
)

type FMEResult struct {
	ID        string `json:"id"`
	ResultURL string `json:"resultUrl"`
}

func NotifyHandler(cmsi cms.Interface, secret string, debug bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var f FMEResult
		if err := c.Bind(&f); err != nil {
			log.Info("sdk notify: invalid payload: %w", err)
			return c.JSON(http.StatusBadRequest, "invalid payload")
		}

		log.Infof("sdk notify: received: %+v", f)

		id, err := fme.ParseID(f.ID, secret)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}

		log.Errorf("sdk notify: validate: itemID=%s, assetID=%s", id.ItemID, id.AssetID)

		aid, err := cmsi.UploadAsset(ctx, id.ProjectID, f.ResultURL)
		if err != nil {
			log.Errorf("sdk notify: failed to update assets: %w", err)

			if _, err := cmsi.UpdateItem(ctx, id.ItemID, Item{
				MaxLODStatus: StatusError,
			}.Fields()); err != nil {
				log.Errorf("sdk notify: failed to update item: %w", err)
			}
			return nil
		}

		if _, err := cmsi.UpdateItem(ctx, id.ItemID, Item{
			MaxLODStatus: StatusOK,
			MaxLOD:       aid,
		}.Fields()); err != nil {
			log.Errorf("sdk notify: failed to update item: %w", err)
			return nil
		}

		log.Infof("sdk notify: done")
		return nil
	}
}
