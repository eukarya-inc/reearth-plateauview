package datacatalogv3

import (
	"fmt"
	"regexp"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

func (i *RelatedItem) toDatasets(area *areaContext, dts []plateauapi.DatasetType) (res []plateauapi.Dataset, warning []string) {
	if area == nil || area.PrefID == nil || area.PrefCode == nil || area.CityID == nil || area.CityCode == nil {
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
			continue
		}

		items, w := assetUrlsToRelatedItems(assets, area.City, area.Wards)
		warning = append(warning, w...)

		for _, item := range items {
			sid := standardItemID(ftcode, item.Area)
			id := plateauapi.NewID(sid, plateauapi.TypeDataset)
			res = append(res, plateauapi.RelatedDataset{
				ID:             id,
				Name:           standardItemName(ftname, item.Area),
				Description:    toPtrIfPresent(i.Desc),
				Year:           area.CityItem.YearInt(),
				PrefectureID:   area.PrefID,
				PrefectureCode: area.PrefCode,
				CityID:         area.CityID,
				CityCode:       area.CityCode,
				WardID:         item.WardID,
				WardCode:       item.WardCode,
				TypeID:         dt.GetID(),
				TypeCode:       ftcode,
				Items: []*plateauapi.RelatedDatasetItem{
					{
						ID:       plateauapi.NewID(sid, plateauapi.TypeDatasetItem),
						Format:   format,
						Name:     ftname,
						URL:      item.URL,
						ParentID: id,
					},
				},
			})
		}
	}

	return
}

type relatedItem struct {
	Area     plateauapi.Area
	WardID   *plateauapi.ID
	WardCode *plateauapi.AreaCode
	URL      string
}

func assetUrlsToRelatedItems(urls []string, city *plateauapi.City, wards []*plateauapi.Ward) (items []relatedItem, warning []string) {
	var res []relatedItem

	wasCityAdded := false
	for _, url := range urls {
		name := nameFromURL(url)
		assetName := ParseRelatedAssetName(name)
		if assetName == nil {
			warning = append(warning, fmt.Sprintf("related %s %s: invalid asset name: %s", city.Code, assetName.Type, name))
			continue
		}

		if assetName.Code == city.Code.String() {
			if wasCityAdded {
				warning = append(warning, fmt.Sprintf("related %s %s: city already added: %s", city.Code, assetName.Type, name))
				continue
			}

			// it's a city
			items = append(items, relatedItem{
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
			warning = append(warning, fmt.Sprintf("related %s %s: ward not found: %s", city.Code, assetName.Type, name))
			continue
		}

		res = append(res, relatedItem{
			Area:     ward,
			WardID:   lo.ToPtr(ward.ID),
			WardCode: lo.ToPtr(ward.Code),
			URL:      url,
		})
	}

	return
}

type RelatedAssetName struct {
	Code string
	Name string
	Type string
	Ext  string
}

var reRelatedAssetName = regexp.MustCompile(`^(\d{5})_([a-zA-Z0-9-]+)_([a-zA-Z0-9-]+)\.([a-z0-9]+)$`)

func ParseRelatedAssetName(name string) *RelatedAssetName {
	m := reRelatedAssetName.FindStringSubmatch(name)
	if m == nil {
		return nil
	}

	return &RelatedAssetName{
		Code: m[1],
		Name: m[2],
		Type: m[3],
		Ext:  m[4],
	}
}
