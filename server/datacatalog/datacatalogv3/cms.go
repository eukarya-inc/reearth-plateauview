package datacatalogv3

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
)

type CMS struct {
	cms cms.Interface
}

func NewCMS(cms cms.Interface) *CMS {
	return &CMS{cms: cms}
}

func (c *CMS) GetAll(ctx context.Context, project string) (*AllData, error) {
	all := AllData{
		Name: project,
	}

	// TODO: get CMSInfo

	specs, err := c.GetPlateauSpecs(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to get plateau specs: %w", err)
	}

	featureTypes, err := c.GetFeatureTypes(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("failed to get feature types: %w", err)
	}

	all.PlateauSpecs = specs
	all.FeatureTypes = featureTypes

	cityItemsChan := lo.Async2(func() ([]*CityItem, error) {
		return c.GetCityItems(ctx, project, featureTypes.Plateau)
	})

	relatedItemsChan := lo.Async2(func() ([]*RelatedItem, error) {
		return c.GetRelatedItems(ctx, project, featureTypes.Related)
	})

	genericItemsChan := lo.Async2(func() ([]*GenericItem, error) {
		return c.GetGenericItems(ctx, project)
	})

	geospatialjpDataItemsChan := lo.Async2(func() ([]*GeospatialjpDataItem, error) {
		return c.GetGeospatialjpDataItemsWithMaxLODContent(ctx, project)
	})

	featureItemsChans := make([]<-chan lo.Tuple3[string, []*PlateauFeatureItem, error], 0, len(all.FeatureTypes.Plateau))
	for _, featureType := range all.FeatureTypes.Plateau {
		featureType := featureType
		featureItemsChan := lo.Async3(func() (string, []*PlateauFeatureItem, error) {
			res, err := c.GetPlateauItems(ctx, project, featureType.Code)
			return featureType.Code, res, err
		})
		featureItemsChans = append(featureItemsChans, featureItemsChan)
	}

	if res := <-cityItemsChan; res.B != nil {
		return nil, fmt.Errorf("failed to get city items: %w", res.B)
	} else {
		all.City = res.A
	}

	if res := <-relatedItemsChan; res.B != nil {
		return nil, fmt.Errorf("failed to get related items: %w", res.B)
	} else {
		all.Related = res.A
	}

	if res := <-genericItemsChan; res.B != nil {
		return nil, fmt.Errorf("failed to get generic items: %w", res.B)
	} else {
		all.Generic = res.A
	}

	if res := <-geospatialjpDataItemsChan; res.B != nil {
		return nil, fmt.Errorf("failed to get geospatialjp data items: %w", res.B)
	} else {
		all.GeospatialjpDataItems = res.A
	}

	all.Plateau = make(map[string][]*PlateauFeatureItem)
	for _, featureItemsChan := range featureItemsChans {
		if res := <-featureItemsChan; res.C != nil {
			return nil, fmt.Errorf("failed to get feature items (%s): %w", res.A, res.C)
		} else {
			all.Plateau[res.A] = append(all.Plateau[res.A], res.B...)
		}
	}

	return &all, nil
}

func (c *CMS) GetPlateauSpecs(ctx context.Context, project string) ([]plateauapi.PlateauSpecSimple, error) {
	// TODO: load specs from CMS
	return plateauSpecs, nil
}

func (c *CMS) GetFeatureTypes(ctx context.Context, project string) (FeatureTypes, error) {
	// TODO: load feature types from CMS
	return FeatureTypes{
		Plateau: plateauFeatureTypes,
		Related: relatedFeatureTypes,
		Generic: genericFeatureTypes,
	}, nil
}

func (c *CMS) GetCityItems(ctx context.Context, project string, featureTypes []FeatureType) ([]*CityItem, error) {
	items, err := getItemsAndConv[CityItem](
		c.cms, ctx, project, modelPrefix+cityModel,
		func(i cms.Item) *CityItem {
			return CityItemFrom(&i, featureTypes)
		},
	)

	// TODO: dynamic year
	for _, item := range items {
		if item.Year == "" {
			item.Year = "令和5年度"
		}
	}

	return items, err
}

func (c *CMS) GetPlateauItems(ctx context.Context, project, feature string) ([]*PlateauFeatureItem, error) {
	items, err := getItemsAndConv[PlateauFeatureItem](
		c.cms, ctx, project, modelPrefix+feature,
		func(i cms.Item) *PlateauFeatureItem {
			return PlateauFeatureItemFrom(&i)
		},
	)
	return items, err
}

func (c *CMS) GetRelatedItems(ctx context.Context, project string, featureTypes []FeatureType) ([]*RelatedItem, error) {
	items, err := getItemsAndConv[RelatedItem](
		c.cms, ctx, project, modelPrefix+relatedModel,
		func(i cms.Item) *RelatedItem {
			return RelatedItemFrom(&i, featureTypes)
		},
	)
	return items, err
}

func (c *CMS) GetGenericItems(ctx context.Context, project string) ([]*GenericItem, error) {
	items, err := getItemsAndConv[GenericItem](
		c.cms, ctx, project, modelPrefix+genericModel,
		func(i cms.Item) *GenericItem {
			return GenericItemFrom(&i)
		},
	)

	for _, item := range items {
		if item.Category == "" {
			item.Category = "ユースケース"
		}
	}

	return items, err
}

func (c *CMS) GetGeospatialjpDataItems(ctx context.Context, project string) ([]*GeospatialjpDataItem, error) {
	items, err := getItemsAndConv[GeospatialjpDataItem](
		c.cms, ctx, project, modelPrefix+geospatialjpDataModel,
		func(i cms.Item) *GeospatialjpDataItem {
			return GeospatialjpDataItemFrom(&i)
		},
	)

	return items, err
}

func (c *CMS) GetGeospatialjpDataItemsWithMaxLODContent(ctx context.Context, project string) ([]*GeospatialjpDataItem, error) {
	items, err := c.GetGeospatialjpDataItems(ctx, project)
	if err != nil {
		return nil, err
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(10)

	for i := 0; i < len(items); i++ {
		i := i
		url := items[i].MaxLOD
		if url == "" {
			continue
		}

		eg.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			if err != nil {
				return fmt.Errorf("items[%d]: failed to create request: %w", i, err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("items[%d]: failed to get max LOD content: %w", i, err)
			}

			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("items[%d]: failed to get max LOD content: status code %d", i, resp.StatusCode)
			}

			c := csv.NewReader(resp.Body)
			records, err := c.ReadAll()
			if err != nil {
				return fmt.Errorf("items[%d]: failed to read max LOD content: %w", i, err)
			}

			items[i].MaxLODContent = records
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return items, nil
}

func getItemsAndConv[T any](cms cms.Interface, ctx context.Context, project, model string, conv func(cms.Item) *T) ([]*T, error) {
	items, err := cms.GetItemsByKeyInParallel(ctx, project, model, true, 100)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return nil, nil
	}

	res := make([]*T, 0, len(items.Items))
	for _, item := range items.Items {
		res = append(res, conv(item))
	}

	return res, nil
}
