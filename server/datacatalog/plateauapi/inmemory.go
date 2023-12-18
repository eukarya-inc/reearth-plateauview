package plateauapi

import (
	"context"
	"slices"

	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
)

type InMemoryRepoContext struct {
	Areas        Areas
	DatasetTypes DatasetTypes
	Datasets     Datasets
	PlateauSpecs []PlateauSpec
	Years        []int
}

// InMemoryRepo is a repository that stores all data in memory.
// Note that it is not thread-safe.
type InMemoryRepo struct {
	ctx               *InMemoryRepoContext
	areasForDataTypes map[string]map[AreaCode]struct{}
	includedStages    []string
}

var _ Repo = (*InMemoryRepo)(nil)

func NewInMemoryRepo(ctx *InMemoryRepoContext) *InMemoryRepo {
	r := &InMemoryRepo{}
	r.SetContext(ctx)
	return r
}

func (c *InMemoryRepo) SetContext(ctx *InMemoryRepoContext) {
	c.ctx = ctx
	c.areasForDataTypes = areasForDatasetTypes(ctx.Datasets.All())
}

func (c *InMemoryRepo) SetIncludeAllStage(stages ...string) {
	c.includedStages = stages
}

func (c *InMemoryRepo) IncludeAllStage() []string {
	return slices.Clone(c.includedStages)
}

func (c *InMemoryRepo) Node(ctx context.Context, id ID) (Node, error) {
	ty := id.Type()
	switch ty {
	case TypeArea:
		if a := c.ctx.Areas.Area(id); a != nil {
			return a, nil
		}
	case TypeDatasetType:
		if dt := c.ctx.DatasetTypes.DatasetType(id); dt != nil {
			return dt, nil
		}
	case TypeDataset:
		if d := c.ctx.Datasets.Dataset(id); d != nil {
			return d, nil
		}
	case TypeDatasetItem:
		if i := c.ctx.Datasets.Item(id); i != nil {
			return i, nil
		}
	case TypePlateauSpec:
		if p, ok := lo.Find(c.ctx.PlateauSpecs, func(p PlateauSpec) bool {
			return p.ID == id || lo.SomeBy(p.MinorVersions, func(v *PlateauSpecMinor) bool {
				return v.ID == id
			})
		}); ok {
			if p.ID != id {
				m, _ := lo.Find(p.MinorVersions, func(v *PlateauSpecMinor) bool {
					return v.ID == id
				})
				return m, nil
			}
			return &p, nil
		}
	}

	return nil, nil
}

func (c *InMemoryRepo) Nodes(ctx context.Context, ids []ID) ([]Node, error) {
	return util.TryMap(ids, func(id ID) (Node, error) {
		return c.Node(ctx, id)
	})
}

func (c *InMemoryRepo) Area(ctx context.Context, code AreaCode) (Area, error) {
	return c.ctx.Areas.Find(func(a Area) bool {
		return a.GetCode() == code
	}), nil
}

func (c *InMemoryRepo) Areas(ctx context.Context, input *AreasInput) (res []Area, _ error) {
	inp := lo.FromPtr(input)
	types := c.getDatasetTypeCodes(inp.DatasetTypes, inp.Categories)

	var codes []AreaCode
	if inp.DatasetTypes != nil {
		for _, t := range types {
			codes = append(codes, maps.Keys(c.areasForDataTypes[t])...)
		}
	}

	res = c.ctx.Areas.Filter(func(a Area) bool {
		if !filterArea(a, inp) {
			return false
		}

		if len(codes) > 0 && !lo.Contains(codes, a.GetCode()) {
			return false
		}

		if inp.ParentCode != nil && lo.IsNotEmpty(inp.ParentCode) && ParentAreaCode(a) != *inp.ParentCode {
			return false
		}

		return true
	})
	return
}

func (c *InMemoryRepo) DatasetTypes(ctx context.Context, input *DatasetTypesInput) (res []DatasetType, _ error) {
	inp := lo.FromPtr(input)
	return c.ctx.DatasetTypes.Filter(func(t DatasetType) bool {
		return filterDatasetType(t, inp)
	}), nil
}

func (c *InMemoryRepo) Datasets(ctx context.Context, input *DatasetsInput) (res []Dataset, _ error) {
	if input == nil {
		input = &DatasetsInput{}
	}
	return c.ctx.Datasets.Filter(func(t Dataset) bool {
		return filterDataset(t, *input, c.includedStages)
	}), nil
}

func (c *InMemoryRepo) PlateauSpecs(ctx context.Context) ([]*PlateauSpec, error) {
	return lo.Map(c.ctx.PlateauSpecs, func(p PlateauSpec, _ int) *PlateauSpec {
		return &p
	}), nil
}

func (c *InMemoryRepo) Years(ctx context.Context) ([]int, error) {
	return slices.Clone(c.ctx.Years), nil
}

func (c *InMemoryRepo) getDatasetTypeCodes(types []string, categories []DatasetTypeCategory) (res []string) {
	if len(categories) == 0 {
		categories = AllDatasetTypeCategory
	}
	dt := c.ctx.DatasetTypes.DatasetTypesByCategories(categories)

	for _, t := range dt {
		code := t.GetCode()
		if len(types) > 0 && !slices.Contains(types, code) {
			continue
		}
		res = append(res, code)
	}
	return res
}

func areasForDatasetTypes(ds []Dataset) map[string]map[AreaCode]struct{} {
	res := make(map[string]map[AreaCode]struct{})

	for _, d := range ds {
		datasetTypeCode := d.GetTypeCode()

		for _, code := range areaCodesFrom(d) {
			if _, ok := res[datasetTypeCode]; !ok {
				res[datasetTypeCode] = make(map[AreaCode]struct{})
			}
			res[datasetTypeCode][code] = struct{}{}
		}
	}

	return res
}
