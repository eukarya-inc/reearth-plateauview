package datacatalogv3

import "github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"

type AllData struct {
	Name                  string
	PlateauSpecs          []plateauapi.PlateauSpecSimple
	FeatureTypes          FeatureTypes
	City                  []*CityItem
	Related               []*RelatedItem
	Generic               []*GenericItem
	GeospatialjpDataItems []*GeospatialjpDataItem
	Plateau               map[string][]*PlateauFeatureItem
	CMSInfo               CMSInfo
}

type FeatureTypes struct {
	Plateau []FeatureType
	Related []FeatureType
	Generic []FeatureType
}

type CMSInfo struct {
	CMSURL         string
	WorkspaceID    string
	ProjectID      string
	PlateauModelID string
	RelatedModelID string
	GenericModelID string
}
