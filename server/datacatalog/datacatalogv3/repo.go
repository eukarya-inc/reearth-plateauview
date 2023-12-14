package datacatalogv3

import (
	"context"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	cms "github.com/reearth/reearth-cms-api/go"
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
		return plateauapi.NewInMemoryRepo(res.Into()), nil
	})
}
