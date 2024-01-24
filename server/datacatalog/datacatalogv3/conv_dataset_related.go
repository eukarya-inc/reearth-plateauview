package datacatalogv3

import (
	"fmt"
	"slices"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

var convIgnored = []string{"border"}

func (i *RelatedItem) toDatasets(area *areaContext, dts []plateauapi.DatasetType, cmsurl string) (res []plateauapi.Dataset, warning []string) {
	if !area.IsValid() {
		warning = append(warning, fmt.Sprintf("related %s: invalid area", i.ID))
		return
	}

	for _, dt := range dts {
		ftname, ftcode := dt.GetName(), dt.GetCode()
		d := i.Items[ftcode]
		if len(d.Asset) == 0 && len(d.Converted) == 0 {
			// warning = append(warning, fmt.Sprintf("related %s: no data for %s", area.CityCode, ftcode))
			continue
		}

		if slices.Contains(convIgnored, ftcode) {
			d.Converted = nil
		}

		seeds, w := assetUrlsToRelatedDatasetSeeds(d.Asset, d.Converted, area.City, area.Wards, area.CityItem.YearInt())
		warning = append(warning, w...)

		for _, seed := range seeds {
			sid := standardItemID(ftcode, seed.Area, "")
			id := plateauapi.NewID(sid, plateauapi.TypeDataset)

			var ou *string
			if seed.OriginalURL != "" && seed.OriginalFormat != nil {
				ou = lo.EmptyableToPtr(assetURLFromFormat(seed.OriginalURL, *seed.OriginalFormat))
			}

			res = append(res, &plateauapi.RelatedDataset{
				ID:             id,
				Name:           standardItemName(ftname, "", seed.Area),
				Description:    toPtrIfPresent(d.Description),
				Year:           area.CityItem.YearInt(),
				PrefectureID:   area.PrefID,
				PrefectureCode: area.PrefCode,
				CityID:         area.CityID,
				CityCode:       area.CityCode,
				WardID:         seed.WardID,
				WardCode:       seed.WardCode,
				TypeID:         dt.GetID(),
				TypeCode:       ftcode,
				Admin:          adminFrom(area.CityItem, cmsurl, "related"),
				Items: []*plateauapi.RelatedDatasetItem{
					{
						ID:             plateauapi.NewID(sid, plateauapi.TypeDatasetItem),
						Name:           ftname,
						Format:         seed.Format,
						URL:            assetURLFromFormat(seed.URL, seed.Format),
						OriginalFormat: seed.OriginalFormat,
						OriginalURL:    ou,
						ParentID:       id,
					},
				},
			})
		}
	}

	return
}

type relatedDatasetSeed struct {
	Area           plateauapi.Area
	WardID         *plateauapi.ID
	WardCode       *plateauapi.AreaCode
	URL            string
	Format         plateauapi.DatasetFormat
	OriginalURL    string
	OriginalFormat *plateauapi.DatasetFormat
}

func assetUrlsToRelatedDatasetSeeds(orig, conv []string, city *plateauapi.City, wards []*plateauapi.Ward, year int) (items []relatedDatasetSeed, warning []string) {
	var assets []OriginalAndConv
	if len(conv) == 0 {
		assets = lo.Map(orig, func(a string, _ int) OriginalAndConv {
			return OriginalAndConv{
				Converted: a,
			}
		})
	} else {
		assets = OriginalAndConvsFrom(orig, conv)
	}

	wasCityAdded := false
	for _, asset := range assets {
		var format plateauapi.DatasetFormat
		var origFormat *plateauapi.DatasetFormat
		if asset.Original != "" {
			format = plateauapi.DatasetFormatCzml
			origFormat = lo.ToPtr(plateauapi.DatasetFormatGeojson)
		} else {
			format = plateauapi.DatasetFormatGeojson
		}

		nameOrig := nameFromURL(asset.Converted)
		assetNameOrig := ParseRelatedAssetName(nameOrig)

		nameConv := nameFromURL(asset.Converted)
		assetNameConv := ParseRelatedAssetName(nameConv)

		if nameOrig != "" && assetNameOrig == nil {
			warning = append(warning, fmt.Sprintf("related %s: invalid asset name: %s", city.Code, nameOrig))
			continue
		}

		if nameConv != "" && assetNameConv == nil {
			warning = append(warning, fmt.Sprintf("related %s: invalid asset name: %s", city.Code, nameConv))
			continue
		}

		if assetNameConv.Year != 0 && assetNameConv.Year != year {
			warning = append(warning, fmt.Sprintf("related %s: invalid year: %s: %d should be %d", city.Code, nameConv, assetNameConv.Year, year))
		}

		if assetNameOrig.Year != 0 && assetNameOrig.Year != year {
			warning = append(warning, fmt.Sprintf("related %s: invalid year: %s: %d should be %d", city.Code, nameConv, assetNameOrig.Year, year))
		}

		// city
		if assetNameConv.Code == city.Code.String() {
			if wasCityAdded {
				warning = append(warning, fmt.Sprintf("related %s: duplicated assets that have the same city code: %s", city.Code, nameConv))
				continue
			}

			// it's a city
			items = append(items, relatedDatasetSeed{
				Area:           city,
				URL:            asset.Converted,
				Format:         format,
				OriginalURL:    asset.Original,
				OriginalFormat: origFormat,
			})
			wasCityAdded = true
			continue
		}

		// wards
		ward, _ := lo.Find(wards, func(w *plateauapi.Ward) bool {
			return w.Code.String() == assetNameConv.Code
		})

		if ward == nil {
			warning = append(warning, fmt.Sprintf("related %s: ward not found: %s", city.Code, nameConv))
			continue
		}

		items = append(items, relatedDatasetSeed{
			Area:           ward,
			WardID:         lo.ToPtr(ward.ID),
			WardCode:       lo.ToPtr(ward.Code),
			URL:            asset.Converted,
			Format:         format,
			OriginalURL:    asset.Original,
			OriginalFormat: origFormat,
		})
	}

	return
}
