package datacatalog

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogv2adapter"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv3"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
	"golang.org/x/sync/errgroup"
)

type reposHandler struct {
	reposv3            *datacatalogv3.Repos
	reposv2            *datacatalogv2adapter.Repos
	pcms               *plateaucms.CMS
	gqlComplexityLimit int
	cacheUpdateKey     string
}

const pidParamName = "pid"
const gqlComplexityLimit = 1000
const cmsSchemaVersion = "v3"
const cmsSchemaVersionV2 = "v2"

func newReposHandler(conf Config) (*reposHandler, error) {
	pcms, err := plateaucms.New(conf.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize plateau cms: %w", err)
	}

	reposv3 := datacatalogv3.NewRepos()
	reposv2 := datacatalogv2adapter.NewRepos()

	if conf.GraphqlMaxComplexity <= 0 {
		conf.GraphqlMaxComplexity = gqlComplexityLimit
	}

	return &reposHandler{
		reposv3:            reposv3,
		reposv2:            reposv2,
		pcms:               pcms,
		gqlComplexityLimit: conf.GraphqlMaxComplexity,
		cacheUpdateKey:     conf.CacheUpdateKey,
	}, nil
}

func (h *reposHandler) Middleware() echo.MiddlewareFunc {
	return h.pcms.AuthMiddleware(plateaucms.AuthMiddlewareConfig{
		Key:             pidParamName,
		FindDataCatalog: true,
		UseDefault:      true,
	})
}

func (h *reposHandler) Handler(admin bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		merged, err := h.prepareMergedRepo(c, admin)
		if err != nil {
			return err
		}

		srv := plateauapi.NewService(merged, plateauapi.FixedComplexityLimit(h.gqlComplexityLimit))
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	}
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
	if md.DataCatalogProjectAlias != pid || !isV3(md) {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	if !md.Auth {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	t := h.reposv3.UpdatedAt(pid)
	res := ""
	if !t.IsZero() {
		res = fmt.Sprintf("updated at: %s\n", t.Format(time.RFC3339))
	}
	res += strings.Join(h.reposv3.Warnings(pid), "\n")
	return c.String(http.StatusOK, res)
}

func (h *reposHandler) UpdateCache(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, p := range h.reposv3.Projects() {
		p := p
		g.Go(func() error {
			return h.updateV3(ctx, p)
		})
	}

	for _, p := range h.reposv2.Projects() {
		p := p
		g.Go(func() error {
			return h.updateV2(ctx, p)
		})
	}

	return g.Wait()
}

func (h *reposHandler) Init(ctx context.Context) error {
	metadata, err := h.pcms.AllMetadata(ctx, true)
	if err != nil {
		return fmt.Errorf("datacatalogv3: failed to get all metadata: %w", err)
	}

	plateauMetadata := metadata.PlateauProjects()
	if err := h.prepareAll(ctx, plateauMetadata); err != nil {
		return err
	}

	return nil
}

func (h *reposHandler) prepareMergedRepo(c echo.Context, admin bool) (plateauapi.Repo, error) {
	ctx := c.Request().Context()
	md := plateaucms.GetCMSMetadataFromContext(ctx)
	if admin && !md.Auth {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	pid := c.Param(pidParamName)
	mds := plateaucms.GetAllCMSMetadataFromContext(ctx)
	merged := h.prepareAndGetMergedRepo(ctx, admin, pid, mds)
	if merged == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "not found")
	}

	log.Debugfc(ctx, "datacatalogv3: use repo for %s: %s", pid, merged.Name())
	return merged, nil
}

func (h *reposHandler) prepareAndGetMergedRepo(ctx context.Context, admin bool, project string, metadata plateaucms.MetadataList) plateauapi.Repo {
	var mds plateaucms.MetadataList
	if project == "" {
		mds = metadata.PlateauProjects()
	} else {
		mds = metadata.FindDataCatalogAndSub(project)
	}

	if err := h.prepareAll(ctx, mds); err != nil {
		log.Errorfc(ctx, "datacatalogv3: failed to prepare repos: %w", err)
	}

	repos := make([]plateauapi.Repo, 0, len(mds))
	for _, s := range mds {
		if r := h.getRepo(admin, s); r != nil {
			repos = append(repos, r)
		}
	}

	if len(repos) == 0 {
		return nil
	}

	if len(repos) == 1 {
		return repos[0]
	}

	merged := plateauapi.NewMerger(repos...)
	if err := merged.Init(ctx); err != nil {
		log.Errorfc(ctx, "datacatalogv3: failed to initialize merged repo: %w", err)
		return nil
	}

	return merged
}

func (h *reposHandler) getRepo(admin bool, md plateaucms.Metadata) (repo plateauapi.Repo) {
	if md.DataCatalogProjectAlias == "" {
		return
	}

	if isV2(md) {
		repo = h.reposv2.Repo(md.DataCatalogProjectAlias, admin)
	} else if isV3(md) {
		repo = h.reposv3.Repo(md.DataCatalogProjectAlias, admin)
	}
	return
}

func (h *reposHandler) prepareAll(ctx context.Context, metadata plateaucms.MetadataList) error {
	errg, ctx := errgroup.WithContext(ctx)
	for _, md := range metadata {
		md := md

		errg.Go(func() error {
			if err := h.prepare(ctx, md); err != nil {
				return fmt.Errorf("datacatalogv3: failed to prepare repo for %s: %w", md.DataCatalogProjectAlias, err)
			}
			return nil
		})
	}
	return errg.Wait()
}

func (h *reposHandler) prepare(ctx context.Context, md plateaucms.Metadata) error {
	if isV2(md) {
		return h.prepareV2(ctx, md)
	}
	return h.prepareV3(ctx, md)
}

func (h *reposHandler) prepareV2(ctx context.Context, md plateaucms.Metadata) error {
	if !isV2(md) {
		return nil
	}

	fetcher, err := datacatalogv2.NewFetcher(md.CMSBaseURL)
	if err != nil {
		return fmt.Errorf("datacatalogv2: failed to create fetcher %s: %w", md.DataCatalogProjectAlias, err)
	}

	if err := h.reposv2.Prepare(ctx, md.DataCatalogProjectAlias, fetcher); err != nil {
		return fmt.Errorf("datacatalogv2: failed to prepare repo for %s: %w", md.DataCatalogProjectAlias, err)
	}

	return nil
}

func (h *reposHandler) prepareV3(ctx context.Context, md plateaucms.Metadata) error {
	if !isV3(md) {
		return nil
	}

	cms, err := md.CMS()
	if err != nil {
		return fmt.Errorf("datacatalogv3: failed to create cms for %s: %w", md.DataCatalogProjectAlias, err)
	}

	if err := h.reposv3.Prepare(ctx, md.DataCatalogProjectAlias, md.PlateauYear(), cms); err != nil {
		return fmt.Errorf("datacatalogv3: failed to prepare repo for %s: %w", md.DataCatalogProjectAlias, err)
	}

	return nil
}

func (h *reposHandler) updateV2(ctx context.Context, prj string) error {
	if _, err := h.reposv2.Update(ctx, prj); err != nil {
		return fmt.Errorf("datacatalogv2: failed to update repo %s: %w", prj, err)
	}
	return nil
}

func (h *reposHandler) updateV3(ctx context.Context, prj string) error {
	if _, err := h.reposv3.Update(ctx, prj); err != nil {
		return fmt.Errorf("datacatalogv3: failed to update repo %s: %w", prj, err)
	}
	return nil
}

func isV2(md plateaucms.Metadata) bool {
	return md.DataCatalogSchemaVersion == "" || md.DataCatalogSchemaVersion == cmsSchemaVersionV2
}

func isV3(md plateaucms.Metadata) bool {
	return md.DataCatalogSchemaVersion == cmsSchemaVersion
}
