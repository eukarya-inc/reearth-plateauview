package sdkapi

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"path"
	"strconv"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/rerror"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

const limit = 10

type CMS struct {
	Project              string
	PublicAPI            bool
	IntegrationAPIClient cms.Interface
	PublicAPIClient      *cms.PublicAPIClient[Item]
}

func NewCMS(c cms.Interface, pac *cms.PublicAPIClient[Item], project string, usePublic bool) *CMS {
	return &CMS{
		Project:              project,
		PublicAPI:            usePublic,
		IntegrationAPIClient: c,
		PublicAPIClient:      pac,
	}
}

func (c *CMS) Datasets(ctx context.Context, model string) (*DatasetResponse, error) {
	if c.PublicAPI {
		return c.DatasetsWithPublicAPI(ctx, model)
	}
	return c.DatasetsWithIntegrationAPI(ctx, model)
}

func (c *CMS) Files(ctx context.Context, model, id string) (FilesResponse, error) {
	if c.PublicAPI {
		return c.FilesWithPublicAPI(ctx, model, id)
	}
	return c.FilesWithIntegrationAPI(ctx, model, id)
}

func (c *CMS) DatasetsWithPublicAPI(ctx context.Context, model string) (*DatasetResponse, error) {
	items, err := c.PublicAPIClient.GetAllItemsInParallel(ctx, c.Project, model, limit)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}

	return Items(items).DatasetResponse(), nil
}

func (c *CMS) FilesWithPublicAPI(ctx context.Context, model, id string) (FilesResponse, error) {
	item, err := c.PublicAPIClient.GetItem(ctx, c.Project, model, id)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}
	if item.CityGML == nil || item.MaxLOD == nil {
		return nil, rerror.ErrNotFound
	}

	asset, err := c.PublicAPIClient.GetAsset(ctx, c.Project, item.CityGML.ID)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}

	maxlod, err := getMaxLOD(ctx, item.MaxLOD.URL)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}

	return MaxLODFiles(maxlod, asset.Files, nil), nil
}

func (c *CMS) DatasetsWithIntegrationAPI(ctx context.Context, model string) (*DatasetResponse, error) {
	items, err := c.IntegrationAPIClient.GetItemsByKey(ctx, c.Project, model, true)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}

	return ItemsFromIntegration(items.Items).DatasetResponse(), nil
}

func (c *CMS) FilesWithIntegrationAPI(ctx context.Context, model, id string) (FilesResponse, error) {
	item, err := c.IntegrationAPIClient.GetItem(ctx, id, true)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}

	iitem := ItemFromIntegration(item)
	if iitem.CityGML == nil || iitem.MaxLOD == nil || !iitem.IsPublic() {
		return nil, rerror.ErrNotFound
	}

	asset, err := c.IntegrationAPIClient.Asset(ctx, iitem.CityGML.ID)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}
	if asset.File == nil {
		return nil, rerror.ErrNotFound
	}

	assetURL, err := url.Parse(asset.URL)
	if asset.File == nil {
		return nil, rerror.ErrInternalBy(fmt.Errorf("failed to parse asset url %s: %w", asset.URL, err))
	}

	assetBase := util.CloneRef(assetURL)
	assetBase.Path = path.Dir(assetBase.Path)

	maxlod, err := getMaxLOD(ctx, iitem.MaxLOD.URL)
	if err != nil {
		return nil, rerror.ErrInternalBy(err)
	}

	return MaxLODFiles(maxlod, asset.File.Paths(), assetBase), nil
}

func ReadMaxLODCSV(b io.Reader) (MaxLODColumns, error) {
	r := csv.NewReader(b)
	r.ReuseRecord = true
	var results MaxLODColumns
	for {
		c, err := r.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read csv: %w", err)
		}

		if len(c) < 3 || !isInt(c[0]) {
			continue
		}

		m, err := strconv.ParseFloat(c[2], 64)
		if err != nil {
			continue
		}

		f := ""
		if len(c) > 3 {
			f = c[3]
		}

		results = append(results, MaxLODColumn{
			Code:   c[0],
			Type:   c[1],
			MaxLOD: m,
			File:   f,
		})
	}

	return results, nil
}

func MaxLODFiles(maxLOD MaxLODColumns, assetPaths []string, assetBase *url.URL) FilesResponse {
	files := lo.FilterMap(assetPaths, func(u string, _ int) (*url.URL, bool) {
		if path.Ext(u) != ".gml" {
			return nil, false
		}

		u2, err := url.Parse(u)
		if err != nil {
			return nil, false
		}

		if assetBase == nil {
			return u2, true
		}

		fu := util.CloneRef(assetBase)
		fu.Path = path.Join(fu.Path, u)
		return fu, true
	})

	return maxLOD.Map().Files(files)
}
