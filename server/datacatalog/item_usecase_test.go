package datacatalog

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/stretchr/testify/assert"
)

func TestUsecaseItem_DataCatalogs(t *testing.T) {
	assert.Equal(t, []DataCatalogItem{{
		ID:          "id",
		Name:        "name",
		Type:        "type",
		Prefecture:  "pref",
		City:        "city",
		Format:      "format",
		Layers:      "layers",
		URL:         "url",
		Description: "desc",
		Year:        "year",
		Config:      map[string]any{"a": "b"},
	}}, UsecaseItem{
		ID:          "id",
		Name:        "name",
		Type:        "type",
		Prefecture:  "pref",
		CityName:    "city",
		OpenDataURL: "https://example.com",
		Description: "desc",
		Year:        "year",
		DataFormat:  "format",
		DataLayers:  "layers",
		DataURL:     "url",
		Config:      `{"a":"b"}`,
	}.DataCatalogs())
	assert.Equal(t, []DataCatalogItem{{
		ID:  "id",
		URL: "url",
	}}, UsecaseItem{
		ID:      "id",
		DataURL: "url2",
		Data: &cms.PublicAsset{
			URL: "url",
		},
	}.DataCatalogs())
}
