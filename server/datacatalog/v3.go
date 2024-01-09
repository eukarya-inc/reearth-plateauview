package datacatalog

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogv2adapter"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv3"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

// TODO
const datacatalogDefaultProject = "plateau-2023"
const datacatalogv2project = "plateau-2022"

func echov3(conf Config, g *echo.Group) (func(ctx context.Context) error, error) {
	h, err := newReposHandler(conf)
	if err != nil {
		return nil, err
	}

	// PLATEAU API
	plateauapig := g.Group("")
	plateauapig.Use(
		middleware.CORS(),
		middleware.Gzip(),
		h.Middleware(),
	)

	// GraphQL playground (all)
	plateauapig.GET("/graphql", gqlPlaygroundHandler(conf.PlaygroundEndpoint, false))

	// GraphQL playground (project)
	plateauapig.GET("/:pid/graphql", gqlPlaygroundHandler(conf.PlaygroundEndpoint, false))

	// GraphQL playground (admin)
	plateauapig.GET("/:pid/admin/graphql", gqlPlaygroundHandler(conf.PlaygroundEndpoint, true))

	// GraphQL API (all)
	plateauapig.POST("/graphql", h.Handler(false))

	// GraphQL API (project)
	plateauapig.POST("/:pid/graphql", h.Handler(false))

	// GraphQL API (admin)
	plateauapig.POST("/:pid/admin/graphql", h.Handler(true))

	// warning API
	plateauapig.GET("/:pid/warnings", h.WarningHandler)

	// cache update API
	g.POST("/update-cache", h.UpdateCacheHandler)

	return func(ctx context.Context) error {
		return h.Init(ctx)
	}, nil
}

type reposHandler struct {
	reposv3            *datacatalogv3.Repos
	repov2             *plateauapi.RepoWrapper
	pcms               *plateaucms.CMS
	gqlComplexityLimit int
	cacheUpdateKey     string
}

const pidParamName = "pid"
const gqlComplexityLimit = 1000
const cmsSchemaVersion = "v3"

func newReposHandler(conf Config) (*reposHandler, error) {
	pcms, err := plateaucms.New(conf.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize plateau cms: %w", err)
	}

	repov2, err := datacatalogv2adapter.New(conf.Config.CMSBaseURL, datacatalogv2project)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize datacatalog v2 repo: %w", err)
	}

	reposv3 := datacatalogv3.NewRepos()

	if conf.GraphqlMaxComplexity <= 0 {
		conf.GraphqlMaxComplexity = gqlComplexityLimit
	}

	return &reposHandler{
		reposv3:            reposv3,
		repov2:             repov2,
		pcms:               pcms,
		gqlComplexityLimit: conf.GraphqlMaxComplexity,
		cacheUpdateKey:     conf.CacheUpdateKey,
	}, nil
}

func (h *reposHandler) Middleware() echo.MiddlewareFunc {
	return h.pcms.AuthMiddleware(pidParamName, nil, true, datacatalogDefaultProject)
}

func (h *reposHandler) Handler(admin bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		pid := c.Param(pidParamName)
		repo, err := h.getRepo(c.Request().Context(), admin, pid == "")
		if err != nil {
			return err
		}
		if repo == nil {
			return echo.NewHTTPError(http.StatusNotFound, "not found")
		}

		srv := plateauapi.NewService(repo, plateauapi.FixedComplexityLimit(h.gqlComplexityLimit))
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func (h *reposHandler) getRepo(ctx context.Context, admin, mergev2 bool) (plateauapi.Repo, error) {
	cms := plateaucms.GetCMSFromContext(ctx)
	if cms == nil {
		return nil, nil
	}

	cmsmd := plateaucms.GetCMSMetadataFromContext(ctx)
	project := cmsmd.DataCatalogProjectAlias
	if project == "" || (admin && !cmsmd.Auth) || cmsmd.DataCatalogSchemaVersion != cmsSchemaVersion {
		return nil, nil
	}

	log.Debugfc(ctx, "datacatalogv3: use CMS project: %s", project)

	if err := h.reposv3.Prepare(ctx, project, cms); err != nil {
		return nil, err
	}

	rw := h.reposv3.Repo(project, admin)

	if mergev2 {
		return plateauapi.NewMerger(rw, h.repov2), nil
	}
	return rw, nil
}

func (h *reposHandler) UpdateCacheHandler(c echo.Context) error {
	if h.cacheUpdateKey != "" {
		b := struct {
			Key string `json:"key"`
		}{}
		if err := c.Bind(&b); err != nil {
			return echo.ErrUnauthorized
		}
		if b.Key != h.cacheUpdateKey {
			return echo.ErrUnauthorized
		}
	}

	if err := h.UpdateCache(c.Request().Context()); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "ok")
}

func (h *reposHandler) WarningHandler(c echo.Context) error {
	pid := c.Param(pidParamName)
	md := plateaucms.GetCMSMetadataFromContext(c.Request().Context())
	if md.DataCatalogProjectAlias != pid {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	if !md.Auth {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	w := strings.Join(h.reposv3.Warnings(pid), "\n")
	return c.String(http.StatusOK, w)
}

func (h *reposHandler) UpdateCache(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return h.repov2.Update(ctx)
	})

	g.Go(func() error {
		return h.reposv3.UpdateAll(ctx)
	})

	return g.Wait()
}

func (h *reposHandler) Init(ctx context.Context) error {
	all, err := h.pcms.AllMetadata(ctx, true)
	if err != nil {
		return fmt.Errorf("failed to get all metadata: %w", err)
	}

	target := lo.Filter(all, func(m plateaucms.Metadata, _ int) bool {
		return m.DataCatalogSchemaVersion == cmsSchemaVersion
	})

	log.Infofc(ctx, "datacatalogv3: initializing repos for %d projects", len(target))

	for _, md := range target {
		cms, err := md.CMS()
		if err != nil {
			log.Errorfc(ctx, "datacatalogv3: failed to create cms for %s: %w", md.DataCatalogProjectAlias, err)
			continue
		}

		if err := h.reposv3.Prepare(ctx, md.DataCatalogProjectAlias, cms); err != nil {
			log.Errorfc(ctx, "datacatalogv3: failed to prepare repo for %s: %w", md.DataCatalogProjectAlias, err)
		}
	}

	return nil
}

func gqlPlaygroundHandler(endpoint string, admin bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		pid := c.Param(pidParamName)

		p := make([]string, 0, 4)
		p = append(p, endpoint)
		if pid != "" {
			p = append(p, pid)
		}
		if admin {
			p = append(p, "admin")
		}
		p = append(p, "graphql")

		h := plateauapi.PlaygroundHandler(
			"PLATEAU GraphQL API Playground",
			path.Join(p...),
		)
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
