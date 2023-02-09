package datacatalog

type PlateauItem struct {
	ID string `json:"id"`
	// TODO
}

func (i PlateauItem) DataCatalogs() []DataCatalogItem {
	// TODO
	return []DataCatalogItem{
		{ID: i.ID},
	}
}
