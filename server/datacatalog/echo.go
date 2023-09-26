package plateauapi

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfig struct {
	plateaucms.Config
	CMSBase      string
	DisableCache bool
	CacheTTL     int
}

func Echo(conf ServerConfig, g *echo.Group) error {
	g.Use(
		middleware.CORS(),
		middleware.Gzip(),
	)

	// PLATEAU API
	srv := plateauapi.NewService(nil) // TODO
	g.GET("/graphql", echo.WrapHandler(playground.Handler("PLATEAU GraphQL API Playground", "/graphql")))
	g.POST("/graphql", echo.WrapHandler(srv))

	// PLATEAU VIEW 3.0 API

	// PLATEAU VIEW 2.0 API
	classic := g.Group("")
	return datacatalogv2.Echo(datacatalogv2.Config{
		Config:       conf.Config,
		CMSBase:      conf.CMSBase,
		DisableCache: conf.DisableCache,
		CacheTTL:     conf.CacheTTL,
	}, classic)
}
