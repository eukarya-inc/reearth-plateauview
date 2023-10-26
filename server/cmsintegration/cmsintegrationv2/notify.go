package cmsintegrationv2

import (
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/fme"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
)

func NotifyHandler(conf Config) (echo.HandlerFunc, error) {
	// s, err := NewServices(conf)
	// if err != nil {
	// 	return nil, err
	// }

	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var f FMEResult
		if err := c.Bind(&f); err != nil {
			log.Info("cmsintegrationv2 notify: invalid payload: %w", err)
			return c.JSON(http.StatusBadRequest, "invalid payload")
		}

		log.Infofc(ctx, "cmsintegrationv2 notify: received: %+v", f)

		_, err := fme.ParseID(f.ID, conf.Secret)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "unauthorized")
		}

		// TODO

		return nil
	}, nil
}
