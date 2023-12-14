package datacatalogv2adapter

import (
	"context"
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"golang.org/x/exp/slices"
)

func New(cmsbase, project string) (*plateauapi.RepoWrapper, error) {
	fetcher, err := datacatalogv2.NewFetcher(cmsbase)
	if err != nil {
		return nil, err
	}
	return From(fetcher, project), nil
}

func From(fetcher datacatalogv2.Fetchable, project string) *plateauapi.RepoWrapper {
	return plateauapi.NewRepoWrapper(func(ctx context.Context) (plateauapi.Repo, error) {
		return fetchAndCreateCache(ctx, project, fetcher, datacatalogv2.FetcherDoOptions{})
	})
}

func fetchAndCreateCache(ctx context.Context, project string, fetcher datacatalogv2.Fetchable, opts datacatalogv2.FetcherDoOptions) (*plateauapi.InMemoryRepo, error) {
	r, err := fetcher.Do(ctx, project, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to update datacatalog cache: %w", err)
	}

	return plateauapi.NewInMemoryRepo(newCache(r)), nil
}

func newCache(r datacatalogv2.ResponseAll) plateauapi.InMemoryRepoContext {
	cache := plateauapi.InMemoryRepoContext{
		PlateauSpecs: plateauSpecs,
	}

	items := r.All()
	areas := make(map[plateauapi.AreaCode]struct{})

	for _, d := range items {
		prefCode := plateauapi.AreaCode(d.PrefCode)
		if _, found := areas[prefCode]; !found {
			if p := prefectureFrom(d); p != nil {
				if cache.Areas == nil {
					cache.Areas = make(plateauapi.Areas)
				}
				cache.Areas.Append(plateauapi.AreaTypePrefecture, []plateauapi.Area{*p})
				areas[prefCode] = struct{}{}
			}
		}

		if d.City != "" {
			areaCode := plateauapi.AreaCode(d.CityCode)
			if _, found := areas[areaCode]; !found {
				if c := cityFrom(d); c != nil {
					if cache.Areas == nil {
						cache.Areas = make(plateauapi.Areas)
					}
					cache.Areas.Append(plateauapi.AreaTypeCity, []plateauapi.Area{*c})
					areas[areaCode] = struct{}{}
				}
			}
		}

		if d.Ward != "" {
			areaCode := plateauapi.AreaCode(d.WardCode)
			if _, found := areas[areaCode]; !found {
				if w := wardFrom(d); w != nil {
					if cache.Areas == nil {
						cache.Areas = make(plateauapi.Areas)
					}
					cache.Areas.Append(plateauapi.AreaTypeWard, []plateauapi.Area{*w})
					areas[areaCode] = struct{}{}
				}
			}
		}

		if ty := plateauDatasetTypeFrom(d); ty.ID != "" {
			if cache.DatasetTypes.DatasetType(ty.ID) == nil {
				if cache.DatasetTypes == nil {
					cache.DatasetTypes = make(plateauapi.DatasetTypes)
				}
				cache.DatasetTypes.Append(plateauapi.DatasetTypeCategoryPlateau, []plateauapi.DatasetType{ty})
			}
		}

		if ty := relatedDatasetTypeFrom(d); ty.ID != "" {
			if cache.DatasetTypes.DatasetType(ty.ID) == nil {
				if cache.DatasetTypes == nil {
					cache.DatasetTypes = make(plateauapi.DatasetTypes)
				}
				cache.DatasetTypes.Append(plateauapi.DatasetTypeCategoryRelated, []plateauapi.DatasetType{ty})
			}
		}

		if ty := genericDatasetTypeFrom(d); ty.ID != "" {
			if cache.DatasetTypes.DatasetType(ty.ID) == nil {
				if cache.DatasetTypes == nil {
					cache.DatasetTypes = make(plateauapi.DatasetTypes)
				}
				cache.DatasetTypes.Append(plateauapi.DatasetTypeCategoryGeneric, []plateauapi.DatasetType{ty})
			}
		}

		if d, ok := plateauDatasetFrom(d); ok {
			if cache.Datasets == nil {
				cache.Datasets = make(plateauapi.Datasets)
			}
			cache.Datasets.Append(plateauapi.DatasetTypeCategoryPlateau, []plateauapi.Dataset{d})
		}
		if d, ok := relatedDatasetFrom(d); ok {
			if cache.Datasets == nil {
				cache.Datasets = make(plateauapi.Datasets)
			}
			cache.Datasets.Append(plateauapi.DatasetTypeCategoryRelated, []plateauapi.Dataset{d})
		}
		if d, ok := genericDatasetFrom(d); ok {
			if cache.Datasets == nil {
				cache.Datasets = make(plateauapi.Datasets)
			}
			cache.Datasets.Append(plateauapi.DatasetTypeCategoryGeneric, []plateauapi.Dataset{d})
		}
		if !slices.Contains(cache.Years, d.Year) {
			cache.Years = append(cache.Years, d.Year)
		}
	}

	slices.SortStableFunc(cache.Years, func(a, b int) int {
		return a - b
	})

	return cache
}
