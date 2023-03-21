package geospatialjp

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/eukarya-inc/reearth-plateauview/server/geospatialjp/ckan"
	"github.com/jarcoal/httpmock"
	"github.com/reearth/reearthx/rerror"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func TestService_CheckCatalog(t *testing.T) {
	ctx := context.Background()
	catalogData := lo.Must(os.ReadFile("testdata/xxxxx_xxx_catalog.xlsx"))

	httpmock.Activate()
	defer httpmock.Deactivate()
	httpmock.RegisterResponder("GET", "https://example.com/catalog.xlsx", httpmock.NewBytesResponder(http.StatusOK, catalogData))

	cmsm := &mockCMS{}
	s := &Services{
		CMS: cmsm,
	}

	assert.NoError(t, s.CheckCatalog(ctx, "prj", Item{
		ID:      "item",
		Catalog: "catalog",
	}))
	assert.Equal(t, cms.Item{
		ID: "item",
		Fields: []cms.Field{
			{Key: "catalog_status", Value: "完了", Type: "select"},
		},
	}, cmsm.item)
}

func TestService_RegisterCkanResources(t *testing.T) {
	ctx := context.Background()
	catalogData := lo.Must(os.ReadFile("testdata/xxxxx_xxx_catalog.xlsx"))
	cf := NewCatalogFile(lo.Must(excelize.OpenReader(bytes.NewReader(catalogData))))
	_ = cf.DeleteSheet()
	catalogData2 := lo.Must(cf.File().WriteToBuffer()).Bytes()

	httpmock.Activate()
	defer httpmock.Deactivate()
	httpmock.RegisterResponder("GET", "https://example.com/catalog.xlsx", httpmock.NewBytesResponder(http.StatusOK, catalogData))
	httpmock.RegisterResponder("GET", "https://example.com/catalog2.xlsx", httpmock.NewBytesResponder(http.StatusOK, catalogData2))

	cms := &mockCMS{}
	ckan := ckan.NewMock("org", nil, nil)
	s := &Services{
		CMS:         cms,
		Ckan:        ckan,
		CkanOrg:     "org",
		CkanPrivate: true,
	}

	// case1: upload all files of 第2.3版
	assert.NoError(t, s.RegisterCkanResources(ctx, Item{
		Specification: "第2.3版",
		CityGML:       "citygml",
		Catalog:       "catalog",
		All:           "all",
	}))

	pkg, err := ckan.ShowPackage(ctx, "plateau-12210-mobara-shi-2022")
	assert.NoError(t, err)
	assert.Equal(t, "plateau-12210-mobara-shi-2022", pkg.Name)
	assert.Equal(t, "TITLE", pkg.Title)
	assert.Greater(t, len(pkg.ThumbnailURL), 100)
	assert.Equal(t, 3, len(pkg.Resources))
	assert.Equal(t, "3D Tiles, MVT（v2）", pkg.Resources[0].Name)
	assert.Equal(t, "https://example.com/all.zip", pkg.Resources[0].URL)
	assert.Equal(t, "CityGML（v2）", pkg.Resources[1].Name)
	assert.Equal(t, "https://example.com/12210_mobara-shi_2022_citygml_1_lsld.zip", pkg.Resources[1].URL)
	assert.Equal(t, "データ目録（v2）", pkg.Resources[2].Name)

	// case2: upload citygml and catalog of 第1版 to an existing package
	// assert.NoError(t, s.RegisterCkanResources(ctx, Item{
	// 	Specification: "第1版",
	// 	CityGML:       "citygml2",
	// 	Catalog:       "catalog2",
	// }))

	// pkg, err = ckan.ShowPackage(ctx, "plateau-12210-mobara-shi-2020")
	// assert.NoError(t, err)
	// assert.Equal(t, "plateau-12210-mobara-shi-2020", pkg.Name)
	// assert.Equal(t, "TITLE", pkg.Title)
	// assert.Greater(t, len(pkg.ThumbnailURL), 100)
	// assert.Equal(t, 2, len(pkg.Resources))
	// assert.Equal(t, "CityGML", pkg.Resources[0].Name)
	// assert.Equal(t, "https://example.com/12210_mobara-shi_2020_citygml_1_lsld.zip", pkg.Resources[0].URL)
	// assert.Equal(t, "データ目録", pkg.Resources[1].Name)
}

type mockCMS struct {
	cms.Interface
	item cms.Item
}

func (c *mockCMS) UpdateItem(ctx context.Context, itemID string, fields []cms.Field) (*cms.Item, error) {
	c.item = cms.Item{
		ID:     itemID,
		Fields: fields,
	}
	return nil, nil
}

func (*mockCMS) Asset(ctx context.Context, id string) (*cms.Asset, error) {
	if id == "catalog" {
		return &cms.Asset{
			ID:  "catalog",
			URL: "https://example.com/catalog.xlsx",
		}, nil
	}
	if id == "catalog2" {
		return &cms.Asset{
			ID:  "catalog",
			URL: "https://example.com/catalog2.xlsx",
		}, nil
	}
	if id == "citygml" {
		return &cms.Asset{
			ID:  "citygml",
			URL: "https://example.com/12210_mobara-shi_2022_citygml_1_lsld.zip",
		}, nil
	}
	if id == "citygml2" {
		return &cms.Asset{
			ID:  "citygml",
			URL: "https://example.com/12210_mobara-shi_2020_citygml_1_lsld.zip",
		}, nil
	}
	if id == "all" {
		return &cms.Asset{
			ID:  "catalog",
			URL: "https://example.com/all.zip",
		}, nil
	}
	return nil, rerror.ErrNotFound
}
