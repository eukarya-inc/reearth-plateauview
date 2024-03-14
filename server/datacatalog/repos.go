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
		pid := c.Param(pidParamName)
		token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
		metadata := plateaucms.GetAllCMSMetadataFromContext(ctx)

		if err := h.auth(ctx, admin, metadata, pid, token); err != nil {
			return err
		}

		merged := h.prepareAndGetMergedRepo(ctx, admin, pid, metadata)
		if merged == nil {
			return echo.NewHTTPError(http.StatusNotFound, "not found")
		}

		log.Debugfc(ctx, "datacatalogv3: use repo for %s: %s", pid, merged.Name())

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

	plateauMetadata := metadata.PlateauProjects()
	if err := h.prepareAll(ctx, plateauMetadata); err != nil {
		return err
	}

	return nil
}

func (*reposHandler) auth(ctx context.Context, admin bool, metadata plateaucms.MetadataList, pid, token string) error {
	if !admin {
		return nil
	}

	if pid == "" {
		plateauMetadata := metadata.PlateauProjects()
		if len(plateauMetadata) == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "not found")
		}

		if token == "" || !plateauMetadata[0].IsValidToken(token) {
			log.Debugfc(ctx, "datacatalogv3: unauthorized access: input_token=%s, project=%s", token, plateauMetadata[0].DataCatalogProjectAlias)
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
	} else {
		md := plateaucms.GetCMSMetadataFromContext(ctx)
		if md.DataCatalogProjectAlias != pid || !isV3(md) {
			return echo.NewHTTPError(http.StatusNotFound, "not found")
		}

		if !md.Auth {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}
	}

	return nil
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
		repo = h.reposv2[md.DataCatalogProjectAlias]
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

	if _, ok := h.reposv2[md.DataCatalogProjectAlias]; ok {
		return nil
	}

	r, err := datacatalogv2adapter.New(md.CMSBaseURL, md.DataCatalogProjectAlias)
	if err != nil {
		return fmt.Errorf("datacatalogv3: failed to create repo(v2) %s: %w", md.DataCatalogProjectAlias, err)
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

	if err := h.reposv3.Prepare(ctx, md.DataCatalogProjectAlias, md.PlateauYear(), cms); err != nil {
		return fmt.Errorf("datacatalogv3: failed to prepare repo for %s: %w", md.DataCatalogProjectAlias, err)
	}

	return nil
}

func (h *reposHandler) updateV2(ctx context.Context, prj string) error {
	r := h.reposv2[prj]
	if r == nil {
		return nil
	}

	log.Infofc(ctx, "datacatalogv3: updating repo(v2) %s", prj)

	if updated, err := r.Update(ctx); err != nil {
		return fmt.Errorf("datacatalogv3: failed to update repo(v2) %s: %w", prj, err)
	} else if !updated {
		log.Infofc(ctx, "datacatalogv3: skip updating repo(v2) %s", prj)
		return nil
	}

	log.Infofc(ctx, "datacatalogv3: updated repo(v2) %s", prj)
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
