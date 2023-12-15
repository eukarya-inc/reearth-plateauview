package datacatalogv3

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
)

func (all *AllData) Into() (res plateauapi.InMemoryRepoContext, warning []string) {
	res.PlateauSpecs = plateauapi.PlateauSpecsFrom(all.PlateauSpecs)
	res.DatasetTypes = all.FeatureTypes.ToDatasetTypes(res.PlateauSpecs)

	ic := newInternalContext()

	for _, cityItem := range all.City {
		pref, city := cityItem.ToPrefecture(), cityItem.ToCity()
		if pref == nil || city == nil {
			continue
		}

		ic.Add(cityItem, pref, city)

		if !ic.HasPref(pref.Code.String()) {
			res.Areas.Append(plateauapi.AreaTypePrefecture, []plateauapi.Area{*pref})
		}

		if !ic.HasCity(city.Code.String()) {
			res.Areas.Append(plateauapi.AreaTypeCity, []plateauapi.Area{*city})
		}
	}

	res.Years = ic.Years()

	// plateau
	for _, ft := range res.DatasetTypes[plateauapi.DatasetTypeCategoryPlateau] {
		datasets, w := convertPlateau(all.Plateau[ft.GetCode()], res.PlateauSpecs, ft, ic)
		warning = append(warning, w...)
		res.Datasets.Append(plateauapi.DatasetTypeCategoryPlateau, datasets)
	}

	// related
	{
		datasets, w := convertRelated(all.Related, res.DatasetTypes[plateauapi.DatasetTypeCategoryRelated], ic)
		warning = append(warning, w...)
		res.Datasets.Append(plateauapi.DatasetTypeCategoryRelated, datasets)
	}

	// generic
	{
		datasets, w := convertGeneric(all.Generic, res.DatasetTypes[plateauapi.DatasetTypeCategoryGeneric], ic)
		warning = append(warning, w...)
		res.Datasets.Append(plateauapi.DatasetTypeCategoryGeneric, datasets)
	}

	return
}

func convertPlateau(items []*PlateauFeatureItem, specs []plateauapi.PlateauSpec, dt plateauapi.DatasetType, ic *internalContext) (res []plateauapi.Dataset, warning []string) {
	pdt, ok := dt.(*plateauapi.PlateauDatasetType)
	if !ok {
		warning = append(warning, fmt.Sprintf("invalid dataset type: %s", dt.GetCode()))
		return
	}

	for _, ds := range items {
		pref, city, cityItem := ic.PrefAndCityFromCityItemID(ds.City)
		if pref == nil || city == nil || cityItem == nil {
			warning = append(warning, fmt.Sprintf("plateau %s: city not found: %s", ds.ID, ds.City))
			continue
		}

		spec := plateauapi.FindSpecMinorByName(specs, cityItem.Spec)
		if spec == nil {
			warning = append(warning, fmt.Sprintf("plateau %s: spec not found: %s", ds.ID, cityItem.Spec))
			continue
		}

		if ds := ds.ToDatasets(
			pref,
			city,
			pdt,
			spec,
		); ds != nil {
			res = append(res, ds...)
		}
	}

	return
}

func convertRelated(items []*RelatedItem, datasetTypes []plateauapi.DatasetType, ic *internalContext) (res []plateauapi.Dataset, warning []string) {
	for _, ds := range items {
		pref, city, cityItem := ic.PrefAndCityFromCityItemID(ds.City)
		if pref == nil || city == nil || cityItem == nil {
			warning = append(warning, fmt.Sprintf("generic %s: city not found: %s", ds.ID, ds.City))
			continue
		}

		if ds := ds.ToDatasets(
			pref,
			city,
			datasetTypes,
			cityItem.YearInt(),
		); ds != nil {
			res = append(res, ds...)
		}
	}

	return
}

func convertGeneric(items []*GenericItem, datasetTypes []plateauapi.DatasetType, ic *internalContext) (res []plateauapi.Dataset, warning []string) {
	for _, ds := range items {
		pref, city, cityItem := ic.PrefAndCityFromCityItemID(ds.City)
		if pref == nil || city == nil || cityItem == nil {
			warning = append(warning, fmt.Sprintf("generic %s: city not found: %s", ds.ID, ds.City))
			continue
		}

		if ds := ds.ToDatasets(
			pref,
			city,
			datasetTypes,
			cityItem.YearInt(),
		); ds != nil {
			res = append(res, ds...)
		}
	}

	return
}
