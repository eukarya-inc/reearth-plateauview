package sdkapiv3

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/reearth/reearthx/log"
)

func Handler(conf Config, g *echo.Group) error {
	if conf.GQLBaseURL == "" || conf.GQLToken == "" {
		return nil
	}

	client, err := NewClient(conf)
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	g.Use(
		auth(conf.Token),
		middleware.Gzip(),
	)

	g.GET("/datasets", func(c echo.Context) error {
		res, err := client.QueryDatasets()
		if err != nil {
			log.Errorfc(c.Request().Context(), "sdkapiv3: error querying datasets: %v", err)
			return c.JSON(http.StatusBadGateway, map[string]any{"error": "bad gateway"})
		}

		return c.JSON(http.StatusOK, res.ToDatasets())
	}, nil)

	g.GET("/datasets/:id/files", func(c echo.Context) error {
		id := strings.TrimPrefix(c.Param("id"), "a_")
		res, err := client.QueryDatasetFiles(id)
		if err != nil {
			log.Errorfc(c.Request().Context(), "sdkapiv3: error querying dataset files: %v", err)
			return c.JSON(http.StatusBadGateway, map[string]any{"error": "bad gateway"})
		}

		return c.JSON(http.StatusOK, res.ToDatasetFiles())
	}, nil)

	log.Infof("sdkapiv3: initialized")
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
