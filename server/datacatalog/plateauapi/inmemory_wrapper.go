package plateauapi

import (
	"context"
	"sync"
)

type RepoUpdater func(ctx context.Context, repo *Repo) error

// RepoWrapper is a thread-safe wrapper of Repo.
type RepoWrapper struct {
	repo    Repo
	lock    sync.RWMutex
	updater RepoUpdater
}

func NewRepoWrapper(repo Repo, updater RepoUpdater) *RepoWrapper {
	return &RepoWrapper{
		repo:    repo,
		updater: updater,
	}
}

func (a *RepoWrapper) GetRepo() Repo {
	return a.repo
}

func (a *RepoWrapper) IsAvailable() bool {
	return a.repo != nil
}

func (a *RepoWrapper) SetRepo(repo Repo) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.repo = repo
}

func (a *RepoWrapper) Update(ctx context.Context) error {
	if a.updater == nil {
		return nil
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	return a.updater(ctx, &a.repo)
}

func (a *RepoWrapper) use(f func(r Repo) error) error {
	if !a.IsAvailable() {
		return ErrDatacatalogUnavailable
	}

	a.lock.RLock()
	defer a.lock.RUnlock()
	return f(a.GetRepo())
}

var _ Repo = (*RepoWrapper)(nil)

func (a *RepoWrapper) Node(ctx context.Context, id ID) (res Node, err error) {
	err = a.use(func(r Repo) (err error) {
		res, err = r.Node(ctx, id)
		return
	})
	return
}

func (a *RepoWrapper) Nodes(ctx context.Context, ids []ID) (res []Node, err error) {
	err = a.use(func(r Repo) (err error) {
		res, err = r.Nodes(ctx, ids)
		return
	})
	return
}

func (a *RepoWrapper) Area(ctx context.Context, code AreaCode) (res Area, err error) {
	err = a.use(func(r Repo) (err error) {
		res, err = r.Area(ctx, code)
		return
	})
	return
}

func (a *RepoWrapper) Areas(ctx context.Context, input *AreasInput) (res []Area, err error) {
	err = a.use(func(r Repo) (err error) {
		res, err = r.Areas(ctx, input)
		return
	})
	return
}

func (a *RepoWrapper) DatasetTypes(ctx context.Context, input *DatasetTypesInput) (res []DatasetType, err error) {
	err = a.use(func(r Repo) (err error) {
		res, err = r.DatasetTypes(ctx, input)
		return
	})
	return
}

func (a *RepoWrapper) Datasets(ctx context.Context, input *DatasetsInput) (res []Dataset, err error) {
	err = a.use(func(r Repo) (err error) {
		res, err = r.Datasets(ctx, input)
		return
	})
	return
}

func (a *RepoWrapper) PlateauSpecs(ctx context.Context) (res []*PlateauSpec, err error) {
	err = a.use(func(r Repo) (err error) {
		res, err = r.PlateauSpecs(ctx)
		return
	})
	return
}

func (a *RepoWrapper) Years(ctx context.Context) (res []int, err error) {
	err = a.use(func(r Repo) (err error) {
		res, err = r.Years(ctx)
		return
	})
	return
}