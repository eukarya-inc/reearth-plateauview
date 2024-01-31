package geospatialjpv3

import (
	"fmt"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/ckan"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogcommon"
	cms "github.com/reearth/reearth-cms-api/go"
)

var featureTypes = []string{
	// *: データカタログ上で複数の項目に分かれて存在
	"bldg", // 建築物モデル
	"tran", // 交通（道路）モデル
	"rwy",  // 交通（鉄道）モデル
	"trk",  // 交通（徒歩道）モデル
	"squr", // 交通（広場）モデル
	"wwy",  // 交通（航路）モデル
	"luse", // 土地利用モデル
	"fld",  // 洪水浸水想定区域モデル*
	"tnm",  // 津波浸水想定区域モデル*
	"htd",  // 高潮浸水想定区域モデル*
	"ifld", // 内水浸水想定区域モデル*
	"lsld", // 土砂災害モデル
	"urf",  // 都市計画決定情報モデル*
	"unf",  // 地下埋設物モデル
	"brid", // 橋梁モデル
	"tun",  // トンネルモデル
	"cons", // その他の構造物モデル
	"frn",  // 都市設備モデル
	"ubld", // 地下街モデル
	"veg",  // 植生モデル
	"dem",  // 地形モデル
	"wtr",  // 水部モデル
	"area", // 区域モデル*
	"gen",  // 汎用都市オブジェクトモデル*
}

type CityItem struct {
	ID                string            `json:"id,omitempty" cms:"id"`
	Prefecture        string            `json:"prefecture,omitempty" cms:"prefecture,select"`
	CityName          string            `json:"city_name,omitempty" cms:"city_name,text"`
	CityNameEn        string            `json:"city_name_en,omitempty" cms:"city_name_en,text"`
	CityCode          string            `json:"city_code,omitempty" cms:"city_code,text"`
	CodeLists         string            `json:"codelists,omitempty" cms:"codelists,asset"`
	Schemas           string            `json:"schemas,omitempty" cms:"schemas,asset"`
	Metadata          string            `json:"metadata,omitempty" cms:"metadata,asset"`
	Specification     string            `json:"specification,omitempty" cms:"specification,asset"`
	Misc              string            `json:"misc,omitempty" cms:"misc,asset"`
	Year              string            `json:"year,omitempty" cms:"year,select"`
	References        map[string]string `json:"references,omitempty" cms:"-"`
	RelatedDataset    string            `json:"related_dataset,omitempty" cms:"related_dataset,reference"`
	GeospatialjpIndex string            `json:"geospatialjp-index,omitempty" cms:"geospatialjp-index,reference"`
	GeospatialjpData  string            `json:"geospatialjp-data,omitempty" cms:"geospatialjp-data,reference"`
}

func CityItemFrom(item *cms.Item) (i *CityItem) {
	i = &CityItem{}
	item.Unmarshal(i)

	references := map[string]string{}
	for _, ft := range featureTypes {
		if ref := item.FieldByKey(ft).GetValue().String(); ref != nil {
			references[ft] = *ref
		}
	}

	i.References = references
	return
}

func (i *CityItem) YearInt() int {
	return datacatalogcommon.YearInt(i.Year)
}

type CMSDataItem struct {
	ID        string         `json:"id,omitempty" cms:"id"`
	CityGML   map[string]any `json:"citygml,omitempty" cms:"citygml,asset"`
	Plateau   map[string]any `json:"plateau,omitempty" cms:"plateau,asset"`
	Related   map[string]any `json:"related,omitempty" cms:"related,asset"`
	DescIndex string         `json:"desc_index,omitempty" cms:"desc_index,markdown"`
}

type CMSIndexItem struct {
	ID          string         `json:"id,omitempty" cms:"id"`
	Thumbnail   map[string]any `json:"thumbnail,omitempty" cms:"thumbnail,asset"`
	Region      string         `json:"region,omitempty" cms:"region,text"`
	Desc        string         `json:"desc,omitempty" cms:"desc,markdown"`
	DescCityGML string         `json:"desc_citygml,omitempty" cms:"desc_citygml,markdown"`
	DescPlateau string         `json:"desc_plateau,omitempty" cms:"desc_plateau,markdown"`
	DescRelated string         `json:"desc_related,omitempty" cms:"desc_related,markdown"`
	Items       []CMSItem      `json:"items,omitempty" cms:"items,items,group"`
}

type CMSItem struct {
	Name  string         `json:"name,omitempty" cms:"name,text"`
	Desc  string         `json:"desc,omitempty" cms:"desc,markdown"`
	Asset map[string]any `json:"asset,omitempty" cms:"asset,asset"`
}

type PackageName struct {
	CityCode, CityNameEn string
	Year                 int
}

func (p PackageName) String() string {
	return datasetName(p.CityCode, p.CityNameEn, p.Year)
}

type PackageSeed struct {
	Name         PackageName
	NameJa       string
	Description  string
	OwnerOrg     string
	Area         string
	ThumbnailURL string
}

func (p PackageSeed) Title() string {
	return fmt.Sprintf("3D都市モデル（Project PLATEAU）%s（%d年度）", p.NameJa, p.Name.Year)
}

func (p PackageSeed) ToPackage() ckan.Package {
	return ckan.Package{
		Name:         p.Name.String(),
		Title:        p.Title(),
		OwnerOrg:     p.OwnerOrg,
		Notes:        p.Description,
		Private:      true,
		Area:         p.Area,
		ThumbnailURL: p.ThumbnailURL,
	}
}
