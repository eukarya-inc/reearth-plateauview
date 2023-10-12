package datacatalogv2adapter

import (
	"context"
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestAdapter_Area(t *testing.T) {
	a := &Adapter{
		prefectures: []plateauapi.Prefecture{
			{Code: "01", Name: "北海道"},
			{Code: "02", Name: "青森県"},
			{Code: "03", Name: "岩手県"},
		},
		cities: []plateauapi.City{
			{Code: "01100", Name: "札幌市", PrefectureCode: "01"},
			{Code: "02100", Name: "青森市", PrefectureCode: "02"},
			{Code: "02101", Name: "弘前市", PrefectureCode: "02"},
			{Code: "03100", Name: "盛岡市", PrefectureCode: "03"},
			{Code: "03101", Name: "花巻市", PrefectureCode: "03"},
		},
		wards: []plateauapi.Ward{
			{Code: "01101", Name: "中央区", CityCode: "01100", PrefectureCode: "01"},
			{Code: "01102", Name: "北区", CityCode: "01100", PrefectureCode: "01"},
		},
	}

	tests := []struct {
		name     string
		code     plateauapi.AreaCode
		expected plateauapi.Area
	}{
		{
			name:     "prefecture",
			code:     "01",
			expected: &plateauapi.Prefecture{Code: "01", Name: "北海道"},
		},
		{
			name:     "city",
			code:     "01100",
			expected: &plateauapi.City{Code: "01100", Name: "札幌市", PrefectureCode: "01"},
		},
		{
			name:     "ward",
			code:     "01101",
			expected: &plateauapi.Ward{Code: "01101", Name: "中央区", CityCode: "01100", PrefectureCode: "01"},
		},
		{
			name:     "not found",
			code:     "99999",
			expected: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.Area(context.Background(), tt.code)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestAdapter_Areas(t *testing.T) {
	a := &Adapter{
		prefectures: []plateauapi.Prefecture{
			{Code: "01", Name: "北海道"},
			{Code: "02", Name: "青森県"},
			{Code: "03", Name: "岩手県"},
		},
		cities: []plateauapi.City{
			{Code: "01100", Name: "札幌市", PrefectureCode: "01"},
			{Code: "02100", Name: "青森市", PrefectureCode: "02"},
			{Code: "02101", Name: "弘前市", PrefectureCode: "02"},
			{Code: "03100", Name: "盛岡市", PrefectureCode: "03"},
			{Code: "03101", Name: "花巻市", PrefectureCode: "03"},
		},
		wards: []plateauapi.Ward{
			{Code: "01101", Name: "中央区", CityCode: "01100", PrefectureCode: "01"},
			{Code: "01102", Name: "北区", CityCode: "01100", PrefectureCode: "01"},
		},
		areasForDataTypes: map[string]map[plateauapi.AreaCode]struct{}{
			"bldg": {
				"01101": {},
			},
		},
	}

	tests := []struct {
		name  string
		input *plateauapi.AreaInput
		want  []plateauapi.Area
	}{
		{
			name:  "no filter",
			input: nil,
			want: []plateauapi.Area{
				&plateauapi.Prefecture{Code: "01", Name: "北海道"},
				&plateauapi.Prefecture{Code: "02", Name: "青森県"},
				&plateauapi.Prefecture{Code: "03", Name: "岩手県"},
				&plateauapi.City{Code: "01100", Name: "札幌市", PrefectureCode: "01"},
				&plateauapi.City{Code: "02100", Name: "青森市", PrefectureCode: "02"},
				&plateauapi.City{Code: "02101", Name: "弘前市", PrefectureCode: "02"},
				&plateauapi.City{Code: "03100", Name: "盛岡市", PrefectureCode: "03"},
				&plateauapi.City{Code: "03101", Name: "花巻市", PrefectureCode: "03"},
				&plateauapi.Ward{Code: "01101", Name: "中央区", CityCode: "01100", PrefectureCode: "01"},
				&plateauapi.Ward{Code: "01102", Name: "北区", CityCode: "01100", PrefectureCode: "01"},
			},
		},
		{
			name: "filter by dataset types",
			input: &plateauapi.AreaInput{
				DatasetTypes: []string{"bldg"},
				ParentCode:   nil,
				SearchTokens: nil,
			},
			want: []plateauapi.Area{
				&plateauapi.Ward{Code: "01101", Name: "中央区", CityCode: "01100", PrefectureCode: "01"},
			},
		},
		{
			name: "filter by prefectures",
			input: &plateauapi.AreaInput{
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("01")),
				DatasetTypes: nil,
				SearchTokens: nil,
			},
			want: []plateauapi.Area{
				&plateauapi.City{Code: "01100", Name: "札幌市", PrefectureCode: "01"},
			},
		},
		{
			name: "filter by cities",
			input: &plateauapi.AreaInput{
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("01100")),
				DatasetTypes: nil,
				SearchTokens: nil,
			},
			want: []plateauapi.Area{
				&plateauapi.Ward{Code: "01101", Name: "中央区", CityCode: "01100", PrefectureCode: "01"},
				&plateauapi.Ward{Code: "01102", Name: "北区", CityCode: "01100", PrefectureCode: "01"},
			},
		},
		{
			name: "filter by search tokens",
			input: &plateauapi.AreaInput{
				ParentCode:   nil,
				DatasetTypes: nil,
				SearchTokens: []string{"弘前"},
			},
			want: []plateauapi.Area{
				&plateauapi.City{Code: "02101", Name: "弘前市", PrefectureCode: "02"},
			},
		},
		{
			name: "filter by search tokens and dataset types",
			input: &plateauapi.AreaInput{
				ParentCode:   nil,
				DatasetTypes: []string{"bldg"},
				SearchTokens: []string{"弘前"},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.Areas(context.Background(), tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAdapter_DatasetTypes(t *testing.T) {
	a := &Adapter{
		plateauDatasetTypes: []plateauapi.PlateauDatasetType{
			{ID: "1", Name: "Plateau Dataset 1", Year: 2022, PlateauSpecName: "第2.3版"},
			{ID: "2", Name: "Plateau Dataset 2", Year: 2022, PlateauSpecName: "第2.3版"},
			{ID: "3", Name: "Plateau Dataset 3", Year: 2023, PlateauSpecName: "第3.0版"},
		},
		relatedDatasetTypes: []plateauapi.RelatedDatasetType{
			{ID: "4", Name: "Related Dataset 1"},
			{ID: "5", Name: "Related Dataset 2"},
			{ID: "6", Name: "Related Dataset 3"},
		},
		genericDatasetTypes: []plateauapi.GenericDatasetType{
			{ID: "7", Name: "Generic Dataset 1"},
			{ID: "8", Name: "Generic Dataset 2"},
			{ID: "9", Name: "Generic Dataset 3"},
		},
	}

	tests := []struct {
		name     string
		input    *plateauapi.DatasetTypeInput
		expected []plateauapi.DatasetType
	}{
		{
			name:  "no filter",
			input: nil,
			expected: []plateauapi.DatasetType{
				&plateauapi.PlateauDatasetType{ID: "1", Name: "Plateau Dataset 1", Year: 2022, PlateauSpecName: "第2.3版"},
				&plateauapi.PlateauDatasetType{ID: "2", Name: "Plateau Dataset 2", Year: 2022, PlateauSpecName: "第2.3版"},
				&plateauapi.PlateauDatasetType{ID: "3", Name: "Plateau Dataset 3", Year: 2023, PlateauSpecName: "第3.0版"},
				&plateauapi.RelatedDatasetType{ID: "4", Name: "Related Dataset 1"},
				&plateauapi.RelatedDatasetType{ID: "5", Name: "Related Dataset 2"},
				&plateauapi.RelatedDatasetType{ID: "6", Name: "Related Dataset 3"},
				&plateauapi.GenericDatasetType{ID: "7", Name: "Generic Dataset 1"},
				&plateauapi.GenericDatasetType{ID: "8", Name: "Generic Dataset 2"},
				&plateauapi.GenericDatasetType{ID: "9", Name: "Generic Dataset 3"},
			},
		},
		{
			name: "filter by year",
			input: &plateauapi.DatasetTypeInput{
				Year: lo.ToPtr(2022),
			},
			expected: []plateauapi.DatasetType{
				&plateauapi.PlateauDatasetType{ID: "1", Name: "Plateau Dataset 1", Year: 2022, PlateauSpecName: "第2.3版"},
				&plateauapi.PlateauDatasetType{ID: "2", Name: "Plateau Dataset 2", Year: 2022, PlateauSpecName: "第2.3版"},
			},
		},
		{
			name: "filter by spec",
			input: &plateauapi.DatasetTypeInput{
				PlateauSpec: lo.ToPtr("2.3"),
			},
			expected: []plateauapi.DatasetType{
				&plateauapi.PlateauDatasetType{ID: "1", Name: "Plateau Dataset 1", Year: 2022, PlateauSpecName: "第2.3版"},
				&plateauapi.PlateauDatasetType{ID: "2", Name: "Plateau Dataset 2", Year: 2022, PlateauSpecName: "第2.3版"},
			},
		},
		{
			name: "filter by category",
			input: &plateauapi.DatasetTypeInput{
				Category: lo.ToPtr(plateauapi.DatasetTypeCategoryGeneric),
			},
			expected: []plateauapi.DatasetType{
				&plateauapi.GenericDatasetType{ID: "7", Name: "Generic Dataset 1"},
				&plateauapi.GenericDatasetType{ID: "8", Name: "Generic Dataset 2"},
				&plateauapi.GenericDatasetType{ID: "9", Name: "Generic Dataset 3"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.DatasetTypes(context.Background(), tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestAdapter_Years(t *testing.T) {
	a := &Adapter{
		years: []int{2020, 2021, 2022, 2023},
	}

	expected := []int{2020, 2021, 2022, 2023}

	years, err := a.Years(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expected, years)
}

func TestAdapter_Datasets(t *testing.T) {
	a := &Adapter{
		plateauDatasets: []plateauapi.PlateauDataset{
			{ID: "1", Name: "Plateau Dataset 1", Year: 2022, CityCode: lo.ToPtr(plateauapi.AreaCode("01100"))},
			{ID: "2", Name: "Plateau Dataset 2", Year: 2022, TypeCode: "bldg"},
			{ID: "3", Name: "Plateau Dataset 3", Year: 2023},
		},
		relatedDatasets: []plateauapi.RelatedDataset{
			{ID: "7", Name: "Related Dataset 1", Year: 2022, Description: lo.ToPtr("desc!")},
			{ID: "8", Name: "Related Dataset 2", Year: 2022},
			{ID: "9", Name: "Related Dataset 3", Year: 2023},
		},
		genericDatasets: []plateauapi.GenericDataset{
			{ID: "10", Name: "Generic Dataset 1", Year: 2022},
			{ID: "11", Name: "Generic Dataset 2", Year: 2022},
			{ID: "12", Name: "Generic Dataset 3", Year: 2023, CityCode: lo.ToPtr(plateauapi.AreaCode("01100")), WardCode: lo.ToPtr(plateauapi.AreaCode("01101"))},
		},
	}

	tests := []struct {
		name  string
		input *plateauapi.DatasetInput
		want  []plateauapi.Dataset
	}{
		{
			name:  "no filter",
			input: nil,
			want: []plateauapi.Dataset{
				&plateauapi.PlateauDataset{ID: "1", Name: "Plateau Dataset 1", Year: 2022, CityCode: lo.ToPtr(plateauapi.AreaCode("01100"))},
				&plateauapi.PlateauDataset{ID: "2", Name: "Plateau Dataset 2", Year: 2022, TypeCode: "bldg"},
				&plateauapi.PlateauDataset{ID: "3", Name: "Plateau Dataset 3", Year: 2023},
				&plateauapi.RelatedDataset{ID: "7", Name: "Related Dataset 1", Year: 2022, Description: lo.ToPtr("desc!")},
				&plateauapi.RelatedDataset{ID: "8", Name: "Related Dataset 2", Year: 2022},
				&plateauapi.RelatedDataset{ID: "9", Name: "Related Dataset 3", Year: 2023},
				&plateauapi.GenericDataset{ID: "10", Name: "Generic Dataset 1", Year: 2022},
				&plateauapi.GenericDataset{ID: "11", Name: "Generic Dataset 2", Year: 2022},
				&plateauapi.GenericDataset{ID: "12", Name: "Generic Dataset 3", Year: 2023, CityCode: lo.ToPtr(plateauapi.AreaCode("01100")), WardCode: lo.ToPtr(plateauapi.AreaCode("01101"))},
			},
		},
		{
			name: "filter by an area code",
			input: &plateauapi.DatasetInput{
				AreaCodes: []plateauapi.AreaCode{"01100"},
			},
			want: []plateauapi.Dataset{
				&plateauapi.PlateauDataset{ID: "1", Name: "Plateau Dataset 1", Year: 2022, CityCode: lo.ToPtr(plateauapi.AreaCode("01100"))},
				&plateauapi.GenericDataset{ID: "12", Name: "Generic Dataset 3", Year: 2023, CityCode: lo.ToPtr(plateauapi.AreaCode("01100")), WardCode: lo.ToPtr(plateauapi.AreaCode("01101"))},
			},
		},
		{
			name: "filter by multiple area codes",
			input: &plateauapi.DatasetInput{
				AreaCodes: []plateauapi.AreaCode{"01100", "01101"},
			},
			want: []plateauapi.Dataset{
				&plateauapi.PlateauDataset{ID: "1", Name: "Plateau Dataset 1", Year: 2022, CityCode: lo.ToPtr(plateauapi.AreaCode("01100"))},
				&plateauapi.GenericDataset{ID: "12", Name: "Generic Dataset 3", Year: 2023, CityCode: lo.ToPtr(plateauapi.AreaCode("01100")), WardCode: lo.ToPtr(plateauapi.AreaCode("01101"))},
			},
		},
		{
			name: "filter by an area code depply",
			input: &plateauapi.DatasetInput{
				AreaCodes: []plateauapi.AreaCode{"01100"},
				Shallow:   lo.ToPtr(true),
			},
			want: []plateauapi.Dataset{
				&plateauapi.PlateauDataset{ID: "1", Name: "Plateau Dataset 1", Year: 2022, CityCode: lo.ToPtr(plateauapi.AreaCode("01100"))},
			},
		},
		{
			name: "filter by included types",
			input: &plateauapi.DatasetInput{
				IncludeTypes: []string{"bldg"},
			},
			want: []plateauapi.Dataset{
				&plateauapi.PlateauDataset{ID: "2", Name: "Plateau Dataset 2", Year: 2022, TypeCode: "bldg"},
			},
		},
		{
			name: "filter by excluded types",
			input: &plateauapi.DatasetInput{
				ExcludeTypes: []string{"bldg"},
			},
			want: []plateauapi.Dataset{
				&plateauapi.PlateauDataset{ID: "1", Name: "Plateau Dataset 1", Year: 2022, CityCode: lo.ToPtr(plateauapi.AreaCode("01100"))},
				&plateauapi.PlateauDataset{ID: "3", Name: "Plateau Dataset 3", Year: 2023},
				&plateauapi.RelatedDataset{ID: "7", Name: "Related Dataset 1", Year: 2022, Description: lo.ToPtr("desc!")},
				&plateauapi.RelatedDataset{ID: "8", Name: "Related Dataset 2", Year: 2022},
				&plateauapi.RelatedDataset{ID: "9", Name: "Related Dataset 3", Year: 2023},
				&plateauapi.GenericDataset{ID: "10", Name: "Generic Dataset 1", Year: 2022},
				&plateauapi.GenericDataset{ID: "11", Name: "Generic Dataset 2", Year: 2022},
				&plateauapi.GenericDataset{ID: "12", Name: "Generic Dataset 3", Year: 2023, CityCode: lo.ToPtr(plateauapi.AreaCode("01100")), WardCode: lo.ToPtr(plateauapi.AreaCode("01101"))},
			},
		},
		{
			name: "filter by search tokens",
			input: &plateauapi.DatasetInput{
				SearchTokens: []string{"desc"},
			},
			want: []plateauapi.Dataset{
				&plateauapi.RelatedDataset{ID: "7", Name: "Related Dataset 1", Year: 2022, Description: lo.ToPtr("desc!")},
			},
		},
		{
			name: "filter by multiple search tokens",
			input: &plateauapi.DatasetInput{
				SearchTokens: []string{"desc", "Related"},
			},
			want: []plateauapi.Dataset{
				&plateauapi.RelatedDataset{ID: "7", Name: "Related Dataset 1", Year: 2022, Description: lo.ToPtr("desc!")},
			},
		},
		{
			name: "filter by non-matched multiple search tokens",
			input: &plateauapi.DatasetInput{
				SearchTokens: []string{"desc", "Related_"},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.Datasets(context.Background(), tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAdapter_Node(t *testing.T) {
	a := &Adapter{
		prefectures: []plateauapi.Prefecture{
			{ID: plateauapi.NewID("01", plateauapi.TypeArea), Name: "北海道"},
			{ID: plateauapi.NewID("02", plateauapi.TypeArea), Name: "青森県"},
		},
		cities: []plateauapi.City{
			{ID: plateauapi.NewID("01100", plateauapi.TypeArea), Name: "札幌市", PrefectureCode: "01"},
		},
		wards: []plateauapi.Ward{
			{ID: plateauapi.NewID("01101", plateauapi.TypeArea), Name: "中央区", CityCode: "01100", PrefectureCode: "01"},
		},
		plateauDatasetTypes: []plateauapi.PlateauDatasetType{
			{ID: plateauapi.NewID("1", plateauapi.TypeDatasetType), Name: "Plateau Dataset 1", Year: 2022},
		},
		relatedDatasetTypes: []plateauapi.RelatedDatasetType{
			{ID: plateauapi.NewID("2", plateauapi.TypeDatasetType), Name: "Related Dataset 1"},
		},
		genericDatasetTypes: []plateauapi.GenericDatasetType{
			{ID: plateauapi.NewID("3", plateauapi.TypeDatasetType), Name: "Generic Dataset 1"},
		},
		plateauDatasets: []plateauapi.PlateauDataset{
			{ID: plateauapi.NewID("1", plateauapi.TypeDataset), Name: "Plateau Dataset 1"},
		},
		relatedDatasets: []plateauapi.RelatedDataset{
			{ID: plateauapi.NewID("3", plateauapi.TypeDataset), Name: "Related Dataset 1"},
		},
		genericDatasets: []plateauapi.GenericDataset{
			{ID: plateauapi.NewID("4", plateauapi.TypeDataset), Name: "Generic Dataset 1"},
		},
	}

	tests := []struct {
		name     string
		id       plateauapi.ID
		expected plateauapi.Node
	}{
		{
			name:     "invalid id",
			id:       plateauapi.NewID("99", plateauapi.TypeArea),
			expected: nil,
		},
		{
			name:     "prefecture",
			id:       plateauapi.NewID("01", plateauapi.TypeArea),
			expected: &plateauapi.Prefecture{ID: plateauapi.NewID("01", plateauapi.TypeArea), Name: "北海道"},
		},
		{
			name:     "city",
			id:       plateauapi.NewID("01100", plateauapi.TypeArea),
			expected: &plateauapi.City{ID: plateauapi.NewID("01100", plateauapi.TypeArea), Name: "札幌市", PrefectureCode: "01"},
		},
		{
			name:     "ward",
			id:       plateauapi.NewID("01101", plateauapi.TypeArea),
			expected: &plateauapi.Ward{ID: plateauapi.NewID("01101", plateauapi.TypeArea), Name: "中央区", CityCode: "01100", PrefectureCode: "01"},
		},
		{
			name:     "plateau dataset type",
			id:       plateauapi.NewID("1", plateauapi.TypeDatasetType),
			expected: &plateauapi.PlateauDatasetType{ID: plateauapi.NewID("1", plateauapi.TypeDatasetType), Name: "Plateau Dataset 1", Year: 2022},
		},
		{
			name:     "related dataset type",
			id:       plateauapi.NewID("2", plateauapi.TypeDatasetType),
			expected: &plateauapi.RelatedDatasetType{ID: plateauapi.NewID("2", plateauapi.TypeDatasetType), Name: "Related Dataset 1"},
		},
		{
			name:     "generic dataset type",
			id:       plateauapi.NewID("3", plateauapi.TypeDatasetType),
			expected: &plateauapi.GenericDatasetType{ID: plateauapi.NewID("3", plateauapi.TypeDatasetType), Name: "Generic Dataset 1"},
		},
		{
			name:     "plateau dataset",
			id:       plateauapi.NewID("1", plateauapi.TypeDataset),
			expected: &plateauapi.PlateauDataset{ID: plateauapi.NewID("1", plateauapi.TypeDataset), Name: "Plateau Dataset 1"},
		},
		{
			name:     "related dataset",
			id:       plateauapi.NewID("3", plateauapi.TypeDataset),
			expected: &plateauapi.RelatedDataset{ID: plateauapi.NewID("3", plateauapi.TypeDataset), Name: "Related Dataset 1"},
		},
		{
			name:     "generic dataset",
			id:       plateauapi.NewID("4", plateauapi.TypeDataset),
			expected: &plateauapi.GenericDataset{ID: plateauapi.NewID("4", plateauapi.TypeDataset), Name: "Generic Dataset 1"},
		},
		{
			name:     "spec",
			id:       plateauapi.NewID("2", plateauapi.TypePlateauSpec),
			expected: plateauSpecs[0],
		},
		{
			name:     "spec minor",
			id:       plateauapi.NewID("2.3", plateauapi.TypePlateauSpec),
			expected: plateauSpecs[0].MinorVersions[0],
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.Node(context.Background(), tt.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestAdapter_Nodes(t *testing.T) {
	a := &Adapter{
		prefectures: []plateauapi.Prefecture{
			{ID: plateauapi.NewID("01", plateauapi.TypeArea), Name: "北海道"},
			{ID: plateauapi.NewID("02", plateauapi.TypeArea), Name: "青森県"},
		},
	}

	tests := []struct {
		name     string
		ids      []plateauapi.ID
		expected []plateauapi.Node
	}{
		{
			name:     "empty ids",
			ids:      []plateauapi.ID{},
			expected: []plateauapi.Node{},
		},
		{
			name: "single id",
			ids:  []plateauapi.ID{plateauapi.NewID("01", plateauapi.TypeArea)},
			expected: []plateauapi.Node{
				&plateauapi.Prefecture{ID: plateauapi.NewID("01", plateauapi.TypeArea), Name: "北海道"},
			},
		},
		{
			name: "multiple ids",
			ids: []plateauapi.ID{
				plateauapi.NewID("01", plateauapi.TypeArea),
				plateauapi.NewID("02", plateauapi.TypeArea),
			},
			expected: []plateauapi.Node{
				&plateauapi.Prefecture{ID: plateauapi.NewID("01", plateauapi.TypeArea), Name: "北海道"},
				&plateauapi.Prefecture{ID: plateauapi.NewID("02", plateauapi.TypeArea), Name: "青森県"},
			},
		},
		{
			name: "multiple ids with an invalid id",
			ids: []plateauapi.ID{
				plateauapi.NewID("99", plateauapi.TypeArea),
				plateauapi.NewID("02", plateauapi.TypeArea),
			},
			expected: []plateauapi.Node{
				nil,
				&plateauapi.Prefecture{ID: plateauapi.NewID("02", plateauapi.TypeArea), Name: "青森県"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.Nodes(context.Background(), tt.ids)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}
