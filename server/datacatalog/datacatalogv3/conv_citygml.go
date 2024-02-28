package datacatalogv3

import (
	"net/url"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"golang.org/x/exp/maps"
)

func toCityGMLs(all *AllData, regYear int) (map[plateauapi.ID]*plateauapi.CityGMLDataset, []string) {
	cities := all.City
	data := all.GeospatialjpDataItems
	cmsurl := all.CMSInfo.CMSURL

	dataMap := make(map[string]*plateauapi.CityGMLDataset)
	featureTypesMap := make(map[string][]string)

	for _, d := range data {
		if d.CityGML == "" || d.MaxLODContent == nil {
			continue
		}

		featureTypesMap[d.City] = citygmlFeatureTypes(d)
		dataMap[d.City] = &plateauapi.CityGMLDataset{
			URL:   d.CityGML,
			Items: toCityGMLItems(d),
			Admin: map[string]any{
				"maxlod": d.MaxLOD,
			},
		}
	}

	for _, city := range cities {
		if _, ok := dataMap[city.ID]; !ok {
			continue
		}

		if !slices.Contains(featureTypesMap[city.ID], "dem") {
			if dem := all.FindPlateauFeatureItemByCityID("dem", city.ID); dem != nil && dem.CityGML != "" {
				featureTypesMap[city.ID] = append(featureTypesMap[city.ID], "dem")
			}
		}

		dataMap[city.ID].ID = plateauapi.CityGMLDatasetIDFrom(plateauapi.AreaCode(city.CityCode))
		dataMap[city.ID].Year = city.YearInt()
		dataMap[city.ID].RegistrationYear = regYear
		dataMap[city.ID].PrefectureCode = plateauapi.AreaCode(plateauapi.AreaCode(city.CityCode).PrefectureCode())
		dataMap[city.ID].PrefectureID = plateauapi.NewID(dataMap[city.ID].PrefectureCode.String(), plateauapi.TypeArea)
		dataMap[city.ID].CityID = plateauapi.NewID(city.CityCode, plateauapi.TypeArea)
		dataMap[city.ID].CityCode = plateauapi.AreaCode(city.CityCode)
		dataMap[city.ID].FeatureTypes = featureTypesMap[city.ID]
		dataMap[city.ID].PlateauSpecMinorID = plateauapi.PlateauSpecIDFrom(city.Spec)
		dataMap[city.ID].Admin = newAdmin(city.ID, city.geospatialjpStage(), cmsurl, dataMap[city.ID].Admin)
	}

	res := make(map[plateauapi.ID]*plateauapi.CityGMLDataset)
	for _, v := range dataMap {
		res[v.ID] = v
	}

	return res, nil
}

func citygmlFeatureTypes(d *GeospatialjpDataItem) []string {
	types := map[string]struct{}{}
	for _, c := range d.MaxLODContent {
		if len(c) == 0 {
			continue
		}

		types[c[1]] = struct{}{}
	}

	res := maps.Keys(types)
	sort.Strings(res)
	return res
}

func toCityGMLItems(data *GeospatialjpDataItem) (res []*plateauapi.CityGMLDatasetItem) {
	for _, d := range data.MaxLODContent {
		if len(d) < 4 || len(d[0]) == 0 || !isNumeric(rune(d[0][0])) {
			continue
		}

		// code,type,maxLod,file
		maxlod, _ := strconv.Atoi(d[2])
		item := &plateauapi.CityGMLDatasetItem{
			MeshCode: d[0],
			TypeCode: d[1],
			MaxLod:   maxlod,
			URL:      citygmlItemURLFrom(data.CityGML, d[3], d[1]),
		}

		if item.URL == "" {
			continue
		}

		res = append(res, item)
	}

	return
}

func citygmlItemURLFrom(base, p, typeCode string) string {
	b := path.Base(base)
	base = strings.TrimSuffix(base, b)
	u, _ := url.JoinPath(base, nameWithoutExt(b), "udx", typeCode, p)
	return u
}

func isNumeric(r rune) bool {
	return strings.ContainsRune("0123456789", r)
}
