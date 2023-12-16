package datacatalogv3

import (
	"fmt"
	"strings"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
)

type plateauDatasetSeed struct {
	AssetURLs  []string
	Assets     []*AssetName
	SubName    string
	Desc       string
	TargetArea plateauapi.Area
	Ward       *plateauapi.Ward
	// common
	DatasetType *plateauapi.PlateauDatasetType
	Dic         Dic
	Area        *areaContext
	Pref        *plateauapi.Prefecture
	City        *plateauapi.City
	Spec        *plateauapi.PlateauSpecMinor
	River       *plateauapi.River
}

func plateauDatasetSeedsFrom(i *PlateauFeatureItem, dt *plateauapi.PlateauDatasetType, area *areaContext, spec *plateauapi.PlateauSpecMinor) (res []plateauDatasetSeed) {
	dic := i.ReadDic()

	if len(i.Items) > 0 {
		for _, item := range i.Items {
			res = append(res, plateauDatasetSeedsFromItem(i, item, dt, dic))
		}
	} else {
		res = plateauDatasetSeedsFromBldg(i, dt, area.Wards)
	}

	for _, seed := range res {
		seed.DatasetType = dt
		seed.Dic = dic
		seed.Area = area
		seed.Pref = area.Pref
		seed.City = area.City
		seed.Spec = spec
		if seed.TargetArea == nil {
			seed.TargetArea = area.City
		}
	}

	return res
}

func plateauDatasetSeedsFromItem(i *PlateauFeatureItem, item PlateauFeatureItemDatum, dt *plateauapi.PlateauDatasetType, dic Dic) (res plateauDatasetSeed) {
	subname := item.Name
	var river *plateauapi.River

	assets := lo.Map(item.Data, func(url string, _ int) *AssetName {
		return ParseAssetName(nameWithoutExt(nameFromURL(url)))
	})

	// TODO: how to get dicName
	dicName := strings.TrimPrefix(item.Key, dt.Code+"/")

	if e := dic.FindEntry(dt.Code, dicName); e != nil {
		if subname == "" {
			subname = datasetSubNameFromDicEntry(e)
		}
		if a := riverAdminFrom(e.Admin); a != nil {
			// it's river
			river = &plateauapi.River{
				Name:  e.Description,
				Admin: *a,
			}
		}
	}

	res = plateauDatasetSeed{
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

func plateauDatasetSeedsFromBldg(i *PlateauFeatureItem, dt *plateauapi.PlateauDatasetType, wards []*plateauapi.Ward) (res []plateauDatasetSeed) {
	if len(wards) == 0 {
		res = append(res, plateauDatasetSeed{
			AssetURLs: i.Data,
			Desc:      i.Desc,
		})
		return
	}

	assets := lo.Zip2(lo.Map(i.Data, func(url string, _ int) *AssetName {
		return ParseAssetName(nameWithoutExt(nameFromURL(url)))
	}), i.Data)

	for _, ward := range wards {
		wardCode := ward.Code.String()
		assets := lo.Filter(assets, func(name lo.Tuple2[*AssetName, string], _ int) bool {
			return name.A != nil && name.A.Ex.Normal != nil && name.A.Ex.Normal.WardCode == wardCode
		})
		if len(assets) == 0 {
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
			Ward:       ward,
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

func (i plateauDatasetItemSeed) GetID(parent string) plateauapi.ID {
	return plateauapi.NewID(fmt.Sprintf("%s_%s", parent, i.ID), plateauapi.TypeDatasetItem)
}

func (i plateauDatasetItemSeed) GetName() string {
	if i.Name != "" {
		return i.Name
	}

	return i.StandardName()
}

func (i plateauDatasetItemSeed) StandardName() string {
	var lod, tex string

	if i.LOD != nil {
		lod = fmt.Sprintf("LOD%d", *i.LOD)
	}

	if i.NoTexture != nil && *i.NoTexture {
		if lod != "" {
			lod += "（"
		}
		tex = "テクスチャなし"
		if lod != "" {
			lod += "）"
		}
	}

	return lod + tex
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
		case assetName.Ex.Tnm != nil:
			item = plateauDatasetItemSeedFromTnm(url, assetName.Ex.Tnm, seed.Dic)
		}

		if item == nil {
			warning = append(warning, fmt.Sprintf("plateau %s %s: dic entry not found: %s", seed.TargetArea.GetCode(), assetName.Ex.Type(), name))
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
	id := fmt.Sprintf("lod%d", ex.LOD)
	if ex.NoTexture {
		id += "_notexture"
	}

	return &plateauDatasetItemSeed{
		ID:        id,
		Name:      "", // use default
		URL:       url,
		Format:    ex.Format,
		LOD:       lo.ToPtr(ex.LOD),
		NoTexture: lo.ToPtr(ex.NoTexture),
	}
}

func plateauDatasetItemSeedFromUrf(url string, ex *AssetNameExUrf, dic Dic) *plateauDatasetItemSeed {
	entry := dic.FindEntry(ex.Type, ex.Name)
	if entry == nil {
		return nil
	}

	return &plateauDatasetItemSeed{
		ID:     entry.Name,
		Name:   entry.Description,
		URL:    url,
		Format: ex.Format,
		LOD:    toPtrIfPresent(ex.LOD),
	}
}

func plateauDatasetItemSeedFromTnm(url string, ex *AssetNameExTnm, dic Dic) *plateauDatasetItemSeed {
	entry := dic.FindEntry(ex.Type, ex.Name)
	if entry == nil {
		return nil
	}

	return &plateauDatasetItemSeed{
		ID:     entry.Name,
		Name:   entry.Description,
		URL:    url,
		Format: ex.Format,
	}
}

func plateauDatasetItemSeedFromFld(url string, ex *AssetNameExFld, dic Dic) *plateauDatasetItemSeed {
	key := fmt.Sprintf("%s_%s_%d", ex.Admin, ex.River, ex.L)
	entry := dic.FindEntry(ex.Type, key)
	if entry == nil {
		return nil
	}

	return &plateauDatasetItemSeed{
		ID:     key,
		Name:   entry.Scale,
		URL:    url,
		Format: ex.Format,
	}
}
