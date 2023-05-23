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
		[]*cms.PublicAsset{
			{
				URL: "https://example.com/01100_sapporo-shi_2020_mvt_op.zip",
			},
		},
		"xxxモデル",
		[]string{"layer"},
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
		[]*cms.PublicAsset{
			{
				URL: "https://example.com/01100_sapporo-shi_2020_mvt_op.zip",
			},
			{
				URL: "https://example.com/01100_sapporo-shi_2020_3dtiles_op.zip",
			},
		},
		"xxxモデル",
		[]string{"layer"},
	))

	// case3: multiple asset with LOD
	assert.Equal(t, DataCatalogItemConfig{
		Data: []DataCatalogItemConfigItem{
			{
				Name:   "LOD1",
				URL:    "https://example.com/01100_sapporo-shi_2020_mvt_op_tran_lod1/{z}/{x}/{y}.mvt",
				Type:   "mvt",
				Layers: []string{"layer"},
			},
			{
				Name:   "LOD2",
				URL:    "https://example.com/01100_sapporo-shi_2020_3dtiles_op_tran_lod2/tileset.json",
				Type:   "3dtiles",
				Layers: []string{"layer"},
			},
		},
	}, multipleLODData(
		[]*cms.PublicAsset{
			{
				URL: "https://example.com/01100_sapporo-shi_2020_mvt_op_tran_lod1.zip",
			},
			{
				URL: "https://example.com/01100_sapporo-shi_2020_3dtiles_op_tran_lod2.zip",
			},
		},
		"xxxモデル",
		[]string{"layer"},
	))
}
