package datacatalogv3

import (
	"fmt"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

const dicKeyAdmin = "admin"

func (i *PlateauFeatureItem) toWards(pref *plateauapi.Prefecture, city *plateauapi.City) (res []*plateauapi.Ward) {
	dic := i.ReadDic()
	if dic == nil || len(dic[dicKeyAdmin]) == 0 {
		return nil
	}

	entries := dic[dicKeyAdmin]
	for _, entry := range entries {
		if entry.Code == "" || entry.Description == "" {
			continue
		}

		_, name, _ := strings.Cut(entry.Description, " ")
		if name == "" {
			name = entry.Description
		}

		ward := &plateauapi.Ward{
			ID:             plateauapi.NewID(entry.Code, plateauapi.TypeArea),
			Name:           name,
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

func (i *PlateauFeatureItem) toDatasets(area *areaContext, dt *plateauapi.PlateauDatasetType, spec *plateauapi.PlateauSpecMinor) (res []plateauapi.Dataset, warning []string) {
	if !area.IsValid() {
		warning = append(warning, fmt.Sprintf("plateau %s: invalid area", i.ID))
		return
	}

	datasetSeeds, w := plateauDatasetSeedsFrom(i, dt, area, spec)
	warning = append(warning, w...)
	for _, seed := range datasetSeeds {
		dataset, w := seedToDataset(seed)
		warning = append(warning, w...)
		if dataset != nil {
			res = append(res, dataset)
		}
	}

	return
}

func seedToDataset(seed plateauDatasetSeed) (res *plateauapi.PlateauDataset, warning []string) {
	if len(seed.AssetURLs) == 0 {
		warning = append(warning, fmt.Sprintf("plateau %s %s: no asset urls", seed.TargetArea.GetCode(), seed.DatasetType.Code))
		return
	}

	sid := seed.GetID()
	id := plateauapi.NewID(sid, plateauapi.TypeDataset)

	seeds, w := plateauDatasetItemSeedFrom(seed)
	warning = append(warning, w...)
	items := lo.FilterMap(seeds, func(s plateauDatasetItemSeed, i int) (*plateauapi.PlateauDatasetItem, bool) {
		item := seedToDatasetItem(s, sid)
		if item == nil {
			warning = append(warning, fmt.Sprintf("plateau %s %s[%d]: unknown dataset format: %s", seed.TargetArea.GetCode(), seed.DatasetType.Code, i, s.URL))
		}
		return item, item != nil
	})

	if len(items) == 0 {
		// warning is already reported by plateauDatasetItemSeedFrom
		warning = append(warning, fmt.Sprintf("plateau %s %s: no items", seed.TargetArea.GetCode(), seed.DatasetType.Code))
		return
	}

	res = &plateauapi.PlateauDataset{
		ID:              id,
		Name:            standardItemName(seed.DatasetType.Name, seed.SubName, seed.TargetArea),
		Description:     toPtrIfPresent(seed.Desc),
		Year:            seed.Area.CityItem.YearInt(),
		PrefectureID:    seed.Area.PrefID,
		PrefectureCode:  seed.Area.PrefCode,
		CityID:          seed.Area.CityID,
		CityCode:        seed.Area.CityCode,
		WardID:          seed.WardID,
		WardCode:        seed.WardCode,
		TypeID:          seed.DatasetType.ID,
		TypeCode:        seed.DatasetType.Code,
		PlateauSpecID:   seed.Spec.ParentID,
		PlateauSpecName: seed.Spec.Name,
		River:           seed.River,
		Items:           items,
	}

	return
}

func seedToDatasetItem(i plateauDatasetItemSeed, parentID string) *plateauapi.PlateauDatasetItem {
	f := datasetFormatFromOrDetect(i.Format, i.URL)
	if f == "" {
		return nil
	}

	return &plateauapi.PlateauDatasetItem{
		ID:       plateauapi.NewID(i.GetID(parentID), plateauapi.TypeDatasetItem),
		Name:     i.GetName(),
		URL:      i.URL,
		Format:   f,
		Lod:      i.LOD,
		Texture:  textureFrom(i.NoTexture),
		ParentID: plateauapi.NewID(parentID, plateauapi.TypeDataset),
	}
}
