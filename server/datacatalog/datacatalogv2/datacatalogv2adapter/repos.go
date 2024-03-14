package datacatalogv2adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/util"
)

type Repos struct {
	fetchers *util.SyncMap[string, datacatalogv2.Fetchable]
	*plateauapi.Repos
}

func NewRepos() *Repos {
	r := &Repos{
		fetchers: util.NewSyncMap[string, datacatalogv2.Fetchable](),
	}
	r.Repos = plateauapi.NewRepos(r.update)
	return r
}

func (r *Repos) Prepare(ctx context.Context, project string, f datacatalogv2.Fetchable) error {
	if _, ok := r.fetchers.Load(project); ok {
		return nil
	}

	r.setCMS(project, f)
	_, err := r.Update(ctx, project)
	return err
}

func (r *Repos) update(ctx context.Context, project string) (*plateauapi.ReposUpdateResult, error) {
	fetcher, ok := r.fetchers.Load(project)
	if !ok {
		return nil, fmt.Errorf("fetcher is not initialized for %s", project)
	}

	updated := r.UpdatedAt(project)
	var updatedStr string
	if !updated.IsZero() {
		updatedStr = updated.Format(time.RFC3339)
	}
	log.Debugfc(ctx, "datacatalogv2: updating repo %s: last_update=%s", project, updatedStr)

	repo, err := fetchAndCreateCache(ctx, project, fetcher, datacatalogv2.FetcherDoOptions{})
	if err != nil {
		return nil, err
	}

	log.Debugfc(ctx, "datacatalogv2: updated repo %s", project)

	return &plateauapi.ReposUpdateResult{
		Repo: repo,
	}, nil
}

func (r *Repos) setCMS(project string, f datacatalogv2.Fetchable) {
	r.fetchers.Store(project, f)
}
