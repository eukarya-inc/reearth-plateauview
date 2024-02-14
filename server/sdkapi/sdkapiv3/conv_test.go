package sdkapiv3

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestQueryToDatasets(t *testing.T) {
	query := &Query{
		Areas: []QueryArea{
			{
				Name: "Prefecture 1",
				Prefecture: QueryPrefecture{
					Cities: []QueryCity{
						{
							ID:   "City1",
							Name: "City 1",
							Datasets: []QueryDataset{
								{
									TypeCode:    "bldg",
									Description: "Description",
									PlateauDataset: QueryPlateauDataset{
										PlateauSpecMinor: QueryPlateauSpecMinor{
											Version: "3.4",
										},
									},
								},
								{
									TypeCode: "DatasetType1",
									PlateauDataset: QueryPlateauDataset{
										PlateauSpecMinor: QueryPlateauSpecMinor{
											Version: "3.5",
										},
									},
								},
							},
						},
						{
							ID:   "City2",
							Name: "City 2",
							Datasets: []QueryDataset{
								{
									TypeCode: "DatasetType2",
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
						Spec:         "3.4",
						Description:  "Description",
						FeatureTypes: []string{"bldg", "DatasetType1"},
					},
					{
						ID:           "City2",
						Title:        "City 2",
						Spec:         "",
						Description:  "",
						FeatureTypes: []string{"DatasetType2"},
					},
				},
			},
		},
	}

	datasets := query.ToDatasets()
	assert.Equal(t, expected, datasets)
}
