package preparegspatialjp

import (
	"regexp"
	"strconv"
	"strings"

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

	if i.Year == "" {
		i.Year = "2023年度"
	}

	i.References = references
	return
}

func (c *CityItem) YearInt() int {
	return YearInt(c.Year)
}

type GspatialjpItem struct {
	ID                 string   `json:"id,omitempty" cms:"id"`
	CityGML            string   `json:"citygml,omitempty" cms:"citygml,asset"`
	Plateau            string   `json:"plateau,omitempty" cms:"plateau,asset"`
	Related            string   `json:"related,omitempty" cms:"related,asset"`
	Generic            []string `json:"generic,omitempty" cms:"generic,asset"`
	MergeCityGMLStatus *cms.Tag `json:"merge_citygml_status" cms:"merge_citygml_status,tag,metadata"`
	MergePlateauStatus *cms.Tag `json:"merge_plateau_status" cms:"merge_plateau_status,tag,metadata"`
	MergeRelatedStatus *cms.Tag `json:"merge_related_status" cms:"merge_related_status,tag,metadata"`
}

var reReiwa = regexp.MustCompile(`令和([0-9]+?)年度?`)

func YearInt(y string) (year int) {
	if ym := reReiwa.FindStringSubmatch(y); len(ym) > 1 {
		yy, _ := strconv.Atoi(ym[1])
		if yy > 0 {
			year = yy + 2018
		}
	} else if yy, err := strconv.Atoi(strings.TrimSuffix(strings.TrimSuffix(y, "度"), "年")); err == nil {
		year = yy
	}
	return year
}

func SpecVersion(version string) string {
	return strings.TrimSuffix(strings.TrimPrefix(version, "第"), "版")
}