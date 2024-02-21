package datacatalogv3

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

func toCityGMLs(cities []*CityItem, data []*GeospatialjpDataItem, dts []plateauapi.DatasetType, cmsurl string) (map[plateauapi.ID]*plateauapi.CityGMLDataset, []string) {
	dataMap := make(map[string]*plateauapi.CityGMLDataset)
	for _, d := range data {
		if d.CityGML == "" || d.MaxLODContent == nil {
			continue
		}

		dataMap[d.City] = &plateauapi.CityGMLDataset{
			URL:   d.CityGML,
			Items: toCityGMLItems(d, dts),
			Admin: map[string]any{
				"maxlod": d.MaxLOD,
			},
		}
	}

	for _, city := range cities {
		if _, ok := dataMap[city.ID]; !ok {
			continue
		}

		dataMap[city.ID].ID = plateauapi.CityGMLDatasetIDFrom(plateauapi.AreaCode(city.CityCode))
		dataMap[city.ID].Year = city.YearInt()
		dataMap[city.ID].PrefectureCode = plateauapi.AreaCode(plateauapi.AreaCode(city.CityCode).PrefectureCode())
		dataMap[city.ID].PrefectureID = plateauapi.NewID(dataMap[city.ID].PrefectureCode.String(), plateauapi.TypeArea)
		dataMap[city.ID].CityID = plateauapi.NewID(city.CityCode, plateauapi.TypeArea)
		dataMap[city.ID].CityCode = plateauapi.AreaCode(city.CityCode)
		dataMap[city.ID].PlateauSpecMinorID = plateauapi.PlateauSpecIDFrom(city.Spec)
		dataMap[city.ID].Admin = newAdmin(city.ID, city.geospatialjpStage(), cmsurl, dataMap[city.ID].Admin)
	}

	res := make(map[plateauapi.ID]*plateauapi.CityGMLDataset)
	for _, v := range dataMap {
		res[v.ID] = v
	}

	return res, nil
}

func toCityGMLItems(data *GeospatialjpDataItem, dts []plateauapi.DatasetType) (res []*plateauapi.CityGMLDatasetItem) {
	for _, d := range data.MaxLODContent {
		if len(d) < 4 || len(d[0]) == 0 || !isNumeric(rune(d[0][0])) {
			continue
		}

		// code,type,maxLod,file
		maxlod, _ := strconv.Atoi(d[2])
		ty, _ := lo.Find(dts, func(dt plateauapi.DatasetType) bool {
			return dt.GetCode() == d[1]
		})
		if ty == nil {
			continue
		}

		item := &plateauapi.CityGMLDatasetItem{
			MeshCode: d[0],
			TypeCode: ty.GetCode(),
			TypeID:   ty.GetID(),
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