package datacatalogv2adapter

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestFilterArea(t *testing.T) {
	testCases := []struct {
		name     string
		area     plateauapi.Area
		input    plateauapi.AreaInput
		expected bool
	}{
		{
			name: "Prefecture with search tokens",
			area: plateauapi.Prefecture{Name: "Tokyo"},
			input: plateauapi.AreaInput{
				SearchTokens: []string{"Tokyo"},
			},
			expected: true,
		},
		{
			name: "Prefecture without search tokens",
			area: plateauapi.Prefecture{Name: "Tokyo"},
			input: plateauapi.AreaInput{
				SearchTokens: []string{},
			},
			expected: true,
		},
		{
			name: "Prefecture without non-matching search tokens",
			area: plateauapi.Prefecture{Name: "Tokyo"},
			input: plateauapi.AreaInput{
				SearchTokens: []string{"Kanagawa"},
			},
			expected: false,
		},
		{
			name: "City with search tokens and matching parent code",
			area: plateauapi.City{Name: "Shinjuku", PrefectureCode: "13"},
			input: plateauapi.AreaInput{
				SearchTokens: []string{"Shinjuku"},
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("13")),
			},
			expected: true,
		},
		{
			name: "City with search tokens and non-matching parent code",
			area: plateauapi.City{Name: "Shinjuku", PrefectureCode: "13"},
			input: plateauapi.AreaInput{
				SearchTokens: []string{"Shinjuku"},
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("14")),
			},
			expected: false,
		},
		{
			name: "Ward with search tokens and matching parent code",
			area: plateauapi.Ward{Name: "Shinjuku", PrefectureCode: "13", CityCode: "13104"},
			input: plateauapi.AreaInput{
				SearchTokens: []string{"Shinjuku"},
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("13104")),
			},
			expected: true,
		},
		{
			name: "Ward with search tokens and non-matching parent code",
			area: plateauapi.Ward{Name: "Shinjuku", PrefectureCode: "13", CityCode: "13104"},
			input: plateauapi.AreaInput{
				SearchTokens: []string{"Shinjuku"},
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("13105")),
			},
			expected: false,
		},
		{
			name: "Ward without search tokens",
			area: plateauapi.Ward{Name: "Shinjuku", PrefectureCode: "13", CityCode: "13104"},
			input: plateauapi.AreaInput{
				SearchTokens: []string{},
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual := filterArea(tc.area, tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestFilterByCode(t *testing.T) {
	assert.True(t, filterByCode("a", []string{"a"}, nil))
	assert.True(t, filterByCode("a", []string{"a", "b"}, nil))
	assert.True(t, filterByCode("b", []string{"a", "b"}, nil))
	assert.False(t, filterByCode("b", []string{"a"}, nil))
	assert.False(t, filterByCode("a", nil, []string{"a"}))
	assert.False(t, filterByCode("a", nil, []string{"a", "b"}))
	assert.False(t, filterByCode("a", []string{"a"}, []string{"a", "b"}))
}

func TestFilterByPlateauSpec(t *testing.T) {
	testCases := []struct {
		name        string
		querySpec   *string
		datasetSpec string
		expected    bool
	}{
		{
			name:        "Nil query spec and empty dataset spec",
			querySpec:   nil,
			datasetSpec: "",
			expected:    true,
		},
		{
			name:        "Empty query spec and empty dataset spec",
			datasetSpec: "",
			querySpec:   lo.ToPtr(""),
			expected:    true,
		},
		{
			name:        "Nil query spec and non-empty dataset spec",
			querySpec:   nil,
			datasetSpec: "1.0",
			expected:    true,
		},
		{
			name:        "Empty query spec and non-empty dataset spec",
			querySpec:   lo.ToPtr(""),
			datasetSpec: "1.0",
			expected:    true,
		},
		{
			name:        "Non-empty query spec and non-empty dataset spec with matching major version",
			querySpec:   lo.ToPtr("1"),
			datasetSpec: "1.2",
			expected:    true,
		},
		{
			name:        "Non-empty query spec and non-empty dataset spec with non-matching major version",
			querySpec:   lo.ToPtr("1"),
			datasetSpec: "2.0",
			expected:    false,
		},
		{
			name:        "Non-empty query spec and non-empty dataset spec with matching major and minor version",
			querySpec:   lo.ToPtr("1.2"),
			datasetSpec: "1.2",
			expected:    true,
		},
		{
			name:        "Non-empty query spec and non-empty dataset spec with non-matching minor version",
			querySpec:   lo.ToPtr("1.2"),
			datasetSpec: "1.3",
			expected:    false,
		},
		{
			name:        "Non-empty query spec and non-empty dataset spec with non-matching major version and matching minor version",
			querySpec:   lo.ToPtr("1.2"),
			datasetSpec: "2.2",
			expected:    false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual := filterByPlateauSpec(tc.querySpec, tc.datasetSpec)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestFilterDataType(t *testing.T) {
	testCases := []struct {
		name     string
		ty       plateauapi.DatasetType
		input    plateauapi.DatasetTypeInput
		expected bool
	}{
		{
			name: "PlateauDatasetType with matching category, year and plateau spec",
			ty: plateauapi.PlateauDatasetType{
				Year:            2021,
				PlateauSpecName: "1.0",
			},
			input: plateauapi.DatasetTypeInput{
				Category:    lo.ToPtr(plateauapi.DatasetTypeCategoryPlateau),
				Year:        lo.ToPtr(2021),
				PlateauSpec: lo.ToPtr("1.0"),
			},
			expected: true,
		},
		{
			name: "PlateauDatasetType with matching plateau spec major version",
			ty: plateauapi.PlateauDatasetType{
				PlateauSpecName: "1.0",
			},
			input: plateauapi.DatasetTypeInput{
				PlateauSpec: lo.ToPtr("1"),
			},
			expected: true,
		},
		{
			name: "PlateauDatasetType with non-matching category",
			ty: plateauapi.PlateauDatasetType{
				Year:            2021,
				PlateauSpecName: "1.0",
			},
			input: plateauapi.DatasetTypeInput{
				Category: lo.ToPtr(plateauapi.DatasetTypeCategoryRelated),
			},
			expected: false,
		},
		{
			name: "PlateauDatasetType with non-matching year",
			ty: plateauapi.PlateauDatasetType{
				Year:            2021,
				PlateauSpecName: "1.0",
			},
			input: plateauapi.DatasetTypeInput{
				Year: lo.ToPtr(2022),
			},
			expected: false,
		},
		{
			name: "PlateauDatasetType with non-matching plateau spec",
			ty: plateauapi.PlateauDatasetType{
				Year:            2021,
				PlateauSpecName: "1.0",
			},
			input: plateauapi.DatasetTypeInput{
				PlateauSpec: lo.ToPtr("2.0"),
			},
			expected: false,
		},
		{
			name: "RelatedDatasetType with matching category",
			ty:   plateauapi.RelatedDatasetType{},
			input: plateauapi.DatasetTypeInput{
				Category: lo.ToPtr(plateauapi.DatasetTypeCategoryRelated),
			},
			expected: true,
		},
		{
			name: "RelatedDatasetType with non-matching category",
			ty:   plateauapi.RelatedDatasetType{},
			input: plateauapi.DatasetTypeInput{
				Category: lo.ToPtr(plateauapi.DatasetTypeCategoryPlateau),
			},
			expected: false,
		},
		{
			name: "GenericDatasetType with matching category",
			ty:   plateauapi.GenericDatasetType{},
			input: plateauapi.DatasetTypeInput{
				Category: lo.ToPtr(plateauapi.DatasetTypeCategoryGeneric),
			},
			expected: true,
		},
		{
			name: "GenericDatasetType with non-matching category",
			ty:   plateauapi.GenericDatasetType{},
			input: plateauapi.DatasetTypeInput{
				Category: lo.ToPtr(plateauapi.DatasetTypeCategoryPlateau),
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actual := filterDataType(tc.ty, tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
