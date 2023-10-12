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
	areas := make(map[plateauapi.AreaCode]struct{})

	for _, d := range items {
		ty := d.TypeEn
		areasForType := a.areasForDataTypes[ty]
		if areasForType == nil {
			areasForType = make(map[plateauapi.AreaCode]struct{})
		}

		prefCode := plateauapi.AreaCode(d.PrefCode)
		areasForType[prefCode] = struct{}{}
		if _, found := areas[prefCode]; !found {
			if p := prefectureFrom(d); p != nil {
				a.prefectures = append(a.prefectures, *p)
				areas[prefCode] = struct{}{}
			}
		}

		if d.City != "" {
			areaCode := plateauapi.AreaCode(d.CityCode)
			areasForType[areaCode] = struct{}{}
			if _, found := areas[areaCode]; !found {
				if c := cityFrom(d); c != nil {
					a.cities = append(a.cities, *c)
					areas[areaCode] = struct{}{}
				}
			}
		}

		if d.Ward != "" {
			areaCode := plateauapi.AreaCode(d.WardCode)
			areasForType[areaCode] = struct{}{}
			if _, found := areas[areaCode]; !found {
				if w := wardFrom(d); w != nil {
					a.wards = append(a.wards, *w)
					areas[areaCode] = struct{}{}
				}
			}
		}

		if ty := plateauDatasetTypeFrom(d); lo.IsNotEmpty(ty) {
			if !lo.Contains(a.plateauDatasetTypes, ty) {
				a.plateauDatasetTypes = append(a.plateauDatasetTypes, ty)
			}
		}

		if ty := relatedDatasetTypeFrom(d); lo.IsNotEmpty(ty) {
			if !lo.Contains(a.relatedDatasetTypes, ty) {
				a.relatedDatasetTypes = append(a.relatedDatasetTypes, ty)
			}
		}

		if ty := genericDatasetTypeFrom(d); lo.IsNotEmpty(ty) {
			if !lo.Contains(a.genericDatasetTypes, ty) {
				a.genericDatasetTypes = append(a.genericDatasetTypes, ty)
			}
		}

		if d, ok := plateauDatasetFrom(d); ok {
			a.plateauDatasets = append(a.plateauDatasets, d)
		}
		if d, ok := relatedDatasetFrom(d); ok {
			a.relatedDatasets = append(a.relatedDatasets, d)
		}
		if d, ok := genericDatasetFrom(d); ok {
			a.genericDatasets = append(a.genericDatasets, d)
		}
		if !slices.Contains(a.years, d.Year) {
			a.years = append(a.years, d.Year)
		}

		a.areasForDataTypes[ty] = areasForType
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
	slices.SortStableFunc(a.years, func(a, b int) bool {
		return a < b
	})

	return nil
}
