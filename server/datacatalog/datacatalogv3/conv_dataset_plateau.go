package datacatalogv3

import (
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

func (i *PlateauFeatureItem) toDatasets(area *areaContext, dt *plateauapi.PlateauDatasetType, spec *plateauapi.PlateauSpecMinor) (res []plateauapi.Dataset, warning []string) {
	if len(i.Items) == 0 || len(i.Data) == 0 || area == nil || area.CityID == nil || area.CityCode == nil || area.PrefID == nil || area.PrefCode == nil {
		return nil, nil
	}

	datasetSeeds := plateauDatasetSeedsFrom(i, dt, area, spec)
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
		return
	}

	sid := standardItemID(seed.DatasetType.Code, seed.Area.City)
	id := plateauapi.NewID(sid, plateauapi.TypeDataset)

	seeds, w := plateauDatasetItemSeedFrom(seed)
	warning = append(warning, w...)
	items := lo.Map(seeds, func(s plateauDatasetItemSeed, _ int) *plateauapi.PlateauDatasetItem {
		return seedToDatasetItem(s, sid)
	})

	if len(items) == 0 {
		// warning is already reported by plateauDatasetItemSeedFrom
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
	return &plateauapi.PlateauDatasetItem{
		ID:      i.GetID(parentID),
		Name:    i.GetName(),
		URL:     i.URL,
		Format:  datasetFormatFrom(i.Format),
		Lod:     i.LOD,
		Texture: textureFrom(i.NoTexture),
	}
}
