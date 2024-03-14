package datacatalogv3

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
)

func toCityGMLs(all *AllData, regYear int) (map[plateauapi.ID]*plateauapi.CityGMLDataset, []string) {
	cities := all.City
	data := all.GeospatialjpDataItems
	cmsurl := all.CMSInfo.CMSURL

	dataMap := make(map[string]*plateauapi.CityGMLDataset)

	for _, d := range data {
		if d.CityGML == "" || d.MaxLOD == "" {
			continue
		}

		dataMap[d.City] = &plateauapi.CityGMLDataset{
			URL: d.CityGML,
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
		dataMap[city.ID].RegistrationYear = regYear
		dataMap[city.ID].PrefectureCode = plateauapi.AreaCode(plateauapi.AreaCode(city.CityCode).PrefectureCode())
		dataMap[city.ID].PrefectureID = plateauapi.NewID(dataMap[city.ID].PrefectureCode.String(), plateauapi.TypePrefecture)
		dataMap[city.ID].CityID = plateauapi.NewID(city.CityCode, plateauapi.TypeCity)
		dataMap[city.ID].CityCode = plateauapi.AreaCode(city.CityCode)
		dataMap[city.ID].FeatureTypes = all.FeatureTypesOf(city.ID)
		dataMap[city.ID].PlateauSpecMinorID = plateauapi.PlateauSpecIDFrom(city.Spec)
		dataMap[city.ID].Admin = newAdmin(city.ID, city.SDKStage(), cmsurl, dataMap[city.ID].Admin)
	}

	res := make(map[plateauapi.ID]*plateauapi.CityGMLDataset)
	for _, v := range dataMap {
		res[v.ID] = v
	}

	return res, nil
}
