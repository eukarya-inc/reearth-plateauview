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

	return r.Update(ctx, project, cms)
}

func (r *Repos) Repo(project string, admin bool) *plateauapi.RepoWrapper {
	if admin {
		return r.adminRepos[project]
	}
	return r.repos[project]
}

func (r *Repos) UpdateAll(ctx context.Context) error {
	projects := lo.Keys(r.repos)
	sort.Strings(projects)

	for _, project := range projects {
		if err := r.Update(ctx, project, nil); err != nil {
			return fmt.Errorf("failed to update project %s: %w", project, err)
		}
	}
	return nil
}

func (r *Repos) Update(ctx context.Context, project string, rawcms cms.Interface) error {
	r.locks.Lock(project)
	defer r.locks.Unlock(project)

	cms := r.cms[project]
	if cms == nil {
		if rawcms == nil {
			return nil
		}
		cms = NewCMS(rawcms)
	}

	log.Infofc(ctx, "datacatalogv3: updating project %s", project)
	data, err := cms.GetAll(ctx, project)
	if err != nil {
		return err
	}

	c, warning := data.Into()
	r.warnings[project] = warning
	r.context[project] = c

	repo := plateauapi.NewInMemoryRepo(c)
	adminRepo := plateauapi.NewInMemoryRepo(c)
	adminRepo.SetIncludedStages(stagesForAdmin...)

	adminRepoWrapper := r.adminRepos[project]
	if adminRepoWrapper == nil {
		adminRepoWrapper = plateauapi.NewRepoWrapper(adminRepo, nil)
		r.adminRepos[project] = adminRepoWrapper
	} else {
		adminRepoWrapper.SetRepo(adminRepo)
	}

	repoWrapper := r.repos[project]
	if repoWrapper == nil {
		repoWrapper = plateauapi.NewRepoWrapper(repo, nil)
		r.repos[project] = repoWrapper
	} else {
		repoWrapper.SetRepo(repo)
	}

	if r.now != nil {
		r.updatedAt[project] = r.now()
	} else {
		r.updatedAt[project] = time.Now()
	}

	log.Infofc(ctx, "datacatalogv3: updated project %s", project)
	return nil
}

func (r *Repos) Warnings(project string) []string {
	return slices.Clone(r.warnings[project])
}

func (r *Repos) UpdatedAt(project string) time.Time {
	return r.updatedAt[project]
}