package datacatalogv3

import (
	"fmt"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

func (i *RelatedItem) ToDatasets(pref *plateauapi.Prefecture, city *plateauapi.City, featureTypes []plateauapi.DatasetType, year int) (res []plateauapi.Dataset) {
	prefID, cityID, prefCode, cityCode := areaInfo(pref, city)
	if prefID == nil || cityID == nil || prefCode == nil || cityCode == nil {
		return nil
	}

	for _, ft := range featureTypes {
		ftname, ftcode := ft.GetName(), ft.GetCode()
		asset := i.ConvertedAssets[ftcode]
		format := plateauapi.DatasetFormatCzml
		if len(asset) == 0 {
			asset = i.Assets[ftcode]
			format = plateauapi.DatasetFormatGeojson
		}
		if len(asset) == 0 {
			continue
		}

		sid := standardItemID(ftcode, city)
		id := plateauapi.NewID(sid, plateauapi.TypeDataset)
		res = append(res, plateauapi.RelatedDataset{
			ID:             id,
			Name:           standardItemName(ftname, city),
			Description:    toPtrIfPresent(i.Desc),
			Year:           year,
			PrefectureID:   prefID,
			PrefectureCode: prefCode,
			CityID:         cityID,
			CityCode:       cityCode,
			TypeID:         ft.GetID(),
			TypeCode:       ftcode,
			Items: []*plateauapi.RelatedDatasetItem{
				{
					ID:       plateauapi.NewID(sid, plateauapi.TypeDatasetItem),
					Format:   format,
					Name:     ftname,
					URL:      asset[0], // TODO
					ParentID: id,
				},
			},
		})
	}

	return
}

func (i *GenericItem) ToDatasets(pref *plateauapi.Prefecture, city *plateauapi.City, dts []plateauapi.DatasetType, year int) []plateauapi.Dataset {
	id := plateauapi.NewID(i.ID, plateauapi.TypeDataset)
	prefID, cityID, prefCode, cityCode := areaInfo(pref, city)
	if prefID == nil || cityID == nil || prefCode == nil || cityCode == nil {
		return nil
	}

	dt, _ := lo.Find(dts, func(dt plateauapi.DatasetType) bool {
		return dt.GetName() == i.Category
	})
	if dt == nil {
		return nil
	}

	items := dropNil(lo.Map(i.Data, func(datum GenericItemDataset, ind int) *plateauapi.GenericDatasetItem {
		if datum.Data == "" {
			return nil
		}

		var inds string
		if len(i.Data) > 1 {
			inds = fmt.Sprintf(" %d", ind+1)
		}

		return &plateauapi.GenericDatasetItem{
			ID:       plateauapi.NewID(datum.ID, plateauapi.TypeDatasetItem),
			Name:     firstNonEmptyValue(datum.Name, fmt.Sprintf("%s%s", i.Name, inds)),
			URL:      datum.Data,
			Format:   datasetFormatFrom(datum.DataFormat),
			Layers:   layerNamesFrom(datum.LayerName),
			ParentID: id,
		}
	}))

	if len(items) == 0 {
		return nil
	}

	res := plateauapi.GenericDataset{
		ID:             id,
		Name:           i.Name,
		Description:    toPtrIfPresent(i.Desc),
		Year:           year,
		PrefectureID:   prefID,
		PrefectureCode: prefCode,
		CityID:         cityID,
		CityCode:       cityCode,
		TypeID:         dt.GetID(),
		TypeCode:       dt.GetCode(),
		Items:          items,
	}

	return []plateauapi.Dataset{&res}
}

// func assetsToWards(assets []string) []string {
// }

func areaInfo(pref *plateauapi.Prefecture, city *plateauapi.City) (prefID, cityID *plateauapi.ID, prefCode, cityCode *plateauapi.AreaCode) {
	if pref != nil {
		prefID = lo.ToPtr(pref.ID)
		prefCode = lo.ToPtr(pref.Code)
	}
	if city != nil {
		cityID = lo.ToPtr(city.ID)
		cityCode = lo.ToPtr(city.Code)
	}
	return
}

func datasetFormatFrom(f string) plateauapi.DatasetFormat {
	switch strings.ToLower(f) {
	case "geojson":
		return plateauapi.DatasetFormatGeojson
	case "3dtiles":
		fallthrough
	case "3d tiles":
		return plateauapi.DatasetFormatCesium3dtiles
	case "czml":
		return plateauapi.DatasetFormatCzml
	case "gtfs":
		fallthrough
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

func standardItemID(name string, area plateauapi.Area) string {
	return fmt.Sprintf("%s_%s", area.GetCode(), name)
}

func standardItemName(name string, area *plateauapi.City) string {
	return fmt.Sprintf("%s (%s)", name, area.Name)
}

func layerNamesFrom(layer string) []string {
	if layer == "" {
		return nil
	}

	return lo.Map(strings.Split(layer, ","), func(s string, _ int) string {
		return strings.TrimSpace(s)
	})
}
