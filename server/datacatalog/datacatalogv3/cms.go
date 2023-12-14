package datacatalogv3

import (
	"context"
	"fmt"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/samber/lo"
)

type CMS struct {
	cms cms.Interface
}

func NewCMS(cms cms.Interface) *CMS {
	return &CMS{cms: cms}
}

func (c *CMS) GetAll(ctx context.Context, project string) (*AllData, error) {
	all := AllData{}

	cityItemsChan := lo.Async2(func() ([]*CityItem, error) {
		return c.GetCityItems(ctx, project)
	})

	featureItemsChan := lo.Async2(func() ([]*FeatureItem, error) {
		return c.GetFeatureItems(ctx, project)
	})

	relatedItemsChan := lo.Async2(func() ([]*RelatedItem, error) {
		return c.GetRelatedItems(ctx, project)
	})

	genericItemsChan := lo.Async2(func() ([]*GenericItem, error) {
		return c.GetGenericItems(ctx, project)
	})

	if res := <-cityItemsChan; res.B != nil {
		return nil, fmt.Errorf("failed to get city items: %w", res.B)
	} else {
		all.Cities = res.A
	}

	if res := <-featureItemsChan; res.B != nil {
		return nil, fmt.Errorf("failed to get feature items: %w", res.B)
	} else {
		all.Features = res.A
	}

	if res := <-relatedItemsChan; res.B != nil {
		return nil, fmt.Errorf("failed to get related items: %w", res.B)
	} else {
		all.Relateds = res.A
	}

	if res := <-genericItemsChan; res.B != nil {
		return nil, fmt.Errorf("failed to get generic items: %w", res.B)
	} else {
		all.Generics = res.A
	}

	return &all, nil
}

func (c *CMS) GetCityItems(ctx context.Context, project string) ([]*CityItem, error) {
	panic("not implemented")
}

func (c *CMS) GetFeatureItems(ctx context.Context, project string) ([]*FeatureItem, error) {
	panic("not implemented")
}

func (c *CMS) GetRelatedItems(ctx context.Context, project string) ([]*RelatedItem, error) {
	panic("not implemented")
}

func (c *CMS) GetGenericItems(ctx context.Context, project string) ([]*GenericItem, error) {
	panic("not implemented")
}
