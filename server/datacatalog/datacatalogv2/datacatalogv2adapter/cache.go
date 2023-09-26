package datacatalogv2adapter

import (
	"context"
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

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
	a.areasForDataTypes = make(map[string]map[plateauapi.AreaCode]struct{})

	for _, d := range items {
		ty := d.TypeEn
		areas := a.areasForDataTypes[ty]
		if areas == nil {
			areas = make(map[plateauapi.AreaCode]struct{})
			a.areasForDataTypes[ty] = areas
		}

		if _, found := areas[plateauapi.AreaCode(d.PrefCode)]; !found {
			a.prefectures = append(a.prefectures, prefectureFrom(d))
			areas[plateauapi.AreaCode(d.PrefCode)] = struct{}{}
		}

		if d.City != "" {
			if _, found := areas[plateauapi.AreaCode(d.CityCode)]; !found {
				a.cities = append(a.cities, cityFrom(d))
				areas[plateauapi.AreaCode(d.CityCode)] = struct{}{}
			}
		}

		if d.Ward != "" {
			if _, found := areas[plateauapi.AreaCode(d.WardCode)]; !found {
				a.wards = append(a.wards, wardFrom(d))
				areas[plateauapi.AreaCode(d.WardCode)] = struct{}{}
			}
		}

		if ty := plateauTypeFrom(d); lo.IsNotEmpty(ty) {
			if !lo.Contains(a.plateauDatasetTypes, ty) {
				a.plateauDatasetTypes = append(a.plateauDatasetTypes, ty)
			}
		}

		if ty := relatedTypeFrom(d); lo.IsNotEmpty(ty) {
			if !lo.Contains(a.relatedDatasetTypes, ty) {
				a.relatedDatasetTypes = append(a.relatedDatasetTypes, ty)
			}
		}

		if ty := genericTypeFrom(d); lo.IsNotEmpty(ty) {
			if !lo.Contains(a.genericDatasetTypes, ty) {
				a.genericDatasetTypes = append(a.genericDatasetTypes, ty)
			}
		}

		if !lo.ContainsBy(a.specs, func(a plateauapi.PlateauSpec) bool {
			return a.Name == d.Spec
		}) {
			if s := specFrom(d); lo.IsNotEmpty(s) {
				a.specs = append(a.specs, s)
			}
		}

		if d, ok := plateauDatasetFrom(d); ok {
			a.plateauDatasets = append(a.plateauDatasets, d)
		}
		if d, ok := plateauFloodingDatasetFrom(d); ok {
			a.plateauFloodingDatasets = append(a.plateauFloodingDatasets, d)
		}
		if d, ok := relatedDatasetFrom(d); ok {
			a.relatedDatasets = append(a.relatedDatasets, d)
		}
		if d, ok := genericDatasetFrom(d); ok {
			a.genericDatasets = append(a.genericDatasets, d)
		}
	}

	slices.SortStableFunc(a.prefectures, func(a, b plateauapi.Prefecture) bool {
		return a.Code < b.Code
	})
	slices.SortStableFunc(a.cities, func(a, b plateauapi.City) bool {
		return a.Code < b.Code
	})
	slices.SortStableFunc(a.wards, func(a, b plateauapi.Ward) bool {
		return a.Code < b.Code
	})
	slices.SortStableFunc(a.plateauDatasetTypes, func(a, b plateauapi.PlateauDatasetType) bool {
		return a.Code < b.Code
	})
	slices.SortStableFunc(a.relatedDatasetTypes, func(a, b plateauapi.RelatedDatasetType) bool {
		return a.Code < b.Code
	})
	slices.SortStableFunc(a.genericDatasetTypes, func(a, b plateauapi.GenericDatasetType) bool {
		return a.Code < b.Code
	})
	slices.SortStableFunc(a.specs, func(a, b plateauapi.PlateauSpec) bool {
		return a.Year < b.Year || a.Name < b.Name
	})

	return nil
}
