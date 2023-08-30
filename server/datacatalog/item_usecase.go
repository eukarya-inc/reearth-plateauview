package datacatalog

import (
	"encoding/json"
	"strings"

	"github.com/eukarya-inc/jpareacode"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogutil"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

const folder = "フォルダ"
const folderEn = "folder"

type UsecaseItem struct {
	ID          string           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Type        string           `json:"type,omitempty"`
	Prefecture  string           `json:"prefecture,omitempty"`
	CityName    string           `json:"city_name,omitempty"`
	WardName    string           `json:"ward_name,omitempty"`
	OpenDataURL string           `json:"opendata_url,omitempty"`
	Description string           `json:"description,omitempty"`
	Year        string           `json:"year,omitempty"`
	Data        *cms.PublicAsset `json:"data,omitempty"`
	DataFormat  string           `json:"data_format,omitempty"`
	DataURL     string           `json:"data_url,omitempty"`
	DataLayers  string           `json:"data_layer,omitempty"`
	Config      string           `json:"config,omitempty"`
	Order       *int             `json:"order,omitempty"`
	Category    string           `json:"category,omitempty"`
}

func (i UsecaseItem) GetCityName() string {
	return i.CityName
}

func (i UsecaseItem) DataCatalogs() []DataCatalogItem {
	pref, prefCodeInt := normalizePref(i.Prefecture)
	prefCode := jpareacode.FormatPrefectureCode(prefCodeInt)

	var city, ward string
	if i.WardName != "" {
		city = i.CityName
		ward = i.WardName
	} else {
		city, ward, _ = strings.Cut(i.CityName, "/")
	}

	cCode := datacatalogutil.CityCode("", city, prefCodeInt)
	wCode := datacatalogutil.CityCode("", ward, prefCodeInt)

	if i.DataFormat == folder {
		return []DataCatalogItem{{
			ID:          i.ID,
			Name:        i.Name,
			Type:        folder,
			TypeEn:      folderEn,
			Pref:        pref,
			PrefCode:    prefCode,
			City:        city,
			CityCode:    cCode,
			Ward:        ward,
			WardCode:    wCode,
			Description: i.Description,
		}}
	}

	var c any
	_ = json.Unmarshal([]byte(i.Config), &c)

	u := ""
	if i.Data != nil && i.Data.URL != "" {
		u = i.Data.URL
	}
	if u == "" {
		u = i.DataURL
	}

	f := formatTypeEn(i.DataFormat)

	var layers []string
	if i.DataLayers != "" {
		layers = lo.Filter(util.Map(strings.Split(i.DataLayers, ","), strings.TrimSpace), func(s string, _ int) bool { return s != "" })
	}

	ty, tye := i.Category, i.Category
	if ty == "" || ty == "ユースケース" {
		ty = "ユースケース"
		tye = "usecase"
	}

	return []DataCatalogItem{{
		ID:          i.ID,
		Name:        i.Name,
		Type:        ty,
		TypeEn:      tye,
		Pref:        pref,
		PrefCode:    prefCode,
		City:        city,
		CityCode:    cCode,
		Ward:        ward,
		WardCode:    wCode,
		Format:      f,
		URL:         datacatalogutil.AssetURLFromFormat(u, f),
		Description: i.Description,
		Config:      c,
		Layers:      layers,
		Year:        yearInt(i.Year),
		OpenDataURL: i.OpenDataURL,
		Order:       i.Order,
		RootType:    pref != zenkyu,
	}}
}
