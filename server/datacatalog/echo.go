package datacatalog

import (
	"context"
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogv2adapter"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
)

type Config struct {
	plateaucms.Config
	CMSBase              string
	CMSToken             string
	DisableCache         bool
	CacheTTL             int
	CacheUpdateKey       string
	PlaygroundEndpoint   string
	GraphqlMaxComplexity int
}

func Echo(conf Config, g *echo.Group) error {
	repov2, err := datacatalogv2adapter.New(conf.Config.CMSBaseURL, "plateau-2022")
	if err != nil {
		return fmt.Errorf("failed to initialize datacatalog repository: %w", err)
	}

	if err := echov3(conf, g, repov2); err != nil {
		return fmt.Errorf("failed to initialize datacatalog v3 repo: %w", err)
	}

	if err := echov2(conf, g, repov2); err != nil {
		return fmt.Errorf("failed to initialize datacatalog v2 repo: %w", err)
	}

	return nil
}

func echov2(conf Config, g *echo.Group, repov2 *plateauapi.RepoWrapper) (err error) {
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

		if err := repov2.Update(c.Request().Context()); err != nil {
			log.Errorfc(c.Request().Context(), "datacatalog: failed to update cache: %v", err)
			return echo.ErrInternalServerError
		}

		return nil
	})

	// first cache update
	if err := repov2.Update(context.Background()); err != nil {
		log.Errorf("datacatalog: failed to update cache: %w", err)
	}

	return nil
}
