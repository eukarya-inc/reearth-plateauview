package datacatalog

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/eukarya-inc/jpareacode"
	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/samber/lo"
	"github.com/spkg/bom"
	"golang.org/x/exp/slices"
)

//go:embed urf.csv
var urfFeatureTypesData []byte
var urfFeatureTypeMap map[string]string
var urfFeatureTypes []string

func init() {
	r := csv.NewReader(bom.NewReader(bytes.NewReader(urfFeatureTypesData)))
	d := lo.Must(r.ReadAll())
	urfFeatureTypes = make([]string, 0, len(d)-1)
	for _, c := range d[1:] {
		urfFeatureTypes = append(urfFeatureTypes, c[0])
	}
	urfFeatureTypeMap = lo.SliceToMap(d[1:], func(c []string) (string, string) {
		return c[0], c[1]
	})
}

type DataCatalogItemBuilder struct {
	Assets           []*cms.PublicAsset
	Descriptions     []string
	IntermediateItem PlateauIntermediateItem
	Options          DataCatalogItemBuilderOption
}

type DataCatalogItemBuilderOption struct {
	ModelName           string
	NameOverride        string
	NameOverrideBy      func(AssetName) (string, string, string)
	MultipleDesc        bool
	ItemID              bool
	LOD                 bool
	Layers              []string
	LayersForLOD        map[string][]string
	UseMaxLODAsDefault  bool
	UseGroupNameAsName  bool
	UseGroupNameAsLayer bool
	GroupBy             func(AssetName) string
	SortGroupBy         func(string, string, AssetName, AssetName) bool
}

func (b DataCatalogItemBuilder) Build() []*DataCatalogItem {
	if len(b.Assets) == 0 {
		return nil
	}

	type asset struct {
		Index int
		URL   string
		Name  AssetName
	}

	type assetGroup struct {
		Name   string
		Assets []asset
	}

	assets := lo.Map(b.Assets, func(a *cms.PublicAsset, i int) asset {
		return asset{
			Index: i,
			URL:   a.URL,
			Name:  AssetNameFrom(a.URL),
		}
	})

	// create groups
	var groups []assetGroup
	if b.Options.GroupBy != nil {
		groups = lo.MapToSlice(lo.GroupBy(assets, func(a asset) string {
			return b.Options.GroupBy(a.Name)
		}), func(k string, a []asset) assetGroup {
			return assetGroup{
				Name:   k,
				Assets: a,
			}
		})

		// sort groups
		sort.SliceStable(groups, func(i, j int) bool {
			if b.Options.SortGroupBy != nil {
				return b.Options.SortGroupBy(groups[i].Name, groups[j].Name, groups[i].Assets[0].Name, groups[j].Assets[0].Name)
			}

			// sort by asset index
			return groups[i].Assets[0].Index < groups[j].Assets[0].Index
		})
	} else {
		groups = []assetGroup{{Assets: assets}}
	}

	results := make([]*DataCatalogItem, 0, len(groups))

	for i, g := range groups {
		itemID := b.Options.ItemID && i == 0

		// desc and name override
		name, desc := "", ""
		if b.Options.NameOverride != "" {
			name = b.Options.NameOverride
		}
		if b.Options.MultipleDesc {
			an, ad := descFromAsset(g.Assets[0].URL, b.Descriptions)
			if an != "" && name == "" {
				name = an
			}
			if name == "" && b.Options.UseGroupNameAsName {
				name = g.Name
			}
			if ad != "" {
				desc = ad
			}
		} else if len(b.Descriptions) > 0 {
			desc = b.Descriptions[0]
		}

		// default asset and its URL
		defaultAsset := g.Assets[0]
		if b.Options.UseMaxLODAsDefault {
			defaultAsset = lo.MaxBy(g.Assets, func(a, b asset) bool {
				return a.Name.LODInt() > b.Name.LODInt()
			})
		}

		// default layers
		layersForLOD := b.Options.LayersForLOD
		defaultLayers := b.Options.Layers
		if b.Options.UseGroupNameAsLayer {
			defaultLayers = []string{g.Name}
		}
		if b.Options.UseGroupNameAsLayer || layersForLOD == nil {
			layersForLOD = map[string][]string{"": defaultLayers}
		}

		var mainDefaultLayers []string
		if isLayerSupported(defaultAsset.Name.Format) {
			mainDefaultLayers = defaultLayers
		}

		dci := b.IntermediateItem.DataCatalogItem(
			b.Options.ModelName,
			defaultAsset.Name,
			defaultAsset.URL,
			desc,
			mainDefaultLayers,
			itemID,
			name,
			b.Options.NameOverrideBy,
		)
		if dci != nil && b.Options.LOD {
			assetURLs := lo.Map(g.Assets, func(a asset, _ int) string {
				return a.URL
			})

			mn := b.Options.ModelName
			if name != "" {
				mn = name
			} else if b.Options.UseGroupNameAsName && g.Name != "" {
				mn = g.Name
			}

			dci.Config = multipleLODData(assetURLs, mn, layersForLOD)
		}

		results = append(results, dci)
	}

	return results
}

type PlateauIntermediateItem struct {
	ID          string
	Prefecture  string
	City        string
	CityEn      string
	CityCode    string
	Dic         Dic
	OpenDataURL string
	Year        int
}

func (i PlateauItem) IntermediateItem() PlateauIntermediateItem {
	au := ""
	if i.CityGML != nil {
		au = i.CityGML.URL
	} else if len(i.Bldg) > 0 {
		au = i.Bldg[0].URL
	}

	if au == "" {
		return PlateauIntermediateItem{}
	}

	an := AssetNameFrom(au)
	dic := Dic{}
	_ = json.Unmarshal(bom.Clean([]byte(i.Dic)), &dic)
	y, _ := strconv.Atoi(an.Year)

	return PlateauIntermediateItem{
		ID:          i.ID,
		Prefecture:  i.Prefecture,
		City:        i.CityName,
		CityEn:      an.CityEn,
		CityCode:    an.CityCode,
		Dic:         dic,
		OpenDataURL: i.OpenDataURL,
		Year:        y,
	}
}

func (i *PlateauIntermediateItem) DataCatalogItem(t string, an AssetName, assetURL, desc string, layers []string, addItemID bool, nameOverride string, nameOverrideBy func(AssetName) (string, string, string)) *DataCatalogItem {
	if i == nil {
		return nil
	}

	id := i.id(an)
	if id == "" {
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

	// name
	name := ""
	t2, t2en := "", ""
	if name == "" && nameOverrideBy != nil {
		name, t2, t2en = nameOverrideBy(an)
	}
	if name == "" {
		name = nameOverride
	}
	if name == "" {
		name = t
	}

	prefCode := jpareacode.PrefectureCodeInt(i.Prefecture)

	var itemID string
	if addItemID {
		itemID = i.ID
	}

	opd := i.OpenDataURL
	if opd == "" {
		opd = openDataURLFromAssetName(an)
	}

	return &DataCatalogItem{
		ID:          id,
		ItemID:      itemID,
		Type:        t,
		TypeEn:      an.Feature,
		Type2:       t2,
		Type2En:     t2en,
		Name:        fmt.Sprintf("%s（%s）", name, cityOrWardName),
		Pref:        i.Prefecture,
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
		Year:        i.Year,
		Layers:      layers,
		OpenDataURL: opd,
	}
}

func (i *PlateauIntermediateItem) id(an AssetName) string {
	return strings.Join(lo.Filter([]string{
		i.CityCode,
		i.CityEn,
		an.WardCode,
		an.WardEn,
		an.Feature,
		an.UrfFeatureType,
		an.FldNameAndCategory(),
		an.GenName,
	}, func(s string, _ int) bool { return s != "" }), "_")
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

var reName = regexp.MustCompile(`^@name:\s*(.+)(?:$|\n)`)

func descFromAsset(assetURL string, descs []string) (string, string) {
	if assetURL == "" || len(descs) == 0 {
		return "", ""
	}

	fn := strings.TrimSuffix(path.Base(assetURL), path.Ext(assetURL))
	for _, desc := range descs {
		b, a, ok := strings.Cut(desc, "\n")
		if ok && strings.Contains(b, fn) {
			return nameFromDescription(strings.TrimSpace(a))
		}
	}

	return "", ""
}

func nameFromDescription(d string) (string, string) {
	if m := reName.FindStringSubmatch(d); len(m) > 0 {
		name := m[1]
		_, n, _ := strings.Cut(d, "\n")
		return name, strings.TrimSpace(n)
	}

	return "", d
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

func urfLayers(ty string) []string {
	if ty == "WaterWay" {
		ty = "Waterway"
	}
	return []string{ty}
}

func openDataURLFromAssetName(a AssetName) string {
	return fmt.Sprintf("https://www.geospatial.jp/ckan/dataset/plateau-%s-%s-%s", a.CityCode, a.CityEn, a.Year)
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

func multipleLODData(assets []string, modelName string, layers map[string][]string) DataCatalogItemConfig {
	data := lo.Map(assets, func(a string, j int) DataCatalogItemConfigItem {
		an := AssetNameFrom(a)
		name := ""
		if an.LOD != "" {
			name = fmt.Sprintf("LOD%s", an.LOD)
		} else if len(assets) == 1 {
			name = modelName
		} else {
			name = fmt.Sprintf("%s%d", modelName, j+1)
		}

		var l []string
		if layers != nil && isLayerSupported(an.Format) {
			l = slices.Clone(layers[an.LOD])
			if l == nil {
				l = slices.Clone(layers[""])
			}
		}

		return DataCatalogItemConfigItem{
			Name:   name,
			URL:    assetURLFromFormat(a, an.Format),
			Type:   an.Format,
			Layers: l,
		}
	})

	sort.Slice(data, func(a, b int) bool {
		return strings.Compare(data[a].Name, data[b].Name) < 0
	})

	return DataCatalogItemConfig{
		Data: data,
	}
}
