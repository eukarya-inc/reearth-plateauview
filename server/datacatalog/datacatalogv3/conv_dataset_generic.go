package datacatalogv3

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

func (i *GenericItem) toDatasets(area *areaContext, dts []plateauapi.DatasetType) []plateauapi.Dataset {
	if area == nil {
		area = &areaContext{}
	}

	id := plateauapi.NewID(i.ID, plateauapi.TypeDataset)

	dt, _ := lo.Find(dts, func(dt plateauapi.DatasetType) bool {
		return dt.GetName() == i.Category
	})
	if dt == nil {
		return nil
	}

	items := lo.FilterMap(i.Data, func(datum GenericItemDataset, ind int) (*plateauapi.GenericDatasetItem, bool) {
		if datum.Data == "" {
			return nil, false
		}

		var inds string
		if len(i.Data) > 1 {
			inds = fmt.Sprintf(" %d", ind+1)
		}

		return &plateauapi.GenericDatasetItem{
			ID:       plateauapi.NewID(datum.ID, plateauapi.TypeDatasetItem),
			Name:     firstNonEmptyValue(datum.Name, fmt.Sprintf("%s%s", i.Name, inds)),
			URL:      datum.Data,
			Format:   datasetFormatFrom(datum.DataFormat),
			Layers:   layerNamesFrom(datum.LayerName),
			ParentID: id,
		}, true
	})

	if len(items) == 0 {
		return nil
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

	return []plateauapi.Dataset{&res}
}
