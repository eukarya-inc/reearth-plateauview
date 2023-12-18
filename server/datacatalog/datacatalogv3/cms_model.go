package datacatalogv3

import (
	"encoding/json"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogcommon"
	cms "github.com/reearth/reearth-cms-api/go"
)

const modelPrefix = "plateau-"
const cityModel = "city"
const relatedModel = "related"
const genericModel = "generic"

type ManagementStatus string

type stage string

const (
	stageAlpha            stage            = "alpha"
	stageBeta             stage            = "beta"
	stageGA               stage            = "ga"
	ManagementStatusReady ManagementStatus = "確認可能"
)

type FeatureType struct {
	Code string `json:"code,omitempty" cms:"code,text"`
	Name string `json:"name,omitempty" cms:"name,text"`
	// for plateau
	SpecMajor int  `json:"spec_major,omitempty" cms:"spec_major,integer"`
	Flood     bool `json:"flood,omitempty" cms:"flood,bool"`
}

type CityItem struct {
	ID             string            `json:"id,omitempty" cms:"id"`
	Prefecture     string            `json:"prefecture,omitempty" cms:"prefecture,select"`
	CityName       string            `json:"city_name,omitempty" cms:"city_name,text"`
	CityNameEn     string            `json:"city_name_en,omitempty" cms:"city_name_en,text"`
	CityCode       string            `json:"city_code,omitempty" cms:"city_code,text"`
	Spec           string            `json:"spec,omitempty" cms:"spec,select"`
	References     map[string]string `json:"references,omitempty" cms:"-"`
	RelatedDataset string            `json:"related_dataset,omitempty" cms:"related_dataset,reference"`
	Year           string            `json:"year,omitempty" cms:"year,select"`
	// meatadata
	PlateauDataStatus ManagementStatus `json:"plateau_data_status,omitempty" cms:"plateau_data_status,select,metadata"`
	SDKPublic         bool             `json:"sdk_public,omitempty" cms:"sdk_public,bool,metadata"`
	RelatedPublic     bool             `json:"related_public,omitempty" cms:"related_public,bool,metadata"`
	Public            map[string]bool  `json:"public,omitempty" cms:"-"`
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

func (i *CityItem) YearInt() int {
	return datacatalogcommon.YearInt(i.Year)
}

type PlateauFeatureItem struct {
	ID      string                    `json:"id,omitempty" cms:"id"`
	City    string                    `json:"city,omitempty" cms:"city,reference"`
	CityGML string                    `json:"citygml,omitempty" cms:"citygml,-"`
	Data    []string                  `json:"data,omitempty" cms:"data,-"`
	Desc    string                    `json:"desc,omitempty" cms:"desc,textarea"`
	Items   []PlateauFeatureItemDatum `json:"items,omitempty" cms:"items,group"`
	Dic     string                    `json:"dic,omitempty" cms:"dic,textarea"`
	MaxLOD  string                    `json:"maxlod,omitempty" cms:"maxlod,-"`
	// metadata
	Status ManagementStatus `json:"status,omitempty" cms:"status,select,metadata"`
}

func (c PlateauFeatureItem) IsPublicForAdmin() bool {
	return c.Status == ManagementStatusReady
}

func (c PlateauFeatureItem) ReadDic() (d Dic, _ error) {
	err := json.Unmarshal([]byte(c.Dic), &d)
	return d, err
}

type PlateauFeatureItemDatum struct {
	ID   string   `json:"id,omitempty" cms:"id"`
	Data []string `json:"data,omitempty" cms:"data,-"`
	Name string   `json:"name,omitempty" cms:"name,text"`
	Desc string   `json:"desc,omitempty" cms:"desc,textarea"`
	Key  string   `json:"key,omitempty" cms:"key,text"`
}

type Dic map[string][]DicEntry // admin, fld. htd, tnm, gen

func (d Dic) FindEntry(key, name string) *DicEntry {
	if d == nil {
		return nil
	}

	if entries, ok := d[key]; ok {
		for _, e := range entries {
			if e.Name == name {
				return &e
			}
		}
	}

	return nil
}

type DicEntry struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Code        string `json:"code,omitempty"`  // bldg only
	Admin       string `json:"admin,omitempty"` // fld only
	Scale       string `json:"scale,omitempty"` // fld only
}

func PlateauFeatureItemFrom(item *cms.Item) (i *PlateauFeatureItem) {
	i = &PlateauFeatureItem{}
	item.Unmarshal(i)

	i.CityGML = valueToAssetURL(item.FieldByKey("citygml").GetValue())
	i.Data = valueToAssetURLs(item.FieldByKey("data").GetValue())
	i.MaxLOD = valueToAssetURL(item.FieldByKey("maxlod").GetValue())
	for ind, d := range i.Items {
		i.Items[ind].Data = valueToAssetURLs(item.FieldByKeyAndGroup("data", d.ID).GetValue())
	}

	return
}

type GenericItem struct {
	ID          string               `json:"id,omitempty" cms:"id"`
	City        string               `json:"city,omitempty" cms:"city,reference"`
	Name        string               `json:"name,omitempty" cms:"name,text"`
	Desc        string               `json:"desc,omitempty" cms:"desc,textarea"`
	Type        string               `json:"type,omitempty" cms:"type,text"`
	TypeEn      string               `json:"type_en,omitempty" cms:"type_en,text"`
	Data        []GenericItemDataset `json:"data,omitempty" cms:"data,group"`
	OpenDataUrl string               `json:"open-data-url,omitempty" cms:"open_data_url,url"`
	Category    string               `json:"category,omitempty" cms:"category,select"`
	// metadata
	Status ManagementStatus `json:"status,omitempty" cms:"status,select,metadata"`
	Public bool             `json:"public,omitempty" cms:"public,bool,metadata"`
	UseAR  bool             `json:"use-ar,omitempty" cms:"use-ar,bool,metadata"`
}

func (c *GenericItem) Stage() stage {
	if c.Public {
		return stageGA
	}
	if c.Status == ManagementStatusReady {
		return stageBeta
	}
	return stageAlpha
}

type GenericItemDataset struct {
	ID         string `json:"id,omitempty" cms:"id"`
	Name       string `json:"name,omitempty" cms:"name,text"`
	Data       string `json:"data,omitempty" cms:"data,-"`
	Desc       string `json:"desc,omitempty" cms:"desc,textarea"`
	DataURL    string `json:"url,omitempty" cms:"data_url,url"`
	DataFormat string `json:"data-format,omitempty" cms:"data_format,select"`
	LayerName  string `json:"layer-name,omitempty" cms:"layer_name,text"`
}

func GenericItemFrom(item *cms.Item) (i *GenericItem) {
	i = &GenericItem{}
	item.Unmarshal(i)

	for ind, d := range i.Data {
		i.Data[ind].Data = valueToAssetURL(item.FieldByKeyAndGroup("data", d.ID).GetValue())
	}

	return
}

type RelatedItem struct {
	ID              string              `json:"id,omitempty" cms:"id"`
	City            string              `json:"city,omitempty" cms:"city,reference"`
	Assets          map[string][]string `json:"assets,omitempty" cms:"-"`
	ConvertedAssets map[string][]string `json:"converted,omitempty" cms:"-"`
	Desc            string              `json:"desc,omitempty" cms:"-"`
	// metadata
	Status ManagementStatus `json:"status,omitempty" cms:"status,select,metadata"`
	Public bool             `json:"public,omitempty" cms:"public,bool,metadata"`
}

func (c *RelatedItem) Stage() stage {
	if c.Public {
		return stageGA
	}
	if c.Status == ManagementStatusReady {
		return stageBeta
	}
	return stageAlpha
}

func RelatedItemFrom(item *cms.Item, featureTypes []FeatureType) (i *RelatedItem) {
	i = &RelatedItem{}
	item.Unmarshal(i)

	for _, t := range featureTypes {
		assets := valueToAssetURLs(item.FieldByKey(t.Code).GetValue())
		conv := valueToAssetURLs(item.FieldByKey(t.Code + "_conv").GetValue())

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

func valueToAssetURL(v *cms.Value) string {
	return anyToAssetURL(v.Interface())
}

func valueToAssetURLs(v *cms.Value) (res []string) {
	i := v.Interface()
	if i == nil {
		return
	}

	values := []any{}
	if s, ok := i.([]any); ok {
		values = s
	} else {
		values = append(values, i)
	}

	for _, v := range values {
		if url := anyToAssetURL(v); url != "" {
			res = append(res, url)
		}
	}

	return
}

func anyToAssetURL(v any) string {
	if v == nil {
		return ""
	}

	m, ok := v.(map[string]any)
	if !ok {
		m2, ok := v.(map[any]any)
		if !ok {
			return ""
		}

		m = map[string]interface{}{}
		for k, v := range m2 {
			if s, ok := k.(string); ok {
				m[s] = v
			}
		}
	}

	url, ok := m["url"].(string)
	if !ok {
		return ""
	}

	return url
}