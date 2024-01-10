package datacatalog

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogv2adapter"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv3"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearthx/log"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type reposHandler struct {
	reposv3            *datacatalogv3.Repos
	reposv2            map[string]*plateauapi.RepoWrapper
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

	if conf.GraphqlMaxComplexity <= 0 {
		conf.GraphqlMaxComplexity = gqlComplexityLimit
	}

	return &reposHandler{
		reposv3:            reposv3,
		reposv2:            map[string]*plateauapi.RepoWrapper{},
		pcms:               pcms,
		gqlComplexityLimit: conf.GraphqlMaxComplexity,
		cacheUpdateKey:     conf.CacheUpdateKey,
	}, nil
}

func (h *reposHandler) Middleware() echo.MiddlewareFunc {
	return h.pcms.AuthMiddleware(pidParamName, nil, true, "")
}

func (h *reposHandler) Handler(admin bool) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var repos []plateauapi.Repo
		pid := c.Param(pidParamName)
		token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")

		if pid == "" {
			metadata := plateaucms.GetAllCMSMetadataFromContext(ctx)
			plateauMetadata := plateaucms.PlateauProjectsFromMetadata(metadata)
			if len(plateauMetadata) == 0 {
				return echo.NewHTTPError(http.StatusNotFound, "not found")
			}

			if admin && (token == "" || !plateauMetadata[0].IsValidToken(token)) {
				log.Debugfc(ctx, "datacatalogv3: unauthorized access: input_token=%s, project=%#v", token, plateauMetadata[0].DataCatalogProjectAlias)
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			repos = h.getAllRepos(c.Request().Context(), admin, plateauMetadata)
		} else {
			md := plateaucms.GetCMSMetadataFromContext(ctx)
			if md.DataCatalogProjectAlias != pid || !isV3(md) {
				return echo.NewHTTPError(http.StatusNotFound, "not found")
			}

			if admin && !md.Auth {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			if err := h.prepare(ctx, md); err != nil {
				return echo.NewHTTPError(http.StatusBadGateway, "failed to prepare")
			}

			repo := h.getRepo(ctx, admin, md)
			if repo == nil {
				return echo.NewHTTPError(http.StatusNotFound, "not found")
			}

			repos = append(repos, repo)
		}

		merged := plateauapi.NewMerger(repos...)
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

	v2prj := lo.Keys(h.reposv2)
	sort.Strings(v2prj)
	for _, prj := range v2prj {
		prj := prj
		g.Go(func() error {
			return h.updateV2(ctx, prj)
		})
	}

	return g.Wait()
}

func (h *reposHandler) Init(ctx context.Context) error {
	metadata, err := h.pcms.AllMetadata(ctx, true)
	if err != nil {
		return fmt.Errorf("datacatalogv3: failed to get all metadata: %w", err)
	}

	if err := h.prepareAllPlateauProjects(ctx, metadata); err != nil {
		return err
	}

	return nil
}

func (h *reposHandler) getAllRepos(ctx context.Context, admin bool, metadata []plateaucms.Metadata) []plateauapi.Repo {
	targets := plateaucms.PlateauProjectsFromMetadata(metadata)

	repos := make([]plateauapi.Repo, 0, len(targets))
	for _, md := range targets {
		r := h.getRepo(ctx, admin, md)
		if r != nil {
			repos = append(repos, r)
			log.Infofc(ctx, "datacatalogv3: found repo for %s", md.DataCatalogProjectAlias)
		}
	}

	return repos
}

func (h *reposHandler) getRepo(ctx context.Context, admin bool, md plateaucms.Metadata) plateauapi.Repo {
	if isV2(md) {
		return h.reposv2[md.DataCatalogProjectAlias]
	} else if isV3(md) {
		return h.reposv3.Repo(md.DataCatalogProjectAlias, admin)
	}

	return nil
}

func (h *reposHandler) prepareAllPlateauProjects(ctx context.Context, metadata []plateaucms.Metadata) error {
	targets := plateaucms.PlateauProjectsFromMetadata(metadata)
	log.Infofc(ctx, "datacatalogv3: preparing repos for %d projects", len(targets))

	errg, ctx := errgroup.WithContext(ctx)
	for _, md := range targets {
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

	if _, ok := h.reposv2[md.DataCatalogProjectAlias]; ok {
		return nil
	}

	r, err := datacatalogv2adapter.New(md.CMSBaseURL, md.DataCatalogProjectAlias)
	if err != nil {
		return fmt.Errorf("datacatalogv3: failed to create repo v2 for %s: %w", md.DataCatalogProjectAlias, err)
	}

	h.reposv2[md.DataCatalogProjectAlias] = r
	return h.updateV2(ctx, md.DataCatalogProjectAlias)
}

func (h *reposHandler) prepareV3(ctx context.Context, md plateaucms.Metadata) error {
	if !isV3(md) {
		return nil
	}

	cms, err := md.CMS()
	if err != nil {
		return fmt.Errorf("datacatalogv3: failed to create cms for %s: %w", md.DataCatalogProjectAlias, err)
	}

	if err := h.reposv3.Prepare(ctx, md.DataCatalogProjectAlias, cms); err != nil {
		return fmt.Errorf("datacatalogv3: failed to prepare repo for %s: %w", md.DataCatalogProjectAlias, err)
	}

	return nil
}

func (h *reposHandler) updateV2(ctx context.Context, prj string) error {
	r := h.reposv2[prj]
	if r == nil {
		return nil
	}

	log.Infofc(ctx, "datacatalogv3: updating repo v2 for %s", prj)

	if err := r.Update(ctx); err != nil {
		return fmt.Errorf("datacatalogv3: failed to update repo v2 for %s: %w", prj, err)
	}

	log.Infofc(ctx, "datacatalogv3: updated repo v2 for %s", prj)

	return nil
}

func (h *reposHandler) updateV3(ctx context.Context, prj string) error {
	if err := h.reposv3.Update(ctx, prj); err != nil {
		return fmt.Errorf("datacatalogv3: failed to update repo for %s: %w", prj, err)
	}
	return nil
}

func isV2(md plateaucms.Metadata) bool {
	return md.DataCatalogSchemaVersion == cmsSchemaVersionV2
}

func isV3(md plateaucms.Metadata) bool {
	return md.DataCatalogSchemaVersion == cmsSchemaVersion
}
