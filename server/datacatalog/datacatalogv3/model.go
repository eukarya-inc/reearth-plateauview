package datacatalogv3

import "github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"

type AllData struct {
	Name         string
	Year         int
	PlateauSpecs []plateauapi.PlateauSpecSimple
	FeatureTypes FeatureTypes
	City         []*CityItem
	Related      []*RelatedItem
	Generic      []*GenericItem
	Plateau      map[string][]*PlateauFeatureItem
	CMSInfo      CMSInfo
}

type FeatureTypes struct {
	Plateau []FeatureType
	Related []FeatureType
	Generic []FeatureType
}

func (ft FeatureTypes) FindPlateauByCode(code string) *FeatureType {
	for _, f := range ft.Plateau {
		if f.Code == code {
			return &f
		}
	}
	return nil
}

type CMSInfo struct {
	CMSURL         string
	WorkspaceID    string
	ProjectID      string
	PlateauModelID string
	RelatedModelID string
	GenericModelID string
}
