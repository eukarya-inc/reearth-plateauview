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

		return c.JSON(http.StatusOK, res)
	}, nil)

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
