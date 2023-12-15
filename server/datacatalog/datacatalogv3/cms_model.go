package datacatalogv3

import (
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationcommon"
	cms "github.com/reearth/reearth-cms-api/go"
)

const modelPrefix = "plateau-"
const cityModel = "city"
const relatedModel = "related"
const genericModel = "generic"

type ManagementStatus string

const (
	ManagementStatusReady ManagementStatus = "確認可能"
)

type AllData struct {
	FeatureTypes FeatureTypes
	City         []*CityItem
	Related      []*RelatedItem
	Generic      []*GenericItem
	Plateau      map[string][]*FeatureItem
}

type FeatureTypes struct {
	Plateau []FeatureType
	Related []FeatureType
	Generic []FeatureType
}

type FeatureType struct {
	Code string `json:"code,omitempty" cms:"code,text"`
	Name string `json:"name,omitempty" cms:"name,text"`
}

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

func CityItemFrom(item *cms.Item, featureTypes []FeatureType) (i *CityItem) {
	i = &CityItem{}
	item.Unmarshal(i)

	references := map[string]string{}
	public := map[string]bool{}
	for _, ft := range featureTypes {
		if ref := item.FieldByKey(ft.Code).GetValue().String(); ref != nil {
			references[ft.Code] = *ref
		}

		if pub := item.MetadataFieldByKey(ft.Code + "_public").GetValue().Bool(); pub != nil {
			public[ft.Code] = *pub
		}
	}

	i.References = references
	i.Public = public
	return
}

func (i *CityItem) CMSItem() *cms.Item {
	item := &cms.Item{}
	cms.Marshal(i, item)

	for ft, ref := range i.References {
		item.Fields = append(item.Fields, &cms.Field{
			Key:   ft,
			Type:  "reference",
			Value: ref,
		})
	}

	for ft, pub := range i.Public {
		item.MetadataFields = append(item.MetadataFields, &cms.Field{
			Key:   ft + "_public",
			Type:  "bool",
			Value: pub,
		})
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
	Status ManagementStatus `json:"status,omitempty" cms:"status,select,metadata"`
}

func (c FeatureItem) IsPublicForAdmin() bool {
	return c.Status == ManagementStatusReady
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

func (c GenericItem) IsPublicForAdmin() bool {
	return c.Status == ManagementStatusReady
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
	Public bool `json:"public,omitempty" cms:"public,bool,metadata"`
}

func RelatedItemFrom(item *cms.Item, featureTypes []FeatureType) (i *RelatedItem) {
	i = &RelatedItem{}
	item.Unmarshal(i)

	for _, t := range featureTypes {
		v := item.FieldByKey(t.Code).GetValue()
		cv := item.FieldByKey(t.Code + "_conv").GetValue()

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
			i.Assets[t.Code] = append(i.Assets[t.Code], assets...)
		}

		if len(conv) > 0 {
			if i.ConvertedAssets == nil {
				i.ConvertedAssets = map[string][]string{}
			}
			i.ConvertedAssets[t.Code] = append(i.ConvertedAssets[t.Code], conv...)
		}
	}

	return
}

func (i *RelatedItem) CMSItem(relatedDataTypes []FeatureType) *cms.Item {
	item := &cms.Item{}
	cms.Marshal(i, item)

	for _, t := range relatedDataTypes {
		if asset, ok := i.Assets[t.Code]; ok {
			item.Fields = append(item.Fields, &cms.Field{
				Key:   t.Code,
				Type:  "asset",
				Value: asset,
			})
		}

		if conv, ok := i.ConvertedAssets[t.Code]; ok {
			item.Fields = append(item.Fields, &cms.Field{
				Key:   t.Code + "_conv",
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
