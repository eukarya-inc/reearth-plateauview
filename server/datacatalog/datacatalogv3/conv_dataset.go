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

func toRiverAdminName(a plateauapi.RiverAdmin) string {
	switch a {
	case plateauapi.RiverAdminNational:
		return "国"
	case plateauapi.RiverAdminPrefecture:
		return "都道府県"
	}
	return ""
}

func textureFrom(notexture *bool) *plateauapi.Texture {
	if notexture == nil {
		return nil
	}
	if *notexture {
		return lo.ToPtr(plateauapi.TextureNone)
	}
	return lo.ToPtr(plateauapi.TextureTexture)
}

func datasetFormatFromOrDetect(f string, url string) plateauapi.DatasetFormat {
	if f != "" {
		return datasetFormatFrom(f)
	}
	return detectDatasetFormatFromURL(url)
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

func detectDatasetFormatFromURL(url string) plateauapi.DatasetFormat {
	name := strings.ToLower(nameFromURL(url))

	switch {
	case strings.HasSuffix(name, ".geojson"):
		return plateauapi.DatasetFormatGeojson
	case strings.HasSuffix(name, ".czml"):
		return plateauapi.DatasetFormatCzml
	case strings.HasSuffix(name, "{z}/{x}/{y}.pbf"):
		fallthrough
	case strings.HasSuffix(name, ".mvt"):
		return plateauapi.DatasetFormatMvt
	case name == "tileset.json":
		return plateauapi.DatasetFormatCesium3dtiles
	case strings.HasSuffix(name, ".csv"):
		return plateauapi.DatasetFormatCSV
	case strings.HasSuffix(name, ".gltf"):
		return plateauapi.DatasetFormatGltf
	case strings.HasSuffix(name, "{z}/{x}/{y}.png"):
		return plateauapi.DatasetFormatTiles
	}

	return ""
}

func standardItemID(name string, area plateauapi.Area, ex string) string {
	if ex != "" {
		ex = fmt.Sprintf("_%s", ex)
	}
	return fmt.Sprintf("%s_%s%s", area.GetCode(), name, ex)
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

func adminFrom(cityItem *CityItem, cmsurl string, ft string) any {
	var stage stage
	if ft == "related" {
		stage = cityItem.relatedStage()
	} else {
		stage = cityItem.plateauStage(ft)
	}

	return newAdmin(cityItem.ID, stage, cmsurl)
}

func newAdmin(id string, stage stage, cmsurl string) any {
	a := map[string]any{}

	if cmsurl != "" && id != "" {
		a["cmsUrl"] = cmsurl + id
	}

	if stage != stageGA {
		if stage == "" {
			stage = stageAlpha
		}
		a["stage"] = string(stage)
	}

	if len(a) == 0 {
		return nil
	}

	return a
}
