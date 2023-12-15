package datacatalogv3

import (
	"fmt"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

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
