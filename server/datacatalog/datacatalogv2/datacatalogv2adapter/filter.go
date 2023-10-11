package datacatalogv2adapter

import (
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

const SpecVersion = "2.3"

func filterDataset(d plateauapi.Dataset, input plateauapi.DatasetInput) bool {
	var dataType string
	var text []string

	switch d2 := d.(type) {
	case plateauapi.PlateauDataset:
		dataType = d2.TypeCode
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	case plateauapi.PlateauFloodingDataset:
		dataType = d2.TypeCode
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	case plateauapi.RelatedDataset:
		dataType = d2.TypeCode
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	case plateauapi.GenericDataset:
		dataType = d2.TypeCode
		text = []string{
			d2.Name,
			lo.FromPtr(d2.Description),
			lo.FromPtr(d2.Subname),
		}
	default:
		return false
	}

	if len(input.AreaCodes) > 0 {
		var areaCodes []plateauapi.AreaCode
		if lo.FromPtr(input.Deep) {
			areaCodes = areaCodesFrom(d)
		} else {
			areaCodes = util.DerefSlice([]*plateauapi.AreaCode{areaCodeFrom(d)})
		}

		if lo.EveryBy(input.AreaCodes, func(code plateauapi.AreaCode) bool {
			return !slices.Contains(areaCodes, code)
		}) {
			return false
		}
	}

	if !filterByCode(dataType, input.IncludeTypes, input.ExcludeTypes) {
		return false
	}

	if len(input.SearchTokens) > 0 {
		// all tokens must be included in at least one of the text
		if lo.SomeBy(input.SearchTokens, func(t string) bool {
			return lo.EveryBy(text, func(t2 string) bool {
				return t2 == "" || !strings.Contains(t2, t)
			})
		}) {
			return false
		}
	}

	return true
}

func filterByCode(code string, includes []string, excludes []string) bool {
	code = strings.ToLower(code)

	if len(excludes) > 0 {
		if lo.SomeBy(excludes, func(t string) bool {
			return lo.SomeBy(strings.Split(t, "_"), func(c string) bool {
				return strings.ToLower(c) == code
			})
		}) {
			return false
		}
	}

	if len(includes) > 0 {
		if lo.EveryBy(includes, func(t string) bool {
			return lo.EveryBy(strings.Split(t, "_"), func(c string) bool {
				return strings.ToLower(c) != code
			})
		}) {
			return false
		}
	}

	return true
}

func areaCodeFrom(d plateauapi.Dataset) *plateauapi.AreaCode {
	switch d2 := d.(type) {
	case plateauapi.PlateauDataset:
		if d2.WardCode != nil {
			return d2.WardCode
		}
		if d2.CityCode != nil {
			return d2.CityCode
		}
		return &d2.PrefectureCode
	case plateauapi.PlateauFloodingDataset:
		if d2.WardCode != nil {
			return d2.WardCode
		}
		if d2.CityCode != nil {
			return d2.CityCode
		}
		return &d2.PrefectureCode
	case plateauapi.RelatedDataset:
		if d2.WardCode != nil {
			return d2.WardCode
		}
		if d2.CityCode != nil {
			return d2.CityCode
		}
		return &d2.PrefectureCode
	case plateauapi.GenericDataset:
		if d2.WardCode != nil {
			return d2.WardCode
		}
		if d2.CityCode != nil {
			return d2.CityCode
		}
		return &d2.PrefectureCode
	}
	return nil
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

func filterArea(area plateauapi.Area, input plateauapi.AreaInput) bool {
	testName := func(name string) bool {
		return len(input.SearchTokens) == 0 || lo.SomeBy(input.SearchTokens, func(t string) bool {
			return strings.Contains(name, t)
		})
	}

	switch area2 := area.(type) {
	case plateauapi.Prefecture:
		if !testName(area2.Name) {
			return false
		}
	case plateauapi.City:
		if !testName(area2.Name) {
			return false
		}

		if input.ParentCode != nil && area2.PrefectureCode != *input.ParentCode {
			return false
		}
	case plateauapi.Ward:
		if !testName(area2.Name) {
			return false
		}

		if input.ParentCode != nil && area2.CityCode != *input.ParentCode {
			return false
		}
	}

	return true
}

func filterDataType(ty plateauapi.DatasetType, input plateauapi.DatasetTypeInput) bool {
	switch ty2 := ty.(type) {
	case plateauapi.PlateauDatasetType:
		if input.Category != nil && *input.Category != plateauapi.DatasetTypeCategoryPlateau {
			return false
		}

		if input.Year != nil && ty2.Year != *input.Year {
			return false
		}

		if input.PlateauSpec != nil && SpecVersion != *input.PlateauSpec {
			return false
		}
	case plateauapi.RelatedDatasetType:
		if input.Category != nil && *input.Category != plateauapi.DatasetTypeCategoryRelated {
			return false
		}

		if input.Year != nil || input.PlateauSpec != nil {
			return false
		}
	case plateauapi.GenericDatasetType:
		if input.Category != nil && *input.Category != plateauapi.DatasetTypeCategoryGeneric {
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
