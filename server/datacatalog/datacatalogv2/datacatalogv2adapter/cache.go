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

	for _, d := range items {
		ty := d.TypeEn
		areas := a.areasForDataTypes[ty]
		if areas == nil {
			areas = make(map[plateauapi.AreaCode]struct{})
			a.areasForDataTypes[ty] = areas
		}

		if d.Pref != "" && !lo.ContainsBy(a.prefectures, func(p plateauapi.Prefecture) bool {
			return p.Code == plateauapi.AreaCode(d.PrefCode)
		}) {
			a.prefectures = append(a.prefectures, prefectureFrom(d))
			areas[plateauapi.AreaCode(d.PrefCode)] = struct{}{}
		}

		if d.City != "" && !lo.ContainsBy(a.cities, func(p plateauapi.City) bool {
			return p.Code == plateauapi.AreaCode(d.CityCode)
		}) {
			a.cities = append(a.cities, cityFrom(d))
			areas[plateauapi.AreaCode(d.CityCode)] = struct{}{}
		}

		if d.Ward != "" && !lo.ContainsBy(a.wards, func(p plateauapi.Ward) bool {
			return p.Code == plateauapi.AreaCode(d.WardCode)
		}) {
			a.wards = append(a.wards, wardFrom(d))
			areas[plateauapi.AreaCode(d.WardCode)] = struct{}{}
		}

		if !lo.ContainsBy(a.plateauDatasetTypes, func(a plateauapi.PlateauDatasetType) bool {
			return a.Name == d.TypeEn
		}) {
			if ty := plateauTypeFrom(d); lo.IsNotEmpty(ty) {
				a.plateauDatasetTypes = append(a.plateauDatasetTypes, ty)
			}
		}

		if !lo.ContainsBy(a.relatedDatasetTypes, func(a plateauapi.RelatedDatasetType) bool {
			return a.Name == d.TypeEn
		}) {
			if ty := relatedTypeFrom(d); lo.IsNotEmpty(ty) {
				a.relatedDatasetTypes = append(a.relatedDatasetTypes, ty)
			}
		}

		if !lo.ContainsBy(a.genericDatasetTypes, func(a plateauapi.GenericDatasetType) bool {
			return a.Name == d.TypeEn
		}) {
			if ty := genericTypeFrom(d); lo.IsNotEmpty(ty) {
				a.genericDatasetTypes = append(a.genericDatasetTypes, ty)
			}
		}

		if !lo.ContainsBy(a.specs, func(a plateauapi.PlateauSpec) bool {
			return a.Name == d.Name
		}) {
			a.specs = append(a.specs, specFrom(d))
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
