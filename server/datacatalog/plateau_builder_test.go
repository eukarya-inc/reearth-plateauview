package datacatalog

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/cms"
	"github.com/stretchr/testify/assert"
)

func TestPlateauItem_IntermediateItem(t *testing.T) {
	assert.Equal(t, PlateauIntermediateItem{
		ID:          "id",
		Prefecture:  "北海道",
		City:        "札幌市",
		CityEn:      "sapporo-shi",
		CityCode:    "01100",
		Dic:         Dic{"dic": []DicEntry{{Name: "name"}}},
		OpenDataURL: "https://example.com",
		Year:        2020,
	}, PlateauItem{
		ID:         "id",
		Prefecture: "北海道",
		CityName:   "札幌市",
		CityGML: &cms.PublicAsset{
			URL: "https://example.com/01100_sapporo-shi_2020_citygml_op.zip",
		},
		Dic:         `{"dic":[{"name":"name"}]}`,
		OpenDataURL: "https://example.com",
	}.IntermediateItem())
}

func TestPlateauIntermediateItem_DataCatalogItem(t *testing.T) {
	// case1: normal
	assert.Equal(
		t,
		&DataCatalogItem{
			ID:          "01100_sapporo-shi_luse",
			Name:        "土地利用モデル（札幌市）",
			Pref:        "北海道",
			PrefCode:    "01",
			City:        "札幌市",
			CityEn:      "sapporo-shi",
			CityCode:    "01100",
			Type:        "土地利用モデル",
			TypeEn:      "luse",
			Description: "説明",
			URL:         "https://example.com/01100_sapporo-shi_2020_mvt_op_luse/{z}/{x}/{y}.mvt",
			OpenDataURL: "https://example.com",
			Format:      "mvt",
			Layers:      []string{"luse"},
			Year:        2020,
		},
		(&PlateauIntermediateItem{
			ID:          "itemid",
			Prefecture:  "北海道",
			City:        "札幌市",
			CityEn:      "sapporo-shi",
			CityCode:    "01100",
			OpenDataURL: "https://example.com",
			Year:        2020,
		}).dataCatalogItem(
			"土地利用モデル",
			AssetName{
				CityCode: "01100",
				CityEn:   "sapporo-shi",
				Format:   "mvt",
				Op:       "op",
				Feature:  "luse",
				Year:     "2020",
			},
			"https://example.com/01100_sapporo-shi_2020_mvt_op_luse.zip",
			"説明",
			[]string{"luse"},
			false,
			"",
			nil,
			nil,
		),
	)

	// case2: ward (not first), no entry in dic
	assert.Equal(
		t,
		&DataCatalogItem{
			ID:          "01100_sapporo-shi_011011_chuo-ku_bldg",
			Name:        "建築物モデル（chuo-ku）",
			Pref:        "北海道",
			PrefCode:    "01",
			City:        "札幌市",
			CityEn:      "sapporo-shi",
			CityCode:    "01100",
			Ward:        "chuo-ku",
			WardEn:      "chuo-ku",
			WardCode:    "011011",
			Type:        "建築物モデル",
			TypeEn:      "bldg",
			Description: "説明",
			URL:         "https://example.com/01100_sapporo-shi_2020_op_bldg_011011_chuo-ku_lod1/tileset.json",
			OpenDataURL: "https://www.geospatial.jp/ckan/dataset/plateau-01100-sapporo-shi-2020",
			Format:      "3dtiles",
			Year:        2020,
		},
		(&PlateauIntermediateItem{
			ID:         "itemid",
			Prefecture: "北海道",
			City:       "札幌市",
			CityEn:     "sapporo-shi",
			CityCode:   "01100",
			Year:       2020,
			Dic: Dic{
				"admin": []DicEntry{
					{Code: "011012", Name: "kita-ku", Description: "北区"},
				},
			},
		}).dataCatalogItem(
			"建築物モデル",
			AssetName{
				CityCode: "01100",
				CityEn:   "sapporo-shi",
				WardCode: "011011",
				WardEn:   "chuo-ku",
				Format:   "3dtiles",
				Op:       "op",
				Feature:  "bldg",
				Year:     "2020",
			},
			"https://example.com/01100_sapporo-shi_2020_op_bldg_011011_chuo-ku_lod1.zip",
			"説明",
			nil,
			false,
			"",
			nil,
			nil,
		),
	)

	// case3: ward (first)
	assert.Equal(
		t,
		&DataCatalogItem{
			ID:          "01100_sapporo-shi_011011_chuo-ku_bldg",
			ItemID:      "itemid",      // diff from case2
			Name:        "建築物モデル（中央区）", // diff from case2
			Pref:        "北海道",
			PrefCode:    "01",
			City:        "札幌市",
			CityEn:      "sapporo-shi",
			CityCode:    "01100",
			Ward:        "中央区", // diff from case2
			WardEn:      "chuo-ku",
			WardCode:    "011011",
			Type:        "建築物モデル",
			TypeEn:      "bldg",
			Description: "説明",
			URL:         "https://example.com/01100_sapporo-shi_2020_op_bldg_011011_chuo-ku_lod1/tileset.json",
			OpenDataURL: "https://example.com",
			Format:      "3dtiles",
			Year:        2020,
		},
		(&PlateauIntermediateItem{
			ID:          "itemid",
			Prefecture:  "北海道",
			City:        "札幌市",
			CityEn:      "sapporo-shi",
			CityCode:    "01100",
			Year:        2020,
			OpenDataURL: "https://example.com",
			Dic: Dic{
				"admin": []DicEntry{
					{Code: "011011", Name: "chuo-ku", Description: "中央区"}, // diff from case2
					{Code: "011012", Name: "kita-ku", Description: "北区"},
				},
			},
		}).dataCatalogItem(
			"建築物モデル",
			AssetName{
				CityCode: "01100",
				CityEn:   "sapporo-shi",
				WardCode: "011011",
				WardEn:   "chuo-ku",
				Format:   "3dtiles",
				Op:       "op",
				Feature:  "bldg",
				Year:     "2020",
			},
			"https://example.com/01100_sapporo-shi_2020_op_bldg_011011_chuo-ku_lod1.zip",
			"説明",
			nil,
			true,
			"",
			nil,
			nil,
		),
	)

	// case4: name override by
	assert.Equal(
		t,
		&DataCatalogItem{
			ID:          "01100_sapporo-shi_urf_UrbanPlanningArea",
			Name:        "1（札幌市）",
			Pref:        "北海道",
			PrefCode:    "01",
			City:        "札幌市",
			CityEn:      "sapporo-shi",
			CityCode:    "01100",
			Type:        "都市計画決定情報モデル",
			TypeEn:      "urf",
			Type2:       "2",
			Type2En:     "3",
			Description: "説明",
			URL:         "https://example.com/01100_sapporo-shi_2020_mvt_op_urf_UrbanPlanningArea/{z}/{x}/{y}.mvt",
			OpenDataURL: "https://www.geospatial.jp/ckan/dataset/plateau-01100-sapporo-shi-2020",
			Format:      "mvt",
			Layers:      []string{"UrbanPlanningArea"},
			Year:        2020,
		},
		(&PlateauIntermediateItem{
			ID:         "itemid",
			Prefecture: "北海道",
			City:       "札幌市",
			CityEn:     "sapporo-shi",
			CityCode:   "01100",
			Year:       2020,
		}).dataCatalogItem(
			"都市計画決定情報モデル",
			AssetName{
				CityCode:       "01100",
				CityEn:         "sapporo-shi",
				Format:         "mvt",
				Op:             "op",
				Feature:        "urf",
				Year:           "2020",
				UrfFeatureType: "UrbanPlanningArea",
			},
			"https://example.com/01100_sapporo-shi_2020_mvt_op_urf_UrbanPlanningArea.zip",
			"説明",
			[]string{"UrbanPlanningArea"},
			false,
			"",
			func(an AssetName) (string, string, string) {
				return "1", "2", "3"
			},
			nil,
		),
	)

	// case5: name override
	assert.Equal(
		t,
		&DataCatalogItem{
			ID:          "01100_sapporo-shi_urf_UrbanPlanningArea",
			Name:        "NAME（札幌市）",
			Pref:        "北海道",
			PrefCode:    "01",
			City:        "札幌市",
			CityEn:      "sapporo-shi",
			CityCode:    "01100",
			Type:        "都市計画決定情報モデル",
			TypeEn:      "urf",
			Description: "説明",
			URL:         "https://example.com/01100_sapporo-shi_2020_mvt_op_urf_UrbanPlanningArea/{z}/{x}/{y}.mvt",
			OpenDataURL: "https://www.geospatial.jp/ckan/dataset/plateau-01100-sapporo-shi-2020",
			Format:      "mvt",
			Layers:      []string{"UrbanPlanningArea"},
			Year:        2020,
		},
		(&PlateauIntermediateItem{
			ID:         "itemid",
			Prefecture: "北海道",
			City:       "札幌市",
			CityEn:     "sapporo-shi",
			CityCode:   "01100",
			Year:       2020,
		}).dataCatalogItem(
			"都市計画決定情報モデル",
			AssetName{
				CityCode:       "01100",
				CityEn:         "sapporo-shi",
				Format:         "mvt",
				Op:             "op",
				Feature:        "urf",
				Year:           "2020",
				UrfFeatureType: "UrbanPlanningArea",
			},
			"https://example.com/01100_sapporo-shi_2020_mvt_op_urf_UrbanPlanningArea.zip",
			"説明",
			[]string{"UrbanPlanningArea"},
			false,
			"NAME",
			nil,
			nil,
		),
	)

	// case6: sub name
	assert.Equal(
		t,
		&DataCatalogItem{
			ID:          "01100_sapporo-shi_urf_UrbanPlanningArea",
			Name:        "都市計画決定情報モデル SUB（札幌市）",
			Pref:        "北海道",
			PrefCode:    "01",
			City:        "札幌市",
			CityEn:      "sapporo-shi",
			CityCode:    "01100",
			Type:        "都市計画決定情報モデル",
			TypeEn:      "urf",
			Description: "説明",
			URL:         "https://example.com/01100_sapporo-shi_2020_mvt_op_urf_UrbanPlanningArea/{z}/{x}/{y}.mvt",
			OpenDataURL: "https://www.geospatial.jp/ckan/dataset/plateau-01100-sapporo-shi-2020",
			Format:      "mvt",
			Layers:      []string{"UrbanPlanningArea"},
			Year:        2020,
		},
		(&PlateauIntermediateItem{
			ID:         "itemid",
			Prefecture: "北海道",
			City:       "札幌市",
			CityEn:     "sapporo-shi",
			CityCode:   "01100",
			Year:       2020,
		}).dataCatalogItem(
			"都市計画決定情報モデル",
			AssetName{
				CityCode:       "01100",
				CityEn:         "sapporo-shi",
				Format:         "mvt",
				Op:             "op",
				Feature:        "urf",
				Year:           "2020",
				UrfFeatureType: "UrbanPlanningArea",
			},
			"https://example.com/01100_sapporo-shi_2020_mvt_op_urf_UrbanPlanningArea.zip",
			"説明",
			[]string{"UrbanPlanningArea"},
			false,
			"",
			nil,
			func(an AssetName, dic Dic) string {
				return "SUB"
			},
		),
	)
}
