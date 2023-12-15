package datacatalogv3

import "github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"

type AllData struct {
	PlateauSpecs []plateauapi.PlateauSpecSimple
	FeatureTypes FeatureTypes
	City         []*CityItem
	Related      []*RelatedItem
	Generic      []*GenericItem
	Plateau      map[string][]*PlateauFeatureItem
}

type FeatureTypes struct {
	Plateau []FeatureType
	Related []FeatureType
	Generic []FeatureType
}
