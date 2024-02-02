package datacatalogv3

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/jarcoal/httpmock"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestRepos(t *testing.T) {
	ctx := context.Background()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mockCMS(t)

	cms := lo.Must(cms.New("https://example.com", "token"))

	repos := NewRepos()
	err := repos.Prepare(ctx, "prj", cms)
	assert.NoError(t, err)
	assert.Nil(t, repos.Warnings("prj"))

	assertRes := func(t *testing.T, r plateauapi.Repo, cityName, cityCode string, stage *string, found bool) {
		t.Helper()

		prefCode := cityCode[:2]
		area, err := r.Area(ctx, plateauapi.AreaCode(cityCode))
		assert.NoError(t, err)
		assert.Equal(t, &plateauapi.City{
			ID:             plateauapi.ID("a_" + cityCode),
			Type:           plateauapi.AreaTypeCity,
			Code:           plateauapi.AreaCode(cityCode),
			Name:           cityName,
			PrefectureID:   plateauapi.ID("a_" + prefCode),
			PrefectureCode: plateauapi.AreaCode(prefCode),
		}, area)

		dataset, err := r.Datasets(ctx, &plateauapi.DatasetsInput{
			AreaCodes: []plateauapi.AreaCode{plateauapi.AreaCode(cityCode)},
		})
		assert.NoError(t, err)

		var admin any
		if stage != nil {
			admin = map[string]any{
				"stage": string(*stage),
			}
		}

		if found {
			assert.Equal(t, []plateauapi.Dataset{
				&plateauapi.PlateauDataset{
					ID:                 plateauapi.ID("d_" + cityCode + "_bldg"),
					Name:               "建築物モデル（" + cityName + "）",
					Year:               2023,
					PrefectureID:       lo.ToPtr(plateauapi.ID("a_00")),
					PrefectureCode:     lo.ToPtr(plateauapi.AreaCode("00")),
					CityID:             lo.ToPtr(plateauapi.ID("a_" + cityCode)),
					CityCode:           lo.ToPtr(plateauapi.AreaCode(cityCode)),
					TypeID:             plateauapi.NewID("bldg_3", plateauapi.TypeDatasetType),
					TypeCode:           "bldg",
					PlateauSpecMinorID: plateauapi.ID("ps_3.2"),
					Items: []*plateauapi.PlateauDatasetItem{
						{
							ID:       plateauapi.ID("di_" + cityCode + "_bldg_lod1"),
							Format:   plateauapi.DatasetFormatCesium3dtiles,
							URL:      "https://example.com/00000_hoge_city_2023_citygml_1_op_bldg_3dtiles_lod1/tileset.json",
							Name:     "LOD1",
							ParentID: plateauapi.ID("d_" + cityCode + "_bldg"),
							Lod:      lo.ToPtr(1),
							Texture:  lo.ToPtr(plateauapi.TextureTexture),
						},
					},
					Admin: admin,
				},
			}, dataset)
		} else {
			assert.Len(t, dataset, 0)
		}
	}

	radmin := repos.Repo("prj", true)
	assertRes(t, radmin, "hoge", "00000", lo.ToPtr(string(stageBeta)), true)
	assertRes(t, radmin, "foo", "00001", nil, true)
	assertRes(t, radmin, "bar", "00002", nil, false)

	rpublic := repos.Repo("prj", false)
	assertRes(t, rpublic, "hoge", "00000", nil, false)
	assertRes(t, rpublic, "foo", "00001", nil, true)
	assertRes(t, rpublic, "bar", "00002", nil, false)

	assert.NoError(t, repos.UpdateAll(ctx))
}

func mockCMS(t *testing.T) {
	t.Helper()
	httpmock.RegisterResponder(
		"GET", "https://example.com/api/projects/prj/models/plateau-city/items",
		httpmock.NewJsonResponderOrPanic(200, cities),
	)
	httpmock.RegisterResponder(
		"GET", "https://example.com/api/projects/prj/models/plateau-related/items",
		httpmock.NewJsonResponderOrPanic(200, empty),
	)
	httpmock.RegisterResponder(
		"GET", "https://example.com/api/projects/prj/models/plateau-generic/items",
		httpmock.NewJsonResponderOrPanic(200, empty),
	)
	for _, ft := range plateauFeatureTypes {
		res := empty
		if ft.Code == "bldg" {
			res = bldg
		}
		httpmock.RegisterResponder(
			"GET", "https://example.com/api/projects/prj/models/plateau-"+ft.Code+"/items",
			httpmock.NewJsonResponderOrPanic(200, res),
		)
	}
}

func j(j string) any {
	var v any
	lo.Must0(json.Unmarshal([]byte(j), &v))
	return v
}

var cities = j(`{
	"totalCount": 1,
	"items": [
		{
			"id": "city0",
			"fields": [
				{
					"key": "prefecture",
					"value": "PREF"
				},
				{
					"key": "city_name",
					"value": "hoge"
				},
				{
					"key": "city_code",
					"value": "00000"
				},
				{
					"key": "bldg",
					"value": "bldg0"
				},
				{
					"key": "spec",
					"value": "第3.2版"
				}
			],
			"metadataFields": [
				{
					"key": "plateau_data_status",
					"value": {
						"name": "確認可能"
					}
				}
			]
		},
		{
			"id": "city1",
			"fields": [
				{
					"key": "prefecture",
					"value": "PREF"
				},
				{
					"key": "city_name",
					"value": "foo"
				},
				{
					"key": "city_code",
					"value": "00001"
				},
				{
					"key": "bldg",
					"value": "bldg1"
				},
				{
					"key": "spec",
					"value": "第3.2版"
				}
			],
			"metadataFields": [
				{
					"key": "bldg_public",
					"value": true
				}
			]
		},
		{
			"id": "city2",
			"fields": [
				{
					"key": "prefecture",
					"value": "PREF"
				},
				{
					"key": "city_name",
					"value": "bar"
				},
				{
					"key": "city_code",
					"value": "00002"
				},
				{
					"key": "bldg",
					"value": "bldg2"
				},
				{
					"key": "spec",
					"value": "第3.2版"
				}
			],
			"metadataFields": [
			]
		}
	]
}`)

var bldg = j(`{
	"totalCount": 1,
	"items": [
		{
			"id": "bldg0",
			"fields": [
				{
					"key": "city",
					"value": "city0"
				},
				{
					"key": "data",
					"value": [{"url": "https://example.com/00000_hoge_city_2023_citygml_1_op_bldg_3dtiles_lod1.zip"}]
				}
			]
		},
		{
			"id": "bldg1",
			"fields": [
				{
					"key": "city",
					"value": "city1"
				},
				{
					"key": "data",
					"value": [{"url": "https://example.com/00000_hoge_city_2023_citygml_1_op_bldg_3dtiles_lod1.zip"}]
				}
			]
		},
		{
			"id": "bldg2",
			"fields": [
				{
					"key": "city",
					"value": "city2"
				},
				{
					"key": "data",
					"value": [{"url": "https://example.com/00000_hoge_city_2023_citygml_1_op_bldg_3dtiles_lod1.zip"}]
				}
			]
		}
	]
}`)

var empty = j(`{
	"totalCount": 0,
	"items": []
}`)
