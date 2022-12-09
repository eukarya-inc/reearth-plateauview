package visualizer

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func Echo(g *echo.Group, c Config) error {
	s, err := NewServices(c)
	if err != nil {
		return err
	}

	initEcho(g, c, s)
	return nil
}

func initEcho(g *echo.Group, c Config, s Services) {
	g.GET("/viz/:pid", fetchRoot(s.CMS), authMiddleware(c.VToken))
	g.GET("/viz/:pid/data", getDataHandler(s.CMS))
	g.POST("/viz/:pid/data", createDataHandler(s.CMS))
	g.PATCH("/viz/:pid/data/:did", updateDataHandler(s.CMS))
	g.DELETE("/vis/:pid/data/:did", exampleHandler)
	g.GET("/viz/:pid/templates", exampleHandler)
	g.POST("/viz/:pid/templates", createTemplateHandler(s.CMS))
	g.PATCH("/viz/:pid/templates/:tid", exampleHandler)
	g.DELETE("/viz/:pid/templates/:tid", exampleHandler)
}

func authMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			ctx := req.Context()

			// TODO: ここでbearer tokenを取得
			header := req.Header.Get("Authorization")
			token := strings.TrimPrefix(header, "bearer ")
			if token != secret {
				return c.JSON(http.StatusUnauthorized, nil)
			}
			return next(c)
		}
	}
}

func exampleHandler(c echo.Context) error {
	return nil
}
