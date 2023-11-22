package datacatalogv2adapter

import (
	"context"
	"fmt"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type cache struct {
	cache               []datacatalogv2.DataCatalogItem
	prefectures         []plateauapi.Prefecture
	cities              []plateauapi.City
	wards               []plateauapi.Ward
	areasForDataTypes   map[string]map[plateauapi.AreaCode]struct{}
	plateauDatasetTypes []plateauapi.PlateauDatasetType
	relatedDatasetTypes []plateauapi.RelatedDatasetType
	genericDatasetTypes []plateauapi.GenericDatasetType
	plateauDatasets     []plateauapi.PlateauDataset
	relatedDatasets     []plateauapi.RelatedDataset
	genericDatasets     []plateauapi.GenericDataset
	years               []int
}

func fetchAndCreateCache(ctx context.Context, project string, fetcher datacatalogv2.Fetchable, opts datacatalogv2.FetcherDoOptions) (*cache, error) {
	r, err := fetcher.Do(ctx, project, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to update datacatalog cache: %w", err)
	}

	return newCache(r)
}

func newCache(r datacatalogv2.ResponseAll) (*cache, error) {
	cache := &cache{}

	items := r.All()
	cache.cache = items
	cache.areasForDataTypes = make(map[string]map[plateauapi.AreaCode]struct{})
	areas := make(map[plateauapi.AreaCode]struct{})

	for _, d := range items {
		ty := d.TypeEn

		areasForType := cache.areasForDataTypes[ty]
		if areasForType == nil {
			areasForType = make(map[plateauapi.AreaCode]struct{})
		}

		prefCode := plateauapi.AreaCode(d.PrefCode)
		areasForType[prefCode] = struct{}{}
		if _, found := areas[prefCode]; !found {
			if p := prefectureFrom(d); p != nil {
				cache.prefectures = append(cache.prefectures, *p)
				areas[prefCode] = struct{}{}
			}
		}

		if d.City != "" {
			areaCode := plateauapi.AreaCode(d.CityCode)
			areasForType[areaCode] = struct{}{}
			if _, found := areas[areaCode]; !found {
				if c := cityFrom(d); c != nil {
					cache.cities = append(cache.cities, *c)
					areas[areaCode] = struct{}{}
				}
			}
		}

		if d.Ward != "" {
			areaCode := plateauapi.AreaCode(d.WardCode)
			areasForType[areaCode] = struct{}{}
			if _, found := areas[areaCode]; !found {
				if w := wardFrom(d); w != nil {
					cache.wards = append(cache.wards, *w)
					areas[areaCode] = struct{}{}
				}
			}
		}

		if ty := plateauDatasetTypeFrom(d); ty.ID != "" {
			if !lo.ContainsBy(cache.plateauDatasetTypes, func(t plateauapi.PlateauDatasetType) bool {
				return t.ID == ty.ID
			}) {
				cache.plateauDatasetTypes = append(cache.plateauDatasetTypes, ty)
			}
		}

		if ty := relatedDatasetTypeFrom(d); ty.ID != "" {
			if !lo.ContainsBy(cache.relatedDatasetTypes, func(t plateauapi.RelatedDatasetType) bool {
				return t.ID == ty.ID
			}) {
				cache.relatedDatasetTypes = append(cache.relatedDatasetTypes, ty)
			}
		}

		if ty := genericDatasetTypeFrom(d); ty.ID != "" {
			if !lo.ContainsBy(cache.genericDatasetTypes, func(t plateauapi.GenericDatasetType) bool {
				return t.ID == ty.ID
			}) {
				cache.genericDatasetTypes = append(cache.genericDatasetTypes, ty)
			}
		}

		if d, ok := plateauDatasetFrom(d); ok {
			cache.plateauDatasets = append(cache.plateauDatasets, d)
		}
		if d, ok := relatedDatasetFrom(d); ok {
			cache.relatedDatasets = append(cache.relatedDatasets, d)
		}
		if d, ok := genericDatasetFrom(d); ok {
			cache.genericDatasets = append(cache.genericDatasets, d)
		}
		if !slices.Contains(cache.years, d.Year) {
			cache.years = append(cache.years, d.Year)
		}

		cache.areasForDataTypes[ty] = areasForType
	}

	slices.SortStableFunc(cache.prefectures, func(a, b plateauapi.Prefecture) int {
		return strings.Compare(string(a.Code), string(b.Code))
	})
	slices.SortStableFunc(cache.cities, func(a, b plateauapi.City) int {
		return strings.Compare(string(a.Code), string(b.Code))
	})
	slices.SortStableFunc(cache.wards, func(a, b plateauapi.Ward) int {
		return strings.Compare(string(a.Code), string(b.Code))
	})
	slices.SortStableFunc(cache.plateauDatasetTypes, func(a, b plateauapi.PlateauDatasetType) int {
		return strings.Compare(a.Code, b.Code)
	})
	slices.SortStableFunc(cache.relatedDatasetTypes, func(a, b plateauapi.RelatedDatasetType) int {
		return strings.Compare(a.Code, b.Code)
	})
	slices.SortStableFunc(cache.genericDatasetTypes, func(a, b plateauapi.GenericDatasetType) int {
		return strings.Compare(a.Code, b.Code)
	})
	slices.SortStableFunc(cache.years, func(a, b int) int {
		return a - b
	})

	return cache, nil
}

var _ plateauapi.Repo = (*cache)(nil)

func (c *cache) Node(ctx context.Context, id plateauapi.ID) (plateauapi.Node, error) {
	i, ty := id.Unwrap()
	switch ty {
	case plateauapi.TypeArea:
		if p, ok := lo.Find(c.prefectures, func(p plateauapi.Prefecture) bool {
			return p.ID == id
		}); ok {
			return &p, nil
		}

		if p, ok := lo.Find(c.cities, func(p plateauapi.City) bool {
			return p.ID == id
		}); ok {
			return &p, nil
		}

		if p, ok := lo.Find(c.wards, func(p plateauapi.Ward) bool {
			return p.ID == id
		}); ok {
			return &p, nil
		}
	case plateauapi.TypeDatasetType:
		if p, ok := lo.Find(c.plateauDatasetTypes, func(p plateauapi.PlateauDatasetType) bool {
			return p.ID == id
		}); ok {
			return &p, nil
		}

		if p, ok := lo.Find(c.relatedDatasetTypes, func(p plateauapi.RelatedDatasetType) bool {
			return p.ID == id
		}); ok {
			return &p, nil
		}

		if p, ok := lo.Find(c.genericDatasetTypes, func(p plateauapi.GenericDatasetType) bool {
			return p.ID == id
		}); ok {
			return &p, nil
		}
	case plateauapi.TypeDataset:
		if p, ok := lo.Find(c.plateauDatasets, func(p plateauapi.PlateauDataset) bool {
			return p.ID == id
		}); ok {
			return &p, nil
		}

		if p, ok := lo.Find(c.relatedDatasets, func(p plateauapi.RelatedDataset) bool {
			return p.ID == id
		}); ok {
			return &p, nil
		}

		if p, ok := lo.Find(c.genericDatasets, func(p plateauapi.GenericDataset) bool {
			return p.ID == id
		}); ok {
			return &p, nil
		}
	case plateauapi.TypeDatasetItem:
		parent, _, _ := cutRight(i, "_")
		parentID := newDatasetID(parent)

		if p, ok := lo.Find(c.plateauDatasets, func(p plateauapi.PlateauDataset) bool {
			return p.ID == parentID
		}); ok {
			item, _ := lo.Find(p.Items, func(i *plateauapi.PlateauDatasetItem) bool {
				return i.ID == id
			})
			return item, nil
		}

		if p, ok := lo.Find(c.relatedDatasets, func(p plateauapi.RelatedDataset) bool {
			return p.ID == parentID
		}); ok {
			item, _ := lo.Find(p.Items, func(i *plateauapi.RelatedDatasetItem) bool {
				return i.ID == id
			})
			return item, nil
		}

		if p, ok := lo.Find(c.genericDatasets, func(p plateauapi.GenericDataset) bool {
			return p.ID == id
		}); ok {
			item, _ := lo.Find(p.Items, func(i *plateauapi.GenericDatasetItem) bool {
				return i.ID == id
			})
			return item, nil
		}
	case plateauapi.TypePlateauSpec:
		if p, ok := lo.Find(plateauSpecs, func(p *plateauapi.PlateauSpec) bool {
			return p.ID == id || lo.SomeBy(p.MinorVersions, func(v *plateauapi.PlateauSpecMinor) bool {
				return v.ID == id
			})
		}); ok {
			if p.ID != id {
				m, _ := lo.Find(p.MinorVersions, func(v *plateauapi.PlateauSpecMinor) bool {
					return v.ID == id
				})
				return m, nil
			}
			return util.CloneRef(p), nil
		}
	}

	return nil, nil
}

func (c *cache) Nodes(ctx context.Context, ids []plateauapi.ID) ([]plateauapi.Node, error) {
	return util.TryMap(ids, func(id plateauapi.ID) (plateauapi.Node, error) {
		return c.Node(ctx, id)
	})
}

func (c *cache) Area(ctx context.Context, code plateauapi.AreaCode) (plateauapi.Area, error) {
	if code.IsPrefectureCode() {
		area, _ := lo.Find(c.prefectures, func(p plateauapi.Prefecture) bool {
			return p.Code == code
		})
		return &area, nil
	}

	if area, ok := lo.Find(c.cities, func(p plateauapi.City) bool {
		return p.Code == code
	}); ok {
		return &area, nil
	}

	if area, ok := lo.Find(c.wards, func(p plateauapi.Ward) bool {
		return p.Code == code
	}); ok {
		return &area, nil
	}

	return nil, nil
}

func (c *cache) Areas(ctx context.Context, input *plateauapi.AreasInput) (res []plateauapi.Area, _ error) {
	inp := lo.FromPtr(input)

	types := c.getDatasetTypes(inp.DatasetTypes, inp.Categories)

	var codes []plateauapi.AreaCode
	if inp.DatasetTypes != nil {
		for _, t := range types {
			codes = append(codes, maps.Keys(c.areasForDataTypes[t])...)
		}
	}

	var prefs []plateauapi.Prefecture
	var cities []plateauapi.City
	var wards []plateauapi.Ward

	if len(inp.AreaTypes) == 0 || slices.Contains(inp.AreaTypes, plateauapi.AreaTypePrefecture) {
		prefs = lo.Filter(c.prefectures, func(t plateauapi.Prefecture, _ int) bool {
			return filterArea(t, inp) && inp.ParentCode == nil && (len(codes) == 0 || lo.Contains(codes, t.Code))
		})
	}

	if len(inp.AreaTypes) == 0 || slices.Contains(inp.AreaTypes, plateauapi.AreaTypeCity) {
		cities = lo.Filter(c.cities, func(t plateauapi.City, _ int) bool {
			if !filterArea(t, inp) {
				return false
			}

			if len(codes) > 0 && !lo.Contains(codes, t.Code) {
				return false
			}

			if inp.ParentCode != nil && t.PrefectureCode != *inp.ParentCode {
				return false
			}

			return true
		})
	}

	if len(inp.AreaTypes) == 0 || slices.Contains(inp.AreaTypes, plateauapi.AreaTypeWard) {
		wards = lo.Filter(c.wards, func(t plateauapi.Ward, _ int) bool {
			if !filterArea(t, inp) {
				return false
			}

			if len(codes) > 0 && !lo.Contains(codes, t.Code) {
				return false
			}

			if inp.ParentCode != nil && t.CityCode != *input.ParentCode {
				return false
			}

			return true
		})
	}

	for _, t := range prefs {
		t := t
		res = append(res, &t)
	}
	for _, t := range cities {
		t := t
		res = append(res, &t)
	}
	for _, t := range wards {
		t := t
		res = append(res, &t)
	}
	return
}

func (c *cache) DatasetTypes(ctx context.Context, input *plateauapi.DatasetTypesInput) (res []plateauapi.DatasetType, _ error) {
	inp := lo.FromPtr(input)
	plateau := lo.Filter(c.plateauDatasetTypes, func(t plateauapi.PlateauDatasetType, _ int) bool {
		return filterDataType(t, inp)
	})
	related := lo.Filter(c.relatedDatasetTypes, func(t plateauapi.RelatedDatasetType, _ int) bool {
		return filterDataType(t, inp)
	})
	generic := lo.Filter(c.genericDatasetTypes, func(t plateauapi.GenericDatasetType, _ int) bool {
		return filterDataType(t, inp)
	})

	for _, t := range plateau {
		t := t
		res = append(res, &t)
	}
	for _, t := range related {
		t := t
		res = append(res, &t)
	}
	for _, t := range generic {
		t := t
		res = append(res, &t)
	}
	return
}

func (c *cache) Datasets(ctx context.Context, input *plateauapi.DatasetsInput) (res []plateauapi.Dataset, _ error) {
	inp := lo.FromPtr(input)
	plateau := lo.Filter(c.plateauDatasets, func(t plateauapi.PlateauDataset, _ int) bool {
		return filterDataset(t, inp)
	})
	related := lo.Filter(c.relatedDatasets, func(t plateauapi.RelatedDataset, _ int) bool {
		return filterDataset(t, inp)
	})
	generic := lo.Filter(c.genericDatasets, func(t plateauapi.GenericDataset, _ int) bool {
		return filterDataset(t, inp)
	})

	for _, t := range plateau {
		t := t
		res = append(res, &t)
	}
	for _, t := range related {
		t := t
		res = append(res, &t)
	}
	for _, t := range generic {
		t := t
		res = append(res, &t)
	}

	return
}

func (c *cache) PlateauSpecs(ctx context.Context) ([]*plateauapi.PlateauSpec, error) {
	return lo.Map(plateauSpecs, func(p *plateauapi.PlateauSpec, _ int) *plateauapi.PlateauSpec {
		return util.CloneRef(p)
	}), nil
}

func (c *cache) Years(ctx context.Context) ([]int, error) {
	return slices.Clone(c.years), nil
}

func (c *cache) getDatasetTypes(types []string, categories []plateauapi.DatasetTypeCategory) (res []string) {
	for _, t := range c.allDatasetTypes(categories) {
		code := t.GetCode()
		if len(types) > 0 && !slices.Contains(types, code) {
			continue
		}
		res = append(res, code)
	}
	return res
}

func (c *cache) allDatasetTypes(categories []plateauapi.DatasetTypeCategory) (res []plateauapi.DatasetType) {
	var plateau []plateauapi.DatasetType
	if len(categories) == 0 || slices.Contains(categories, plateauapi.DatasetTypeCategoryPlateau) {
		for _, t := range c.plateauDatasetTypes {
			t := t
			plateau = append(plateau, &t)
		}
	}
	slices.SortStableFunc(plateau, func(a, b plateauapi.DatasetType) int {
		return strings.Compare(a.GetCode(), b.GetCode())
	})

	var related []plateauapi.DatasetType
	if len(categories) == 0 || slices.Contains(categories, plateauapi.DatasetTypeCategoryRelated) {
		for _, t := range c.relatedDatasetTypes {
			t := t
			related = append(related, &t)
		}
	}
	slices.SortStableFunc(related, func(a, b plateauapi.DatasetType) int {
		return strings.Compare(a.GetCode(), b.GetCode())
	})

	var generic []plateauapi.DatasetType
	if len(categories) == 0 || slices.Contains(categories, plateauapi.DatasetTypeCategoryGeneric) {
		for _, t := range c.genericDatasetTypes {
			t := t
			generic = append(generic, &t)
		}
	}
	slices.SortStableFunc(generic, func(a, b plateauapi.DatasetType) int {
		return strings.Compare(a.GetCode(), b.GetCode())
	})

	return append(append(plateau, related...), generic...)
}
