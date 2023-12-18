package datacatalogv3

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestPlateauDataset_ToWards(t *testing.T) {
	dic := `{
		"admin": [
			{
				"code": "11112",
				"name": "bar-shi_bar-ku",
				"description": "bar市 hoge区"
			},
			{
				"code": "11113",
				"name": "foo-shi_foo-ku",
				"description": "foo区"
			}
		]
	}`

	item := &PlateauFeatureItem{
		ID:  "id",
		Dic: dic,
	}

	expected := []*plateauapi.Ward{
		{
			ID:             plateauapi.NewID("11112", plateauapi.TypeArea),
			Name:           "hoge区",
			Type:           plateauapi.AreaTypeWard,
			Code:           plateauapi.AreaCode("11112"),
			PrefectureID:   plateauapi.NewID("11", plateauapi.TypeArea),
			PrefectureCode: plateauapi.AreaCode("11"),
			CityID:         plateauapi.NewID("11111", plateauapi.TypeArea),
			CityCode:       plateauapi.AreaCode("11111"),
		},
		{
			ID:             plateauapi.NewID("11113", plateauapi.TypeArea),
			Name:           "foo区",
			Type:           plateauapi.AreaTypeWard,
			Code:           plateauapi.AreaCode("11113"),
			PrefectureID:   plateauapi.NewID("11", plateauapi.TypeArea),
			PrefectureCode: plateauapi.AreaCode("11"),
			CityID:         plateauapi.NewID("11111", plateauapi.TypeArea),
			CityCode:       plateauapi.AreaCode("11111"),
		},
	}

	pref := &plateauapi.Prefecture{
		ID:   plateauapi.NewID("11", plateauapi.TypeArea),
		Code: plateauapi.AreaCode("11"),
	}

	city := &plateauapi.City{
		ID:   plateauapi.NewID("11111", plateauapi.TypeArea),
		Code: plateauapi.AreaCode("11111"),
	}

	res := item.toWards(pref, city)
	assert.Equal(t, expected, res)
}

func TestPlateauDataset_ToDatasets_Bldg(t *testing.T) {
	item := &PlateauFeatureItem{
		ID:   "id",
		Desc: "desc",
		Data: []string{
			"https://example.com/11111_bar-shi_city_2023_citygml_1_op_bldg_3dtiles_11112_hoge-ku_lod1.zip",
			"https://example.com/11111_bar-shi_city_2023_citygml_1_op_bldg_3dtiles_11112_hoge-ku_lod1_no_texture.zip",
			"https://example.com/11111_bar-shi_city_2023_citygml_1_op_bldg_3dtiles_11112_hoge-ku_lod2.zip",
			"https://example.com/11111_bar-shi_city_2023_citygml_1_op_bldg_3dtiles_11113_foo-ku_lod1.zip",
		},
	}

	expected := []plateauapi.Dataset{
		&plateauapi.PlateauDataset{
			ID:                 plateauapi.NewID("11112_bldg", plateauapi.TypeDataset),
			Name:               "建築物モデル（hoge区）",
			Description:        lo.ToPtr("desc"),
			Year:               2023,
			PrefectureID:       lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode:     lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:             lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:           lo.ToPtr(plateauapi.AreaCode("11111")),
			WardID:             lo.ToPtr(plateauapi.NewID("11112", plateauapi.TypeArea)),
			WardCode:           lo.ToPtr(plateauapi.AreaCode("11112")),
			TypeID:             plateauapi.NewID("bldg", plateauapi.TypeDatasetType),
			TypeCode:           "bldg",
			PlateauSpecMinorID: plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
			Stage:              lo.ToPtr(string(stageAlpha)),
			Items: []*plateauapi.PlateauDatasetItem{
				{
					ID:       plateauapi.NewID("11112_bldg_lod1", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "LOD1",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_bldg_3dtiles_11112_hoge-ku_lod1.zip",
					ParentID: plateauapi.NewID("11112_bldg", plateauapi.TypeDataset),
					Lod:      lo.ToPtr(1),
					Texture:  lo.ToPtr(plateauapi.TextureTexture),
				},
				{
					ID:       plateauapi.NewID("11112_bldg_lod1_no_texture", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "LOD1（テクスチャなし）",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_bldg_3dtiles_11112_hoge-ku_lod1_no_texture.zip",
					ParentID: plateauapi.NewID("11112_bldg", plateauapi.TypeDataset),
					Lod:      lo.ToPtr(1),
					Texture:  lo.ToPtr(plateauapi.TextureNone),
				},
				{
					ID:       plateauapi.NewID("11112_bldg_lod2", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "LOD2",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_bldg_3dtiles_11112_hoge-ku_lod2.zip",
					ParentID: plateauapi.NewID("11112_bldg", plateauapi.TypeDataset),
					Lod:      lo.ToPtr(2),
					Texture:  lo.ToPtr(plateauapi.TextureTexture),
				},
			},
		},
		&plateauapi.PlateauDataset{
			ID:                 plateauapi.NewID("11113_bldg", plateauapi.TypeDataset),
			Name:               "建築物モデル（foo区）",
			Description:        lo.ToPtr("desc"),
			Year:               2023,
			PrefectureID:       lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode:     lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:             lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:           lo.ToPtr(plateauapi.AreaCode("11111")),
			WardID:             lo.ToPtr(plateauapi.NewID("11113", plateauapi.TypeArea)),
			WardCode:           lo.ToPtr(plateauapi.AreaCode("11113")),
			TypeID:             plateauapi.NewID("bldg", plateauapi.TypeDatasetType),
			TypeCode:           "bldg",
			PlateauSpecMinorID: plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
			Stage:              lo.ToPtr(string(stageAlpha)),
			Items: []*plateauapi.PlateauDatasetItem{
				{
					ID:       plateauapi.NewID("11113_bldg_lod1", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "LOD1",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_bldg_3dtiles_11113_foo-ku_lod1.zip",
					ParentID: plateauapi.NewID("11113_bldg", plateauapi.TypeDataset),
					Lod:      lo.ToPtr(1),
					Texture:  lo.ToPtr(plateauapi.TextureTexture),
				},
			},
		},
	}

	area := &areaContext{
		Pref: &plateauapi.Prefecture{},
		City: &plateauapi.City{
			Name: "bar市",
			Code: "11111",
		},
		PrefID:   lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
		CityID:   lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
		PrefCode: lo.ToPtr(plateauapi.AreaCode("11")),
		CityCode: lo.ToPtr(plateauapi.AreaCode("11111")),
		CityItem: &CityItem{
			Year: "2023年",
		},
		Wards: []*plateauapi.Ward{
			{
				ID:   plateauapi.NewID("11112", plateauapi.TypeArea),
				Code: plateauapi.AreaCode("11112"),
				Name: "hoge区",
			},
			{
				ID:   plateauapi.NewID("11113", plateauapi.TypeArea),
				Code: plateauapi.AreaCode("11113"),
				Name: "foo区",
			},
		},
	}

	dts := &plateauapi.PlateauDatasetType{
		ID:   plateauapi.NewID("bldg", plateauapi.TypeDatasetType),
		Code: "bldg",
		Name: "建築物モデル",
	}

	spec := &plateauapi.PlateauSpecMinor{
		ID:           plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
		Name:         "第3.2版",
		MajorVersion: 3,
		Version:      "3.2",
		Year:         2023,
		ParentID:     plateauapi.NewID("3", plateauapi.TypePlateauSpec),
	}

	res, warning := item.toDatasets(area, dts, spec)
	assert.Nil(t, warning)
	assert.Equal(t, expected, res)
}

func TestPlateauDataset_ToDatasets_Tnm(t *testing.T) {
	item := &PlateauFeatureItem{
		ID:   "id",
		Desc: "desc",
		Items: []PlateauFeatureItemDatum{
			{
				ID: "id1",
				Data: []string{
					"https://example.com/11111_bar-shi_city_2023_citygml_1_op_tnm_AAA_3dtiles.zip",
					"https://example.com/11111_bar-shi_city_2023_citygml_1_op_tnm_AAA_3dtiles_no_texture.zip",
				},
				Desc: "desc1",
			},
			{
				ID: "id2",
				Data: []string{
					"https://example.com/11111_bar-shi_city_2023_citygml_1_op_tnm_BBB_3dtiles.zip",
				},
				Desc: "desc2",
			},
		},
		Dic: `{
			"tnm": [
				{ "name": "AAA", "description": "AAA!" },
				{ "name": "BBB", "description": "BBB!" }
			]
		}`,
	}

	expected := []plateauapi.Dataset{
		&plateauapi.PlateauDataset{
			ID:                 plateauapi.NewID("11111_tnm_AAA", plateauapi.TypeDataset),
			Name:               "津波浸水想定区域モデル AAA!（bar市）",
			Description:        lo.ToPtr("desc1"),
			Year:               2023,
			PrefectureID:       lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode:     lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:             lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:           lo.ToPtr(plateauapi.AreaCode("11111")),
			TypeID:             plateauapi.NewID("tnm", plateauapi.TypeDatasetType),
			TypeCode:           "tnm",
			PlateauSpecMinorID: plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
			Items: []*plateauapi.PlateauDatasetItem{
				{
					ID:       plateauapi.NewID("11111_tnm_AAA", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "AAA!",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_tnm_AAA_3dtiles.zip",
					Texture:  lo.ToPtr(plateauapi.TextureTexture),
					ParentID: plateauapi.NewID("11111_tnm_AAA", plateauapi.TypeDataset),
				},
				{
					ID:       plateauapi.NewID("11111_tnm_AAA_no_texture", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "AAA!（テクスチャなし）",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_tnm_AAA_3dtiles_no_texture.zip",
					Texture:  lo.ToPtr(plateauapi.TextureNone),
					ParentID: plateauapi.NewID("11111_tnm_AAA", plateauapi.TypeDataset),
				},
			},
		},
		&plateauapi.PlateauDataset{
			ID:                 plateauapi.NewID("11111_tnm_BBB", plateauapi.TypeDataset),
			Name:               "津波浸水想定区域モデル BBB!（bar市）",
			Description:        lo.ToPtr("desc2"),
			Year:               2023,
			PrefectureID:       lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode:     lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:             lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:           lo.ToPtr(plateauapi.AreaCode("11111")),
			TypeID:             plateauapi.NewID("tnm", plateauapi.TypeDatasetType),
			TypeCode:           "tnm",
			PlateauSpecMinorID: plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
			Items: []*plateauapi.PlateauDatasetItem{
				{
					ID:       plateauapi.NewID("11111_tnm_BBB", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "BBB!",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_tnm_BBB_3dtiles.zip",
					Texture:  lo.ToPtr(plateauapi.TextureTexture),
					ParentID: plateauapi.NewID("11111_tnm_BBB", plateauapi.TypeDataset),
				},
			},
		},
	}

	area := &areaContext{
		Pref: &plateauapi.Prefecture{},
		City: &plateauapi.City{
			Name: "bar市",
			Code: "11111",
		},
		PrefID:   lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
		CityID:   lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
		PrefCode: lo.ToPtr(plateauapi.AreaCode("11")),
		CityCode: lo.ToPtr(plateauapi.AreaCode("11111")),
		CityItem: &CityItem{
			Year: "2023年",
			Public: map[string]bool{
				"tnm": true,
			},
		},
	}

	dts := &plateauapi.PlateauDatasetType{
		ID:   plateauapi.NewID("tnm", plateauapi.TypeDatasetType),
		Code: "tnm",
		Name: "津波浸水想定区域モデル",
	}

	spec := &plateauapi.PlateauSpecMinor{
		ID:           plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
		Name:         "第3.2版",
		MajorVersion: 3,
		Version:      "3.2",
		Year:         2023,
		ParentID:     plateauapi.NewID("3", plateauapi.TypePlateauSpec),
	}

	res, warning := item.toDatasets(area, dts, spec)
	assert.Nil(t, warning)
	assert.Equal(t, expected, res)
}

func TestPlateauDataset_ToDatasets_Fld(t *testing.T) {
	item := &PlateauFeatureItem{
		ID:   "id",
		Desc: "desc",
		Items: []PlateauFeatureItemDatum{
			{
				ID: "id1",
				Data: []string{
					"https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_natl_yabegawa_haegawa_3dtiles_l1.zip",
					"https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_natl_yabegawa_haegawa_3dtiles_l1_no_texture.zip",
					"https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_natl_yabegawa_haegawa_3dtiles_l2.zip",
				},
				Desc: "desc1",
			},
			{
				ID: "id2",
				Data: []string{
					"https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_natl_yodogawa_ujigawa_3dtiles_l1.zip",
				},
				Desc: "desc2",
			},
			{
				ID: "id3",
				Data: []string{
					"https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_pref_yodogawa_ujigawa_3dtiles_l1.zip",
				},
				Desc: "desc3",
			},
		},
		Dic: `{
			"fld": [
				{ "name": "yabegawa_haegawa_l1", "description": "矢部川水系八重川", "scale": "計画規模", "admin": "国" },
				{ "name": "yabegawa_haegawa_l2", "description": "矢部川水系八重川", "scale": "想定最大規模", "admin": "国" },
				{ "name": "yodogawa_ujigawa_l1", "description": "淀川水系宇治川", "scale": "計画規模", "admin": "国" },
				{ "name": "yodogawa_ujigawa_l1", "description": "淀川水系宇治川", "scale": "計画規模", "admin": "都道府県" }
			]
		}`,
	}

	expected := []plateauapi.Dataset{
		&plateauapi.PlateauDataset{
			ID:                 plateauapi.NewID("11111_fld_natl_yabegawa_haegawa", plateauapi.TypeDataset),
			Name:               "洪水浸水想定区域モデル 矢部川水系八重川（国管理区間）（bar市）",
			Description:        lo.ToPtr("desc1"),
			Year:               2023,
			PrefectureID:       lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode:     lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:             lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:           lo.ToPtr(plateauapi.AreaCode("11111")),
			TypeID:             plateauapi.NewID("fld", plateauapi.TypeDatasetType),
			TypeCode:           "fld",
			PlateauSpecMinorID: plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
			Stage:              lo.ToPtr(string(stageBeta)),
			Items: []*plateauapi.PlateauDatasetItem{
				{
					ID:       plateauapi.NewID("11111_fld_natl_yabegawa_haegawa_l1", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "計画規模",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_natl_yabegawa_haegawa_3dtiles_l1.zip",
					Texture:  lo.ToPtr(plateauapi.TextureTexture),
					ParentID: plateauapi.NewID("11111_fld_natl_yabegawa_haegawa", plateauapi.TypeDataset),
				},
				{
					ID:       plateauapi.NewID("11111_fld_natl_yabegawa_haegawa_l1_no_texture", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "計画規模（テクスチャなし）",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_natl_yabegawa_haegawa_3dtiles_l1_no_texture.zip",
					Texture:  lo.ToPtr(plateauapi.TextureNone),
					ParentID: plateauapi.NewID("11111_fld_natl_yabegawa_haegawa", plateauapi.TypeDataset),
				},
				{
					ID:       plateauapi.NewID("11111_fld_natl_yabegawa_haegawa_l2", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "想定最大規模",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_natl_yabegawa_haegawa_3dtiles_l2.zip",
					Texture:  lo.ToPtr(plateauapi.TextureTexture),
					ParentID: plateauapi.NewID("11111_fld_natl_yabegawa_haegawa", plateauapi.TypeDataset),
				},
			},
			River: &plateauapi.River{
				Name:  "矢部川水系八重川",
				Admin: plateauapi.RiverAdminNational,
			},
		},
		&plateauapi.PlateauDataset{
			ID:                 plateauapi.NewID("11111_fld_natl_yodogawa_ujigawa", plateauapi.TypeDataset),
			Name:               "洪水浸水想定区域モデル 淀川水系宇治川（国管理区間）（bar市）",
			Description:        lo.ToPtr("desc2"),
			Year:               2023,
			PrefectureID:       lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode:     lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:             lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:           lo.ToPtr(plateauapi.AreaCode("11111")),
			TypeID:             plateauapi.NewID("fld", plateauapi.TypeDatasetType),
			TypeCode:           "fld",
			PlateauSpecMinorID: plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
			Stage:              lo.ToPtr(string(stageBeta)),
			Items: []*plateauapi.PlateauDatasetItem{
				{
					ID:       plateauapi.NewID("11111_fld_natl_yodogawa_ujigawa_l1", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "計画規模",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_natl_yodogawa_ujigawa_3dtiles_l1.zip",
					Texture:  lo.ToPtr(plateauapi.TextureTexture),
					ParentID: plateauapi.NewID("11111_fld_natl_yodogawa_ujigawa", plateauapi.TypeDataset),
				},
			},
			River: &plateauapi.River{
				Name:  "淀川水系宇治川",
				Admin: plateauapi.RiverAdminNational,
			},
		},
		&plateauapi.PlateauDataset{
			ID:                 plateauapi.NewID("11111_fld_pref_yodogawa_ujigawa", plateauapi.TypeDataset),
			Name:               "洪水浸水想定区域モデル 淀川水系宇治川（都道府県管理区間）（bar市）",
			Description:        lo.ToPtr("desc3"),
			Year:               2023,
			PrefectureID:       lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode:     lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:             lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:           lo.ToPtr(plateauapi.AreaCode("11111")),
			TypeID:             plateauapi.NewID("fld", plateauapi.TypeDatasetType),
			TypeCode:           "fld",
			PlateauSpecMinorID: plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
			Stage:              lo.ToPtr(string(stageBeta)),
			Items: []*plateauapi.PlateauDatasetItem{
				{
					ID:       plateauapi.NewID("11111_fld_pref_yodogawa_ujigawa_l1", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					Name:     "計画規模",
					URL:      "https://example.com/11111_bar-shi_city_2023_citygml_1_op_fld_pref_yodogawa_ujigawa_3dtiles_l1.zip",
					Texture:  lo.ToPtr(plateauapi.TextureTexture),
					ParentID: plateauapi.NewID("11111_fld_pref_yodogawa_ujigawa", plateauapi.TypeDataset),
				},
			},
			River: &plateauapi.River{
				Name:  "淀川水系宇治川",
				Admin: plateauapi.RiverAdminPrefecture,
			},
		},
	}

	area := &areaContext{
		Pref: &plateauapi.Prefecture{},
		City: &plateauapi.City{
			Name: "bar市",
			Code: "11111",
		},
		PrefID:   lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
		CityID:   lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
		PrefCode: lo.ToPtr(plateauapi.AreaCode("11")),
		CityCode: lo.ToPtr(plateauapi.AreaCode("11111")),
		CityItem: &CityItem{
			Year:              "2023年",
			PlateauDataStatus: ManagementStatusReady,
		},
	}

	dts := &plateauapi.PlateauDatasetType{
		ID:   plateauapi.NewID("fld", plateauapi.TypeDatasetType),
		Code: "fld",
		Name: "洪水浸水想定区域モデル",
	}

	spec := &plateauapi.PlateauSpecMinor{
		ID:           plateauapi.NewID("3.2", plateauapi.TypePlateauSpec),
		Name:         "第3.2版",
		MajorVersion: 3,
		Version:      "3.2",
		Year:         2023,
		ParentID:     plateauapi.NewID("3", plateauapi.TypePlateauSpec),
	}

	res, warning := item.toDatasets(area, dts, spec)
	assert.Nil(t, warning)
	assert.Equal(t, expected, res)
}