package datacatalogv3

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
)

func (i *PlateauFeatureItem) ToDatasets(pref *plateauapi.Prefecture, city *plateauapi.City, ft *plateauapi.PlateauDatasetType, spec *plateauapi.PlateauSpecMinor) []plateauapi.Dataset {
	if len(i.Items) == 0 || len(i.Data) == 0 {
		return nil
	}

	sid := standardItemID(ft.Code, city)
	id := plateauapi.NewID(sid, plateauapi.TypeDataset)
	prefID, cityID, prefCode, cityCode := areaInfo(pref, city)
	if prefID == nil || cityID == nil || prefCode == nil || cityCode == nil {
		return nil
	}

	var river *plateauapi.River                // TODO
	var items []*plateauapi.PlateauDatasetItem // TODO

	data := i.Items
	if len(data) == 0 && len(i.Data) > 0 {
		data = append(data, PlateauFeatureItemDatum{
			Data: i.Data,
			Desc: i.Desc,
		})
	}

	for _, d := range data {
		if len(d.Data) == 0 {
			continue
		}

		items = append(items, &plateauapi.PlateauDatasetItem{
			// TODO
			// ID:   plateauapi.NewID(fmt.Sprintf("%s_%d", sid, d.Index), plateauapi.TypeDatasetItem),
			// Name: firstNonEmptyValue(d.Name, fmt.Sprintf("%s%s", i.Name, inds)),
			// URL:      d.Data,
			// Format:   datasetFormatFrom(d.DataFormat),
			// Layers:   layerNamesFrom(d.LayerName),
			ParentID: id,
		})
	}

	if len(items) == 0 {
		return nil
	}

	res := plateauapi.PlateauDataset{
		ID:              id,
		Name:            standardItemName(ft.Name, city),
		Description:     toPtrIfPresent(i.Desc),
		Year:            ft.Year,
		PrefectureID:    prefID,
		PrefectureCode:  prefCode,
		CityID:          cityID,
		CityCode:        cityCode,
		TypeID:          ft.ID,
		TypeCode:        ft.Code,
		PlateauSpecID:   spec.ParentID,
		PlateauSpecName: spec.Name,
		River:           river,
		Items:           items,
	}

	return []plateauapi.Dataset{&res}
}
