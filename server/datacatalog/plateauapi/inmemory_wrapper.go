package plateauapi

import (
	"context"
	"fmt"
	"sync"

	"github.com/samber/lo"
)

type RepoWrappers []*RepoWrapper

func (a RepoWrappers) Update(ctx context.Context) error {
	results := make([]<-chan error, 0, len(a))
	for _, w := range a {
		w := w
		ch := lo.Async1(func() error {
			return w.Update(ctx)
		})
		results = append(results, ch)
	}

	for i, ch := range results {
		if err := <-ch; err != nil {
			return fmt.Errorf("repo %d: %w", i, err)
		}
	}

	return nil
}

type RepoUpdater func(ctx context.Context, repo *Repo) error

type RepoWrapper struct {
	repo    Repo
	lock    sync.RWMutex
	updater RepoUpdater
}

func NewRepoWrapper(updater RepoUpdater) *RepoWrapper {
	return &RepoWrapper{
		updater: updater,
	}
}

func (a *RepoWrapper) GetRepo() Repo {
	return a.repo
}

func (a *RepoWrapper) IsAvailable() bool {
	return a.repo != nil
}

func (a *RepoWrapper) Update(ctx context.Context) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	err := a.updater(ctx, &a.repo)
	if err != nil {
		return err
	}

	return nil
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