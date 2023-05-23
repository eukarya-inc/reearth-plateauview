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
		}).DataCatalogItem(
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
		}).DataCatalogItem(
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
		}).DataCatalogItem(
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
		),
	)

	// case4: urf
	assert.Equal(
		t,
		&DataCatalogItem{
			ID:          "01100_sapporo-shi_urf_UrbanPlanningArea",
			Name:        "都市計画区域モデル（札幌市）",
			Pref:        "北海道",
			PrefCode:    "01",
			City:        "札幌市",
			CityEn:      "sapporo-shi",
			CityCode:    "01100",
			Type:        "都市計画決定情報モデル",
			TypeEn:      "urf",
			Type2:       "都市計画区域",
			Type2En:     "UrbanPlanningArea",
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
		}).DataCatalogItem(
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
			Type2:       "都市計画区域",
			Type2En:     "UrbanPlanningArea",
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
		}).DataCatalogItem(
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
		),
	)
}

func TestMultipleLODData(t *testing.T) {
	// case1: single asset without LOD
	assert.Equal(t, DataCatalogItemConfig{
		Data: []DataCatalogItemConfigItem{
			{
				Name:   "xxxモデル",
				URL:    "https://example.com/01100_sapporo-shi_2020_mvt_op/{z}/{x}/{y}.mvt",
				Type:   "mvt",
				Layers: []string{"layer"},
			},
		},
	}, multipleLODData(
		[]string{
			"https://example.com/01100_sapporo-shi_2020_mvt_op.zip",
		},
		"xxxモデル",
		map[string][]string{"": {"layer"}},
	))

	// case2: multiple asset without LOD
	assert.Equal(t, DataCatalogItemConfig{
		Data: []DataCatalogItemConfigItem{
			{
				Name:   "xxxモデル1",
				URL:    "https://example.com/01100_sapporo-shi_2020_mvt_op/{z}/{x}/{y}.mvt",
				Type:   "mvt",
				Layers: []string{"layer"},
			},
			{
				Name:   "xxxモデル2",
				URL:    "https://example.com/01100_sapporo-shi_2020_3dtiles_op/tileset.json",
				Type:   "3dtiles",
				Layers: []string{"layer"},
			},
		},
	}, multipleLODData(
		[]string{
			"https://example.com/01100_sapporo-shi_2020_mvt_op.zip",
			"https://example.com/01100_sapporo-shi_2020_3dtiles_op.zip",
		},
		"xxxモデル",
		map[string][]string{"": {"layer"}},
	))

	// case3: multiple asset with LOD
	assert.Equal(t, DataCatalogItemConfig{
		Data: []DataCatalogItemConfigItem{
			{
				Name:   "LOD1",
				URL:    "https://example.com/01100_sapporo-shi_2020_mvt_op_tran_lod1/{z}/{x}/{y}.mvt",
				Type:   "mvt",
				Layers: []string{"layer1"},
			},
			{
				Name:   "LOD2",
				URL:    "https://example.com/01100_sapporo-shi_2020_3dtiles_op_tran_lod2/tileset.json",
				Type:   "3dtiles",
				Layers: []string{"layer"},
			},
		},
	}, multipleLODData(
		[]string{
			"https://example.com/01100_sapporo-shi_2020_mvt_op_tran_lod1.zip",
			"https://example.com/01100_sapporo-shi_2020_3dtiles_op_tran_lod2.zip",
		},
		"xxxモデル",
		map[string][]string{"": {"layer"}, "1": {"layer1"}},
	))
}
