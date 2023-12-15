package datacatalogv3

import (
	"sort"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
)

func (all *AllData) Into() (res plateauapi.InMemoryRepoContext) {
	res.PlateauSpecs = plateauapi.PlateauSpecsFrom(all.PlateauSpecs)
	res.DatasetTypes = all.FeatureTypes.ToDatasetTypes(res.PlateauSpecs)

	years := map[int]struct{}{}
	prefs := map[plateauapi.AreaCode]struct{}{}
	cities := map[plateauapi.AreaCode]struct{}{}

	for _, city := range all.City {
		pref := city.ToPrefecture()
		if pref == nil {
			continue
		}

		if _, ok := prefs[pref.Code]; !ok {
			prefs[pref.Code] = struct{}{}
			res.Areas.Append(plateauapi.AreaTypePrefecture, []plateauapi.Area{*pref})
		}

		if _, ok := cities[plateauapi.AreaCode(city.CityCode)]; !ok {
			prefs[pref.Code] = struct{}{}
			res.Areas.Append(plateauapi.AreaTypePrefecture, []plateauapi.Area{*pref})
		}

		if y := city.YearInt(); y != 0 {
			years[y] = struct{}{}
		}
	}

	for _, ft := range all.FeatureTypes.Plateau {
		for _, f := range all.Plateau[ft.Code] {
			if ds := f.ToDatasets(); ds != nil {
				res.Datasets.Append(plateauapi.DatasetTypeCategoryPlateau, ds)
			}
		}
	}

	for _, ds := range all.Related {
		if ds := ds.ToDatasets(); ds != nil {
			res.Datasets.Append(plateauapi.DatasetTypeCategoryRelated, ds)
		}
	}

	for _, ds := range all.Generic {
		if ds := ds.ToDatasets(); ds != nil {
			res.Datasets.Append(plateauapi.DatasetTypeCategoryGeneric, ds)
		}
	}

	for y := range years {
		res.Years = append(res.Years, y)
	}
	sort.Ints(res.Years)

	return
}

func (all *AllData) FindPlateauByCityID(id, ft string) (res *PlateauFeatureItem) {
	features, ok := all.Plateau[ft]
	if !ok {
		return
	}

	for _, r := range features {
		if r.City == id {
			res = r
			return
		}
	}
	return
}

func (all *AllData) FindRelatedByCityID(id string) (res *RelatedItem) {
	for _, r := range all.Related {
		if r.City == id {
			res = r
			return
		}
	}
	return
}

func (all *AllData) FindGenericByCityID(id string) (res *GenericItem) {
	for _, r := range all.Generic {
		if r.City == id {
			res = r
			return
		}
	}
	return
}
