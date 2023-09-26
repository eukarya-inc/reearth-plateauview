package datacatalogv2

import (
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogutil"
	"github.com/samber/lo"
)

type ItemCommon interface {
	GetCityName() string
	DataCatalogs() []DataCatalogItem
}

type DataCatalogItem struct {
	ID          string                                 `json:"id,omitempty"`
	ItemID      string                                 `json:"itemId,omitempty"`
	Name        string                                 `json:"name,omitempty"`
	Pref        string                                 `json:"pref,omitempty"`
	PrefCode    string                                 `json:"pref_code,omitempty"`
	City        string                                 `json:"city,omitempty"`
	CityEn      string                                 `json:"city_en,omitempty"`
	CityCode    string                                 `json:"city_code,omitempty"`
	Ward        string                                 `json:"ward,omitempty"`
	WardEn      string                                 `json:"ward_en,omitempty"`
	WardCode    string                                 `json:"ward_code,omitempty"`
	Type        string                                 `json:"type,omitempty"`
	Type2       string                                 `json:"type2,omitempty"`
	TypeEn      string                                 `json:"type_en,omitempty"`
	Type2En     string                                 `json:"type2_en,omitempty"`
	Format      string                                 `json:"format,omitempty"`
	Layers      []string                               `json:"layers,omitempty"`
	URL         string                                 `json:"url,omitempty"`
	Description string                                 `json:"desc,omitempty"`
	SearchIndex string                                 `json:"search_index,omitempty"`
	Year        int                                    `json:"year,omitempty"`
	OpenDataURL string                                 `json:"openDataUrl,omitempty"`
	Config      *datacatalogutil.DataCatalogItemConfig `json:"config,omitempty"`
	Order       *int                                   `json:"order,omitempty"`
	// force not creatign a type folder
	Root bool `json:"root,omitempty"`
	// force creating folder on root
	RootType bool   `json:"root_type,omitempty"`
	Group    string `json:"group,omitempty"`
	Infobox  bool   `json:"infobox,omitempty"`
	// alias of type that is used as a folder name
	Category string `json:"category,omitempty"`
	// internal
	Spec    string `json:"-"`
	Family  string `json:"-"`
	Edition string `json:"-"`
}

func (i DataCatalogItem) ConfigItems() []datacatalogutil.DataCatalogItemConfigItem {
	if i.Config == nil {
		return nil
	}
	return i.Config.Data
}

type DataCatalogGroup struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Prefecture string `json:"pref,omitempty"`
	City       string `json:"city,omitempty"`
	CityEn     string `json:"cityEn,omitempty"`
	Type       string `json:"type,omitempty"`
	Children   []any  `json:"children"`
}

type ResponseAll struct {
	Plateau []PlateauItem
	Dataset []DatasetItem
	Usecase []UsecaseItem
}

func (d ResponseAll) All() []DataCatalogItem {
	return append(append(d.PlateauItems(), d.DatasetItems()...), d.UsecaseItems()...)
}

func (d ResponseAll) PlateauItems() []DataCatalogItem {
	return items(d.Plateau, true)
}

func (d ResponseAll) DatasetItems() []DataCatalogItem {
	return items(d.Dataset, false)
}

func (d ResponseAll) UsecaseItems() []DataCatalogItem {
	return items(d.Usecase, false)
}

func items[T ItemCommon](data []T, omitOldItems bool) []DataCatalogItem {
	items := lo.FlatMap(data, func(i T, _ int) []DataCatalogItem {
		return i.DataCatalogs()
	})

	if !omitOldItems {
		return items
	}

	m := map[string]int{}
	for _, i := range items {
		m[i.CityCode] = i.Year
	}
	return lo.Filter(items, func(i DataCatalogItem, _ int) bool {
		y, ok := m[i.CityCode]
		return ok && y == i.Year
	})
}
