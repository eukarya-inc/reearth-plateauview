package datacatalogv2

import (
	"context"
	"net/http"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/eukarya-inc/reearth-plateauview/server/putil"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/reearth/reearthx/log"
)

type Config struct {
	plateaucms.Config
	CMSBase      string
	DisableCache bool
	CacheTTL     int
}

func Echo(conf Config, g *echo.Group) error {
	pcms, err := plateaucms.New(conf.Config)
	if err != nil {
		return err
	}

	f, err := NewFetcher(conf.CMSBase)
	if err != nil {
		return err
	}

	g.Use(
		middleware.CORS(),
		middleware.Gzip(),
		putil.NewCacheMiddleware(putil.CacheConfig{
			Disabled:     conf.DisableCache,
			TTL:          time.Duration(conf.CacheTTL) * time.Second,
			CacheControl: true,
		}).Middleware(),
		pcms.AuthMiddleware("pid", nil, true, ""),
	)

	g.GET("/:pid", func(c echo.Context) error {
		ctx := c.Request().Context()
		prj := c.Param("pid")
		res, err := f.Do(ctx, prj, options(ctx, prj))
		if err != nil {
			log.Errorfc(ctx, "datacatalog: %v", err)
			return c.JSON(http.StatusInternalServerError, "error")
		}
		return c.JSON(http.StatusOK, res.All())
	})

	return nil
}

func options(ctx context.Context, prj string) FetcherDoOptions {
	md := plateaucms.GetCMSMetadataFromContext(ctx)
	if md.Name == "" {
		return FetcherDoOptions{}
	}

	return FetcherDoOptions{
		Subproject: md.SubPorjectAlias,
		CityName:   md.Name,
	}
}
