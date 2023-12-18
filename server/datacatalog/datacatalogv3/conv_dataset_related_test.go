package datacatalogv3

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestRelatedItem_ToDatasets(t *testing.T) {
	item := &RelatedItem{
		ID: "id",
		Assets: map[string][]string{
			"landmark": {
				"https://example.com/11112_hoge-ku_city_2023_landmark.geojson",
				"https://example.com/11113_foo-ku_city_2023_landmark.geojson",
			},
			"border": {"https://example.com/11111_bar-shi_city_2023_border.geojson"},
		},
		ConvertedAssets: map[string][]string{
			"landmark": {
				"https://example.com/11112_hoge-ku_city_2023_landmark.czml",
				"https://example.com/11113_foo-ku_city_2023_landmark.czml",
			},
		},
		Desc: "desc",
	}

	expected := []plateauapi.Dataset{
		&plateauapi.RelatedDataset{
			ID:             plateauapi.NewID("11112_landmark", plateauapi.TypeDataset),
			Name:           "ランドマーク情報（hoge区）",
			Description:    toPtrIfPresent("desc"),
			Year:           2023,
			PrefectureID:   lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode: lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:         lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:       lo.ToPtr(plateauapi.AreaCode("11111")),
			WardID:         lo.ToPtr(plateauapi.NewID("11112", plateauapi.TypeArea)),
			WardCode:       lo.ToPtr(plateauapi.AreaCode("11112")),
			TypeID:         plateauapi.NewID("landmark", plateauapi.TypeDatasetType),
			TypeCode:       "landmark",
			Stage:          lo.ToPtr(string(stageAlpha)),
			Items: []*plateauapi.RelatedDatasetItem{
				{
					ID:       plateauapi.NewID("11112_landmark", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCzml,
					Name:     "ランドマーク情報",
					URL:      "https://example.com/11112_hoge-ku_city_2023_landmark.czml",
					ParentID: plateauapi.NewID("11112_landmark", plateauapi.TypeDataset),
				},
			},
		},
		&plateauapi.RelatedDataset{
			ID:             plateauapi.NewID("11113_landmark", plateauapi.TypeDataset),
			Name:           "ランドマーク情報（foo区）",
			Description:    toPtrIfPresent("desc"),
			Year:           2023,
			PrefectureID:   lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode: lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:         lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:       lo.ToPtr(plateauapi.AreaCode("11111")),
			WardID:         lo.ToPtr(plateauapi.NewID("11113", plateauapi.TypeArea)),
			WardCode:       lo.ToPtr(plateauapi.AreaCode("11113")),
			TypeID:         plateauapi.NewID("landmark", plateauapi.TypeDatasetType),
			TypeCode:       "landmark",
			Stage:          lo.ToPtr(string(stageAlpha)),
			Items: []*plateauapi.RelatedDatasetItem{
				{
					ID:       plateauapi.NewID("11113_landmark", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatCzml,
					Name:     "ランドマーク情報",
					URL:      "https://example.com/11113_foo-ku_city_2023_landmark.czml",
					ParentID: plateauapi.NewID("11113_landmark", plateauapi.TypeDataset),
				},
			},
		},
		&plateauapi.RelatedDataset{
			ID:             plateauapi.NewID("11111_border", plateauapi.TypeDataset),
			Name:           "行政界情報（bar市）",
			Description:    toPtrIfPresent("desc"),
			Year:           2023,
			PrefectureID:   lo.ToPtr(plateauapi.NewID("11", plateauapi.TypeArea)),
			PrefectureCode: lo.ToPtr(plateauapi.AreaCode("11")),
			CityID:         lo.ToPtr(plateauapi.NewID("11111", plateauapi.TypeArea)),
			CityCode:       lo.ToPtr(plateauapi.AreaCode("11111")),
			TypeID:         plateauapi.NewID("border", plateauapi.TypeDatasetType),
			TypeCode:       "border",
			Stage:          lo.ToPtr(string(stageAlpha)),
			Items: []*plateauapi.RelatedDatasetItem{
				{
					ID:       plateauapi.NewID("11111_border", plateauapi.TypeDatasetItem),
					Format:   plateauapi.DatasetFormatGeojson,
					Name:     "行政界情報",
					URL:      "https://example.com/11111_bar-shi_city_2023_border.geojson",
					ParentID: plateauapi.NewID("11111_border", plateauapi.TypeDataset),
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

	dts := []plateauapi.DatasetType{
		&plateauapi.RelatedDatasetType{
			ID:   plateauapi.NewID("landmark", plateauapi.TypeDatasetType),
			Code: "landmark",
			Name: "ランドマーク情報",
		},
		&plateauapi.RelatedDatasetType{
			ID:   plateauapi.NewID("border", plateauapi.TypeDatasetType),
			Code: "border",
			Name: "行政界情報",
		},
	}

	res, warnings := item.toDatasets(area, dts)
	assert.Nil(t, warnings)
	assert.Equal(t, expected, res)
}
