package plateauv2

import (
	"testing"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestDatasetItem_DataCatalogs(t *testing.T) {
	assert.Equal(t, []DataCatalogItem{{
		ID:          "id",
		Name:        "公園情報（北区）",
		Type:        "公園情報",
		TypeEn:      "park",
		Pref:        "大阪府",
		PrefCode:    "27",
		City:        "大阪市",
		CityCode:    "27100",
		Ward:        "北区",
		WardCode:    "27146",
		Format:      "3dtiles",
		Layers:      []string{"layers", "layers2"},
		URL:         "https://example.com/aaaaa/tileset.json",
		Description: "desc",
		Year:        2021,
		Config:      map[string]any{"a": "b"},
		OpenDataURL: "https://example.com",
		Order:       lo.ToPtr(100),
	}}, DatasetItem{
		ID:          "id",
		Name:        "name",
		Type:        "公園",
		Prefecture:  "大阪府",
		CityName:    "大阪市/北区",
		OpenDataURL: "https://example.com",
		Description: "desc",
		Year:        "令和3年度以前",
		DataFormat:  "3D Tiles",
		DataLayers:  "layers, layers2",
		DataURL:     "https://example.com/aaaaa.zip",
		Config:      `{"a":"b"}`,
		Order:       lo.ToPtr(100),
	}.DataCatalogs())

	assert.Equal(t, []DataCatalogItem{{
		ID:   "id",
		URL:  "url",
		City: "city",
		Ward: "ward",
		Year: 2023,
	}}, DatasetItem{
		ID:       "id",
		DataURL:  "url2",
		CityName: "city",
		WardName: "ward",
		Data: &cms.PublicAsset{
			URL: "url",
		},
		Year: "2023年度",
	}.DataCatalogs())
}
