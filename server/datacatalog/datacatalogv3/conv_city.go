package datacatalogv3

import "github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"

func (city *CityItem) ToPrefecture() *plateauapi.Prefecture {
	prefCode := city.CityCode[:2]
	if prefCode == "" {
		return nil
	}

	return &plateauapi.Prefecture{
		ID:   plateauapi.NewID(prefCode, plateauapi.TypeArea),
		Name: city.Prefecture,
		Code: plateauapi.AreaCode(prefCode),
		Type: plateauapi.AreaTypePrefecture,
	}
}

func (city *CityItem) ToCity() *plateauapi.City {
	return &plateauapi.City{
		ID:   plateauapi.NewID(city.CityCode, plateauapi.TypeArea),
		Name: city.CityName,
		Code: plateauapi.AreaCode(city.CityCode),
		Type: plateauapi.AreaTypeCity,
	}
}
