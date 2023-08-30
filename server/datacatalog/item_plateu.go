package datacatalog

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauv2"
	"github.com/samber/lo"
)

type PlateauItem plateauv2.CMSItem

var _ ItemCommon = &PlateauItem{}

func (i PlateauItem) GetCityName() string {
	return i.CityName
}

func (i PlateauItem) DataCatalogs() []DataCatalogItem {
	c := plateauv2.CMSItem(i).IntermediateItem()
	if c.Year == 0 {
		return nil
	}
	return lo.Map(plateauv2.CMSItem(i).AllDataCatalogItems(c), dataCatalogItemFromPlateauV2)
}

func dataCatalogItemFromPlateauV2(i plateauv2.DataCatalogItem, _ int) DataCatalogItem {
	return DataCatalogItem(i)
}
