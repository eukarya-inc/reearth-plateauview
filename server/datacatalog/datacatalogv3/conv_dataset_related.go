package datacatalogv3

import "github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"

func (i *RelatedItem) ToDatasets(pref *plateauapi.Prefecture, city *plateauapi.City, featureTypes []plateauapi.DatasetType, year int) (res []plateauapi.Dataset) {
	prefID, cityID, prefCode, cityCode := areaInfo(pref, city)
	if prefID == nil || cityID == nil || prefCode == nil || cityCode == nil {
		return nil
	}

	for _, ft := range featureTypes {
		ftname, ftcode := ft.GetName(), ft.GetCode()
		asset := i.ConvertedAssets[ftcode]
		format := plateauapi.DatasetFormatCzml
		if len(asset) == 0 {
			asset = i.Assets[ftcode]
			format = plateauapi.DatasetFormatGeojson
		}
		if len(asset) == 0 {
			continue
		}

		sid := standardItemID(ftcode, city)
		id := plateauapi.NewID(sid, plateauapi.TypeDataset)
		res = append(res, plateauapi.RelatedDataset{
			ID:             id,
			Name:           standardItemName(ftname, city),
			Description:    toPtrIfPresent(i.Desc),
			Year:           year,
			PrefectureID:   prefID,
			PrefectureCode: prefCode,
			CityID:         cityID,
			CityCode:       cityCode,
			TypeID:         ft.GetID(),
			TypeCode:       ftcode,
			Items: []*plateauapi.RelatedDatasetItem{
				{
					ID:       plateauapi.NewID(sid, plateauapi.TypeDatasetItem),
					Format:   format,
					Name:     ftname,
					URL:      asset[0], // TODO
					ParentID: id,
				},
			},
		})
	}

	return
}
