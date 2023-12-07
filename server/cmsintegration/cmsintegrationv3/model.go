package cmsintegrationv3

import (
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationcommon"
	cms "github.com/reearth/reearth-cms-api/go"
)

const modelPrefix = "plateau-"
const cityModel = "city"
const relatedModel = "related"
const geospatialjpIndex = "geospatialjp-index"
const geospatialjpData = "geospatialjp-data"

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

var featureTypesWithItems = []string{
	"fld",  // 洪水浸水想定区域モデル*
	"tnm",  // 津波浸水想定区域モデル*
	"htd",  // 高潮浸水想定区域モデル*
	"ifld", // 内水浸水想定区域モデル*
	"urf",  // 都市計画決定情報モデル*
	"gen",  // 汎用都市オブジェクトモデル*
}

var relatedDataTypes = []string{
	"shelter",
	"park",
	"landmark",
	"station",
	"railway",
	"emergency_route",
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
	ID                   string                    `json:"id,omitempty" cms:"id"`
	Prefecture           string                    `json:"prefecture,omitempty" cms:"prefecture,select"`
	CityName             string                    `json:"city_name,omitempty" cms:"city_name,text"`
	CityNameEn           string                    `json:"city_name_en,omitempty" cms:"city_name_en,text"`
	CityCode             string                    `json:"city_code,omitempty" cms:"city_code,text"`
	SpecificationVersion string                    `json:"spec,omitempty" cms:"spec,select"`
	OpenDataUrl          string                    `json:"open_data_url,omitempty" cms:"open_data_url,url"`
	PRCS                 cmsintegrationcommon.PRCS `json:"prcs,omitempty" cms:"prcs,select"`
	CodeLists            string                    `json:"codelists,omitempty" cms:"codelists,asset"`
	Schemas              string                    `json:"schemas,omitempty" cms:"schemas,asset"`
	Metadata             string                    `json:"metadata,omitempty" cms:"metadata,asset"`
	Specification        string                    `json:"specification,omitempty" cms:"specification,asset"`
	References           map[string]string         `json:"references,omitempty" cms:"-"`
	RelatedDataset       string                    `json:"related_dataset,omitempty" cms:"related_dataset,reference"`
	GeospatialjpIndex    string                    `json:"geospatialjp-index,omitempty" cms:"geospatialjp-index,reference"`
	GeospatialjpData     string                    `json:"geospatialjp-data,omitempty" cms:"geospatialjp-data,reference"`
	// meatadata
	PlateauDataStatus string          `json:"plateau_data_status,omitempty" cms:"plateau_data_status,select,metadata"`
	CityPublic        bool            `json:"city_public,omitempty" cms:"city_public,bool,metadata"`
	SDKPublic         bool            `json:"sdk_public,omitempty" cms:"sdk_public,bool,metadata"`
	Public            map[string]bool `json:"public,omitempty" cms:"-"`
}

func CityItemFrom(item *cms.Item) (i *CityItem) {
	i = &CityItem{}
	item.Unmarshal(i)

	references := map[string]string{}
	public := map[string]bool{}
	for _, ft := range featureTypes {
		if ref := item.FieldByKey(ft).GetValue().String(); ref != nil {
			references[ft] = *ref
		}

		if pub := item.MetadataFieldByKey(ft + "_public").GetValue().Bool(); pub != nil {
			public[ft] = *pub
		}
	}

	i.References = references
	i.Public = public
	return
}

func (i *CityItem) CMSItem() *cms.Item {
	item := &cms.Item{}
	cms.Marshal(i, item)

	for _, ft := range featureTypes {
		if ref, ok := i.References[ft]; ok {
			item.Fields = append(item.Fields, &cms.Field{
				Key:   ft,
				Type:  "reference",
				Value: ref,
			})
		}

		if pub, ok := i.Public[ft]; ok {
			item.MetadataFields = append(item.MetadataFields, &cms.Field{
				Key:   ft + "_public",
				Type:  "bool",
				Value: pub,
			})
		}
	}

	return item
}

type FeatureItem struct {
	ID          string             `json:"id,omitempty" cms:"id"`
	City        string             `json:"city,omitempty" cms:"city,reference"`
	CityGML     string             `json:"citygml,omitempty" cms:"citygml,asset"`
	Data        []string           `json:"data,omitempty" cms:"data,asset"`
	Desc        string             `json:"desc,omitempty" cms:"desc,textarea"`
	Items       []FeatureItemDatum `json:"items,omitempty" cms:"items,group"`
	QCResult    string             `json:"qcresult,omitempty" cms:"qc_result,asset"`
	SearchIndex string             `json:"search_index,omitempty" cms:"search_index,asset"`
	Dic         string             `json:"dic,omitempty" cms:"dic,textarea"`
	MaxLOD      string             `json:"maxlod,omitempty" cms:"maxlod,asset"`
	// metadata
	Status            ManagementStatus `json:"status,omitempty" cms:"status,select,metadata"`
	SkipQC            bool             `json:"skip_qc,omitempty" cms:"skip_qc,bool,metadata"`
	SkipConvert       bool             `json:"skip_conv,omitempty" cms:"skip_conv,bool,metadata"`
	ConvertionStatus  ConvertionStatus `json:"conv_status,omitempty" cms:"conv_status,select,metadata"`
	QCStatus          ConvertionStatus `json:"qc_status,omitempty" cms:"qc_status,select,metadata"`
	SearchIndexStatus ConvertionStatus `json:"search_index_status,omitempty" cms:"search_index_status,select,metadata"`
}

type FeatureItemDatum struct {
	ID   string   `json:"id,omitempty" cms:"id"`
	Data []string `json:"data,omitempty" cms:"data,asset"`
	Name string   `json:"name,omitempty" cms:"name,text"`
	Desc string   `json:"desc,omitempty" cms:"desc,textarea"`
	Key  string   `json:"key,omitempty" cms:"key,text"`
}

func FeatureItemFrom(item *cms.Item) (i *FeatureItem) {
	i = &FeatureItem{}
	item.Unmarshal(i)
	return
}

func (i *FeatureItem) CMSItem() *cms.Item {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item
}

type GenericItem struct {
	ID          string               `json:"id,omitempty" cms:"id"`
	City        string               `json:"city,omitempty" cms:"city,reference"`
	Name        string               `json:"name,omitempty" cms:"name,text"`
	Type        string               `json:"type,omitempty" cms:"type,text"`
	TypeEn      string               `json:"type_en,omitempty" cms:"type_en,text"`
	Datasets    []GenericItemDataset `json:"datasets,omitempty" cms:"datasets,group"`
	OpenDataUrl string               `json:"open-data-url,omitempty" cms:"open_data_url,url"`
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
	DataURL    string `json:"url,omitempty" cms:"data_url,url"`
	DataFormat string `json:"data-format,omitempty" cms:"data_format,select"`
	LayerName  string `json:"layer-name,omitempty" cms:"layer_name,text"`
}

func GenericItemFrom(item *cms.Item) (i *GenericItem) {
	i = &GenericItem{}
	item.Unmarshal(i)
	return
}

func (i *GenericItem) CMSItem() *cms.Item {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item
}

type RelatedItem struct {
	ID              string              `json:"id,omitempty" cms:"id"`
	City            string              `json:"city,omitempty" cms:"city,reference"`
	Assets          map[string][]string `json:"assets,omitempty" cms:"-"`
	ConvertedAssets map[string][]string `json:"converted,omitempty" cms:"-"`
	Merged          string              `json:"merged,omitempty" cms:"merged,asset"`
	// metadata
	ConvertStatus ConvertionStatus `json:"conv_status,omitempty" cms:"conv_status,select,metadata"`
	MergeStatus   ConvertionStatus `json:"merge_status,omitempty" cms:"merge_status,select,metadata"`
	Public        bool             `json:"public,omitempty" cms:"public,bool,metadata"`
}

func RelatedItemFrom(item *cms.Item) (i *RelatedItem) {
	i = &RelatedItem{}
	item.Unmarshal(i)

	for _, t := range relatedDataTypes {
		v := item.FieldByKey(t).GetValue()
		cv := item.FieldByKey(t + "_conv").GetValue()

		var assets []string
		if s := v.String(); s != nil {
			assets = []string{*s}
		} else if s := v.Strings(); s != nil {
			assets = s
		}

		var conv []string
		if s := cv.String(); s != nil {
			conv = []string{*s}
		} else if s := cv.Strings(); s != nil {
			conv = s
		}

		if len(assets) > 0 {
			if i.Assets == nil {
				i.Assets = map[string][]string{}
			}
			i.Assets[t] = append(i.Assets[t], assets...)
		}

		if len(conv) > 0 {
			if i.ConvertedAssets == nil {
				i.ConvertedAssets = map[string][]string{}
			}
			i.ConvertedAssets[t] = append(i.ConvertedAssets[t], conv...)
		}
	}

	return
}

func (i *RelatedItem) CMSItem() *cms.Item {
	item := &cms.Item{}
	cms.Marshal(i, item)

	for _, t := range relatedDataTypes {
		if asset, ok := i.Assets[t]; ok {
			item.Fields = append(item.Fields, &cms.Field{
				Key:   t,
				Type:  "asset",
				Value: asset,
			})
		}

		if conv, ok := i.ConvertedAssets[t]; ok {
			item.Fields = append(item.Fields, &cms.Field{
				Key:   t + "_conv",
				Type:  "asset",
				Value: conv,
			})
		}

		// if pub, ok := i.Public[t]; ok {
		// 	item.MetadataFields = append(item.MetadataFields, &cms.Field{
		// 		Key:   t + "_public",
		// 		Type:  "bool",
		// 		Value: pub,
		// 	})
		// }
	}

	return item
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

func GeospatialjpIndexItemFrom(item *cms.Item) (i *GeospatialjpIndexItem) {
	i = &GeospatialjpIndexItem{}
	item.Unmarshal(i)
	return
}

func (i *GeospatialjpIndexItem) CMSItem() *cms.Item {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item
}

type GeospatialjpDataItem struct {
	ID          string `json:"id,omitempty" cms:"id"`
	City        string `json:"city,omitempty" cms:"city,reference"`
	CityGML     string `json:"citygml,omitempty" cms:"citygml,asset"`
	PlateauData string `json:"converted-data,omitempty" cms:"plateau_data,asset"`
	RelatedData string `json:"related-data,omitempty" cms:"related_data,asset"`
	GenericData string `json:"generic-data,omitempty" cms:"generic_data,asset"`
	// metadata
	CityGMLMergeStatus        ConvertionStatus `json:"citygml_merge_status,omitempty" cms:"citygml_merge_status,select,metadata"`
	ConvertedDataMergeSatatus ConvertionStatus `json:"plateau_merge_status,omitempty" cms:"plateau_merge_status,select,metadata"`
	RelatedDataMergeStatus    ConvertionStatus `json:"related_merge_status,omitempty" cms:"related_merge_status,select,metadata"`
}

func GeospatialjpDataItemFrom(item *cms.Item) (i *GeospatialjpDataItem) {
	i = &GeospatialjpDataItem{}
	item.Unmarshal(i)
	return
}

func (i *GeospatialjpDataItem) CMSItem() *cms.Item {
	item := &cms.Item{}
	cms.Marshal(i, item)
	return item
}
