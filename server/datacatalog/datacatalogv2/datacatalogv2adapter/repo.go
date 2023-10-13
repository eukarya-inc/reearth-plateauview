package datacatalogv2adapter

import (
	"context"
	"sync"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
)

type Adapter struct {
	fetcher       datacatalogv2.Fetchable
	project       string
	lock          sync.Mutex
	updatingCache bool
	cache         *cache
}

func New(cmsbase, project string) (*Adapter, error) {
	f, err := datacatalogv2.NewFetcher(cmsbase)
	return &Adapter{
		fetcher: f,
		project: project,
	}, err
}

func From(proejct string, fetcher datacatalogv2.Fetchable) *Adapter {
	return &Adapter{
		fetcher: fetcher,
		project: proejct,
	}
}

func (a *Adapter) UpdateCache(ctx context.Context) error {
	updating := a.updatingCache
	a.lock.Lock()
	defer a.lock.Unlock()
	if updating {
		return nil
	}

	a.updatingCache = true
	defer func() {
		a.updatingCache = false
	}()

	c, err := fetchAndCreateCache(ctx, a.project, a.fetcher, datacatalogv2.FetcherDoOptions{})
	if err != nil {
		return err
	}

	a.cache = c
	return nil
}

func (a *Adapter) IsAvailable() bool {
	return a.cache != nil
}

func (a *Adapter) getCache() *cache {
	return a.cache
}

var _ plateauapi.Repo = (*Adapter)(nil)

func (a *Adapter) Node(ctx context.Context, id plateauapi.ID) (plateauapi.Node, error) {
	if !a.IsAvailable() {
		return nil, plateauapi.ErrDatacatalogUnavailable
	}
	return a.getCache().Node(ctx, id)
}

func (a *Adapter) Nodes(ctx context.Context, ids []plateauapi.ID) ([]plateauapi.Node, error) {
	if !a.IsAvailable() {
		return nil, plateauapi.ErrDatacatalogUnavailable
	}
	return a.getCache().Nodes(ctx, ids)
}

func (a *Adapter) Area(ctx context.Context, code plateauapi.AreaCode) (plateauapi.Area, error) {
	if !a.IsAvailable() {
		return nil, plateauapi.ErrDatacatalogUnavailable
	}
	return a.getCache().Area(ctx, code)
}

func (a *Adapter) Areas(ctx context.Context, input *plateauapi.AreaInput) (res []plateauapi.Area, _ error) {
	if !a.IsAvailable() {
		return nil, plateauapi.ErrDatacatalogUnavailable
	}
	return a.getCache().Areas(ctx, input)
}

func (a *Adapter) DatasetTypes(ctx context.Context, input *plateauapi.DatasetTypeInput) (res []plateauapi.DatasetType, _ error) {
	if !a.IsAvailable() {
		return nil, plateauapi.ErrDatacatalogUnavailable
	}
	return a.getCache().DatasetTypes(ctx, input)
}

func (a *Adapter) Datasets(ctx context.Context, input *plateauapi.DatasetInput) (res []plateauapi.Dataset, _ error) {
	if !a.IsAvailable() {
		return nil, plateauapi.ErrDatacatalogUnavailable
	}
	return a.getCache().Datasets(ctx, input)
}

func (a *Adapter) PlateauSpecs(ctx context.Context) ([]*plateauapi.PlateauSpec, error) {
	if !a.IsAvailable() {
		return nil, plateauapi.ErrDatacatalogUnavailable
	}
	return a.getCache().PlateauSpecs(ctx)
}

func (a *Adapter) Years(ctx context.Context) ([]int, error) {
	if !a.IsAvailable() {
		return nil, plateauapi.ErrDatacatalogUnavailable
	}
	return a.getCache().Years(ctx)
}
