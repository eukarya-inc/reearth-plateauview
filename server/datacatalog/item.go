package datacatalog

type DataCatalogItem struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	CityCode    string `json:"city_code,omitempty"`
	Prefecture  string `json:"pref,omitempty"`
	City        string `json:"city,omitempty"`
	CityEn      string `json:"cityEn,omitempty"`
	Ward        string `json:"ward,omitempty"`
	WardEn      string `json:"wardEn,omitempty"`
	Type        string `json:"type,omitempty"`
	Format      string `json:"format,omitempty"`
	Layers      string `json:"layers,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"desc,omitempty"`
	SearchIndex string `json:"searchIndex,omitempty"`
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

func (d ResponseAll) Merge() []DataCatalogItem {
	// TODO
	return nil
}

func (d ResponseAll) MergeByCities() []DataCatalogGroup {
	// r := d.Merge()
	// TODO
	return nil
}

func (d ResponseAll) MergeByTypes() []DataCatalogGroup {
	// r := d.Merge()
	// TODO
	return nil
}
