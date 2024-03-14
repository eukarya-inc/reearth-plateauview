package datacatalogv3

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/util"
)

var stagesForAdmin = []string{string(stageBeta)}

type Repos struct {
	cms *util.SyncMap[string, *CMS]
	*plateauapi.Repos
}

func NewRepos() *Repos {
	r := &Repos{
		cms: util.NewSyncMap[string, *CMS](),
	}
	r.Repos = plateauapi.NewRepos(r.update)
	return r
}

func (r *Repos) Prepare(ctx context.Context, project string, year int, cms cms.Interface) error {
	if _, ok := r.cms.Load(project); ok {
		return nil
	}

	r.setCMS(project, year, cms)
	_, err := r.Update(ctx, project)
	return err
}

func (r *Repos) update(ctx context.Context, project string) (*plateauapi.ReposUpdateResult, error) {
	cms, ok := r.cms.Load(project)
	if !ok {
		return nil, fmt.Errorf("cms is not initialized for %s", project)
	}

	updated := r.UpdatedAt(project)
	var updatedStr string
	if !updated.IsZero() {
		updatedStr = updated.Format(time.RFC3339)
	}
	log.Debugfc(ctx, "datacatalogv3: updating repo %s: last_update=%s", project, updatedStr)

	data, err := cms.GetAll(ctx, project)
	if err != nil {
		return nil, err
	}

	c, warning := data.Into()
	sort.Strings(warning)

	repo := plateauapi.NewInMemoryRepo(c)

	adminRepo := plateauapi.NewInMemoryRepo(c)
	adminRepo.SetAdmin(true)
	adminRepo.SetIncludedStages(stagesForAdmin...)

	log.Debugfc(ctx, "datacatalogv3: updated repo %s", project)

	return &plateauapi.ReposUpdateResult{
		Repo:      repo,
		AdminRepo: adminRepo,
		Warnings:  warning,
	}, nil
}

func (r *Repos) setCMS(project string, year int, cms cms.Interface) {
	c := NewCMS(cms, year)
	r.cms.Store(project, c)
}
