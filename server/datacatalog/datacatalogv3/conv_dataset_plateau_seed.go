package datacatalogv3

import (
	"fmt"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

type plateauDatasetSeed struct {
	IDEx       string
	AssetURLs  []string
	Assets     []*AssetName
	SubName    string
	Desc       string
	TargetArea plateauapi.Area
	WardID     *plateauapi.ID
	WardCode   *plateauapi.AreaCode
	// common
	DatasetType *plateauapi.PlateauDatasetType
	Dic         Dic
	Area        *areaContext
	Pref        *plateauapi.Prefecture
	City        *plateauapi.City
	Spec        *plateauapi.PlateauSpecMinor
	River       *plateauapi.River
}

func (seed plateauDatasetSeed) GetID() string {
	return standardItemID(seed.DatasetType.Code, seed.TargetArea, seed.IDEx)
}

func plateauDatasetSeedsFrom(i *PlateauFeatureItem, dt *plateauapi.PlateauDatasetType, area *areaContext, spec *plateauapi.PlateauSpecMinor) (res []plateauDatasetSeed, warning []string) {
	dic := i.ReadDic()

	if len(i.Items) > 0 {
		for _, item := range i.Items {
			seeds, w := plateauDatasetSeedsFromItem(i, item, dt, dic)
			warning = append(warning, w...)
			res = append(res, seeds)
		}
	} else {
		seeds, w := plateauDatasetSeedsFromBldg(i, dt, area.Wards)
		warning = append(warning, w...)
		res = append(res, seeds...)
	}

	for i := range res {
		res[i].DatasetType = dt
		res[i].Dic = dic
		res[i].Area = area
		res[i].Pref = area.Pref
		res[i].City = area.City
		res[i].Spec = spec
		if res[i].TargetArea == nil {
			res[i].TargetArea = area.City
		}
	}

	return
}

func plateauDatasetSeedsFromItem(i *PlateauFeatureItem, item PlateauFeatureItemDatum, dt *plateauapi.PlateauDatasetType, dic Dic) (res plateauDatasetSeed, warning []string) {
	assets := lo.Map(item.Data, func(url string, _ int) *AssetName {
		n := nameWithoutExt(nameFromURL(url))
		an := ParseAssetName(n)
		if an == nil || !an.Ex.IsValid() {
			warning = append(warning, fmt.Sprintf("plateau %s %s: invalid asset name: %s", i.ID, dt.Code, n))
		}
		return an
	})
	if len(assets) == 0 {
		warning = append(warning, fmt.Sprintf("plateau %s %s: no assets", i.ID, dt.Code))
		return
	}

	key := assets[0].Ex.Key()
	if key == "" {
		warning = append(warning, fmt.Sprintf("plateau %s %s: invalid asset name key: %s", i.ID, dt.Code, assets[0].Ex.Ex))
		return
	}

	e := dic.FindEntry(dt.Code, key)
	if e == nil {
		warning = append(warning, fmt.Sprintf("plateau %s %s: dic entry not found: %s", i.ID, dt.Code, key))
		return
	}

	var river *plateauapi.River
	if a := riverAdminFrom(e.Admin); a != nil {
		if e.Description == "" {
			warning = append(warning, fmt.Sprintf("plateau %s %s: dic entry has no description: %s", i.ID, dt.Code, key))
		} else {
			river = &plateauapi.River{
				Name:  e.Description,
				Admin: *a,
			}
		}
	}

	subname := item.Name
	if subname == "" {
		subname = datasetSubNameFromDicEntry(e)
	}
	if subname == "" && e != nil {
		warning = append(warning, fmt.Sprintf("plateau %s %s: invalid dic entry: %s", i.ID, dt.Code, key))
	}

	res = plateauDatasetSeed{
		IDEx:      key,
		AssetURLs: item.Data,
		Assets:    assets,
		SubName:   subname,
		Desc:      item.Desc,
		River:     river,
	}
	return
}

func datasetSubNameFromDicEntry(e *DicEntry) string {
	if e == nil {
		return ""
	}
	if e.Admin != "" {
		// it's river
		return fmt.Sprintf("%s（%s管理区間）", e.Description, e.Admin)
	}
	return e.Description
}

func plateauDatasetSeedsFromBldg(i *PlateauFeatureItem, dt *plateauapi.PlateauDatasetType, wards []*plateauapi.Ward) (res []plateauDatasetSeed, warning []string) {
	if len(wards) == 0 {
		res = append(res, plateauDatasetSeed{
			AssetURLs: i.Data,
			Desc:      i.Desc,
		})
		return
	}

	assets := lo.Zip2(lo.Map(i.Data, func(url string, ind int) *AssetName {
		n := nameWithoutExt(nameFromURL(url))
		an := ParseAssetName(n)
		if an == nil || an.Ex.Normal == nil {
			warning = append(warning, fmt.Sprintf("plateau %s %s[%d]: invalid asset name: %s", i.ID, dt.Code, ind, n))
		}
		return an
	}), i.Data)
	if len(assets) == 0 {
		warning = append(warning, fmt.Sprintf("plateau %s %s: no assets", i.ID, dt.Code))
		return
	}

	for _, ward := range wards {
		wardCode := ward.Code.String()
		assets := lo.Filter(assets, func(name lo.Tuple2[*AssetName, string], _ int) bool {
			return name.A != nil && name.A.Ex.Normal != nil && name.A.Ex.Normal.WardCode == wardCode
		})
		if len(assets) == 0 {
			warning = append(warning, fmt.Sprintf("plateau %s %s: no assets for ward %s", i.ID, dt.Code, wardCode))
			continue
		}

		res = append(res, plateauDatasetSeed{
			AssetURLs: lo.Map(assets, func(name lo.Tuple2[*AssetName, string], _ int) string {
				return name.B
			}),
			Assets: lo.Map(assets, func(name lo.Tuple2[*AssetName, string], _ int) *AssetName {
				return name.A
			}),
			Desc:       i.Desc,
			WardID:     lo.ToPtr(ward.ID),
			WardCode:   lo.ToPtr(ward.Code),
			TargetArea: ward,
		})
	}

	return
}

type plateauDatasetItemSeed struct {
	ID        string
	Name      string
	URL       string
	Format    string
	LOD       *int
	NoTexture *bool
}

func (i plateauDatasetItemSeed) GetID(parentID string) string {
	ids := []string{parentID, i.ID}

	if i.LOD != nil {
		ids = append(ids, fmt.Sprintf("lod%d", *i.LOD))
	}

	if i.NoTexture != nil && *i.NoTexture {
		ids = append(ids, "no_texture")
	}

	return strings.Join(lo.Filter(ids, func(s string, _ int) bool {
		return s != ""
	}), "_")
}

func (i plateauDatasetItemSeed) GetName() string {
	name := i.Name
	var lod, tex string

	if i.LOD != nil {
		lod = fmt.Sprintf("LOD%d", *i.LOD)
	}

	if i.NoTexture != nil && *i.NoTexture {
		if name != "" || lod != "" {
			tex += "（"
		}
		tex += "テクスチャなし"
		if name != "" || lod != "" {
			tex += "）"
		}
	}

	if name != "" && lod != "" {
		name += " "
	}
	return name + lod + tex
}

func plateauDatasetItemSeedFrom(seed plateauDatasetSeed) (items []plateauDatasetItemSeed, warning []string) {
	for i, url := range seed.AssetURLs {
		name := nameWithoutExt(nameFromURL(url))
		assetName := seed.Assets[i]
		if assetName == nil {
			warning = append(warning, fmt.Sprintf("plateau %s %s: invalid asset name: %s", seed.TargetArea.GetCode(), seed.DatasetType.Code, name))
			continue
		}

		var item *plateauDatasetItemSeed
		var w []string

		switch {
		case assetName.Ex.Normal != nil:
			item = plateauDatasetItemSeedFromNormal(url, assetName.Ex.Normal)
		case assetName.Ex.Urf != nil:
			item = plateauDatasetItemSeedFromUrf(url, assetName.Ex.Urf, seed.Dic)
		case assetName.Ex.Fld != nil:
			item = plateauDatasetItemSeedFromFld(url, assetName.Ex.Fld, seed.Dic)
		default:
			warning = append(warning, fmt.Sprintf("plateau %s %s[%d]: invalid asset name ex: %s", seed.TargetArea.GetCode(), seed.DatasetType.Code, i, name))
			return
		}

		if item == nil {
			warning = append(warning, fmt.Sprintf("plateau %s %s: dic entry not found: %s", seed.TargetArea.GetCode(), seed.DatasetType.Code, assetName.Ex.Key()))
			continue
		}

		warning = append(warning, w...)
		if item != nil {
			items = append(items, *item)
		}
	}

	return
}

func plateauDatasetItemSeedFromNormal(url string, ex *AssetNameExNormal) *plateauDatasetItemSeed {
	return &plateauDatasetItemSeed{
		// ID:        id,
		Name:      "", // use default
		URL:       url,
		Format:    ex.Format,
		LOD:       lo.ToPtr(ex.LOD),
		NoTexture: lo.ToPtr(ex.NoTexture),
	}
}

func plateauDatasetItemSeedFromUrf(url string, ex *AssetNameExUrf, dic Dic) *plateauDatasetItemSeed {
	entry := dic.FindEntry(ex.Type, ex.Key())
	if entry == nil {
		return nil
	}

	var notexture *bool
	if ex.Format == "3dtiles" {
		notexture = &ex.NoTexture
	}

	return &plateauDatasetItemSeed{
		// ID:        entry.Name,
		Name:      entry.Description,
		URL:       url,
		Format:    ex.Format,
		LOD:       toPtrIfPresent(ex.LOD),
		NoTexture: notexture,
	}
}

func plateauDatasetItemSeedFromFld(url string, ex *AssetNameExFld, dic Dic) *plateauDatasetItemSeed {
	key := fmt.Sprintf("%s_%s_%d", ex.Admin, ex.River, ex.L)
	entry := dic.FindEntry(ex.Type, key)
	if entry == nil {
		return nil
	}

	return &plateauDatasetItemSeed{
		// ID:     key,
		Name:   entry.Scale,
		URL:    url,
		Format: ex.Format,
	}
}
