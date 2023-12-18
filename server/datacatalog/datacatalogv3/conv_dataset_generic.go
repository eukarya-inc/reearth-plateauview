package datacatalogv3

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

func (i *GenericItem) toDatasets(area *areaContext, dts []plateauapi.DatasetType) (_ []plateauapi.Dataset, warning []string) {
	if area == nil {
		area = &areaContext{}
	}

	id := plateauapi.NewID(i.ID, plateauapi.TypeDataset)

	dt, _ := lo.Find(dts, func(dt plateauapi.DatasetType) bool {
		return dt.GetName() == i.Category
	})
	if dt == nil {
		warning = append(warning, fmt.Sprintf("generic %s: dataset type not found: %s", i.ID, i.Category))
		return
	}

	items := lo.FilterMap(i.Data, func(datum GenericItemDataset, ind int) (*plateauapi.GenericDatasetItem, bool) {
		url := datum.DataURL
		if url == "" {
			url = datum.Data
		}
		f := datasetFormatFromOrDetect(datum.DataFormat, url)
		if url == "" || f == "" {
			warning = append(warning, fmt.Sprintf("generic %s[%d]: invalid url: %s", i.ID, ind, url))
			return nil, false
		}

		var inds string
		if len(i.Data) > 1 {
			inds = fmt.Sprintf(" %d", ind+1)
		}

		return &plateauapi.GenericDatasetItem{
			ID:       plateauapi.NewID(datum.ID, plateauapi.TypeDatasetItem),
			Name:     firstNonEmptyValue(datum.Name, fmt.Sprintf("%s%s", i.Name, inds)),
			URL:      url,
			Format:   f,
			Layers:   layerNamesFrom(datum.LayerName),
			ParentID: id,
		}, true
	})

	if len(items) == 0 {
		warning = append(warning, fmt.Sprintf("generic %s: no items", i.ID))
		return
	}

	res := plateauapi.GenericDataset{
		ID:             id,
		Name:           i.Name,
		Description:    toPtrIfPresent(i.Desc),
		Year:           area.CityItem.YearInt(),
		PrefectureID:   area.PrefID,
		PrefectureCode: area.PrefCode,
		CityID:         area.CityID,
		CityCode:       area.CityCode,
		TypeID:         dt.GetID(),
		TypeCode:       dt.GetCode(),
		Items:          items,
	}

	return []plateauapi.Dataset{&res}, warning
}
