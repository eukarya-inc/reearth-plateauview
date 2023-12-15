package datacatalogv3

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
)

func (i *PlateauFeatureItem) toWards(pref *plateauapi.Prefecture, city *plateauapi.City) (res []*plateauapi.Ward) {
	dic := i.ReadDic()
	if dic == nil || len(dic["admin"]) == 0 {
		return nil
	}

	entries := dic["admin"]
	for _, entry := range entries {
		if entry.Code == "" || entry.Description == "" {
			continue
		}

		ward := &plateauapi.Ward{
			ID:             plateauapi.NewID(entry.Code, plateauapi.TypeArea),
			Name:           entry.Description,
			Type:           plateauapi.AreaTypeWard,
			Code:           plateauapi.AreaCode(entry.Code),
			PrefectureID:   pref.ID,
			PrefectureCode: pref.Code,
			CityID:         city.ID,
			CityCode:       city.Code,
		}

		res = append(res, ward)
	}

	return
}

func (i *PlateauFeatureItem) toDatasets(area *areaContext, dt *plateauapi.PlateauDatasetType, spec *plateauapi.PlateauSpecMinor) []plateauapi.Dataset {
	if len(i.Items) == 0 || len(i.Data) == 0 || area == nil || area.CityID == nil || area.CityCode == nil || area.PrefID == nil || area.PrefCode == nil {
		return nil
	}

	sid := standardItemID(dt.Code, area.City)
	id := plateauapi.NewID(sid, plateauapi.TypeDataset)

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
		Name:            standardItemName(dt.Name, area.City),
		Description:     toPtrIfPresent(i.Desc),
		Year:            area.CityItem.YearInt(),
		PrefectureID:    area.PrefID,
		PrefectureCode:  area.PrefCode,
		CityID:          area.CityID,
		CityCode:        area.CityCode,
		TypeID:          dt.ID,
		TypeCode:        dt.Code,
		PlateauSpecID:   spec.ParentID,
		PlateauSpecName: spec.Name,
		River:           river,
		Items:           items,
	}

	return []plateauapi.Dataset{&res}
}
