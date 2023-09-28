package datacatalogv2adapter

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

func filterDataset(d plateauapi.Dataset, input plateauapi.DatasetQuery) bool {
	var dataType string
	var text []string

	switch d2 := d.(type) {
	case plateauapi.PlateauDataset:
		dataType = dataTypeCodeFromDataTypeID(d2.TypeID)
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	case plateauapi.PlateauFloodingDataset:
		dataType = dataTypeCodeFromDataTypeID(d2.TypeID)
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	case plateauapi.RelatedDataset:
		dataType = dataTypeCodeFromDataTypeID(d2.TypeID)
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	case plateauapi.GenericDataset:
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
		areaCodes := areaCodesFrom(d)
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

func areaCodesFrom(d plateauapi.Dataset) []plateauapi.AreaCode {
	switch d2 := d.(type) {
	case plateauapi.PlateauDataset:
		return util.DerefSlice([]*plateauapi.AreaCode{
			lo.ToPtr(d2.PrefectureCode),
			d2.CityCode,
			d2.WardCode,
		})
	case plateauapi.PlateauFloodingDataset:
		return util.DerefSlice([]*plateauapi.AreaCode{
			lo.ToPtr(d2.PrefectureCode),
			d2.CityCode,
			d2.WardCode,
		})
	case plateauapi.RelatedDataset:
		return util.DerefSlice([]*plateauapi.AreaCode{
			lo.ToPtr(d2.PrefectureCode),
			d2.CityCode,
			d2.WardCode,
		})
	case plateauapi.GenericDataset:
		return util.DerefSlice([]*plateauapi.AreaCode{
			lo.ToPtr(d2.PrefectureCode),
			d2.CityCode,
			d2.WardCode,
		})
	}
	return nil
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
		if len(input.SearchTokens) > 0 {
			if !slices.Contains(input.SearchTokens, area2.Name) {
				return false
			}
		}
	case plateauapi.City:
		if len(input.SearchTokens) > 0 {
			if !slices.Contains(input.SearchTokens, area2.Name) {
				return false
			}
		}

		if input.ParentCode != nil && area2.PrefectureCode != *input.ParentCode {
			return false
		}
	case plateauapi.Ward:
		if len(input.SearchTokens) > 0 {
			if !slices.Contains(input.SearchTokens, area2.Name) {
				return false
			}
		}

		if input.ParentCode != nil && (area2.PrefectureCode != *input.ParentCode || area2.CityCode != *input.ParentCode) {
			return false
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

		if input.PlateauSpec != nil && "2.2" != *input.PlateauSpec {
			return false
		}
	case plateauapi.RelatedDatasetType:
		if input.Category != nil && ty2.Category != *input.Category {
			return false
		}

		if input.Year != nil || input.PlateauSpec != nil {
			return false
		}
	case plateauapi.GenericDatasetType:
		if input.Category != nil && ty2.Category != *input.Category {
			return false
		}

		if input.Year != nil || input.PlateauSpec != nil {
			return false
		}
	default:
		return false
	}

	return true
}
