package datacatalogv2adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"golang.org/x/exp/slices"
)

const minCacheDuration = 30 * time.Second

func New(cmsbase, project string) (*plateauapi.RepoWrapper, error) {
	fetcher, err := datacatalogv2.NewFetcher(cmsbase)
	if err != nil {
		return nil, err
	}
	return From(fetcher, project), nil
}

func From(fetcher datacatalogv2.Fetchable, project string) *plateauapi.RepoWrapper {
	r := plateauapi.NewRepoWrapper(nil, func(ctx context.Context, repo *plateauapi.Repo) error {
		r, err := fetchAndCreateCache(ctx, project, fetcher, datacatalogv2.FetcherDoOptions{})
		if err != nil {
			return err
		}
		*repo = r
		return nil
	})
	r.SetName(fmt.Sprintf("%s(v2)", project))
	r.SetMinCacheDuration(minCacheDuration)
	return r
}

func fetchAndCreateCache(ctx context.Context, project string, fetcher datacatalogv2.Fetchable, opts datacatalogv2.FetcherDoOptions) (*plateauapi.InMemoryRepo, error) {
	r, err := fetcher.Do(ctx, project, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to update datacatalog cache: %w", err)
	}

	all := r.All()

	if err := fetchMaxLOD(ctx, all); err != nil {
		return nil, fmt.Errorf("failed to fetch max lod: %w", err)
	}

	return plateauapi.NewInMemoryRepo(newCache(all)), nil
}

func newCache(items []datacatalogv2.DataCatalogItem) *plateauapi.InMemoryRepoContext {
	cache := &plateauapi.InMemoryRepoContext{
		PlateauSpecs: plateauSpecs,
	}

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
				cache.DatasetTypes.Append(plateauapi.DatasetTypeCategoryPlateau, []plateauapi.DatasetType{&ty})
			}
		}

		if ty := relatedDatasetTypeFrom(d); ty.ID != "" {
			if cache.DatasetTypes.DatasetType(ty.ID) == nil {
				if cache.DatasetTypes == nil {
					cache.DatasetTypes = make(plateauapi.DatasetTypes)
				}
				cache.DatasetTypes.Append(plateauapi.DatasetTypeCategoryRelated, []plateauapi.DatasetType{&ty})
			}
		}

		if ty := genericDatasetTypeFrom(d); ty.ID != "" {
			if cache.DatasetTypes.DatasetType(ty.ID) == nil {
				if cache.DatasetTypes == nil {
					cache.DatasetTypes = make(plateauapi.DatasetTypes)
				}
				cache.DatasetTypes.Append(plateauapi.DatasetTypeCategoryGeneric, []plateauapi.DatasetType{&ty})
			}
		}

		if d := plateauDatasetFrom(d); d != nil {
			if cache.Datasets == nil {
				cache.Datasets = make(plateauapi.Datasets)
			}
			cache.Datasets.Append(plateauapi.DatasetTypeCategoryPlateau, []plateauapi.Dataset{d})
		}
		if d := relatedDatasetFrom(d); d != nil {
			if cache.Datasets == nil {
				cache.Datasets = make(plateauapi.Datasets)
			}
			cache.Datasets.Append(plateauapi.DatasetTypeCategoryRelated, []plateauapi.Dataset{d})
		}
		if d := genericDatasetFrom(d); d != nil {
			if cache.Datasets == nil {
				cache.Datasets = make(plateauapi.Datasets)
			}
			cache.Datasets.Append(plateauapi.DatasetTypeCategoryGeneric, []plateauapi.Dataset{d})
		}
		if !slices.Contains(cache.Years, d.Year) {
			cache.Years = append(cache.Years, d.Year)
		}

		if citygml := citygmlFrom(d); citygml != nil {
			if cache.CityGML == nil {
				cache.CityGML = map[plateauapi.ID]*plateauapi.CityGMLDataset{}
			}

			cg := cache.CityGML[citygml.ID]
			if cg == nil || cg.Year < citygml.Year {
				cache.CityGML[citygml.ID] = citygml
			}
		}
	}

	slices.SortStableFunc(cache.Years, func(a, b int) int {
		return a - b
	})

	return cache
}
