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

func plateauDatasetFrom(d datacatalogv2.DataCatalogItem) (plateauapi.PlateauDataset, bool) {
	if d.Family != "plateau" || slices.Contains(floodingTypes, d.TypeEn) {
		return plateauapi.PlateauDataset{}, false
	}

	id := datasetIDFrom(d, nil)
	return plateauapi.PlateauDataset{
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
	} else if format == plateauapi.DatasetFormatCesium3DTiles {
		texture = lo.ToPtr(plateauapi.TextureTexture)
	}

	return &plateauapi.PlateauDatasetItem{
		ID:       plateauapi.NewID(fmt.Sprintf("%s_%s", parentID.ID(), c.Name), plateauapi.TypeDatasetItem),
		Name:     c.Name,
		URL:      c.URL,
		Format:   format,
		Layers:   c.Layers,
		ParentID: parentID,
		Lod:      lod,
		Texture:  texture,
	}
}

var reBrackets = regexp.MustCompile(`（[^（]*）`)

func plateauFloodingDatasetFrom(d datacatalogv2.DataCatalogItem) (plateauapi.PlateauFloodingDataset, bool) {
	if d.Family != "plateau" || !slices.Contains(floodingTypes, d.TypeEn) {
		return plateauapi.PlateauFloodingDataset{}, false
	}

	var subname *string
	var river *plateauapi.River

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

	id := datasetIDFrom(d, subname)
	return plateauapi.PlateauFloodingDataset{
		ID:             id,
		Name:           d.Name,
		Subname:        subname,
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
		River:          river,
		Items: lo.Map(d.MainOrConfigItems(), func(c datacatalogutil.DataCatalogItemConfigItem, _ int) *plateauapi.PlateauFloodingDatasetItem {
			return plateauFloodingDatasetItemFrom(c, id)
		}),
	}, true
}

func plateauFloodingDatasetItemFrom(c datacatalogutil.DataCatalogItemConfigItem, parentID plateauapi.ID) *plateauapi.PlateauFloodingDatasetItem {
	var id string
	var floodingScale plateauapi.FloodingScale
	if strings.Contains(c.Name, "想定最大規模") || strings.Contains(c.Name, "L2") {
		floodingScale = plateauapi.FloodingScaleExpectedMaximum
		id = "l2"
	} else if strings.Contains(c.Name, "計画規模") || strings.Contains(c.Name, "L1") {
		floodingScale = plateauapi.FloodingScalePlanned
		id = "l1"
	}

	return &plateauapi.PlateauFloodingDatasetItem{
		ID:            plateauapi.NewID(fmt.Sprintf("%s_%s", parentID.ID(), id), plateauapi.TypeDatasetItem),
		Name:          c.Name,
		URL:           c.URL,
		Format:        datasetFormatFrom(c.Type),
		Layers:        c.Layers,
		ParentID:      parentID,
		FloodingScale: floodingScale,
	}
}

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
		ID:           id,
		Name:         d.Name,
		Subname:      nil,
		Description:  lo.ToPtr(d.Description),
		PrefectureID: prefectureIDFrom(d),
		CityID:       cityIDFrom(d),
		WardID:       wardIDFrom(d),
		Year:         d.Year,
		TypeID:       datasetTypeIDFrom(d),
		TypeCode:     datasetTypeCodeFrom(d),
		Groups:       groupsFrom(d),
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
	switch f {
	case "geojson":
		return plateauapi.DatasetFormatGeoJSON
	case "3dtiles":
		return plateauapi.DatasetFormatCesium3DTiles
	case "czml":
		return plateauapi.DatasetFormatCzml
	case "gtfs":
	case "gtfs-realtime":
		return plateauapi.DatasetFormatGTFSRelatime
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
		isEx := strings.Contains(d.ID, "_ex_")

		if isFlood || d.TypeEn == "gen" || isEx {
			if isEx {
				typeCode = "ex"
			}

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
		// TODO: refer to a major version of the plateau spec
		return plateauapi.NewID(fmt.Sprintf("%s_%s", code, d.Spec), plateauapi.TypeDatasetType)
	}
	return plateauapi.NewID(code, plateauapi.TypeDatasetType)
}

func datasetTypeCodeFrom(d datacatalogv2.DataCatalogItem) string {
	if d.Family == "plateau" {
		return d.TypeEn
	}
	if d.Family == "related" {
		return d.TypeEn
	}
	if d.Family == "generic" && d.Category != "" {
		return d.Category
	}
	return usecaseID
}

func specIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	return plateauapi.NewID(specNumber(d.Spec), plateauapi.TypePlateauSpec)
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

func plateauTypeFrom(d datacatalogv2.DataCatalogItem) plateauapi.PlateauDatasetType {
	if d.Family != "plateau" {
		return plateauapi.PlateauDatasetType{}
	}

	year, _ := strconv.Atoi(d.Edition)
	return plateauapi.PlateauDatasetType{
		ID:            datasetTypeIDFrom(d),
		Name:          d.Type,
		Code:          datasetTypeCodeFrom(d),
		Year:          year,
		EnglishName:   d.TypeEn,
		Category:      plateauapi.DatasetTypeCategoryPlateau,
		PlateauSpecID: specIDFrom(d),
		Flood:         slices.Contains(floodingTypes, d.TypeEn),
	}
}

func relatedTypeFrom(d datacatalogv2.DataCatalogItem) plateauapi.RelatedDatasetType {
	if d.Family != "related" {
		return plateauapi.RelatedDatasetType{}
	}

	return plateauapi.RelatedDatasetType{
		ID:          datasetTypeIDFrom(d),
		Name:        d.Type,
		Code:        datasetTypeCodeFrom(d),
		EnglishName: d.TypeEn,
		Category:    plateauapi.DatasetTypeCategoryRelated,
	}
}

func genericTypeFrom(d datacatalogv2.DataCatalogItem) plateauapi.GenericDatasetType {
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
		ID:          datasetTypeIDFrom(d),
		Name:        "ユースケース",
		Code:        "usecase",
		EnglishName: "usecase",
		Category:    plateauapi.DatasetTypeCategoryGeneric,
	}
}

func specFrom(d datacatalogv2.DataCatalogItem) *plateauapi.PlateauSpec {
	if d.Spec == "" {
		return nil
	}
	return &plateauapi.PlateauSpec{
		ID:   specIDFrom(d),
		Name: d.Spec,
		Year: d.Year,
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
