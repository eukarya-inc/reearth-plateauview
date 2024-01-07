package datacatalogv3

import (
	"fmt"
	"sort"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

type internalContext struct {
	years         map[int]struct{}
	cityItems     map[string]*CityItem
	prefs         map[string]*plateauapi.Prefecture
	cities        map[string]*plateauapi.City
	wards         map[string][]*plateauapi.Ward
	plateauCMSURL string
	relatedCMSURL string
	genericCMSURL string
}

func newInternalContext() *internalContext {
	return &internalContext{
		years:     map[int]struct{}{},
		cityItems: map[string]*CityItem{},
		prefs:     map[string]*plateauapi.Prefecture{},
		cities:    map[string]*plateauapi.City{},
		wards:     map[string][]*plateauapi.Ward{},
	}
}

func (c *internalContext) CityItem(id string) *CityItem {
	return c.cityItems[id]
}

func (c *internalContext) Wards(cityCode string) []*plateauapi.Ward {
	return c.wards[cityCode]
}

func (c *internalContext) HasPref(prefCode string) bool {
	_, ok := c.prefs[prefCode]
	return ok
}

func (c *internalContext) HasCity(prefCode string) bool {
	_, ok := c.prefs[prefCode]
	return ok
}

func (c *internalContext) Years() []int {
	res := make([]int, 0, len(c.years))
	for y := range c.years {
		res = append(res, y)
	}
	sort.Ints(res)
	return res
}

func (c *internalContext) Add(cityItem *CityItem, pref *plateauapi.Prefecture, city *plateauapi.City) {
	c.cityItems[cityItem.ID] = cityItem
	c.prefs[pref.Code.String()] = pref
	c.cities[city.Code.String()] = city

	if y := cityItem.YearInt(); y != 0 {
		c.years[y] = struct{}{}
	}
}

func (c *internalContext) AddWards(wards []*plateauapi.Ward) {
	for _, w := range wards {
		cityCode := w.CityCode.String()
		c.wards[cityCode] = append(c.wards[cityCode], w)
	}
}

func (c *internalContext) SetURL(t, cmsurl, ws, prj, modelID string) {
	if cmsurl == "" || ws == "" || prj == "" || modelID == "" {
		return
	}

	url := fmt.Sprintf("%s/workspace/%s/project/%s/content/%s/details/", cmsurl, ws, prj, modelID)

	switch t {
	case "plateau":
		c.plateauCMSURL = url
	case "related":
		c.plateauCMSURL = url
	case "generic":
		c.plateauCMSURL = url
	}
}

type areaContext struct {
	Pref               *plateauapi.Prefecture
	City               *plateauapi.City
	CityItem           *CityItem
	Wards              []*plateauapi.Ward
	PrefID, CityID     *plateauapi.ID
	PrefCode, CityCode *plateauapi.AreaCode
}

func (c *areaContext) IsValid() bool {
	return c.Pref != nil && c.City != nil && c.CityItem != nil && c.PrefID != nil && c.CityID != nil && c.PrefCode != nil && c.CityCode != nil
}

func (c *internalContext) AreaContext(cityItemID string) *areaContext {
	var prefID, cityID *plateauapi.ID
	var prefCode, cityCode *plateauapi.AreaCode

	cityItem := c.CityItem(cityItemID)
	if cityItem == nil {
		return nil
	}

	city := c.cities[cityItem.CityCode]
	if city != nil {
		cityID = lo.ToPtr(city.ID)
		cityCode = lo.ToPtr(city.Code)
	}

	pref := c.prefs[city.Code.PrefectureCode()]
	if pref != nil {
		prefID = lo.ToPtr(pref.ID)
		prefCode = lo.ToPtr(pref.Code)
	}

	return &areaContext{
		CityItem: cityItem,
		City:     city,
		Pref:     pref,
		Wards:    c.Wards(cityItem.CityCode),
		PrefID:   prefID,
		CityID:   cityID,
		PrefCode: prefCode,
		CityCode: cityCode,
	}
}
