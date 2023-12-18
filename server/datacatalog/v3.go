package datacatalog

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv3"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/eukarya-inc/reearth-plateauview/server/plateaucms"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func echov3(conf Config, g *echo.Group, repov2 *plateauapi.RepoWrapper) error {
	if conf.GraphqlMaxComplexity <= 0 {
		conf.GraphqlMaxComplexity = 1000
	}

	h, err := newReposHandler(conf, repov2)
	if err != nil {
		return err
	}

	// PLATEAU API
	plateauapig := g.Group(":alias")
	plateauapig.Use(
		middleware.CORS(),
		middleware.Gzip(),
		h.Middleware(),
	)

	plateauapig.GET("/graphql", func(c echo.Context) error {
		project, admin := getProjectAndAdmin(c.Path())
		pa := getGqlPathFromProjectAndAdmin(project, admin)

		p := plateauapi.PlaygroundHandler(
			"PLATEAU GraphQL API Playground",
			path.Join(conf.PlaygroundEndpoint, pa),
		)
		p.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	plateauapig.POST("/graphql", func(c echo.Context) error {
		project, admin := getProjectAndAdmin(c.Path())

		repo, err := h.GerRepo(c.Request().Context(), project, admin)
		if err != nil {
			return err
		}
		if repo == nil {
			return echo.NewHTTPError(http.StatusNotFound, "not found")
		}

		srv := plateauapi.NewService(repo, plateauapi.FixedComplexityLimit(conf.GraphqlMaxComplexity))
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	return nil
}

type reposHandler struct {
	reposv3 *datacatalogv3.Repos
	repov2  *plateauapi.RepoWrapper
	pcms    *plateaucms.CMS
}

func newReposHandler(conf Config, repov2 *plateauapi.RepoWrapper) (*reposHandler, error) {
	pcms, err := plateaucms.New(conf.Config)
	if err != nil {
		return nil, err
	}

	reposv3 := datacatalogv3.NewRepos(conf.Config.CMSBaseURL)

	return &reposHandler{
		reposv3: reposv3,
		repov2:  repov2,
		pcms:    pcms,
	}, nil
}

func (h *reposHandler) GerRepo(ctx context.Context, project string, admin bool) (plateauapi.Repo, error) {
	cmsmd, err := h.pcms.Metadata(ctx, project)
	if err != nil {
		return nil, err
	}
	if cmsmd.ProjectAlias == "" || cmsmd.CMSAPIKey == "" || admin && !cmsmd.Auth {
		return nil, nil
	}

	projectalias := cmsmd.ProjectAlias
	if err := h.reposv3.Prepare(ctx, cmsmd.CMSAPIKey, projectalias); err != nil {
		return nil, err
	}

	return plateauapi.NewMerger(h.reposv3.Repo(projectalias, admin), h.repov2), nil
}

func (h *reposHandler) Middleware() echo.MiddlewareFunc {
	return h.pcms.AuthMiddleware(true)
}

func getProjectAndAdmin(path string) (project string, admin bool) {
	parts := strings.Split(path, "/")
	if parts[len(parts)-1] != "graphql" || len(parts) < 2 {
		return "", false
	}

	if parts[len(parts)-2] == "admin" {
		admin = true
		project = parts[len(parts)-3]
	} else {
		project = parts[len(parts)-2]
	}

	return
}

func getGqlPathFromProjectAndAdmin(project string, admin bool) string {
	if admin {
		return fmt.Sprintf("/%s/admin/graphql", project)
	}
	return fmt.Sprintf("/%s/graphql", project)
}
