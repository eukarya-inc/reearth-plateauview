package datacatalogv3

import (
	"context"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
)

func New(cmsbase, token, project string) (*plateauapi.RepoWrapper, error) {
	cms, err := cms.New(cmsbase, token)
	if err != nil {
		return nil, err
	}
	return From(NewCMS(cms), project), nil
}

func From(cms *CMS, project string) *plateauapi.RepoWrapper {
	return plateauapi.NewRepoWrapper(func(ctx context.Context) (plateauapi.Repo, error) {
		res, err := cms.GetAll(ctx, project)
		if err != nil {
			return nil, err
		}

		c, warning := res.Into()
		if len(warning) > 0 {
			log.Warnfc(ctx, "datacatalogv3: warning: \n%s", strings.Join(warning, "\n"))
		}

		return plateauapi.NewInMemoryRepo(c), nil
	})
}
