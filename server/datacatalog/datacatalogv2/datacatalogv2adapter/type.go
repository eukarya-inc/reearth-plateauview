package datacatalogv2adapter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogutil"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

const usecaseID = "usecase"

var floodingTypes = []string{"fld", "htd", "tnm", "ifld"}

var plateauSpecs = []*plateauapi.PlateauSpec{
	{
		ID:           plateauSpecIDFrom("2"),
		MajorVersion: 2,
		Year:         2022,
		MinorVersions: []*plateauapi.PlateauSpecMinor{
			{
				ID:           plateauSpecIDFrom("2.3"),
				Version:      "2.3",
				Name:         "第2.3版",
				MajorVersion: 2,
				Year:         2022,
				ParentID:     plateauSpecIDFrom("2"),
			},
		},
	},
	{
		ID:           plateauSpecIDFrom("3"),
		MajorVersion: 3,
		Year:         2022,
		MinorVersions: []*plateauapi.PlateauSpecMinor{
			{
				ID:           plateauSpecIDFrom("3.0"),
				Version:      "3.0",
				Name:         "第3.0版",
				MajorVersion: 3,
				Year:         2023,
				ParentID:     plateauSpecIDFrom("3"),
			},
		},
	},
}

func plateauDatasetFrom(d datacatalogv2.DataCatalogItem) (plateauapi.PlateauDataset, bool) {
	if d.Family != "plateau" {
		return plateauapi.PlateauDataset{}, false
	}

	plateauSpecVersion := d.Spec
	if isEx(d) {
		plateauSpecVersion = "3.0"
	}

	var subname *string
	var river *plateauapi.River
	if slices.Contains(floodingTypes, d.TypeEn) {
		if d.TypeEn == "fld" {
			var admin plateauapi.RiverAdmin
			if strings.Contains(d.Name, "（国管理区間）") {
				admin = plateauapi.RiverAdminNational
			} else if strings.Contains(d.Name, "（都道府県管理区間）") {
				admin = plateauapi.RiverAdminPrefecture
			}

			names := strings.Split(reBrackets.ReplaceAllString(d.Name, ""), " ")
			name, _ := lo.Find(names, func(s string) bool {
				return strings.HasSuffix(s, "川")
			})

			river = &plateauapi.River{
				Name:  name,
				Admin: admin,
			}
		} else {
			names := strings.Split(d.Name, " ")
			if len(names) > 1 {
				subname = lo.ToPtr(names[1])
			}
		}
	}

	id := datasetIDFrom(d, nil)
	return plateauapi.PlateauDataset{
		ID:              id,
		Name:            d.Name,
		Subname:         subname,
		Description:     lo.ToPtr(d.Description),
		PrefectureID:    prefectureIDFrom(d),
		PrefectureCode:  prefectureCodeFrom(d),
		CityID:          cityIDFrom(d),
		CityCode:        cityCodeFrom(d),
		WardID:          wardIDFrom(d),
		WardCode:        wardCodeFrom(d),
		Year:            d.Year,
		TypeID:          datasetTypeIDFrom(d),
		TypeCode:        datasetTypeCodeFrom(d),
		Groups:          groupsFrom(d),
		PlateauSpecID:   plateauSpecIDFrom(plateauSpecVersion),
		PlateauSpecName: plateauSpecVersion,
		River:           river,
		Items: lo.Map(d.MainOrConfigItems(), func(c datacatalogutil.DataCatalogItemConfigItem, _ int) *plateauapi.PlateauDatasetItem {
			return plateauDatasetItemFrom(c, id)
		}),
	}, true
}

func plateauDatasetItemFrom(c datacatalogutil.DataCatalogItemConfigItem, parentID plateauapi.ID) *plateauapi.PlateauDatasetItem {
	var lod *int
	if strings.HasPrefix(c.Name, "LOD") {
		l, _, _ := strings.Cut(c.Name[3:], "（")
		lodf, err := strconv.Atoi(l)
		if err == nil {
			lod = &lodf
		}
	}

	format := datasetFormatFrom(c.Type)

	var texture *plateauapi.Texture
	if strings.Contains(c.Name, "（テクスチャなし）") {
		texture = lo.ToPtr(plateauapi.TextureNone)
	} else if format == plateauapi.DatasetFormatCesium3dtiles {
		texture = lo.ToPtr(plateauapi.TextureTexture)
	}

	id := c.Name
	var floodingScale *plateauapi.FloodingScale
	if strings.Contains(c.Name, "想定最大規模") || strings.Contains(c.Name, "L2") {
		floodingScale = lo.ToPtr(plateauapi.FloodingScaleExpectedMaximum)
		id = "l2"
	} else if strings.Contains(c.Name, "計画規模") || strings.Contains(c.Name, "L1") {
		floodingScale = lo.ToPtr(plateauapi.FloodingScalePlanned)
		id = "l1"
	}

	return &plateauapi.PlateauDatasetItem{
		ID:            plateauapi.NewID(fmt.Sprintf("%s_%s", parentID.ID(), id), plateauapi.TypeDatasetItem),
		Name:          c.Name,
		URL:           c.URL,
		Format:        format,
		Layers:        c.Layers,
		ParentID:      parentID,
		Lod:           lod,
		Texture:       texture,
		FloodingScale: floodingScale,
	}
}

var reBrackets = regexp.MustCompile(`（[^（]*）`)

func relatedDatasetFrom(d datacatalogv2.DataCatalogItem) (plateauapi.RelatedDataset, bool) {
	if d.Family != "related" {
		return plateauapi.RelatedDataset{}, false
	}

	id := datasetIDFrom(d, nil)
	return plateauapi.RelatedDataset{
		ID:             id,
		Name:           d.Name,
		Subname:        nil,
		Description:    lo.ToPtr(d.Description),
		PrefectureID:   prefectureIDFrom(d),
		PrefectureCode: prefectureCodeFrom(d),
		CityID:         cityIDFrom(d),
		CityCode:       cityCodeFrom(d),
		WardID:         wardIDFrom(d),
		WardCode:       wardCodeFrom(d),
		Year:           d.Year,
		TypeID:         datasetTypeIDFrom(d),
		TypeCode:       datasetTypeCodeFrom(d),
		Groups:         groupsFrom(d),
		Items: lo.Map(d.MainOrConfigItems(), func(c datacatalogutil.DataCatalogItemConfigItem, _ int) *plateauapi.RelatedDatasetItem {
			return &plateauapi.RelatedDatasetItem{
				ID:       plateauapi.NewID(id.ID(), plateauapi.TypeDatasetItem), // RelatedDatasetItem should be single
				Name:     c.Name,
				URL:      c.URL,
				Format:   datasetFormatFrom(c.Type),
				Layers:   c.Layers,
				ParentID: id,
			}
		}),
	}, true
}

func genericDatasetFrom(d datacatalogv2.DataCatalogItem) (plateauapi.GenericDataset, bool) {
	if d.Family != "generic" {
		return plateauapi.GenericDataset{}, false
	}

	id := datasetIDFrom(d, nil)
	return plateauapi.GenericDataset{
		ID:             id,
		Name:           d.Name,
		Subname:        nil,
		Description:    lo.ToPtr(d.Description),
		PrefectureID:   prefectureIDFrom(d),
		PrefectureCode: prefectureCodeFrom(d),
		CityID:         cityIDFrom(d),
		CityCode:       cityCodeFrom(d),
		WardID:         wardIDFrom(d),
		WardCode:       wardCodeFrom(d),
		Year:           d.Year,
		TypeID:         datasetTypeIDFrom(d),
		TypeCode:       datasetTypeCodeFrom(d),
		Groups:         groupsFrom(d),
		Items: lo.Map(d.MainOrConfigItems(), func(c datacatalogutil.DataCatalogItemConfigItem, _ int) *plateauapi.GenericDatasetItem {
			return &plateauapi.GenericDatasetItem{
				ID:       plateauapi.NewID(fmt.Sprintf("%s:%s", d.ID, c.Name), plateauapi.TypeDatasetItem),
				Name:     c.Name,
				URL:      c.URL,
				Format:   datasetFormatFrom(c.Type),
				Layers:   c.Layers,
				ParentID: id,
			}
		}),
	}, true
}

func datasetFormatFrom(f string) plateauapi.DatasetFormat {
	switch strings.ToLower(f) {
	case "geojson":
		return plateauapi.DatasetFormatGeojson
	case "3dtiles":
		return plateauapi.DatasetFormatCesium3dtiles
	case "czml":
		return plateauapi.DatasetFormatCzml
	case "gtfs":
	case "gtfs-realtime":
		return plateauapi.DatasetFormatGtfsRealtime
	case "gltf":
		return plateauapi.DatasetFormatGltf
	case "mvt":
		return plateauapi.DatasetFormatMvt
	case "tiles":
		return plateauapi.DatasetFormatTiles
	case "tms":
		return plateauapi.DatasetFormatTms
	case "wms":
		return plateauapi.DatasetFormatWms
	case "csv":
		return plateauapi.DatasetFormatCSV
	}
	return ""
}

func prefectureIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	return plateauapi.NewID(d.PrefCode, plateauapi.TypeArea)
}

func cityIDFrom(d datacatalogv2.DataCatalogItem) *plateauapi.ID {
	if d.CityCode == "" {
		return nil
	}
	return lo.ToPtr(plateauapi.NewID(d.CityCode, plateauapi.TypeArea))
}

func wardIDFrom(d datacatalogv2.DataCatalogItem) *plateauapi.ID {
	if d.WardCode == "" {
		return nil
	}
	return lo.ToPtr(plateauapi.NewID(d.WardCode, plateauapi.TypeArea))
}

func prefectureCodeFrom(d datacatalogv2.DataCatalogItem) plateauapi.AreaCode {
	return plateauapi.AreaCode(d.PrefCode)
}

func cityCodeFrom(d datacatalogv2.DataCatalogItem) *plateauapi.AreaCode {
	if d.CityCode == "" {
		return nil
	}
	return lo.ToPtr(plateauapi.AreaCode(d.CityCode))
}

func wardCodeFrom(d datacatalogv2.DataCatalogItem) *plateauapi.AreaCode {
	if d.WardCode == "" {
		return nil
	}
	return lo.ToPtr(plateauapi.AreaCode(d.WardCode))
}

func datasetIDFrom(d datacatalogv2.DataCatalogItem, subname *string) plateauapi.ID {
	if d.Family == "plateau" || d.Family == "related" {
		invalid := false
		areaCode := d.WardCode
		if areaCode == "" {
			areaCode = d.CityCode
		}
		if areaCode == "" {
			areaCode = d.PrefCode
		}

		sub := ""
		typeCode := datasetTypeCodeFrom(d)
		isFlood := slices.Contains(floodingTypes, d.TypeEn)

		if isFlood || d.TypeEn == "gen" || isEx(d) {
			if _, after, found := strings.Cut(d.ID, "_"+typeCode+"_"); found {
				if isFlood {
					after = strings.TrimSuffix(after, "_l1")
					after = strings.TrimSuffix(after, "_l2")
				}
				sub = fmt.Sprintf("_%s", after)
			} else {
				invalid = true
			}
		} else if d.TypeEn == "urf" {
			sub = fmt.Sprintf("_%s", d.Type2En)
		}

		if !invalid {
			return newDatasetID(fmt.Sprintf("%s_%s%s", areaCode, typeCode, sub))
		}
	}

	return newDatasetID(d.ID)
}

func newDatasetID(id string) plateauapi.ID {
	return plateauapi.NewID(id, plateauapi.TypeDataset)
}

func datasetTypeIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	code := datasetTypeCodeFrom(d)
	if d.Family == "plateau" {
		spec := d.Spec
		if isEx(d) {
			spec = "3.0"
		}
		return plateauapi.NewID(fmt.Sprintf("%s_%s", code, majorVersion(spec)), plateauapi.TypeDatasetType)
	}
	return plateauapi.NewID(code, plateauapi.TypeDatasetType)
}

func datasetTypeCodeFrom(d datacatalogv2.DataCatalogItem) string {
	if d.Family == "plateau" {
		if strings.HasPrefix(d.TypeEn, "urf_") {
			return "urf"
		}
		return d.TypeEn
	}
	if d.Family == "related" {
		return d.TypeEn
	}
	if d.Family == "generic" && d.Category != "" {
		if d.Category == "サンプルデータ" {
			return "sample"
		}
		return d.Category
	}
	return usecaseID
}

func plateauSpecIDFrom(version string) plateauapi.ID {
	return plateauapi.NewID(specNumber(version), plateauapi.TypePlateauSpec)
}

func plateauSpecMajorIDFrom(version string) plateauapi.ID {
	return plateauapi.NewID(majorVersion(version), plateauapi.TypePlateauSpec)
}

func specNumber(spec string) string {
	return strings.TrimSuffix(strings.TrimPrefix(spec, "第"), "版")
}

func prefectureFrom(d datacatalogv2.DataCatalogItem) *plateauapi.Prefecture {
	if d.PrefCode == "" {
		return nil
	}

	return &plateauapi.Prefecture{
		ID:   prefectureIDFrom(d),
		Code: prefectureCodeFrom(d),
		Name: d.Pref,
	}
}

func cityFrom(d datacatalogv2.DataCatalogItem) *plateauapi.City {
	id, code := cityIDFrom(d), cityCodeFrom(d)
	if id == nil || code == nil {
		return nil
	}

	return &plateauapi.City{
		ID:             *id,
		Code:           *code,
		Name:           d.City,
		PrefectureID:   prefectureIDFrom(d),
		PrefectureCode: prefectureCodeFrom(d),
	}
}

func wardFrom(d datacatalogv2.DataCatalogItem) *plateauapi.Ward {
	id, code := wardIDFrom(d), wardCodeFrom(d)
	if id == nil || code == nil {
		return nil
	}

	cityid, citycode := cityIDFrom(d), cityCodeFrom(d)
	if cityid == nil || citycode == nil {
		return nil
	}

	return &plateauapi.Ward{
		ID:             *id,
		Code:           *code,
		Name:           d.Ward,
		PrefectureID:   prefectureIDFrom(d),
		PrefectureCode: prefectureCodeFrom(d),
		CityID:         *cityid,
		CityCode:       *citycode,
	}
}

func plateauDatasetTypeFrom(d datacatalogv2.DataCatalogItem) plateauapi.PlateauDatasetType {
	if d.Family != "plateau" {
		return plateauapi.PlateauDatasetType{}
	}

	name := d.Type
	if strings.HasPrefix(d.TypeEn, "urf_") {
		name = "都市計画決定情報モデル"
	}
	spec := d.Spec
	if isEx(d) {
		spec = "第3.0版"
	}

	year, _ := strconv.Atoi(d.Edition)
	return plateauapi.PlateauDatasetType{
		ID:              datasetTypeIDFrom(d),
		Name:            name,
		Code:            datasetTypeCodeFrom(d),
		Year:            year,
		Category:        plateauapi.DatasetTypeCategoryPlateau,
		PlateauSpecID:   plateauSpecMajorIDFrom(spec),
		PlateauSpecName: spec,
		Flood:           slices.Contains(floodingTypes, d.TypeEn),
	}
}

func relatedDatasetTypeFrom(d datacatalogv2.DataCatalogItem) plateauapi.RelatedDatasetType {
	if d.Family != "related" {
		return plateauapi.RelatedDatasetType{}
	}

	return plateauapi.RelatedDatasetType{
		ID:       datasetTypeIDFrom(d),
		Name:     d.Type,
		Code:     datasetTypeCodeFrom(d),
		Category: plateauapi.DatasetTypeCategoryRelated,
	}
}

func genericDatasetTypeFrom(d datacatalogv2.DataCatalogItem) plateauapi.GenericDatasetType {
	if d.Family != "generic" {
		return plateauapi.GenericDatasetType{}
	}

	if d.Category != "" {
		return plateauapi.GenericDatasetType{
			ID:       datasetTypeIDFrom(d),
			Name:     d.Category,
			Code:     datasetTypeCodeFrom(d),
			Category: plateauapi.DatasetTypeCategoryGeneric,
		}
	}

	return plateauapi.GenericDatasetType{
		ID:       datasetTypeIDFrom(d),
		Name:     "ユースケース",
		Code:     "usecase",
		Category: plateauapi.DatasetTypeCategoryGeneric,
	}
}

func groupsFrom(d datacatalogv2.DataCatalogItem) []string {
	if d.Group == "" {
		return nil
	}
	return strings.Split(d.Group, "/")
}

// cut string with sep from right
func cutRight(s, sep string) (string, string, bool) {
	i := strings.LastIndex(s, sep)
	if i < 0 {
		return s, "", false
	}
	return s[:i], s[i+len(sep):], true
}

func majorVersion(version string) string {
	v := specNumber(version)
	i := strings.Index(v, ".")
	if i < 0 {
		return version
	}
	return v[:i]
}

func isEx(d datacatalogv2.DataCatalogItem) bool {
	return strings.Contains(d.ID, "_ex_")
}
