package datacatalog

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/eukarya-inc/jpareacode"
	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
	"github.com/spkg/bom"
)

//go:embed urf.csv
var urfFeatureTypesData []byte
var urfFeatureTypes map[string]string

func init() {
	r := csv.NewReader(bom.NewReader(bytes.NewReader(urfFeatureTypesData)))
	d := lo.Must(r.ReadAll())
	urfFeatureTypes = lo.SliceToMap(d[1:], func(c []string) (string, string) {
		return c[0], c[1]
	})
}

type PlateauItem struct {
	ID              string             `json:"id"`
	Prefecture      string             `json:"prefecture"`
	CityName        string             `json:"city_name"`
	Specification   string             `json:"specification"`
	CityGML         *cms.PublicAsset   `json:"citygml"`
	DescriptionBldg string             `json:"description_bldg"`
	DescriptionTran string             `json:"description_tran"`
	DescriptionFrn  string             `json:"description_frn"`
	DescriptionVeg  string             `json:"description_veg"`
	DescriptionLuse string             `json:"description_luse"`
	DescriptionLsld string             `json:"description_lsld"`
	DescriptionUrf  []string           `json:"description_urf"`
	DescriptionFld  []string           `json:"description_fld"`
	DescriptionHtd  []string           `json:"description_htd"`
	DescriptionIfld []string           `json:"description_ifld"`
	DescriptionTnm  []string           `json:"description_tnm"`
	DescriptionBrid string             `json:"description_brid"`
	Bldg            []*cms.PublicAsset `json:"bldg"`
	Tran            []*cms.PublicAsset `json:"tran"`
	Frn             []*cms.PublicAsset `json:"frn"`
	Veg             []*cms.PublicAsset `json:"veg"`
	Luse            []*cms.PublicAsset `json:"luse"`
	Lsld            []*cms.PublicAsset `json:"lsld"`
	Urf             []*cms.PublicAsset `json:"urf"`
	Fld             []*cms.PublicAsset `json:"fld"`
	Htd             []*cms.PublicAsset `json:"htd"`
	Ifld            []*cms.PublicAsset `json:"ifld"`
	Tnm             []*cms.PublicAsset `json:"tnm"`
	Brid            []*cms.PublicAsset `json:"brid"`
	Dictionary      *cms.PublicAsset   `json:"dictionary"`
	Dic             string             `json:"dic"`
	SearchIndex     []*cms.PublicAsset `json:"search_index"`
	OpenDataURL     string             `json:"opendata_url"`
}

func (i PlateauItem) FrnItem(c PlateauIntermediateItem) *DataCatalogItem {
	if len(i.Frn) == 0 {
		return nil
	}

	a := i.Frn[0]
	return c.DataCatalogItem("都市設備モデル", AssetNameFrom(a.URL), a.URL, i.DescriptionFrn, nil, false)
}

func (i PlateauItem) VegItem(c PlateauIntermediateItem) *DataCatalogItem {
	if len(i.Veg) == 0 {
		return nil
	}

	a := i.Veg[0]
	return c.DataCatalogItem("植生モデル", AssetNameFrom(a.URL), a.URL, i.DescriptionVeg, nil, false)
}

func (i PlateauItem) LuseItem(c PlateauIntermediateItem) *DataCatalogItem {
	if i.Luse == nil {
		return nil
	}

	a := i.Luse[0]
	return c.DataCatalogItem("土地利用モデル", AssetNameFrom(a.URL), a.URL, i.DescriptionLuse, []string{"luse"}, false)
}

func (i PlateauItem) LsldItem(c PlateauIntermediateItem) *DataCatalogItem {
	if i.Lsld == nil {
		return nil
	}

	a := i.Lsld[0]
	return c.DataCatalogItem("土砂災害警戒区域モデル", AssetNameFrom(a.URL), a.URL, i.DescriptionLsld, []string{"lsld"}, false)
}

func (i PlateauItem) UrfItems(c PlateauIntermediateItem) []*DataCatalogItem {
	if len(i.Urf) == 0 {
		return nil
	}

	return lo.Map(i.Urf, func(a *cms.PublicAsset, _ int) *DataCatalogItem {
		an := AssetNameFrom(a.URL)
		return c.DataCatalogItem("都市計画決定情報モデル", an, a.URL, descFromAsset(a, i.DescriptionUrf), []string{an.UrfFeatureType}, false)
	})
}

func (i PlateauItem) HtdItems(c PlateauIntermediateItem) []*DataCatalogItem {
	if len(i.Htd) == 0 {
		return nil
	}

	return lo.Map(i.Htd, func(a *cms.PublicAsset, _ int) *DataCatalogItem {
		an := AssetNameFrom(a.URL)
		dci := c.DataCatalogItem("高潮浸水想定区域モデル", an, a.URL, descFromAsset(a, i.DescriptionHtd), nil, false)
		dci.Name = htdTnmIfldName("高潮浸水想定区域モデル", i.CityName, an.FldName, c.Dic.Htd(an.FldName))
		return dci
	})
}

func (i PlateauItem) IfldItems(c PlateauIntermediateItem) []*DataCatalogItem {
	if len(i.Ifld) == 0 {
		return nil
	}

	return lo.Map(i.Ifld, func(a *cms.PublicAsset, _ int) *DataCatalogItem {
		an := AssetNameFrom(a.URL)
		desc := descFromAsset(a, i.DescriptionIfld)
		dci := c.DataCatalogItem("内水浸水想定区域モデル", an, a.URL, desc, nil, false)
		dci.Name = htdTnmIfldName("内水浸水想定区域モデル", i.CityName, an.FldName, c.Dic.Ifld(an.FldName))
		return dci
	})
}

func (i PlateauItem) TnmItems(c PlateauIntermediateItem) []*DataCatalogItem {
	if len(i.Tnm) == 0 {
		return nil
	}

	return lo.Map(i.Tnm, func(a *cms.PublicAsset, _ int) *DataCatalogItem {
		an := AssetNameFrom(a.URL)
		dci := c.DataCatalogItem("津波浸水想定区域モデル", an, a.URL, descFromAsset(a, i.DescriptionTnm), nil, false)
		dci.Name = htdTnmIfldName("津波浸水想定区域モデル", i.CityName, an.FldName, c.Dic.Tnm(an.FldName))
		return dci
	})
}

func (i PlateauItem) BridItem(c PlateauIntermediateItem) *DataCatalogItem {
	if len(i.Brid) == 0 {
		return nil
	}

	a := i.Brid[0]
	return c.DataCatalogItem("橋梁モデル", AssetNameFrom(a.URL), a.URL, i.DescriptionBrid, nil, false)
}

func (i PlateauItem) DataCatalogItems() []DataCatalogItem {
	c := i.IntermediateItem()
	if c.ID == "" {
		return nil
	}

	return util.DerefSlice(lo.Filter(
		append(append(append(append(append(append(append(
			i.BldgItems(c),
			i.TranItem(c),
			i.FrnItem(c),
			i.VegItem(c),
			i.LuseItem(c),
			i.LsldItem(c)),
			i.UrfItems(c)...),
			i.FldItems(c)...),
			i.TnmItems(c)...),
			i.HtdItems(c)...),
			i.IfldItems(c)...),
			i.BridItem(c),
		),
		func(i *DataCatalogItem, _ int) bool {
			return i != nil
		},
	))
}

func (i PlateauItem) IntermediateItem() PlateauIntermediateItem {
	if i.CityGML == nil {
		return PlateauIntermediateItem{}
	}

	an := AssetNameFrom(i.CityGML.URL)
	dic := Dic{}
	_ = json.Unmarshal(bom.Clean([]byte(i.Dic)), &dic)

	return PlateauIntermediateItem{
		ID:          i.ID,
		Prefecture:  i.Prefecture,
		City:        i.CityName,
		CityEn:      an.CityEn,
		CityCode:    an.CityCode,
		Dic:         dic,
		OpenDataURL: i.OpenDataURL,
	}
}

type PlateauIntermediateItem struct {
	ID          string
	Prefecture  string
	City        string
	CityEn      string
	CityCode    string
	Dic         Dic
	OpenDataURL string
}

func (i *PlateauIntermediateItem) DataCatalogItem(t string, an AssetName, assetURL, desc string, layers []string, firstWard bool) *DataCatalogItem {
	if i == nil {
		return nil
	}

	wardName := i.Dic.WardName(an.WardCode)
	if wardName == "" && an.WardCode != "" {
		wardName = an.WardEn
	}

	cityOrWardName := i.City
	if wardName != "" {
		cityOrWardName = wardName
	}

	var name, t2, t2en string
	if an.Feature == "urf" {
		urfName := urfFeatureTypes[an.UrfFeatureType]
		if urfName == "" {
			urfName = an.UrfFeatureType
		}
		name = fmt.Sprintf("%sモデル（%s）", urfName, cityOrWardName)
		t2 = urfName
		t2en = an.UrfFeatureType
	} else {
		name = fmt.Sprintf("%s（%s）", t, cityOrWardName)
	}

	id := strings.Join(lo.Filter([]string{
		i.CityCode,
		i.CityEn,
		an.WardCode,
		an.WardEn,
		an.Feature,
		an.UrfFeatureType,
		an.FldNameAndCategory(),
	}, func(s string, _ int) bool { return s != "" }), "_")

	if id == "" {
		return nil
	}

	y, _ := strconv.Atoi(an.Year)
	pref, prefCode := normalizePref(i.Prefecture)

	var itemID string
	if an.Feature == "bldg" && (an.WardCode == "" || firstWard) {
		itemID = i.ID
	}

	return &DataCatalogItem{
		ID:          id,
		ItemID:      itemID,
		Type:        t,
		TypeEn:      an.Feature,
		Type2:       t2,
		Type2En:     t2en,
		Name:        name,
		Pref:        pref,
		PrefCode:    jpareacode.FormatPrefectureCode(prefCode),
		City:        i.City,
		CityEn:      i.CityEn,
		CityCode:    cityCode(i.CityCode, i.City, prefCode),
		Ward:        wardName,
		WardEn:      an.WardEn,
		WardCode:    cityCode(an.WardCode, wardName, prefCode),
		Description: desc,
		URL:         assetURLFromFormat(assetURL, an.Format),
		Format:      an.Format,
		Year:        y,
		Layers:      layers,
		OpenDataURL: i.OpenDataURL,
	}
}

func assetsByWards(a []*cms.PublicAsset) map[string][]*cms.PublicAsset {
	if len(a) == 0 {
		return nil
	}

	r := map[string][]*cms.PublicAsset{}
	for _, a := range a {
		if a == nil {
			continue
		}

		an := AssetNameFrom(a.URL)
		k := an.WardCode
		if _, ok := r[k]; !ok {
			r[k] = []*cms.PublicAsset{a}
		} else {
			r[k] = append(r[k], a)
		}
	}
	return r
}

func descFromAsset(a *cms.PublicAsset, descs []string) string {
	if a == nil || len(descs) == 0 {
		return ""
	}

	fn := strings.TrimSuffix(path.Base(a.URL), path.Ext(a.URL))
	for _, desc := range descs {
		b, a, ok := strings.Cut(desc, "\n")
		if ok && strings.Contains(b, fn) {
			return strings.TrimSpace(a)
		}
	}
	return ""
}

func searchIndexURLFrom(assets []*cms.PublicAsset, wardCode string) string {
	a, found := lo.Find(assets, func(a *cms.PublicAsset) bool {
		if wardCode == "" {
			return true
		}
		return AssetNameFrom(a.URL).WardCode == wardCode
	})
	if !found {
		return ""
	}

	u, err := url.Parse(a.URL)
	if err != nil {
		return ""
	}

	u.Path = path.Join(assetRootPath(u.Path), "indexRoot.json")
	return u.String()
}

type DataCatalogItemConfig struct {
	Data []DataCatalogItemConfigItem `json:"data,omitempty"`
}

type DataCatalogItemConfigItem struct {
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	Type   string   `json:"type"`
	Layers []string `json:"layer,omitempty"`
}

type assetWithLOD struct {
	A   *cms.PublicAsset
	F   AssetName
	LOD int
}

func assetWithLODFromList(a []*cms.PublicAsset) ([]assetWithLOD, int) {
	maxLOD := 0
	return lo.FilterMap(a, func(a *cms.PublicAsset, _ int) (assetWithLOD, bool) {
		l := assetWithLODFrom(a)
		if l != nil && maxLOD < l.LOD {
			maxLOD = l.LOD
		}
		return *l, l != nil
	}), maxLOD
}

func assetWithLODFrom(a *cms.PublicAsset) *assetWithLOD {
	if a == nil {
		return nil
	}
	f := AssetNameFrom(a.URL)
	l, _ := strconv.Atoi(f.LOD)
	return &assetWithLOD{A: a, LOD: l, F: f}
}

func htdTnmIfldName(t, cityName, raw string, e *DicEntry) string {
	if e == nil {
		return raw
	}
	return fmt.Sprintf("%s %s（%s）", t, e.Description, cityName)
}
