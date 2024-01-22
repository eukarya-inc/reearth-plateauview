package datacatalogv3

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/log"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

const cacheUpdateDuration = 10 * time.Second

var stagesForAdmin = []string{string(stageBeta)}

type Repos struct {
	locks      util.LockMap[string]
	cms        map[string]*CMS
	context    map[string]*plateauapi.InMemoryRepoContext
	repos      map[string]*plateauapi.RepoWrapper
	adminRepos map[string]*plateauapi.RepoWrapper
	warnings   map[string][]string
	updatedAt  map[string]time.Time
	now        func() time.Time
}

func NewRepos() *Repos {
	return &Repos{
		locks:      util.LockMap[string]{},
		cms:        map[string]*CMS{},
		context:    map[string]*plateauapi.InMemoryRepoContext{},
		repos:      map[string]*plateauapi.RepoWrapper{},
		adminRepos: map[string]*plateauapi.RepoWrapper{},
		warnings:   map[string][]string{},
		updatedAt:  map[string]time.Time{},
	}
}

func (r *Repos) Prepare(ctx context.Context, project string, cms cms.Interface) error {
	if r.cms[project] != nil {
		return nil
	}

	r.setCMS(project, cms)
	return r.Update(ctx, project)
}

func (r *Repos) Repo(project string, admin bool) *plateauapi.RepoWrapper {
	if admin {
		return r.adminRepos[project]
	}
	return r.repos[project]
}

func (r *Repos) Projects() []string {
	keys := lo.Keys(r.repos)
	sort.Strings(keys)
	return keys
}

func (r *Repos) UpdateAll(ctx context.Context) error {
	projects := r.Projects()
	for _, project := range projects {
		if err := r.Update(ctx, project); err != nil {
			return fmt.Errorf("failed to update project %s: %w", project, err)
		}
	}
	return nil
}

func (r *Repos) Update(ctx context.Context, project string) error {
	r.locks.Lock(project)
	defer r.locks.Unlock(project)

	updated := r.UpdatedAt(project)
	var updatedStr string
	if !updated.IsZero() {
		updatedStr = updated.Format(time.RFC3339)
	}

	// avoid too frequent updates
	since := r.getNow().Sub(updated)
	if !updated.IsZero() && since < cacheUpdateDuration {
		log.Infofc(ctx, "datacatalogv3: skip updating repo %s: last_update=%s, since=%s", project, updatedStr, since)
		return nil
	}

	cms := r.cms[project]
	if cms == nil {
		return fmt.Errorf("cms is not initialized for %s", project)
	}

	log.Infofc(ctx, "datacatalogv3: updating repo %s: last_update=%s", project, updatedStr)
	data, err := cms.GetAll(ctx, project)
	if err != nil {
		return err
	}

	c, warning := data.Into()
	sort.Strings(warning)
	r.warnings[project] = warning
	r.context[project] = c

	repo := plateauapi.NewInMemoryRepo(c)
	adminRepo := plateauapi.NewInMemoryRepo(c)
	adminRepo.SetAdmin(true)
	adminRepo.SetIncludedStages(stagesForAdmin...)

	adminRepoWrapper := r.adminRepos[project]
	if adminRepoWrapper == nil {
		adminRepoWrapper = plateauapi.NewRepoWrapper(adminRepo, nil)
		adminRepoWrapper.SetName(fmt.Sprintf("%s(admin)", project))
		r.adminRepos[project] = adminRepoWrapper
	} else {
		adminRepoWrapper.SetRepo(adminRepo)
	}

	repoWrapper := r.repos[project]
	if repoWrapper == nil {
		repoWrapper = plateauapi.NewRepoWrapper(repo, nil)
		repoWrapper.SetName(project)
		r.repos[project] = repoWrapper
	} else {
		repoWrapper.SetRepo(repo)
	}

	r.updatedAt[project] = r.getNow()

	log.Infofc(ctx, "datacatalogv3: updated repo %s", project)
	return nil
}

func (r *Repos) Warnings(project string) []string {
	if r.UpdatedAt(project).IsZero() {
		return []string{"project is not initialized"}
	}
	return slices.Clone(r.warnings[project])
}

func (r *Repos) UpdatedAt(project string) time.Time {
	return r.updatedAt[project]
}

func (r *Repos) setCMS(project string, cms cms.Interface) {
	r.locks.Lock(project)
	defer r.locks.Unlock(project)

	c := NewCMS(cms)
	r.cms[project] = c
}

func (r *Repos) getNow() time.Time {
	if r.now != nil {
		return r.now()
	}
	return time.Now()
}
