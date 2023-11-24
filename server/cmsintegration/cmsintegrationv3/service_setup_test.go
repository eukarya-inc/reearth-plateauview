package cmsintegrationv3

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSetupCSV(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		name     string
		csvData  string
		expected []SetupCSVItem
		err      error
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
			err: nil,
		},
		{
			name:     "empty csv",
			csvData:  "",
			expected: nil,
			err:      io.EOF,
		},
		{
			name: "invalid header",
			csvData: `Name,NameEn,Code,Prefecture
Tokyo,東京,13,Tokyo`,
			expected: nil,
			err:      fmt.Errorf("invalid header: [Name NameEn Code Prefecture]"),
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
			items, err := parseSetupCSV(ctx, r)

			if tt.err != nil {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, items)
		})
	}
}
