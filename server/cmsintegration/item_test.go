package cmsintegration

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/stretchr/testify/assert"
)

var item = Item{
	ID:                  "xxx",
	Prefecture:          "東京都",
	CityName:            "千代田区",
	Specification:       "第2.0版",
	CityGML:             "citygml_assetid",
	CityGMLGeoSpatialJP: "citygml_geospatialjp_assetid",
	Catalog:             "catalog_assetid",
	ConversionEnabled:   "変換する",
	PRCS:                "第1系",
	QualityCheckParams:  "qcp_assetid",
	DevideODC:           "分割する",
	Bldg:                []string{"bldg_assetid", "bldg_assetid2"},
	Tran:                "tran_assetid",
	Frn:                 "frn_assetid",
	Veg:                 "veg_assetid",
	Luse:                "luse_assetid",
	Lsld:                "lsld_assetid",
	Urf:                 "urf_assetid",
	Fld:                 []string{"fld_assetid", "fld_assetid2"},
	Tnum:                []string{"tnum_assetid", "tnum_assetid2"},
	Htd:                 []string{"htd_assetid", "htd_assetid2"},
	Ifld:                []string{"ifld_assetid", "ifld_assetid2"},
	All:                 "all_assetid",
	Dictionary:          "dictionary_assetid",
	ConversionStatus:    "実行中",
	CatalogStatus:       "完了",
	MaxLOD:              "maxlod_assetid",
	MaxLODStatus:        "未実行",
	SearchIndex:         "searchindex_assetid",
	SeatchIndexStatus:   "エラー",
}

var cmsitem = cms.Item{
	ID: "xxx",
	Fields: []cms.Field{
		{Key: "prefecture", Type: "select", Value: "東京都"},
		{Key: "city_name", Type: "text", Value: "千代田区"},
		{Key: "specification", Type: "text", Value: "第2.0版"},
		{Key: "citygml", Type: "asset", Value: "citygml_assetid"},
		{Key: "citygml_geospatialjp", Type: "asset", Value: "citygml_geospatialjp_assetid"},
		{Key: "catalog", Type: "asset", Value: "catalog_assetid"},
		{Key: "conversion_enabled", Type: "select", Value: "変換する"},
		{Key: "prcs", Type: "select", Value: "第1系"},
		{Key: "quality_check_params", Type: "asset", Value: "qcp_assetid"},
		{Key: "devide_odc", Type: "select", Value: "分割する"},
		{Key: "bldg", Type: "asset", Value: []string{"bldg_assetid", "bldg_assetid2"}},
		{Key: "tran", Type: "asset", Value: "tran_assetid"},
		{Key: "frn", Type: "asset", Value: "frn_assetid"},
		{Key: "veg", Type: "asset", Value: "veg_assetid"},
		{Key: "luse", Type: "asset", Value: "luse_assetid"},
		{Key: "lsld", Type: "asset", Value: "lsld_assetid"},
		{Key: "urf", Type: "asset", Value: "urf_assetid"},
		{Key: "fld", Type: "asset", Value: []string{"fld_assetid", "fld_assetid2"}},
		{Key: "tnum", Type: "asset", Value: []string{"tnum_assetid", "tnum_assetid2"}},
		{Key: "htd", Type: "asset", Value: []string{"htd_assetid", "htd_assetid2"}},
		{Key: "ifld", Type: "asset", Value: []string{"ifld_assetid", "ifld_assetid2"}},
		{Key: "all", Type: "asset", Value: "all_assetid"},
		{Key: "dictionary", Type: "asset", Value: "dictionary_assetid"},
		{Key: "conversion_status", Type: "select", Value: "実行中"},
		{Key: "catalog_status", Type: "select", Value: "完了"},
		{Key: "max_lod", Type: "asset", Value: "maxlod_assetid"},
		{Key: "max_lod_status", Type: "select", Value: "未実行"},
		{Key: "search_index", Type: "asset", Value: "searchindex_assetid"},
		{Key: "search_index_status", Type: "select", Value: "エラー"},
	},
}

func TestItem(t *testing.T) {
	assert.Equal(t, item, ItemFrom(cmsitem))
	assert.Equal(t, Item{}, ItemFrom(cms.Item{}))
	assert.Equal(t, cmsitem.Fields, item.Fields())
	assert.Equal(t, []cms.Field(nil), Item{}.Fields())
}

func TestPRCS_Code(t *testing.T) {
	assert.Equal(t, "6669", PRCS("第1系").ESPGCode())
	assert.Equal(t, "6670", PRCS("第2系").ESPGCode())
	assert.Equal(t, "6686", PRCS("第18系").ESPGCode())
	assert.Equal(t, "6687", PRCS("第19系").ESPGCode())
}
