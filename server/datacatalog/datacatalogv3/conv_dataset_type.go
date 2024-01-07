package datacatalogv3

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

func (ft FeatureTypes) ToDatasetTypes(specs []plateauapi.PlateauSpec) plateauapi.DatasetTypes {
	res := make(plateauapi.DatasetTypes)
	res[plateauapi.DatasetTypeCategoryPlateau] = lo.Map(ft.Plateau, func(f FeatureType, _ int) plateauapi.DatasetType {
		spec, _ := lo.Find(specs, func(s plateauapi.PlateauSpec) bool {
			return s.MajorVersion == f.SpecMajor
		})
		return f.ToPlateauDatasetType(spec)
	})
	res[plateauapi.DatasetTypeCategoryRelated] = lo.Map(ft.Related, func(f FeatureType, _ int) plateauapi.DatasetType {
		return f.ToRelatedDatasetType()
	})
	res[plateauapi.DatasetTypeCategoryGeneric] = lo.Map(ft.Generic, func(f FeatureType, _ int) plateauapi.DatasetType {
		return f.ToGenericDatasetType()
	})
	return res
}

func (f *FeatureType) ToPlateauDatasetType(spec plateauapi.PlateauSpec) *plateauapi.PlateauDatasetType {
	return &plateauapi.PlateauDatasetType{
		Category:      plateauapi.DatasetTypeCategoryPlateau,
		ID:            plateauapi.NewID(f.Code, plateauapi.TypeDatasetType),
		Name:          f.Name,
		Code:          f.Code,
		Flood:         f.Flood,
		PlateauSpecID: spec.ID,
		Year:          spec.Year,
	}
}

func (f *FeatureType) ToRelatedDatasetType() *plateauapi.RelatedDatasetType {
	return &plateauapi.RelatedDatasetType{
		Category: plateauapi.DatasetTypeCategoryPlateau,
		ID:       plateauapi.NewID(f.Code, plateauapi.TypeDatasetType),
		Name:     f.Name,
		Code:     f.Code,
	}
}

func (f *FeatureType) ToGenericDatasetType() *plateauapi.GenericDatasetType {
	return &plateauapi.GenericDatasetType{
		Category: plateauapi.DatasetTypeCategoryPlateau,
		ID:       plateauapi.NewID(f.Code, plateauapi.TypeDatasetType),
		Name:     f.Name,
		Code:     f.Code,
	}
}
