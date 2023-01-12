package cmsintegration

import (
	"context"
	"fmt"
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
)

func NotifyHandler(cmsi cms.Interface, secret string) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var f fmeResult
		if err := c.Bind(&f); err != nil {
			log.Info("cmsintegration notify: invalid payload: %w", err)
			return c.JSON(http.StatusBadRequest, "invalid payload")
		}

		log.Infof("cmsintegration notify: received: %+v", f)

		id, err := ParseID(f.ID, secret)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}

		log.Errorf("cmsintegration notify: validate: itemID=%s, assetID=%s", id.ItemID, id.AssetID)

		if f.Status != "ok" && f.Status != "error" {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid type: %s", f.Type))
		}

		if err := c.JSON(http.StatusOK, "ok"); err != nil {
			return err
		}

		cc := commentContent(f.Status, f.Type, f.LogURL)
		if err := cmsi.CommentToItem(c.Request().Context(), id.ItemID, cc); err != nil {
			log.Errorf("cmsintegration notify: failed to comment: %w", err)
			return nil
		}

		if f.Type == "error" {
			if _, err := cmsi.UpdateItem(ctx, id.ItemID, Item{
				ConversionStatus:  StatusError,
				ConversionEnabled: ConversionDisabled,
			}.Fields()); err != nil {
				log.Errorf("cmsintegration notify: failed to update item: %w", err)
				return nil
			}

			return nil
		}

		if _, err := cmsi.UpdateItem(ctx, id.ItemID, Item{
			ConversionStatus: StatusOK,
		}.Fields()); err != nil {
			log.Errorf("cmsintegration notify: failed to update item: %w", err)
			return nil
		}

		r, err := uploadAssets(ctx, cmsi, id.ProjectID, f)
		if err != nil {
			return err
		}

		if _, err := cmsi.UpdateItem(ctx, id.ItemID, r.Fields()); err != nil {
			log.Errorf("cmsintegration notify: failed to update item: %w", err)
			return nil
		}

		log.Infof("cmsintegration notify: done")
		return nil
	}
}

func commentContent(s, t, logURL string) string {
	var log string
	if logURL != "" {
		log = fmt.Sprintf(" ログ: %s", logURL)
	}

	var tt string
	if t == "qualityCheck" {
		tt = "品質検査"
	} else if t == "conversion" {
		tt = "3D Tiles への変換"
	}

	if s == "ok" {
		return fmt.Sprintf("%sに成功しました。%s", tt, log)
	}

	return fmt.Sprintf("%sでエラーが発生しました。%s", tt, log)
}

func uploadAssets(ctx context.Context, c cms.Interface, pid string, f fmeResult) (r Item, _ error) {
	// TODO2: support multiple files
	// TODO2: add retry

	bldg := f.GetResultFromAllLOD("bldg")

	assetID, err := c.UploadAsset(ctx, pid, bldg)
	if err != nil {
		log.Errorf("cmsintegration notify: failed to upload asset: %w", err)
		return r, nil
	}

	log.Infof("cmsintegration notify: asset uploaded: %s", assetID)

	r.Bldg = []string{assetID}
	return
}
