package datacatalog

type DataCatalogItem struct {
	// TODO
}

type DataCatalogItems struct {
	Plateau []DataCatalogItem `json:"plateau"`
	Usecase []DataCatalogItem `json:"usecase"`
	Dataset []DataCatalogItem `json:"dataset"`
}

func (d DataCatalogItems) Merge() []DataCatalogItem {
	// TODO
	return nil
}
