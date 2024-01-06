package datacatalog

import (
	"context"
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/labstack/echo/v4"
)

type Config struct {
	plateaucms.Config
	// v3
	CacheUpdateKey       string
	PlaygroundEndpoint   string
	GraphqlMaxComplexity int
	// v2
	DisableCache bool
	CacheTTL     int
}

func Echo(conf Config, g *echo.Group) error {
	// data catalog API
	updateCache, err := echov3(conf, g)
	if err != nil {
		return fmt.Errorf("failed to initialize datacatalog v3 repo: %w", err)
	}

	// compat: PLATEAU VIEW 2.0 data catalog API
	err = datacatalogv2.Echo(datacatalogv2.Config{
		Config:       conf.Config,
		DisableCache: conf.DisableCache,
		CacheTTL:     conf.CacheTTL,
	}, g)
	if err != nil {
		return fmt.Errorf("failed to initialize datacatalog v2 API: %w", err)
	}

	// first cache update
	return updateCache(context.Background())
}
