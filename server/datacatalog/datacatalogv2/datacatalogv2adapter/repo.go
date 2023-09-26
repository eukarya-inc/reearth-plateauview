package datacatalogv2adapter

import (
	"context"
	"strings"
	"sync"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

type Adapter struct {
	fetcher *datacatalogv2.Fetcher
	project string

	// cache
	lock                    sync.Mutex
	updatingCache           bool
	cache                   []datacatalogv2.DataCatalogItem
	prefectures             []plateauapi.Prefecture
	municipalities          []plateauapi.Municipality
	plateauDatasetTypes     []plateauapi.PlateauDatasetType
	relatedDatasetTypes     []plateauapi.RelatedDatasetType
	genericDatasetTypes     []plateauapi.GenericDatasetType
	plateauDatasets         []plateauapi.PlateauDataset
	plateauFloodingDatasets []plateauapi.PlateauFloodingDataset
	relatedDatasets         []plateauapi.RelatedDataset
	genericDatasets         []plateauapi.GenericDataset
	specs                   []plateauapi.PlateauSpec
}

func New(cmsbase, project string) (*Adapter, error) {
	f, err := datacatalogv2.NewFetcher(cmsbase)
	return &Adapter{
		fetcher: f,
		project: project,
	}, err
}

var _ plateauapi.Repo = (*Adapter)(nil)

func (a *Adapter) Node(ctx context.Context, id plateauapi.ID) (plateauapi.Node, error) {
	i, ty := id.Unwrap()
	switch ty {
	case plateauapi.TypePrefecture:
		if p, ok := lo.Find(a.prefectures, func(p plateauapi.Prefecture) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}
	case plateauapi.TypeMunicipality:
		if p, ok := lo.Find(a.municipalities, func(p plateauapi.Municipality) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}
	case plateauapi.TypeDatasetType:
		if p, ok := lo.Find(a.plateauDatasetTypes, func(p plateauapi.PlateauDatasetType) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}

		if p, ok := lo.Find(a.relatedDatasetTypes, func(p plateauapi.RelatedDatasetType) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}

		if p, ok := lo.Find(a.genericDatasetTypes, func(p plateauapi.GenericDatasetType) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}
	case plateauapi.TypeDataset:
		if p, ok := lo.Find(a.plateauDatasets, func(p plateauapi.PlateauDataset) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}

		if p, ok := lo.Find(a.plateauFloodingDatasets, func(p plateauapi.PlateauFloodingDataset) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}

		if p, ok := lo.Find(a.relatedDatasets, func(p plateauapi.RelatedDataset) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}

		if p, ok := lo.Find(a.genericDatasets, func(p plateauapi.GenericDataset) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}
	case plateauapi.TypeDatasetItem:
		parent, _, _ := strings.Cut(i, ":")
		parentID := newDatasetID(parent)

		if p, ok := lo.Find(a.plateauDatasets, func(p plateauapi.PlateauDataset) bool {
			return p.ID == parentID
		}); ok {
			item, _ := lo.Find(p.Data, func(i *plateauapi.PlateauDatasetItem) bool {
				return i.ID == id
			})
			return item, nil
		}

		if p, ok := lo.Find(a.plateauFloodingDatasets, func(p plateauapi.PlateauFloodingDataset) bool {
			return p.ID == parentID
		}); ok {
			item, _ := lo.Find(p.Data, func(i *plateauapi.PlateauFloodingDatasetItem) bool {
				return i.ID == id
			})
			return item, nil
		}

		if p, ok := lo.Find(a.relatedDatasets, func(p plateauapi.RelatedDataset) bool {
			return p.ID == parentID
		}); ok {
			item, _ := lo.Find(p.Data, func(i *plateauapi.RelatedDatasetItem) bool {
				return i.ID == id
			})
			return item, nil
		}

		if p, ok := lo.Find(a.genericDatasets, func(p plateauapi.GenericDataset) bool {
			return p.ID == id
		}); ok {
			item, _ := lo.Find(p.Data, func(i *plateauapi.GenericDatasetItem) bool {
				return i.ID == id
			})
			return item, nil
		}
	case plateauapi.TypePlateauSpec:
		if p, ok := lo.Find(a.specs, func(p plateauapi.PlateauSpec) bool {
			return p.ID == id
		}); ok {
			return p, nil
		}
	}

	return nil, nil
}

func (a *Adapter) Nodes(ctx context.Context, ids []plateauapi.ID) ([]plateauapi.Node, error) {
	return util.TryMap(ids, func(id plateauapi.ID) (plateauapi.Node, error) {
		return a.Node(ctx, id)
	})
}

func (a *Adapter) Area(ctx context.Context, code plateauapi.AreaCode) (plateauapi.Area, error) {
	if code.IsPrefectureCode() {
		area, _ := lo.Find(a.prefectures, func(p plateauapi.Prefecture) bool {
			return p.Code == code
		})
		return area, nil
	}

	area, _ := lo.Find(a.municipalities, func(p plateauapi.Municipality) bool {
		return p.Code == code
	})
	return area, nil
}

func (a *Adapter) Areas(ctx context.Context, input plateauapi.AreaQuery) (res []plateauapi.Area, _ error) {
	prefs := lo.Filter(a.prefectures, func(t plateauapi.Prefecture, _ int) bool {
		return filterArea(t, input)
	})
	municipalities := lo.Filter(a.municipalities, func(t plateauapi.Municipality, _ int) bool {
		return filterArea(t, input)
	})

	for _, t := range prefs {
		res = append(res, t)
	}
	for _, t := range municipalities {
		res = append(res, t)
	}
	return
}

func (a *Adapter) DatasetTypes(ctx context.Context, input plateauapi.DatasetTypeQuery) (res []plateauapi.DatasetType, _ error) {
	plateau := lo.Filter(a.plateauDatasetTypes, func(t plateauapi.PlateauDatasetType, _ int) bool {
		return filterDataType(t, input)
	})
	related := lo.Filter(a.relatedDatasetTypes, func(t plateauapi.RelatedDatasetType, _ int) bool {
		return filterDataType(t, input)
	})
	generic := lo.Filter(a.genericDatasetTypes, func(t plateauapi.GenericDatasetType, _ int) bool {
		return filterDataType(t, input)
	})

	for _, t := range plateau {
		res = append(res, t)
	}
	for _, t := range related {
		res = append(res, t)
	}
	for _, t := range generic {
		res = append(res, t)
	}
	return
}

func (a *Adapter) Datasets(ctx context.Context, input plateauapi.DatasetQuery) (res []plateauapi.Dataset, _ error) {
	plateau := lo.Filter(a.plateauDatasets, func(t plateauapi.PlateauDataset, _ int) bool {
		return filterDataset(t, input)
	})
	flooding := lo.Filter(a.plateauFloodingDatasets, func(t plateauapi.PlateauFloodingDataset, _ int) bool {
		return filterDataset(t, input)
	})
	related := lo.Filter(a.relatedDatasets, func(t plateauapi.RelatedDataset, _ int) bool {
		return filterDataset(t, input)
	})
	generic := lo.Filter(a.genericDatasets, func(t plateauapi.GenericDataset, _ int) bool {
		return filterDataset(t, input)
	})

	for _, t := range plateau {
		res = append(res, t)
	}
	for _, t := range flooding {
		res = append(res, t)
	}
	for _, t := range related {
		res = append(res, t)
	}
	for _, t := range generic {
		res = append(res, t)
	}
	return
}

func (a *Adapter) PlateauSpecs(ctx context.Context) ([]*plateauapi.PlateauSpec, error) {
	return lo.ToSlicePtr(a.specs), nil
}

func (a *Adapter) Years(ctx context.Context) ([]int, error) {
	return lo.Map(a.specs, func(s plateauapi.PlateauSpec, _ int) int {
		return s.Year
	}), nil
}
