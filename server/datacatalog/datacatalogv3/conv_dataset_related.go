package datacatalogv3

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

func (i *RelatedItem) toDatasets(area *areaContext, dts []plateauapi.DatasetType) (res []plateauapi.Dataset, warning []string) {
	if !area.IsValid() {
		warning = append(warning, fmt.Sprintf("related %s: invalid area", i.ID))
		return
	}

	for _, dt := range dts {
		ftname, ftcode := dt.GetName(), dt.GetCode()
		assets := i.ConvertedAssets[ftcode]
		format := plateauapi.DatasetFormatCzml
		if len(assets) == 0 {
			assets = i.Assets[ftcode]
			format = plateauapi.DatasetFormatGeojson
		}
		if len(assets) == 0 {
			warning = append(warning, fmt.Sprintf("related %s: no assets for %s", area.CityCode, ftcode))
			continue
		}

		seeds, w := assetUrlsToRelatedDatasetSeeds(assets, area.City, area.Wards)
		warning = append(warning, w...)

		for _, seed := range seeds {
			sid := standardItemID(ftcode, seed.Area)
			id := plateauapi.NewID(sid, plateauapi.TypeDataset)
			res = append(res, &plateauapi.RelatedDataset{
				ID:             id,
				Name:           standardItemName(ftname, "", seed.Area),
				Description:    toPtrIfPresent(i.Desc),
				Year:           area.CityItem.YearInt(),
				PrefectureID:   area.PrefID,
				PrefectureCode: area.PrefCode,
				CityID:         area.CityID,
				CityCode:       area.CityCode,
				WardID:         seed.WardID,
				WardCode:       seed.WardCode,
				TypeID:         dt.GetID(),
				TypeCode:       ftcode,
				Items: []*plateauapi.RelatedDatasetItem{
					{
						ID:       plateauapi.NewID(sid, plateauapi.TypeDatasetItem),
						Format:   format,
						Name:     ftname,
						URL:      seed.URL,
						ParentID: id,
					},
				},
			})
		}
	}

	return
}

type relatedDatasetSeed struct {
	Area     plateauapi.Area
	WardID   *plateauapi.ID
	WardCode *plateauapi.AreaCode
	URL      string
}

func assetUrlsToRelatedDatasetSeeds(urls []string, city *plateauapi.City, wards []*plateauapi.Ward) (items []relatedDatasetSeed, warning []string) {
	wasCityAdded := false
	for _, url := range urls {
		name := nameFromURL(url)
		assetName := ParseRelatedAssetName(name)
		if assetName == nil {
			warning = append(warning, fmt.Sprintf("related %s: invalid asset name: %s", city.Code, name))
			continue
		}

		if assetName.Code == city.Code.String() {
			if wasCityAdded {
				warning = append(warning, fmt.Sprintf("related %s: city already added: %s", city.Code, name))
				continue
			}

			// it's a city
			items = append(items, relatedDatasetSeed{
				Area: city,
				URL:  url,
			})
			wasCityAdded = true
			continue
		}

		ward, _ := lo.Find(wards, func(w *plateauapi.Ward) bool {
			return w.Code.String() == assetName.Code
		})

		if ward == nil {
			warning = append(warning, fmt.Sprintf("related %s: ward not found: %s", city.Code, name))
			continue
		}

		items = append(items, relatedDatasetSeed{
			Area:     ward,
			WardID:   lo.ToPtr(ward.ID),
			WardCode: lo.ToPtr(ward.Code),
			URL:      url,
		})
	}

	return
}
