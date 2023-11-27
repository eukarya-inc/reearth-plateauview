package cmsintegrationv3

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/stretchr/testify/assert"
)

func TestParseSetupCSV(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		name             string
		csvData          string
		expected         []SetupCSVItem
		expectedFeatures []string
		err              error
	}{
		{
			name: "valid csv",
			csvData: `Name,NameEn,Code,Prefecture,Feature1,Feature2,Feature3
Tokyo,東京,13,Tokyo,Yes,,Yes
Osaka,大阪,27,Osaka,Yes,Yes,`,
			expected: []SetupCSVItem{
				{
					Name:       "Tokyo",
					NameEn:     "東京",
					Code:       "13",
					Prefecture: "Tokyo",
					Features:   []string{"Feature1", "Feature3"},
				},
				{
					Name:       "Osaka",
					NameEn:     "大阪",
					Code:       "27",
					Prefecture: "Osaka",
					Features:   []string{"Feature1", "Feature2"},
				},
			},
			expectedFeatures: []string{"Feature1", "Feature2", "Feature3"},
			err:              nil,
		},
		{
			name:             "empty csv",
			csvData:          "",
			expected:         nil,
			expectedFeatures: nil,
			err:              io.EOF,
		},
		{
			name: "invalid header",
			csvData: `Name,NameEn,Code,Prefecture
Tokyo,東京,13,Tokyo`,
			expected:         nil,
			expectedFeatures: nil,
			err:              fmt.Errorf("invalid header: [Name NameEn Code Prefecture]"),
		},
		{
			name: "invalid row",
			csvData: `Name,NameEn,Code,Prefecture,Feature1
Tokyo,東京,13,Tokyo`,
			expected: nil,
			err:      fmt.Errorf("record on line 2: wrong number of fields"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.csvData)
			items, features, err := parseSetupCSV(ctx, r)

			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, items)
			assert.Equal(t, tt.expectedFeatures, features)
		})
	}
}

func TestSetupCityItems(t *testing.T) {
	ctx := context.Background()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://example.com/data.csv",
		httpmock.NewStringResponder(200, `Name,NameEn,Code,Prefecture,bldg,tran,luse
八王子市,hachioji-shi,13201,東京都,Yes,,Yes
東村山市,higashimurayama-shi,13213,東京都,Yes,Yes,`))

	var createdItems []*cms.Item
	var updateditems []*cms.Item

	s := &Services{
		CMS: &cmsMock{
			getModels: func(ctx context.Context, projectID string) (*cms.Models, error) {
				return &cms.Models{
					Models: []cms.Model{
						{
							ID:  "city",
							Key: "plateau-city",
						},
						{
							ID:  "bldg",
							Key: "plateau-bldg",
						},
						{
							ID:  "tran",
							Key: "plateau-tran",
						},
						{
							ID:  "luse",
							Key: "plateau-luse",
						},
						{
							ID:  "related",
							Key: "plateau-related",
						},
					},
				}, nil
			},
			createItem: func(ctx context.Context, modelID string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
				item := &cms.Item{
					ID:             fmt.Sprintf("item%d", len(createdItems)),
					ModelID:        modelID,
					Fields:         fields,
					MetadataFields: metadataFields,
				}
				createdItems = append(createdItems, item)
				return item, nil
			},
			updateItem: func(ctx context.Context, itemID string, fields []*cms.Field, metadataFields []*cms.Field) (*cms.Item, error) {
				item := &cms.Item{
					ID:             itemID,
					Fields:         fields,
					MetadataFields: metadataFields,
				}
				updateditems = append(updateditems, item)
				return item, nil
			},
		},
		HTTP: http.DefaultClient,
	}

	inp := SetupCityItemsInput{
		ProjectID: "project123",
		DataURL:   "https://example.com/data.csv",
	}

	onprogress := func(i, l int) {}

	t.Run("success", func(t *testing.T) {
		createdItems = nil
		updateditems = nil
		err := SetupCityItems(ctx, s, inp, onprogress)
		assert.NoError(t, err)

		assert.Equal(t, 10, len(createdItems))
		assertCityItem(t, &CityItem{
			ID:         "item0",
			CityName:   "八王子市",
			CityNameEn: "hachioji-shi",
			CityCode:   "13201",
			Prefecture: "東京都",
		}, createdItems[0])
		assertFeatureItem(t, "related", "item0", "", createdItems[1])
		assertFeatureItem(t, "bldg", "item0", "", createdItems[2])
		assertFeatureItem(t, "tran", "item0", ManagementStatusSkip, createdItems[3])
		assertFeatureItem(t, "luse", "item0", "", createdItems[4])
		assertCityItem(t, &CityItem{
			ID:         "item5",
			CityName:   "東村山市",
			CityNameEn: "higashimurayama-shi",
			CityCode:   "13213",
			Prefecture: "東京都",
		}, createdItems[5])
		assertFeatureItem(t, "related", "item5", "", createdItems[6])
		assertFeatureItem(t, "bldg", "item5", "", createdItems[7])
		assertFeatureItem(t, "tran", "item5", "", createdItems[8])
		assertFeatureItem(t, "luse", "item5", ManagementStatusSkip, createdItems[9])

		assert.Equal(t, 2, len(updateditems))
		assertUpdatedCityItem(t, &CityItem{
			ID: "item0",
			References: map[string]string{
				"bldg": "item2",
				"tran": "item3",
				"luse": "item4",
			},
			RelatedDataset: "item1",
		}, updateditems[0])
		assertUpdatedCityItem(t, &CityItem{
			ID: "item5",
			References: map[string]string{
				"bldg": "item7",
				"tran": "item8",
				"luse": "item9",
			},
			RelatedDataset: "item6",
		}, updateditems[1])
	})
}

func assertCityItem(t *testing.T, expected *CityItem, actual *cms.Item) {
	assert.Equal(t, "city", actual.ModelID)
	a := CityItemFrom(actual)
	am := &CityItem{
		ID:         a.ID,
		CityName:   a.CityName,
		CityNameEn: a.CityNameEn,
		CityCode:   a.CityCode,
		Prefecture: a.Prefecture,
	}
	assert.Equal(t, expected, am)
}

func assertFeatureItem(t *testing.T, expectedModel, expectedCity string, status ManagementStatus, actual *cms.Item) {
	assert.Equal(t, expectedModel, actual.ModelID, "model of "+actual.ID)
	assert.Equal(t, expectedCity, actual.FieldByKey("city").GetValue().Interface(), "city of "+actual.ID)
	statusv := actual.MetadataFieldByKey("status").GetValue().Interface()
	if status == "" {
		assert.Nil(t, statusv, "status of "+actual.ID)
		return
	}
	assert.Equal(t, status, statusv, "status of "+actual.ID)
}

func assertUpdatedCityItem(t *testing.T, expected *CityItem, actual *cms.Item) {
	a := CityItemFrom(actual)
	am := &CityItem{
		ID:             a.ID,
		References:     a.References,
		RelatedDataset: a.RelatedDataset,
	}
	assert.Equal(t, expected, am)
}
