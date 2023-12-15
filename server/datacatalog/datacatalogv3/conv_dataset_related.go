package datacatalogv3

import "github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"

func (i *RelatedItem) toDatasets(area *areaContext, dts []plateauapi.DatasetType) (res []plateauapi.Dataset) {
	if area == nil || area.PrefID == nil || area.PrefCode == nil || area.CityID == nil || area.CityCode == nil {
		return nil
	}

	for _, dt := range dts {
		ftname, ftcode := dt.GetName(), dt.GetCode()
		asset := i.ConvertedAssets[ftcode]
		format := plateauapi.DatasetFormatCzml
		if len(asset) == 0 {
			asset = i.Assets[ftcode]
			format = plateauapi.DatasetFormatGeojson
		}
		if len(asset) == 0 {
			continue
		}

		assets := []string{asset[0]}

		for _, asset := range assets {
			sid := standardItemID(ftcode, area.City)
			id := plateauapi.NewID(sid, plateauapi.TypeDataset)
			res = append(res, plateauapi.RelatedDataset{
				ID:             id,
				Name:           standardItemName(ftname, area.City),
				Description:    toPtrIfPresent(i.Desc),
				Year:           area.CityItem.YearInt(),
				PrefectureID:   area.PrefID,
				PrefectureCode: area.PrefCode,
				CityID:         area.CityID,
				CityCode:       area.CityCode,
				TypeID:         dt.GetID(),
				TypeCode:       ftcode,
				Items: []*plateauapi.RelatedDatasetItem{
					{
						ID:       plateauapi.NewID(sid, plateauapi.TypeDatasetItem),
						Format:   format,
						Name:     ftname,
						URL:      asset,
						ParentID: id,
					},
				},
			})
		}
	}

	return
}

// func wardCodesFromAssetURLs(urls []string) []string {
// 	hit := false
// 	res := make([]string, 0, len(urls))
// 	for _, url := range urls {
// 		name := nameWithoutExt(nameFromURL(url))
// 		// TODO
// 		res = append(res, name)
// 	}

// 	if !hit {
// 		return nil
// 	}

// 	panic("not implemented")
// 	// return res
// }
