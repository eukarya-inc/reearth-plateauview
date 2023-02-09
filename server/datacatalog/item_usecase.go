package datacatalog

const TypeUsecase = "usecase"

type UsecaseItem struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	// TODO
}

func (i UsecaseItem) DataCatalogs() []DataCatalogItem {
	// TODO
	return []DataCatalogItem{
		{ID: i.ID, Type: i.Type},
	}
}
