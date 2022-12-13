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
	h, err := NewHandler(s.CMS, c.DataModelKey, c.TemplateModelKey)
	if err != nil {
		panic("failed to init echo")
	}
	g.GET("/viz/:pid", h.fetchRoot(s.CMS), authMiddleware(c.VToken))
	g.GET("/viz/:pid/data", h.getAllDataHandler(s.CMS), authMiddleware(c.VToken))
	g.GET("/viz/:pid/data/:iid", h.getDataHandler(s.CMS), authMiddleware(c.VToken))
	g.POST("/viz/:pid/data", h.createDataHandler(s.CMS), authMiddleware(c.VToken))
	g.PATCH("/viz/:pid/data/:iid", h.updateDataHandler(s.CMS), authMiddleware(c.VToken))
	g.DELETE("/vis/:pid/data/:iid", h.deleteDataHandler(s.CMS), authMiddleware(c.VToken))
	g.GET("/viz/:pid/templates", h.fetchTemplate(s.CMS), authMiddleware(c.VToken))
	g.POST("/viz/:pid/templates", h.createTemplateHandler(s.CMS), authMiddleware(c.VToken))
	g.PATCH("/viz/:pid/templates/:iid", h.updateTemplateHandler(s.CMS), authMiddleware(c.VToken))
	g.DELETE("/viz/:pid/templates/:iid", h.deleteTemplateHandler(s.CMS), authMiddleware(c.VToken))
}

func authMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			header := req.Header.Get("Authorization")
			token := strings.TrimPrefix(header, "Bearer ")
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
