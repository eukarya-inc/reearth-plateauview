package plateauv2

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/eukarya-inc/jpareacode"
	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2/datacatalogutil"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearthx/util"
	"github.com/samber/lo"
)

var datasetTypes = map[string]string{
	"公園":     "park",
	"避難施設":   "shelter",
	"鉄道":     "railway",
	"鉄道駅":    "station",
	"行政界":    "border",
	"ランドマーク": "landmark",
	"緊急輸送道路": "emergency_route",
}

type DatasetItem struct {
	ID           string             `json:"id,omitempty"`
	Name         string             `json:"name,omitempty"`
	Type         string             `json:"type,omitempty"`
	Prefecture   string             `json:"prefecture,omitempty"`
	CityName     string             `json:"city_name,omitempty"`
	WardName     string             `json:"ward_name,omitempty"`
	OpenDataURL  string             `json:"opendata_url,omitempty"`
	Description  string             `json:"description,omitempty"`
	Year         string             `json:"year,omitempty"`
	Data         *cms.PublicAsset   `json:"data,omitempty"`
	DataFormat   string             `json:"data_format,omitempty"`
	DataURL      string             `json:"data_url,omitempty"`
	DataLayers   string             `json:"data_layer,omitempty"`
	Config       string             `json:"config,omitempty"`
	Order        *int               `json:"order,omitempty"`
	OriginalData []*cms.PublicAsset `json:"data_orig,omitempty"`
}

func (i DatasetItem) GetCityName() string {
	return i.CityName
}

func (i DatasetItem) DataCatalogs() []DataCatalogItem {
	pref, prefCodeInt := datacatalogutil.NormalizePref(i.Prefecture)
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

	var c *datacatalogutil.DataCatalogItemConfig
	_ = json.Unmarshal([]byte(i.Config), &c)

	u := ""
	if i.Data != nil && i.Data.URL != "" {
		u = i.Data.URL
	}
	if u == "" {
		u = i.DataURL
	}

	f := datacatalogutil.FormatTypeEn(i.DataFormat)

	var ou, of string
	if len(i.OriginalData) > 0 {
		ou = i.OriginalData[0].URL
		of = datacatalogutil.FormatTypeEn(strings.TrimPrefix(path.Ext(ou), "."))
	}

	t := i.Type
	if t != "" && !strings.HasSuffix(t, "情報") {
		t += "情報"
	}

	var layers []string
	if i.DataLayers != "" {
		layers = lo.Filter(util.Map(strings.Split(i.DataLayers, ","), strings.TrimSpace), func(s string, _ int) bool { return s != "" })
	}

	name := i.Name
	if t != "" {
		cityOrWard := city
		if ward != "" {
			cityOrWard = ward
		}
		name = fmt.Sprintf("%s（%s）", t, cityOrWard)
	}

	return []DataCatalogItem{{
		ID:             i.ID,
		Name:           name,
		Type:           t,
		TypeEn:         datasetTypes[i.Type],
		Pref:           pref,
		PrefCode:       prefCode,
		City:           city,
		CityCode:       cCode,
		Ward:           ward,
		WardCode:       wCode,
		Format:         f,
		URL:            datacatalogutil.AssetURLFromFormat(u, f),
		OriginalFormat: of,
		OriginalURL:    datacatalogutil.AssetURLFromFormat(ou, of),
		Description:    i.Description,
		Config:         c,
		Layers:         layers,
		Year:           datacatalogutil.YearInt(i.Year),
		OpenDataURL:    i.OpenDataURL,
		Order:          i.Order,
		Family:         "related",
		Edition:        "2022",
	}}
}
