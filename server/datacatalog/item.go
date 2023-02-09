package datacatalog

type DataCatalogItem struct {
	ID          string `json:"id"`
	CityCode    string `json:"cityCode"`
	Prefecture  string `json:"pref"`
	City        string `json:"city"`
	CityEn      string `json:"cityEn"`
	Ward        string `json:"ward,omitempty"`
	WardEn      string `json:"wardEn,omitempty"`
	Type        string `json:"type"`
	DataType    string `json:"dataType"`
	URL         string `json:"url"`
	Description string `json:"desc"`
	SearchIndex string `json:"searchIndex,omitempty"`
	// TODO
}

type DataCatalogGroup struct {
	ID         string `json:"id,omitempty"`
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
