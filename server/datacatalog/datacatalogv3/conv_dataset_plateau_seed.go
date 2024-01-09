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
	Admin       any
}

func (seed plateauDatasetSeed) GetID() string {
	return standardItemID(seed.DatasetType.Code, seed.TargetArea, seed.IDEx)
}

func plateauDatasetSeedsFrom(i *PlateauFeatureItem, dt *plateauapi.PlateauDatasetType, area *areaContext, spec *plateauapi.PlateauSpecMinor, cmsurl string) (res []plateauDatasetSeed, warning []string) {
	cityCode := lo.FromPtr(area.CityCode).String()

	dic, err := i.ReadDic()
	if err != nil && i.Dic != "" {
		warning = append(warning, fmt.Sprintf("plateau %s %s: invalid dic: %s", cityCode, dt.Code, err))
		return
	}

	items := i.Items
	if len(i.Data) > 0 {
		items = append(items, lo.Map(i.Data, func(url string, _ int) PlateauFeatureItemDatum {
			return PlateauFeatureItemDatum{
				Data: []string{url},
				Desc: i.Desc,
			}
		})...)
	}

	if dt.Code == "bldg" {
		seeds, w := plateauDatasetSeedsFromBldg(i, dt, cityCode, area.Wards)
		warning = append(warning, w...)
		res = append(res, seeds...)
	} else {
		for _, item := range items {
			seeds, w := plateauDatasetSeedsFromItem(i, item, dt, dic, cityCode)
			warning = append(warning, w...)
			res = append(res, seeds)
		}
	}

	for i := range res {
		res[i].DatasetType = dt
		res[i].Dic = dic
		res[i].Area = area
		res[i].Pref = area.Pref
		res[i].City = area.City
		res[i].Spec = spec
		res[i].Admin = adminFrom(area.CityItem, cmsurl, dt.Code)
		if res[i].TargetArea == nil {
			res[i].TargetArea = area.City
		}
	}

	return
}

func plateauDatasetSeedsFromItem(i *PlateauFeatureItem, item PlateauFeatureItemDatum, dt *plateauapi.PlateauDatasetType, dic Dic, cityCode string) (res plateauDatasetSeed, warning []string) {
	assets := lo.FilterMap(item.Data, func(url string, _ int) (*AssetName, bool) {
		n := nameWithoutExt(nameFromURL(url))
		an := ParseAssetName(n)
		if an == nil || !an.Ex.IsValid() {
			warning = append(warning, fmt.Sprintf("plateau %s %s: invalid asset name: %s", cityCode, dt.Code, n))
		}
		return an, an != nil
	})
	if len(assets) == 0 {
		// warning = append(warning, fmt.Sprintf("plateau %s %s: no assets", cityCode, dt.Code))
		return
	}

	assetName := assets[0]
	key, dickey := assetName.Ex.ItemKey(), assetName.Ex.DicKey()
	var e *DicEntry

	if dickey != "" {
		var found bool
		e, found = dic.FindEntryOrDefault(dt.Code, dickey)
		if !found {
			warning = append(warning, fmt.Sprintf("plateau %s %s: unknown dic key: %s", cityCode, dt.Code, dickey))
			if e == nil {
				return
			}
		}
	}

	var river *plateauapi.River
	if assetName.Ex.Fld != nil {
		if a := riverAdminFrom(assetName.Ex.Fld.Admin); a != nil {
			if e == nil || e.Description == "" {
				warning = append(warning, fmt.Sprintf("plateau %s %s: dic entry has no description or entry not found: %s", cityCode, dt.Code, key))
			} else {
				river = &plateauapi.River{
					Name:  e.Description,
					Admin: *a,
				}
			}
		}
	}

	subname := item.Name
	if subname == "" && e != nil {
		if river != nil {
			// fld
			subname = fmt.Sprintf("%s（%s管理区間）", e.Description, toRiverAdminName(river.Admin))
		} else {
			subname = e.Description
		}
	}
	if subname == "" && e != nil {
		warning = append(warning, fmt.Sprintf("plateau %s %s: invalid dic entry: %s", cityCode, dt.Code, key))
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

func plateauDatasetSeedsFromBldg(i *PlateauFeatureItem, dt *plateauapi.PlateauDatasetType, cityCode string, wards []*plateauapi.Ward) (res []plateauDatasetSeed, warning []string) {
	assets := lo.Zip2(lo.Map(i.Data, func(url string, ind int) *AssetName {
		n := nameWithoutExt(nameFromURL(url))
		an := ParseAssetName(n)
		if an == nil || an.Ex.Normal == nil {
			warning = append(warning, fmt.Sprintf("plateau %s %s[%d]: invalid asset name: %s", cityCode, dt.Code, ind, n))
		}
		return an
	}), i.Data)
	if len(assets) == 0 {
		// warning = append(warning, fmt.Sprintf("plateau %s %s: no assets", cityCode, dt.Code))
		return
	}

	if len(wards) == 0 {
		res = append(res, plateauDatasetSeed{
			AssetURLs: i.Data,
			Assets: lo.Map(assets, func(name lo.Tuple2[*AssetName, string], _ int) *AssetName {
				return name.A
			}),
			Desc: i.Desc,
		})
		return
	}

	for _, ward := range wards {
		wardCode := ward.Code.String()
		assets := lo.Filter(assets, func(name lo.Tuple2[*AssetName, string], _ int) bool {
			return name.A != nil && name.A.Ex.Normal != nil && name.A.Ex.Normal.WardCode == wardCode
		})
		if len(assets) == 0 {
			warning = append(warning, fmt.Sprintf("plateau %s %s: no assets for ward %s", cityCode, dt.Code, wardCode))
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
	cityCode := seed.TargetArea.GetCode().String()

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
			item, w = plateauDatasetItemSeedFromUrf(url, assetName.Ex.Urf, seed.Dic, cityCode)
		case assetName.Ex.Fld != nil:
			item, w = plateauDatasetItemSeedFromFld(url, assetName.Ex.Fld, seed.Dic, cityCode)
		default:
			warning = append(warning, fmt.Sprintf("plateau %s %s[%d]: invalid asset name ex: %s", cityCode, seed.DatasetType.Code, i, name))
			return
		}

		if item == nil {
			warning = append(warning, fmt.Sprintf("plateau %s %s: invalid asset name ex dic key: %s", cityCode, seed.DatasetType.Code, assetName.Ex.DicKey()))
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
		Name:      "", // use default
		URL:       url,
		Format:    ex.Format,
		LOD:       lo.ToPtr(ex.LOD),
		NoTexture: lo.ToPtr(ex.NoTexture),
	}
}

func plateauDatasetItemSeedFromUrf(url string, ex *AssetNameExUrf, dic Dic, cityCode string) (_ *plateauDatasetItemSeed, w []string) {
	entry, found := dic.FindEntryOrDefault(ex.Type, ex.DicKey())
	if !found {
		w = append(w, fmt.Sprintf("plateau %s %s: unknown dic key: %s", cityCode, ex.Type, ex.DicKey()))
	}
	if entry == nil {
		return
	}

	var notexture *bool
	if ex.Format == "3dtiles" {
		notexture = &ex.NoTexture
	}

	return &plateauDatasetItemSeed{
		Name:      entry.Description,
		URL:       url,
		Format:    ex.Format,
		LOD:       toPtrIfPresent(ex.LOD),
		NoTexture: notexture,
	}, w
}

func plateauDatasetItemSeedFromFld(url string, ex *AssetNameExFld, dic Dic, cityCode string) (_ *plateauDatasetItemSeed, w []string) {
	key := ex.Key()
	entry, found := dic.FindEntryOrDefault(ex.Type, ex.DicKey())
	if !found {
		w = append(w, fmt.Sprintf("plateau %s %s: unknown dic key: %s", cityCode, ex.Type, ex.DicKey()))
	}
	if key == "" || entry == nil {
		return
	}

	return &plateauDatasetItemSeed{
		ID:        key,
		Name:      entry.Scale,
		URL:       url,
		Format:    ex.Format,
		NoTexture: &ex.NoTexture,
	}, w
}
