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
	Subname    string
	Subcode    string
	Suborder   *int
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
	LayerNames  LayerNames
	Year        int
}

func (seed plateauDatasetSeed) GetID() string {
	return standardItemID(seed.DatasetType.Code, seed.TargetArea.GetCode(), seed.Subcode)
}

func plateauDatasetSeedsFrom(i *PlateauFeatureItem, dt *plateauapi.PlateauDatasetType, area *areaContext, spec *plateauapi.PlateauSpecMinor, layerNames LayerNames, cmsurl string) (res []plateauDatasetSeed, warning []string) {
	cityCode := lo.FromPtr(area.CityCode).String()
	year := area.CityItem.YearInt()

	dic, err := i.ReadDic()
	if err != nil && i.Dic != "" {
		warning = append(warning, fmt.Sprintf("plateau %s %s: invalid dic: %s", cityCode, dt.Code, err))
		return
	}

	if dt.Code == "bldg" {
		seeds, w := plateauDatasetSeedsFromBldg(i, dt, cityCode, area.Wards)
		warning = append(warning, w...)
		res = append(res, seeds...)
	} else {
		items := i.Items
		if len(i.Data) > 0 {
			items = []PlateauFeatureItemDatum{
				{
					Data:   i.Data,
					Desc:   i.Desc,
					Simple: true,
				},
			}
		}

		for _, item := range items {
			seeds, w := plateauDatasetSeedsFromItem(item, dt, dic, cityCode)
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
		res[i].Admin = newAdmin(area.CityItem.ID, area.CityItem.plateauStage(dt.Code), cmsurl)
		res[i].LayerNames = layerNames
		res[i].Year = year
		if res[i].TargetArea == nil {
			res[i].TargetArea = area.City
		}
	}

	return
}

func plateauDatasetSeedsFromItem(item PlateauFeatureItemDatum, dt *plateauapi.PlateauDatasetType, dic Dic, cityCode string) (res plateauDatasetSeed, warning []string) {
	assets := make([]lo.Tuple2[string, *AssetName], 0, len(item.Data))
	for _, url := range item.Data {
		n := nameWithoutExt(nameFromURL(url))
		an := ParseAssetName(n)
		if an == nil || !an.Ex.IsValid() {
			warning = append(warning, fmt.Sprintf("plateau %s %s: invalid asset name: %s", cityCode, dt.Code, n))
			continue
		}

		assets = append(assets, lo.Tuple2[string, *AssetName]{A: url, B: an})
	}

	if len(assets) == 0 {
		if len(item.Data) > 0 {
			warning = append(warning, fmt.Sprintf("plateau %s %s: some invalid assets", cityCode, dt.Code))
		}
		return
	}

	assetUrls := lo.Map(assets, func(a lo.Tuple2[string, *AssetName], _ int) string {
		return a.A
	})
	assetNames := lo.Map(assets, func(a lo.Tuple2[string, *AssetName], _ int) *AssetName {
		return a.B
	})

	res = plateauDatasetSeed{
		AssetURLs: assetUrls,
		Assets:    assetNames,
		Desc:      item.Desc,
	}

	if !item.Simple {
		assetName := res.Assets[0]
		key, dickey := assetName.Ex.DatasetKey(), assetName.Ex.DicKey()
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

		var suborder *int
		if e != nil {
			suborder = e.Order
		}

		res.Subcode = key
		res.Suborder = suborder
		res.Subname = subname
		res.River = river
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
	ID                  string
	Name                string
	URL                 string
	Format              plateauapi.DatasetFormat
	LOD                 *int
	NoTexture           *bool
	Layers              []string
	FloodingScale       *plateauapi.FloodingScale
	FloodingScaleSuffix *string
}

func (i plateauDatasetItemSeed) GetID(parentID string) string {
	if i.ID != "" {
		parentID = strings.TrimSuffix(parentID, "_"+i.ID)
	}

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
		if assetName.Year != seed.Year {
			warning = append(warning, fmt.Sprintf("plateau %s %s: invalid asset name year: %s: %d should be %d", seed.TargetArea.GetCode(), seed.DatasetType.Code, name, assetName.Year, seed.Year))
		}

		var item *plateauDatasetItemSeed
		var w []string

		switch {
		case assetName.Ex.Normal != nil:
			item, w = plateauDatasetItemSeedFromNormal(url, assetName.Ex.Normal, seed.LayerNames, cityCode)
		case assetName.Ex.Urf != nil:
			item, w = plateauDatasetItemSeedFromUrf(url, assetName.Ex.Urf, seed.Dic, seed.LayerNames, cityCode)
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

func plateauDatasetItemSeedFromNormal(url string, ex *AssetNameExNormal, layerNames LayerNames, cityCode string) (res *plateauDatasetItemSeed, w []string) {
	format := datasetFormatFrom(ex.Format)
	if format == "" {
		w = append(w, fmt.Sprintf("plateau %s %s: invalid format: %s", cityCode, ex.Type, ex.Format))
		return
	}

	return &plateauDatasetItemSeed{
		ID:        "",
		Name:      "", // use default
		URL:       assetURLFromFormat(url, format),
		Format:    format,
		LOD:       lo.ToPtr(ex.LOD),
		NoTexture: lo.ToPtr(ex.NoTexture),
		Layers:    layerNames.LayerName(nil, ex.LOD, format),
	}, nil
}

func plateauDatasetItemSeedFromUrf(url string, ex *AssetNameExUrf, dic Dic, layerNames LayerNames, cityCode string) (_ *plateauDatasetItemSeed, w []string) {
	format := datasetFormatFrom(ex.Format)
	if format == "" {
		w = append(w, fmt.Sprintf("plateau %s %s: unknown format: %s", cityCode, ex.Type, ex.Format))
		return
	}

	key := ex.DicKey()

	entry, found := dic.FindEntryOrDefault(ex.Type, ex.DicKey())
	if !found {
		w = append(w, fmt.Sprintf("plateau %s %s: unknown dic key: %s", cityCode, ex.Type, key))
	}
	if entry == nil {
		return
	}

	var notexture *bool
	if ex.Format == "3dtiles" {
		notexture = &ex.NoTexture
	}

	return &plateauDatasetItemSeed{
		ID:        key,
		Name:      entry.Description,
		URL:       assetURLFromFormat(url, format),
		Format:    format,
		LOD:       lo.EmptyableToPtr(ex.LOD),
		NoTexture: notexture,
		Layers:    layerNames.LayerName([]string{key}, ex.LOD, format),
	}, w
}

func plateauDatasetItemSeedFromFld(url string, ex *AssetNameExFld, dic Dic, cityCode string) (_ *plateauDatasetItemSeed, w []string) {
	format := datasetFormatFrom(ex.Format)
	if format == "" {
		w = append(w, fmt.Sprintf("plateau %s %s: unknown format: %s", cityCode, ex.Type, ex.Format))
		return
	}

	key := ex.Key()
	entry, found := dic.FindEntryOrDefault(ex.Type, ex.DicKey())
	if !found {
		w = append(w, fmt.Sprintf("plateau %s %s: unknown dic key: %s", cityCode, ex.Type, ex.DicKey()))
	}
	if key == "" || entry == nil {
		return
	}

	return &plateauDatasetItemSeed{
		ID:                  key,
		Name:                fldItemName(entry),
		URL:                 assetURLFromFormat(url, format),
		Format:              format,
		NoTexture:           &ex.NoTexture,
		FloodingScale:       toFloodingScale(entry.Scale),
		FloodingScaleSuffix: lo.EmptyableToPtr(entry.SuffixDescription),
	}, w
}

func fldItemName(e *DicEntry) string {
	suffix := ""
	if e.SuffixDescription != "" {
		suffix = e.SuffixDescription
	} else if e.Suffix != "" {
		suffix = e.Suffix
	}
	if suffix != "" {
		suffix = fmt.Sprintf("（%s）", suffix)
	}
	return e.Scale + suffix
}

func toFloodingScale(s string) *plateauapi.FloodingScale {
	switch s {
	case "計画規模":
		fallthrough
	case "l1":
		fallthrough
	case "L1":
		return lo.ToPtr(plateauapi.FloodingScalePlanned)
	case "想定最大規模":
		fallthrough
	case "l2":
		fallthrough
	case "L2":
		return lo.ToPtr(plateauapi.FloodingScaleExpectedMaximum)
	}
	return nil
}
