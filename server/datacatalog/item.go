package datacatalog

import "github.com/samber/lo"

type DataCatalogItem struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Prefecture  string `json:"pref,omitempty"`
	City        string `json:"city,omitempty"`
	CityEn      string `json:"city_en,omitempty"`
	CityCode    string `json:"city_code,omitempty"`
	Ward        string `json:"ward,omitempty"`
	WardEn      string `json:"ward_en,omitempty"`
	WardCode    string `json:"ward_code,omitempty"`
	Type        string `json:"type,omitempty"`
	TypeEn      string `json:"type_en,omitempty"`
	Format      string `json:"format,omitempty"`
	Layers      string `json:"layers,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"desc,omitempty"`
	SearchIndex string `json:"search_index,omitempty"`
	Year        string `json:"year,omitempty"`
	Config      any    `json:"config,omitempty"`
}

type DataCatalogGroup struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Prefecture string `json:"pref,omitempty"`
	City       string `json:"city,omitempty"`
	CityEn     string `json:"cityEn,omitempty"`
	Type       string `json:"type,omitempty"`
	Children   []any  `json:"children"`
}

type ResponseAll struct {
	Plateau []PlateauItem
	Usecase []UsecaseItem
}

func (d ResponseAll) All() []DataCatalogItem {
	return append(d.plateau(), d.usecase()...)
}

func (d ResponseAll) ByCities() []DataCatalogGroup {
	// r := d.All()
	// TODO
	return nil
}

func (d ResponseAll) ByTypes() []DataCatalogGroup {
	// r := d.All()
	// TODO
	return nil
}

func (d ResponseAll) plateau() []DataCatalogItem {
	return lo.FlatMap(d.Plateau, func(i PlateauItem, _ int) []DataCatalogItem {
		return i.DataCatalogs()
	})
}

func (d ResponseAll) usecase() []DataCatalogItem {
	return lo.FlatMap(d.Usecase, func(i UsecaseItem, _ int) []DataCatalogItem {
		return i.DataCatalogs()
	})
}
