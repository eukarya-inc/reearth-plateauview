package sdk

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/rerror"
)

func RequestHandler(conf Config, g *echo.Group) error {
	s, err := NewServices(conf)
	if err != nil {
		return err
	}

	return requestHandler(conf, g, s)
}

func requestHandler(conf Config, g *echo.Group, s *Services) error {
	g.Use(auth(conf.APIToken))

	g.POST("/request_max_lod", func(c echo.Context) error {
		ctx := c.Request().Context()
		q := struct {
			IDs []string `json:"ids"`
		}{}

		if err := c.Bind(&q); err != nil {
			return err
		}

		log.Infof("sdk: request max lod extraction for %d items: %v", len(q.IDs), q.IDs)

		for _, id := range q.IDs {
			log.Infof("sdk:	request max lod extraction for %s", id)

			i, err := s.CMS.GetItem(ctx, id, false)
			if i == nil || err != nil {
				if err == rerror.ErrNotFound {
					return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
				}
				if err != nil {
					log.Errorf("sdk: failed to get item: %v", err)
				} else {
					log.Errorf("sdk: item is nil: %s", id)
				}
				return echo.NewHTTPError(http.StatusInternalServerError, "internal")
			}

			item := ItemFrom(*i)
			// project id can be empty
			s.RequestMaxLODExtraction(c.Request().Context(), item, true)
		}

		return c.JSON(http.StatusOK, "ok")
	})

	return nil
}

func auth(expected string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if expected != "" {
				token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
				if token != expected {
					return echo.ErrUnauthorized
				}
			}

			return next(c)
		}
	}
}
