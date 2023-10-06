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
		input    plateauapi.AreaQuery
		expected bool
	}{
		{
			name: "Prefecture with search tokens",
			area: plateauapi.Prefecture{Name: "Tokyo"},
			input: plateauapi.AreaQuery{
				SearchTokens: []string{"Tokyo"},
			},
			expected: true,
		},
		{
			name: "Prefecture without search tokens",
			area: plateauapi.Prefecture{Name: "Tokyo"},
			input: plateauapi.AreaQuery{
				SearchTokens: []string{},
			},
			expected: true,
		},
		{
			name: "Prefecture without non-matching search tokens",
			area: plateauapi.Prefecture{Name: "Tokyo"},
			input: plateauapi.AreaQuery{
				SearchTokens: []string{"Kanagawa"},
			},
			expected: false,
		},
		{
			name: "City with search tokens and matching parent code",
			area: plateauapi.City{Name: "Shinjuku", PrefectureCode: "13"},
			input: plateauapi.AreaQuery{
				SearchTokens: []string{"Shinjuku"},
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("13")),
			},
			expected: true,
		},
		{
			name: "City with search tokens and non-matching parent code",
			area: plateauapi.City{Name: "Shinjuku", PrefectureCode: "13"},
			input: plateauapi.AreaQuery{
				SearchTokens: []string{"Shinjuku"},
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("14")),
			},
			expected: false,
		},
		{
			name: "Ward with search tokens and matching parent code",
			area: plateauapi.Ward{Name: "Shinjuku", PrefectureCode: "13", CityCode: "13104"},
			input: plateauapi.AreaQuery{
				SearchTokens: []string{"Shinjuku"},
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("13104")),
			},
			expected: true,
		},
		{
			name: "Ward with search tokens and non-matching parent code",
			area: plateauapi.Ward{Name: "Shinjuku", PrefectureCode: "13", CityCode: "13104"},
			input: plateauapi.AreaQuery{
				SearchTokens: []string{"Shinjuku"},
				ParentCode:   lo.ToPtr(plateauapi.AreaCode("13105")),
			},
			expected: false,
		},
		{
			name: "Ward without search tokens",
			area: plateauapi.Ward{Name: "Shinjuku", PrefectureCode: "13", CityCode: "13104"},
			input: plateauapi.AreaQuery{
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
