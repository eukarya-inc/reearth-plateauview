package datacatalogv3

import (
	"sort"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
)

type internalContext struct {
	years     map[int]struct{}
	cityItems map[string]*CityItem
	prefs     map[string]*plateauapi.Prefecture
	cities    map[string]*plateauapi.City
}

func newInternalContext() *internalContext {
	return &internalContext{
		years:     map[int]struct{}{},
		cityItems: map[string]*CityItem{},
		prefs:     map[string]*plateauapi.Prefecture{},
		cities:    map[string]*plateauapi.City{},
	}
}

func (c *internalContext) CityItem(id string) *CityItem {
	return c.cityItems[id]
}

func (c *internalContext) PrefAndCityFromCityItemID(id string) (*plateauapi.Prefecture, *plateauapi.City, *CityItem) {
	cityItem := c.CityItem(id)
	if cityItem == nil {
		return nil, nil, nil
	}

	city := c.cities[cityItem.CityCode]
	if city == nil {
		return nil, nil, nil
	}

	pref := c.prefs[city.Code.PrefectureCode()]
	if pref == nil {
		return nil, nil, nil
	}

	return pref, city, cityItem
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
	c.cityItems[cityItem.CityCode] = cityItem
	c.prefs[pref.Code.String()] = pref
	c.cities[city.Code.String()] = city

	if y := cityItem.YearInt(); y != 0 {
		c.years[y] = struct{}{}
	}
}
