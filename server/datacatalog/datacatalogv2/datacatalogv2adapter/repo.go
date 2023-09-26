package datacatalogv2adapter

import (
	"context"
	"fmt"
	"sync"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

type Adapter struct {
	fetcher *datacatalogv2.Fetcher
	project string

	// cache
	lock          sync.Mutex
	updatingCache bool
	cache         []datacatalogv2.DataCatalogItem
	areas         []area
	specs         []spec
}

type area struct {
	Name string
	Code string
}

func (a area) Into() plateauapi.Area {
	if a.Name == "" {
		return nil
	}
	if len(a.Code) == 2 {
		return plateauapi.Prefecture{
			Code: plateauapi.AreaCode(a.Code),
			Name: a.Name,
		}
	}
	return plateauapi.Municipality{
		Code: plateauapi.AreaCode(a.Code),
		Name: a.Name,
	}
}

type spec struct {
	Name string
	Year int
}

func New(cmsbase, project string) (*Adapter, error) {
	f, err := datacatalogv2.NewFetcher(cmsbase)
	return &Adapter{
		fetcher: f,
		project: project,
	}, err
}

func (a *Adapter) UpdateCache(ctx context.Context, opts datacatalogv2.FetcherDoOptions) error {
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

	r, err := a.fetcher.Do(ctx, a.project, opts)
	if err != nil {
		return fmt.Errorf("failed to update datacatalog cache: %w", err)
	}

	items := r.All()
	a.cache = items

	for _, d := range items {
		if d.Spec != "" {
			if d.Pref != "" {
				a.areas = append(a.areas, area{
					Code: d.PrefCode,
					Name: d.Pref,
				})
			}
			if d.City != "" {
				a.areas = append(a.areas, area{
					Code: d.CityCode,
					Name: d.City,
				})
			}
			if d.Ward != "" {
				a.areas = append(a.areas, area{
					Code: d.WardCode,
					Name: d.Ward,
				})
			}
			a.specs = append(a.specs, spec{
				Name: d.Spec,
				Year: d.Year,
			})
		}
	}

	slices.SortStableFunc(a.areas, func(a, b area) bool {
		return a.Code < b.Code
	})
	slices.SortStableFunc(a.specs, func(a, b spec) bool {
		return a.Year < b.Year || a.Name < b.Name
	})
	lo.Uniq(a.areas)
	lo.Uniq(a.specs)
	return nil
}

var _ plateauapi.Repo = (*Adapter)(nil)

func (a *Adapter) Node(ctx context.Context, id plateauapi.ID) (plateauapi.Node, error) {
	// id, ty := id.Unwrap()
	// switch ty {
	// case "usecase":
	// }
	panic("implement me")
}

func (a *Adapter) Nodes(ctx context.Context, ids []plateauapi.ID) ([]plateauapi.Node, error) {
	return util.TryMap(ids, func(id plateauapi.ID) (plateauapi.Node, error) {
		return a.Node(ctx, id)
	})
}

func (a *Adapter) Area(ctx context.Context, code plateauapi.AreaCode) (plateauapi.Area, error) {
	area, _ := lo.Find(a.areas, func(a area) bool {
		return a.Code == string(code)
	})
	return area.Into(), nil
}

func (a *Adapter) Areas(ctx context.Context, input plateauapi.AreaQuery) ([]plateauapi.Area, error) {
	return lo.Map(a.areas, func(a area, _ int) plateauapi.Area {
		return a.Into()
	}), nil
}

func (a *Adapter) DatasetTypes(ctx context.Context, input plateauapi.DatasetTypeQuery) ([]plateauapi.DatasetType, error) {
	panic("implement me")
}

func (a *Adapter) Datasets(ctx context.Context, input plateauapi.DatasetQuery) ([]plateauapi.Dataset, error) {
	panic("implement me")
}

func (a *Adapter) PlateauSpecs(ctx context.Context) ([]*plateauapi.PlateauSpec, error) {
	specs := lo.Map(a.specs, func(s spec, _ int) *plateauapi.PlateauSpec {
		return &plateauapi.PlateauSpec{
			ID:   plateauapi.NewID(s.Name, "spec"),
			Name: s.Name,
			Year: s.Year,
		}
	})
	return specs, nil
}

func (a *Adapter) Years(ctx context.Context) ([]int, error) {
	return lo.Map(a.specs, func(s spec, _ int) int {
		return s.Year
	}), nil
}
