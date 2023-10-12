package datacatalog

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogv2adapter"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/reearth/reearthx/log"
)

type Config struct {
	plateaucms.Config
	CMSBase        string
	DisableCache   bool
	CacheTTL       int
	CacheUpdateKey string
}

func Echo(conf Config, g *echo.Group) error {
	// TODO: merge 2022 and later 2023 projects
	repo, err := datacatalogv2adapter.New(conf.Config.CMSBaseURL, "plateau-2022")
	if err != nil {
		return fmt.Errorf("failed to initialize datacatalog repository: %w", err)
	}

	// PLATEAU API
	plateauapig := g.Group("")
	plateauapig.Use(
		middleware.CORS(),
		middleware.Gzip(),
	)

	srv := plateauapi.NewService(repo)
	plateauapig.GET("/graphql", echo.WrapHandler(playground.Handler("PLATEAU GraphQL API Playground", "/datacatalog/graphql")))
	plateauapig.POST("/graphql", echo.WrapHandler(srv))

	// compat: PLATEAU VIEW 2.0 API
	v2apig := g.Group("")
	err = datacatalogv2.Echo(datacatalogv2.Config{
		Config:       conf.Config,
		CMSBase:      conf.CMSBase,
		DisableCache: conf.DisableCache,
		CacheTTL:     conf.CacheTTL,
	}, v2apig)
	if err != nil {
		return fmt.Errorf("failed to initialize datacatalog v2 API: %w", err)
	}

	// cache update API
	g.POST("/update-cache", func(c echo.Context) error {
		if conf.CacheUpdateKey != "" {
			b := struct {
				Key string `json:"key"`
			}{}
			if err := c.Bind(&b); err != nil {
				return echo.ErrUnauthorized
			}
			if b.Key != conf.CacheUpdateKey {
				return echo.ErrUnauthorized
			}
		}

		if err := repo.UpdateCache(c.Request().Context(), datacatalogv2.FetcherDoOptions{}); err != nil {
			log.Errorfc(c.Request().Context(), "datacatalog: failed to update cache: %w", err)
			return echo.ErrInternalServerError
		}

		return nil
	})

	// first cache update
	if err := repo.UpdateCache(context.Background(), datacatalogv2.FetcherDoOptions{}); err != nil {
		log.Fatalf("datacatalog: failed to update cache: %w", err)
	}

	return nil
}
