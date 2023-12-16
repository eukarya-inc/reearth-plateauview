package datacatalogv3

import (
	"fmt"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

func riverAdminFrom(admin string) *plateauapi.RiverAdmin {
	switch admin {
	case "国":
		fallthrough
	case "natl":
		return lo.ToPtr(plateauapi.RiverAdminNational)
	case "都道府県":
		fallthrough
	case "pref":
		return lo.ToPtr(plateauapi.RiverAdminPrefecture)
	}
	return nil
}

func textureFrom(texture *bool) *plateauapi.Texture {
	if texture == nil {
		return nil
	}
	if !*texture {
		return lo.ToPtr(plateauapi.TextureNone)
	}
	return lo.ToPtr(plateauapi.TextureTexture)
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

func standardItemName(dtname, subname string, area plateauapi.Area) string {
	space := ""
	if subname != "" {
		space = " "
	}
	return fmt.Sprintf("%s%s%s（%s）", dtname, space, subname, area.GetName())
}

func layerNamesFrom(layer string) []string {
	if layer == "" {
		return nil
	}

	return lo.Map(strings.Split(layer, ","), func(s string, _ int) string {
		return strings.TrimSpace(s)
	})
}
