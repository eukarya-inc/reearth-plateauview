package datacatalogv3

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

func (all *AllData) Into() (res plateauapi.InMemoryRepoContext, warning []string) {
	res.PlateauSpecs = plateauapi.PlateauSpecsFrom(all.PlateauSpecs)
	res.DatasetTypes = all.FeatureTypes.ToDatasetTypes(res.PlateauSpecs)

	ic := newInternalContext()

	// pref and city
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

	// wards
	for _, ft := range res.DatasetTypes[plateauapi.DatasetTypeCategoryPlateau] {
		wards, w := getWards(all.Plateau[ft.GetCode()], ic)
		warning = append(warning, w...)
		ic.AddWards(wards)
		res.Areas.Append(
			plateauapi.AreaTypeWard,
			lo.Map(wards, func(w *plateauapi.Ward, _ int) plateauapi.Area { return w }),
		)
	}

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

func getWards(items []*PlateauFeatureItem, ic *internalContext) (res []*plateauapi.Ward, warning []string) {
	for _, ds := range items {
		area := ic.AreaContext(ds.City)
		if area == nil {
			warning = append(warning, fmt.Sprintf("plateau %s: city not found: %s", ds.ID, ds.City))
			continue
		}

		wards := ds.toWards(area.Pref, area.City)
		res = append(res, wards...)
	}

	return
}

func convertPlateau(items []*PlateauFeatureItem, specs []plateauapi.PlateauSpec, dt plateauapi.DatasetType, ic *internalContext) (res []plateauapi.Dataset, warning []string) {
	pdt, ok := dt.(*plateauapi.PlateauDatasetType)
	if !ok {
		warning = append(warning, fmt.Sprintf("plateau %s: invalid dataset type: %s", dt.GetCode(), dt.GetName()))
		return
	}

	for _, ds := range items {
		area := ic.AreaContext(ds.City)
		if area == nil {
			warning = append(warning, fmt.Sprintf("plateau %s: city not found: %s", ds.ID, ds.City))
			continue
		}

		spec := plateauapi.FindSpecMinorByName(specs, area.CityItem.Spec)
		if spec == nil {
			warning = append(warning, fmt.Sprintf("plateau %s: spec not found: %s", ds.ID, area.CityItem.Spec))
			continue
		}

		ds, w := ds.toDatasets(area, pdt, spec)
		warning = append(warning, w...)
		if ds != nil {
			res = append(res, ds...)
		}
	}

	return
}

func convertRelated(items []*RelatedItem, datasetTypes []plateauapi.DatasetType, ic *internalContext) (res []plateauapi.Dataset, warning []string) {
	for _, ds := range items {
		area := ic.AreaContext(ds.City)
		if area == nil {
			warning = append(warning, fmt.Sprintf("related %s: city not found: %s", ds.ID, ds.City))
			continue
		}

		ds, w := ds.toDatasets(area, datasetTypes)
		warning = append(warning, w...)
		if ds != nil {
			res = append(res, ds...)
		}
	}

	return
}

func convertGeneric(items []*GenericItem, datasetTypes []plateauapi.DatasetType, ic *internalContext) (res []plateauapi.Dataset, warning []string) {
	for _, ds := range items {
		area := ic.AreaContext(ds.City)
		ds, w := ds.toDatasets(area, datasetTypes)
		warning = append(warning, w...)
		if ds != nil {
			res = append(res, ds...)
		}
	}

	return
}
