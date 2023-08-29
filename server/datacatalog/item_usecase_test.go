package datacatalog

import (
	"testing"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/stretchr/testify/assert"
)

func TestUsecaseItem_DataCatalogs(t *testing.T) {
	assert.Equal(t, []DataCatalogItem{{
		ID:       "id",
		Type:     "ユースケース",
		TypeEn:   "usecase",
		URL:      "url",
		City:     "city",
		Ward:     "ward",
		Year:     2023,
		Category: "ユースケース",
	}}, UsecaseItem{
		ID:       "id",
		DataURL:  "url2",
		CityName: "city",
		WardName: "ward",
		Data: &cms.PublicAsset{
			URL: "url",
		},
		Year: "2023年度",
	}.DataCatalogs())

	assert.Equal(t, []DataCatalogItem{{
		ID:       "id",
		Type:     "フォルダ",
		TypeEn:   "folder",
		Name:     "name",
		Pref:     "大阪府",
		PrefCode: "27",
		City:     "大阪市",
		CityCode: "27100",
		Ward:     "北区",
		WardCode: "27146",
		Category: "ユースケース",
	}}, UsecaseItem{
		ID:         "id",
		Name:       "name",
		Prefecture: "大阪府",
		CityName:   "大阪市/北区",
		DataFormat: "フォルダ",
	}.DataCatalogs())
}
