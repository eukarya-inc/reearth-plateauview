package datacatalogv2adapter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogutil"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

var floodingTypes = []string{"fld", "htd", "tnm", "ifld"}

func plateauDatasetFrom(d datacatalogv2.DataCatalogItem) (plateauapi.PlateauDataset, bool) {
	if d.Category != "plateau" || slices.Contains(floodingTypes, d.TypeEn) {
		return plateauapi.PlateauDataset{}, false
	}

	id := datasetIDFrom(d)
	return plateauapi.PlateauDataset{
		ID:           id,
		Name:         d.Name,
		Subname:      nil,
		Description:  lo.ToPtr(d.Description),
		PrefectureID: prefectureIDFrom(d),
		CityID:       cityIDFrom(d),
		WardID:       wardIDFrom(d),
		Year:         d.Year,
		TypeID:       plateauapi.NewID(d.TypeEn, "type"),
		Groups:       strings.Split(d.Group, "/"),
		Data: lo.Map(d.ConfigItems(), func(c datacatalogutil.DataCatalogItemConfigItem, _ int) *plateauapi.PlateauDatasetItem {
			return plateauDatasetItemFrom(c, d.ID, id)
		}),
	}, true
}

func plateauDatasetItemFrom(c datacatalogutil.DataCatalogItemConfigItem, parent string, parentID plateauapi.ID) *plateauapi.PlateauDatasetItem {
	var lod *float64
	if strings.HasPrefix(c.Name, "LOD") {
		l, _, _ := strings.Cut(c.Name[3:], "（")
		lodf, err := strconv.ParseFloat(l, 64)
		if err == nil {
			lod = &lodf
		}
	}

	var texture *plateauapi.Texture
	if strings.Contains(c.Name, "（テクスチャなし）") {
		texture = lo.ToPtr(plateauapi.TextureNone)
	}

	return &plateauapi.PlateauDatasetItem{
		ID:       plateauapi.NewID(fmt.Sprintf("%s:%s", parent, c.Name), plateauapi.TypeDatasetItem),
		Name:     c.Name,
		URL:      c.URL,
		Format:   datasetFormatFrom(c.Type),
		Layers:   c.Layers,
		ParentID: parentID,
		Lod:      lod,
		Texture:  texture,
	}
}

func plateauFloodingDatasetFrom(d datacatalogv2.DataCatalogItem) (plateauapi.PlateauFloodingDataset, bool) {
	if categoryFrom(d) != plateauapi.DatasetTypeCategoryPlateau || !slices.Contains(floodingTypes, d.TypeEn) {
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

		names := strings.Split(d.Name, " ")
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

	id := datasetIDFrom(d)
	return plateauapi.PlateauFloodingDataset{
		ID:           id,
		Name:         d.Name,
		Subname:      subname,
		Description:  lo.ToPtr(d.Description),
		PrefectureID: prefectureIDFrom(d),
		CityID:       cityIDFrom(d),
		WardID:       wardIDFrom(d),
		Year:         d.Year,
		TypeID:       plateauapi.NewID(d.TypeEn, "type"),
		Groups:       strings.Split(d.Group, "/"),
		River:        river,
		Data: lo.Map(d.ConfigItems(), func(c datacatalogutil.DataCatalogItemConfigItem, _ int) *plateauapi.PlateauFloodingDatasetItem {
			return plateauFloodingDatasetItemFrom(c, d.ID, id)
		}),
	}, true
}

func plateauFloodingDatasetItemFrom(c datacatalogutil.DataCatalogItemConfigItem, parent string, parentID plateauapi.ID) *plateauapi.PlateauFloodingDatasetItem {
	var floodingScale plateauapi.FloodingScale
	if strings.Contains(c.Name, "想定最大規模") {
		floodingScale = plateauapi.FloodingScaleExpectedMaximum
	} else if strings.Contains(c.Name, "計画規模") {
		floodingScale = plateauapi.FloodingScalePlanned
	}

	return &plateauapi.PlateauFloodingDatasetItem{
		ID:            plateauapi.NewID(fmt.Sprintf("%s:%s", parent, c.Name), plateauapi.TypeDatasetItem),
		Name:          c.Name,
		URL:           c.URL,
		Format:        datasetFormatFrom(c.Type),
		Layers:        c.Layers,
		ParentID:      parentID,
		FloodingScale: floodingScale,
	}
}

func relatedDatasetFrom(d datacatalogv2.DataCatalogItem) (plateauapi.RelatedDataset, bool) {
	if categoryFrom(d) != plateauapi.DatasetTypeCategoryRelated {
		return plateauapi.RelatedDataset{}, false
	}

	id := datasetIDFrom(d)
	return plateauapi.RelatedDataset{
		ID:           id,
		Name:         d.Name,
		Subname:      nil,
		Description:  lo.ToPtr(d.Description),
		PrefectureID: prefectureIDFrom(d),
		CityID:       cityIDFrom(d),
		WardID:       wardIDFrom(d),
		Year:         d.Year,
		TypeID:       plateauapi.NewID(d.TypeEn, "type"),
		Groups:       strings.Split(d.Group, "/"),
		Data: lo.Map(d.ConfigItems(), func(c datacatalogutil.DataCatalogItemConfigItem, _ int) *plateauapi.RelatedDatasetItem {
			return &plateauapi.RelatedDatasetItem{
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

func genericDatasetFrom(d datacatalogv2.DataCatalogItem) (plateauapi.GenericDataset, bool) {
	if cat := categoryFrom(d); cat == plateauapi.DatasetTypeCategoryPlateau || cat == plateauapi.DatasetTypeCategoryRelated {
		return plateauapi.GenericDataset{}, false
	}

	id := datasetIDFrom(d)
	return plateauapi.GenericDataset{
		ID:           datasetIDFrom(d),
		Name:         d.Name,
		Subname:      nil,
		Description:  lo.ToPtr(d.Description),
		PrefectureID: prefectureIDFrom(d),
		CityID:       cityIDFrom(d),
		WardID:       wardIDFrom(d),
		Year:         d.Year,
		TypeID:       plateauapi.NewID(d.TypeEn, "type"),
		Groups:       strings.Split(d.Group, "/"),
		Data: lo.Map(d.ConfigItems(), func(c datacatalogutil.DataCatalogItemConfigItem, _ int) *plateauapi.GenericDatasetItem {
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
	if d.WardCode == "" {
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
	if d.WardCode == "" {
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

func datasetIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	return newDatasetID(d.ID)
}

func newDatasetID(id string) plateauapi.ID {
	return plateauapi.NewID(id, plateauapi.TypeDataset)
}

func datasetTypeIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	if d.Family == "generic" && d.Category != "" {
		return plateauapi.NewID(fmt.Sprintf("%s:%s", d.Edition, d.Category), plateauapi.TypeDatasetType)
	}
	return plateauapi.NewID(fmt.Sprintf("%s:usecase", d.Edition), plateauapi.TypeDatasetType)
}

func specIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	return plateauapi.NewID(specNumber(d.Spec), plateauapi.TypePlateauSpec)
}

func specNumber(spec string) string {
	return strings.TrimSuffix(strings.TrimPrefix(spec, "第"), "版")
}

func prefectureFrom(d datacatalogv2.DataCatalogItem) plateauapi.Prefecture {
	if d.PrefCode == "" {
		return plateauapi.Prefecture{}
	}

	return plateauapi.Prefecture{
		ID:   prefectureIDFrom(d),
		Code: prefectureCodeFrom(d),
		Name: d.Pref,
	}
}

func cityFrom(d datacatalogv2.DataCatalogItem) plateauapi.City {
	id, code := cityIDFrom(d), cityCodeFrom(d)
	if id == nil || code == nil {
		return plateauapi.City{}
	}

	return plateauapi.City{
		ID:             *id,
		Code:           *code,
		Name:           d.City,
		PrefectureID:   prefectureIDFrom(d),
		PrefectureCode: prefectureCodeFrom(d),
	}
}

func wardFrom(d datacatalogv2.DataCatalogItem) plateauapi.Ward {
	id, code := wardIDFrom(d), wardCodeFrom(d)
	if id == nil || code == nil {
		return plateauapi.Ward{}
	}

	cityid, citycode := cityIDFrom(d), cityCodeFrom(d)
	if cityid == nil || citycode == nil {
		return plateauapi.Ward{}
	}

	return plateauapi.Ward{
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
		Code:          d.TypeEn,
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
		Code:        d.TypeEn,
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
			Code:     d.Category,
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

func specFrom(d datacatalogv2.DataCatalogItem) plateauapi.PlateauSpec {
	if d.Spec == "" {
		return plateauapi.PlateauSpec{}
	}
	return plateauapi.PlateauSpec{
		ID:   specIDFrom(d),
		Name: d.Spec,
		Year: d.Year,
	}
}

func categoryFrom(d datacatalogv2.DataCatalogItem) plateauapi.DatasetTypeCategory {
	switch d.Family {
	case "plateau":
		return plateauapi.DatasetTypeCategoryPlateau
	case "related":
		return plateauapi.DatasetTypeCategoryRelated
	}
	return plateauapi.DatasetTypeCategoryGeneric
}
