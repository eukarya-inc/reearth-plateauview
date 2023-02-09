package datacatalog

import (
	"encoding/json"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
)

const TypeUsecase = "usecase"

type UsecaseItem struct {
	ID          string           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Type        string           `json:"type,omitempty"`
	Prefecture  string           `json:"prefecture,omitempty"`
	CityName    string           `json:"city_name,omitempty"`
	OpenDataURL string           `json:"opendata_url,omitempty"`
	Description string           `json:"description,omitempty"`
	Year        string           `json:"year,omitempty"`
	Data        *cms.PublicAsset `json:"data,omitempty"`
	DataFormat  string           `json:"data_format,omitempty"`
	DataURL     string           `json:"data_url,omitempty"`
	DataLayers  string           `json:"data_layer,omitempty"`
	Config      string           `json:"config,omitempty"`
}

func (i UsecaseItem) DataCatalogs() []DataCatalogItem {
	var c any
	_ = json.Unmarshal([]byte(i.Config), &c)

	u := i.DataURL
	if u == "" || i.Data != nil && i.Data.URL != "" {
		u = i.Data.URL
	}

	return []DataCatalogItem{{
		ID:          i.ID,
		Name:        i.Name,
		Type:        i.Type,
		Prefecture:  i.Prefecture,
		City:        i.CityName,
		Format:      i.DataFormat,
		URL:         u,
		Description: i.Description,
		Config:      c,
		Layers:      i.DataLayers,
		Year:        i.Year,
	}}
}
