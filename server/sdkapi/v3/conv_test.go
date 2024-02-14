package sdkapiv3

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestQueryToDatasets(t *testing.T) {
	query := &Query{
		Areas: []*QueryArea{
			{
				Name: "Prefecture 1",
				Prefecture: &QueryPrefecture{
					Cities: []*QueryCity{
						{
							ID:   "City1",
							Name: "City 1",
							Datasets: []*QueryCityDataset{
								{
									TypeCode: "DatasetType1",
								},
								{
									TypeCode: "DatasetType2",
								},
							},
						},
						{
							ID:   "City2",
							Name: "City 2",
							Datasets: []*QueryCityDataset{
								{
									TypeCode: "DatasetType3",
								},
							},
						},
					},
				},
			},
		},
	}

	expected := &DatasetsResponse{
		Data: []*DatasetPrefectureResponse{
			{
				Title: "Prefecture 1",
				Data: []*DatasetCityResponse{
					{
						ID:           "City1",
						Title:        "City 1",
						Spec:         "",
						Description:  "",
						FeatureTypes: []string{"DatasetType1", "DatasetType2"},
					},
					{
						ID:           "City2",
						Title:        "City 2",
						Spec:         "",
						Description:  "",
						FeatureTypes: []string{"DatasetType3"},
					},
				},
			},
		},
	}

	datasets := query.ToDatasets()
	assert.Equal(t, expected, datasets)
}
