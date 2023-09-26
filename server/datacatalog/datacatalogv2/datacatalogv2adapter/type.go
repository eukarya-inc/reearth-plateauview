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
		ID:          id,
		Name:        d.Name,
		Subname:     nil,
		Description: lo.ToPtr(d.Description),
		AreaID:      areaIDFrom(d),
		Year:        d.Year,
		TypeID:      plateauapi.NewID(d.TypeEn, "type"),
		Groups:      strings.Split(d.Group, "/"),
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
	if d.Category != "plateau" || !slices.Contains(floodingTypes, d.TypeEn) {
		return plateauapi.PlateauFloodingDataset{}, false
	}

	var subname *string
	var river *plateauapi.River

	if d.TypeEn == "fld" {
		var admin plateauapi.RiverAdmin
		if strings.Contains(d.Name, "（国管理区間）") {
			admin = plateauapi.RiverAdminGovernment
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
		ID:          id,
		Name:        d.Name,
		Subname:     subname,
		Description: lo.ToPtr(d.Description),
		AreaID:      areaIDFrom(d),
		Year:        d.Year,
		TypeID:      plateauapi.NewID(d.TypeEn, "type"),
		Groups:      strings.Split(d.Group, "/"),
		River:       river,
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
	if d.Category != "related" {
		return plateauapi.RelatedDataset{}, false
	}

	id := datasetIDFrom(d)
	return plateauapi.RelatedDataset{
		ID:          id,
		Name:        d.Name,
		Subname:     nil,
		Description: lo.ToPtr(d.Description),
		AreaID:      areaIDFrom(d),
		Year:        d.Year,
		TypeID:      plateauapi.NewID(d.TypeEn, "type"),
		Groups:      strings.Split(d.Group, "/"),
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
	if d.Category == "plateau" || d.Category == "related" {
		return plateauapi.GenericDataset{}, false
	}

	id := datasetIDFrom(d)
	return plateauapi.GenericDataset{
		ID:          datasetIDFrom(d),
		Name:        d.Name,
		Subname:     nil,
		Description: lo.ToPtr(d.Description),
		AreaID:      areaIDFrom(d),
		Year:        d.Year,
		TypeID:      plateauapi.NewID(d.TypeEn, "type"),
		Groups:      strings.Split(d.Group, "/"),
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

func areaIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	if id := wardIDFrom(d); id != "" {
		return id
	}

	if id := cityIDFrom(d); id != "" {
		return id
	}

	return prefectureIDFrom(d)
}

func prefectureIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	return plateauapi.NewID(d.PrefCode, plateauapi.TypePrefecture)
}

func cityIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	if d.WardCode == "" {
		return ""
	}
	return plateauapi.NewID(d.CityCode, plateauapi.TypeMunicipality)
}

func wardIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	if d.WardCode == "" {
		return ""
	}
	return plateauapi.NewID(d.WardCode, plateauapi.TypeMunicipality)
}

func datasetIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	return newDatasetID(d.ID)
}

func newDatasetID(id string) plateauapi.ID {
	return plateauapi.NewID(id, plateauapi.TypeDataset)
}

func datasetTypeIDFrom(d datacatalogv2.DataCatalogItem) plateauapi.ID {
	return plateauapi.NewID(d.TypeEn, plateauapi.TypeDatasetType)
}

func prefectureFrom(d datacatalogv2.DataCatalogItem) plateauapi.Prefecture {
	if d.PrefCode == "" {
		return plateauapi.Prefecture{}
	}

	return plateauapi.Prefecture{
		ID:   prefectureIDFrom(d),
		Code: plateauapi.AreaCode(d.PrefCode),
		Name: d.Pref,
	}
}

func municipalityFrom(d datacatalogv2.DataCatalogItem) plateauapi.Municipality {
	if d.CityCode == "" {
		return plateauapi.Municipality{}
	}

	return plateauapi.Municipality{
		ID:   cityIDFrom(d),
		Code: plateauapi.AreaCode(d.CityCode),
		Name: d.City,
	}
}

func wardMunicipalityFrom(d datacatalogv2.DataCatalogItem) plateauapi.Municipality {
	if d.WardCode == "" {
		return plateauapi.Municipality{}
	}

	return plateauapi.Municipality{
		ID:   wardIDFrom(d),
		Code: plateauapi.AreaCode(d.WardCode),
		Name: d.Ward,
	}
}

func plateauTypeFrom(d datacatalogv2.DataCatalogItem) plateauapi.PlateauDatasetType {
	if d.Category != "plateau" {
		return plateauapi.PlateauDatasetType{}
	}

	return plateauapi.PlateauDatasetType{
		ID:          datasetTypeIDFrom(d),
		Name:        d.Type,
		Code:        d.TypeEn,
		Year:        d.Year,
		EnglishName: d.TypeEn,
		// PlateauSpec: d.Spec, // TODO
		Category: plateauapi.DatasetTypeCategoryPlateau,
	}
}

func relatedTypeFrom(d datacatalogv2.DataCatalogItem) plateauapi.RelatedDatasetType {
	return plateauapi.RelatedDatasetType{
		ID:          datasetTypeIDFrom(d),
		Name:        d.Type,
		Code:        d.TypeEn,
		EnglishName: d.TypeEn,
		// PlateauSpec: d.Spec, // TODO
		Category: plateauapi.DatasetTypeCategoryRelated,
	}
}

func genericTypeFrom(d datacatalogv2.DataCatalogItem) plateauapi.GenericDatasetType {
	return plateauapi.GenericDatasetType{
		ID:          datasetTypeIDFrom(d),
		Name:        d.Type,
		Code:        d.TypeEn,
		EnglishName: d.TypeEn,
		// PlateauSpec: d.Spec, // TODO
		Category: plateauapi.DatasetTypeCategoryGeneric,
	}
}

func specFrom(d datacatalogv2.DataCatalogItem) plateauapi.PlateauSpec {
	return plateauapi.PlateauSpec{
		ID:   plateauapi.NewID(d.Name, plateauapi.TypePlateauSpec),
		Name: d.Spec,
		Year: d.Year,
	}
}
