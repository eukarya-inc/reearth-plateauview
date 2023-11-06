package cmsintegrationv3

import (
	cms "github.com/reearth/reearth-cms-api/go"
)

const modelPrefix = "plateau-"
const cityModel = "city"
const relatedModel = "relevant"
const genericModel = "usecase"
const gspatialjpIndexModel = "g-center-index"
const gspatialjpDataModel = "g-center-data"

var featureTypes = []string{
	"bldg", // 建築物モデル
	"tran", // 交通（道路）モデル
	"rwy",  // 交通（鉄道）モデル
	"trk",  // 交通（徒歩道）モデル
	"squr", // 交通（広場）モデル
	"wwy",  // 交通（航路）モデル
	"luse", // 土地利用モデル
	"fld",  // 洪水浸水想定区域モデル
	"tnm",  // 津波浸水想定区域モデル
	"htd",  // 高潮浸水想定区域モデル
	"ifld", // 内水浸水想定区域モデル
	"lsld", // 災害リスク（土砂災害）モデル
	"urf",  // 都市計画決定情報モデル
	"brid", // 橋梁モデル
	"tun",  // トンネルモデル
	"cons", // その他の構造物モデル
	"frn",  // 都市設備モデル
	"ubld", // 地下街モデル
	"veg",  // 植生モデル
	"dem",  // 地形モデル
	"wtr",  // 水部モデル
	"area", // 区域モデル
	"gen",  // 汎用都市オブジェクトモデル
	"app",  // アピアランスモデル
}

var relatedDataTypes = []string{
	"shelter",
	"park",
	"landmark",
	"station",
	"railway",
	"emergency_route",
	"railway",
	"border",
}

type ManagementStatus string

const (
	ManagementStatusNotStarted ManagementStatus = "登録未着手"
	ManagementStatusRunning    ManagementStatus = "新規登録中"
	ManagementStatusSkip       ManagementStatus = "対象外"
	ManagementStatusDone       ManagementStatus = "登録済み"
	ManagementStatusReady      ManagementStatus = "確認可能"
)

type ConvertionStatus string

const (
	ConvertionStatusNotStarted ConvertionStatus = "未実行"
	ConvertionStatusRunning    ConvertionStatus = "実行中"
	ConvertionStatusError      ConvertionStatus = "エラー"
	ConvertionStatusSuccess    ConvertionStatus = "成功"
)

type CityItem struct {
	ID                   string            `json:"id,omitempty" cms:"id"`
	Prefecture           string            `json:"prefecture,omitempty" cms:"prefecture,select"`
	CityName             string            `json:"city_name,omitempty" cms:"city_name,text"`
	SpecificationVersion string            `json:"spec,omitempty" cms:"spec,select"`
	OpenDataUrl          string            `json:"open_data_url,omitempty" cms:"open_data_url,url"`
	PRCS                 string            `json:"prcs,omitempty" cms:"prcs,select"`
	CodeList             string            `json:"code_list,omitempty" cms:"code_list,asset"`
	Schemas              string            `json:"schemas,omitempty" cms:"schemas,asset"`
	Metadata             string            `json:"metadata,omitempty" cms:"metadata,asset"`
	Specification        string            `json:"specification,omitempty" cms:"specification,asset"`
	References           map[string]string `json:"references,omitempty" cms:"-"`
	// meatadata
	PlateauDataStatus string          `json:"plateau_data_status,omitempty" cms:"plateau-data-status,select,metadata"`
	CityPublic        bool            `json:"city_public,omitempty" cms:"city-public,bool,metadata"`
	SDKPublic         bool            `json:"sdk_public,omitempty" cms:"sdk-public,bool,metadata"`
	Public            map[string]bool `json:"public,omitempty" cms:"-"`
}

func CityItemFrom(item cms.Item) (i CityItem) {
	item.Unmarshal(&i)

	references := map[string]string{}
	public := map[string]bool{}
	for _, ft := range featureTypes {
		if asset := item.MetadataFieldByKey(ft).GetValue().String(); asset != nil {
			references[ft] = *asset
		}

		if pub := item.MetadataFieldByKey(ft + "-public").GetValue().Bool(); pub != nil {
			public[ft] = *pub
		}
	}

	i.References = references
	i.Public = public
	return
}

func (i CityItem) Fields() (fields []*cms.Field) {
	item := &cms.Item{}
	cms.Marshal(i, item)

	for _, ft := range featureTypes {
		if ref, ok := i.References[ft]; ok {
			item.MetadataFields = append(item.MetadataFields, &cms.Field{
				Key:   ft,
				Type:  "reference",
				Value: ref,
			})
		}

		if pub, ok := i.Public[ft]; ok {
			item.MetadataFields = append(item.MetadataFields, &cms.Field{
				Key:   ft + "-public",
				Type:  "bool",
				Value: pub,
			})
		}
	}

	return item.Fields
}

type Item struct {
	ID          string   `json:"id,omitempty" cms:"id"`
	City        string   `json:"city,omitempty" cms:"city,reference"`
	CityGML     string   `json:"citygml,omitempty" cms:"citygml,asset"`
	Data        []string `json:"data,omitempty" cms:"data,asset"`
	Desc        string   `json:"desc,omitempty" cms:"desc,textarea"`
	Rivers      []River  `json:"rivers,omitempty" cms:"rivers,group"`
	QAResult    string   `json:"qaresult,omitempty" cms:"qa-result,asset"`
	SearchIndex string   `json:"searchindex,omitempty" cms:"search-index,asset"`
	MaxLOD      string   `json:"maxlod,omitempty" cms:"maxlod,asset"`
	// metadata
	Status            ManagementStatus `json:"status,omitempty" cms:"status,select,metadata"`
	SkipQA            bool             `json:"skip_qa,omitempty" cms:"skip-qa,bool,metadata"`
	SkipCovnert       bool             `json:"skip_convert,omitempty" cms:"skip-convert,bool,metadata"`
	QAStatus          ConvertionStatus `json:"qa_status,omitempty" cms:"qa-status,select,metadata"`
	ConvertStatus     ConvertionStatus `json:"convert_status,omitempty" cms:"convert-status,select,metadata"`
	SearchIndexStatus ConvertionStatus `json:"search_index_status,omitempty" cms:"search-index-status,select,metadata"`
}

type River struct {
	ID   string   `json:"id,omitempty" cms:"id"`
	Data []string `json:"data,omitempty" cms:"data,asset"`
	Desc string   `json:"desc,omitempty" cms:"desc,textarea"`
}

func ItemFrom(item cms.Item) (i Item) {
	item.Unmarshal(&i)
	return
}

func (i Item) Fields() (fields []*cms.Field) {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item.Fields
}

type GenericItem struct {
	ID          string               `json:"id,omitempty" cms:"id"`
	Prefecture  string               `json:"prefecture,omitempty" cms:"prefecture,select"`
	CityName    string               `json:"city-name,omitempty" cms:"city_name,text"`
	Name        string               `json:"name,omitempty" cms:"name,text"`
	Type        string               `json:"type,omitempty" cms:"type,text"`
	TypeEn      string               `json:"type_en,omitempty" cms:"type-en,text"`
	Datasets    []GenericItemDataset `json:"datasets,omitempty" cms:"datasets,group"`
	OpenDataUrl string               `json:"open-data-url,omitempty" cms:"open-data-url,url"`
	Year        string               `json:"year,omitempty" cms:"year,select"`
	// metadata
	Status ManagementStatus `json:"status,omitempty" cms:"status,select,metadata"`
	Public bool             `json:"public,omitempty" cms:"public,bool,metadata"`
	UseAR  bool             `json:"use-ar,omitempty" cms:"use-ar,bool,metadata"`
}

type GenericItemDataset struct {
	ID         string `json:"id,omitempty" cms:"id"`
	Data       string `json:"data,omitempty" cms:"data,asset"`
	Desc       string `json:"desc,omitempty" cms:"desc,textarea"`
	DataURL    string `json:"url,omitempty" cms:"data-url,url"`
	LayerName  string `json:"layer-name,omitempty" cms:"layer-name,text"`
	DataFormat string `json:"data-format,omitempty" cms:"data-format,select"`
}

func GenericItemFrom(item cms.Item) (i GenericItem) {
	item.Unmarshal(&i)
	return
}

func (i GenericItem) Fields() (fields []*cms.Field) {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item.Fields
}

type RelatedItem struct {
	ID         string            `json:"id,omitempty" cms:"id"`
	Prefecture string            `json:"prefecture,omitempty" cms:"prefecture,select"`
	CityName   string            `json:"city-name,omitempty" cms:"city_name,text"`
	Assets     map[string]string `json:"assets,omitempty" cms:"-"`
	// metadata
	Status map[string]ManagementStatus `json:"status,omitempty" cms:"-"`
	Public map[string]bool             `json:"public,omitempty" cms:"-"`
}

func RelatedItemFrom(item cms.Item) (i RelatedItem) {
	item.Unmarshal(&i)

	for _, t := range relatedDataTypes {
		if asset := item.MetadataFieldByKey(t).GetValue().String(); asset != nil {
			i.Assets[t] = *asset
		}

		if pub := item.MetadataFieldByKey(t + "-public").GetValue().Bool(); pub != nil {
			i.Public[t] = *pub
		}
	}

	return
}

func (i RelatedItem) Fields() (fields []*cms.Field) {
	item := &cms.Item{}
	cms.Marshal(i, item)

	for _, t := range relatedDataTypes {
		if asset, ok := i.Assets[t]; ok {
			item.MetadataFields = append(item.MetadataFields, &cms.Field{
				Key:   t,
				Type:  "asset",
				Value: asset,
			})
		}

		if pub, ok := i.Public[t]; ok {
			item.MetadataFields = append(item.MetadataFields, &cms.Field{
				Key:   t + "-public",
				Type:  "bool",
				Value: pub,
			})
		}
	}

	return item.Fields
}

type GeospatialjpIndexItem struct {
	ID        string `json:"id,omitempty" cms:"id"`
	City      string `json:"city,omitempty" cms:"city,reference"`
	Title     string `json:"title,omitempty" cms:"title,text"`
	Desc      string `json:"desc,omitempty" cms:"desc,markdown"`
	Region    string `json:"region,omitempty" cms:"region,text"`
	Thumbnail string `json:"thumbnail,omitempty" cms:"thumbnail,asset"`
	// metadata
	Status ManagementStatus `json:"status,omitempty" cms:"status,select,metadata"`
}

func GeospatialjpIndexItemFrom(item cms.Item) (i GeospatialjpIndexItem) {
	item.Unmarshal(&i)
	return
}

func (i GeospatialjpIndexItem) Fields() (fields []*cms.Field) {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item.Fields
}

type GeospatialjpDataItem struct {
	ID            string `json:"id,omitempty" cms:"id"`
	City          string `json:"city,omitempty" cms:"city,reference"`
	CityGML       string `json:"citygml,omitempty" cms:"citygml,asset"`
	ConvertedData string `json:"converted-data,omitempty" cms:"converted-data,asset"`
	RelatedData   string `json:"related-data,omitempty" cms:"related-data,asset"`
	GenericData   string `json:"generic-data,omitempty" cms:"generic-data,asset"`
	// metadata
	CityGMLMergeStatus        ConvertionStatus `json:"citygml-merge-status,omitempty" cms:"citygml-merge-status,select,metadata"`
	ConvertedDataMergeSatatus ConvertionStatus `json:"converted-data-merge-status,omitempty" cms:"converted-data-merge-status,select,metadata"`
	RelatedDataMergeStatus    ConvertionStatus `json:"related-data-merge-status,omitempty" cms:"related-data-merge-status,select,metadata"`
}

func GeospatialjpDataItemFrom(item cms.Item) (i GeospatialjpDataItem) {
	item.Unmarshal(&i)
	return
}

func (i GeospatialjpDataItem) Fields() (fields []*cms.Field) {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item.Fields
}
