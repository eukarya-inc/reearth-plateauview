package datacatalogv3

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func Test_GenericItem_ToDatasets(t *testing.T) {
	item := &GenericItem{
		ID:   "id",
		Name: "name",
		Desc: "desc",
		Data: []GenericItemDataset{
			{
				ID:         "id1",
				Name:       "name1",
				Data:       "url1",
				Desc:       "desc1",
				DataFormat: "3D Tiles",
			},
			{
				ID:        "id2",
				Data:      "url2",
				DataURL:   "https://example.com/{z}/{x}/{y}.mvt",
				Desc:      "desc2",
				LayerName: "layer1, layer2",
			},
			// invalid item
			{
				ID: "id3",
			},
		},
		Category: "ユースケース",
	}

	expected := []plateauapi.Dataset{
		&plateauapi.GenericDataset{
			ID:             plateauapi.NewID("id", plateauapi.TypeDataset),
			Name:           "name",
			Description:    toPtrIfPresent("desc"),
			Year:           2023,
			PrefectureID:   lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode: lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:         lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:       lo.ToPtr(plateauapi.AreaCode("11111")),
			TypeID:         plateauapi.NewID("usecase", plateauapi.TypeDatasetType),
			TypeCode:       "usecase",
			Stage:          lo.ToPtr(string(stageAlpha)),
			Items: []*plateauapi.GenericDatasetItem{
				{
					ID:       plateauapi.NewID("id1", plateauapi.TypeDatasetItem),
					Name:     "name1",
					URL:      "url1",
					Format:   plateauapi.DatasetFormatCesium3dtiles,
					ParentID: plateauapi.NewID("id", plateauapi.TypeDataset),
				},
				{
					ID:       plateauapi.NewID("id2", plateauapi.TypeDatasetItem),
					Name:     "name 2",
					URL:      "https://example.com/{z}/{x}/{y}.mvt",
					Format:   plateauapi.DatasetFormatMvt,
					Layers:   []string{"layer1", "layer2"},
					ParentID: plateauapi.NewID("id", plateauapi.TypeDataset),
				},
			},
		},
	}

	area := &areaContext{
		PrefID:   lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
		PrefCode: lo.ToPtr(plateauapi.AreaCode("11")),
		CityID:   lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
		CityCode: lo.ToPtr(plateauapi.AreaCode("11111")),
		CityItem: &CityItem{
			Year: "令和5年度",
		},
	}

	dts := []plateauapi.DatasetType{
		&plateauapi.GenericDatasetType{
			ID:   plateauapi.NewID("usecase", plateauapi.TypeDatasetType),
			Code: "usecase",
			Name: "ユースケース",
		},
	}

	res, warning := item.toDatasets(area, dts)
	assert.Equal(t, []string{"generic id[2]: invalid url: "}, warning)
	assert.Equal(t, expected, res)
}
