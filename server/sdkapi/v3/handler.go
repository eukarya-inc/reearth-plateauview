package sdkapiv3

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Handler(conf Config, g *echo.Group) {

	g.Use(
		auth(conf.Token),
		middleware.Gzip(),
	)

	g.GET("/datasets", func(c echo.Context) error {
		client := NewClient(conf)

		_, err := client.QueryDatasets()
		if err != nil {
			return fmt.Errorf("error querying datasets: %w", err)
		}

		return c.JSON(http.StatusOK, nil)

	}, nil)

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
