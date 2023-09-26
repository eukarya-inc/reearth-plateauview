package datacatalogv2adapter

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

func filterDataset(d plateauapi.Dataset, input plateauapi.DatasetQuery) bool {
	var areaCodes []plateauapi.AreaCode
	var dataType string
	var text []string

	switch d2 := d.(type) {
	case plateauapi.PlateauDataset:
		areaCodes = areaCodesFromAreaID(d2.AreaID)
		dataType = dataTypeCodeFromDataTypeID(d2.TypeID)
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	case plateauapi.PlateauFloodingDataset:
		areaCodes = areaCodesFromAreaID(d2.AreaID)
		dataType = dataTypeCodeFromDataTypeID(d2.TypeID)
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	case plateauapi.RelatedDataset:
		areaCodes = areaCodesFromAreaID(d2.AreaID)
		dataType = dataTypeCodeFromDataTypeID(d2.TypeID)
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	case plateauapi.GenericDataset:
		areaCodes = areaCodesFromAreaID(d2.AreaID)
		dataType = dataTypeCodeFromDataTypeID(d2.TypeID)
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	default:
		return false
	}

	if len(input.AreaCodes) > 0 {
		if lo.EveryBy(input.AreaCodes, func(code plateauapi.AreaCode) bool {
			return !slices.Contains(areaCodes, code)
		}) {
			return false
		}
	}

	if len(input.ExcludeTypes) > 0 {
		if lo.SomeBy(input.ExcludeTypes, func(t string) bool {
			return t == dataType
		}) {
			return false
		}
	}

	if len(input.IncludeTypes) > 0 {
		if lo.EveryBy(input.IncludeTypes, func(t string) bool {
			return t != dataType
		}) {
			return false
		}
	}

	if len(input.SearchTokens) > 0 {
		if lo.EveryBy(input.SearchTokens, func(t string) bool {
			return !lo.Contains(text, t)
		}) {
			return false
		}
	}

	return true
}

func areaCodesFromAreaID(id plateauapi.ID) []plateauapi.AreaCode {
	i, ty := id.Unwrap()
	if ty != plateauapi.TypePrefecture && ty != plateauapi.TypeMunicipality {
		return nil
	}

	var pref, city, ward string
	switch ty {
	case plateauapi.TypePrefecture:
		pref = i
	case plateauapi.TypeMunicipality:
		city = i
		// ward = i
		// pref = i
	}

	return lo.Filter([]plateauapi.AreaCode{
		plateauapi.AreaCode(pref),
		plateauapi.AreaCode(city),
		plateauapi.AreaCode(ward),
	}, func(code plateauapi.AreaCode, _ int) bool {
		return code != ""
	})
}

func dataTypeCodeFromDataTypeID(id plateauapi.ID) string {
	i, ty := id.Unwrap()
	if ty != plateauapi.TypeDatasetType {
		return ""
	}

	return i
}

func filterArea(area plateauapi.Area, input plateauapi.AreaQuery) bool {
	switch area2 := area.(type) {
	case plateauapi.Prefecture:
		// TODO: datasetType

		if len(input.SearchTokens) > 0 {
			if !slices.Contains(input.SearchTokens, area2.Name) {
				return false
			}
		}
	case plateauapi.Municipality:
		// TODO: parentCode
		// if input.ParentCode != nil && area2.ParentCode != *input.PrefectureCode {
		// 	return false
		// }

		// TODO: datasetType

		if len(input.SearchTokens) > 0 {
			if !slices.Contains(input.SearchTokens, area2.Name) {
				return false
			}
		}
	}

	return true
}

func filterDataType(ty plateauapi.DatasetType, input plateauapi.DatasetTypeQuery) bool {
	switch ty2 := ty.(type) {
	case plateauapi.PlateauDatasetType:
		if input.Category != nil && ty2.Category != *input.Category {
			return false
		}

		if input.Year != nil && ty2.Year != *input.Year {
			return false
		}

		if input.PlateauSpec != nil && ty2.PlateauSpec.Name != *input.PlateauSpec {
			return false
		}
	case plateauapi.RelatedDatasetType:
		if input.Category != nil && ty2.Category != *input.Category {
			return false
		}
	case plateauapi.GenericDatasetType:
		if input.Category != nil && ty2.Category != *input.Category {
			return false
		}
	default:
		return false
	}

	return true
}
