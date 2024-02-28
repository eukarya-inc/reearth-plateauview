package preparegspatialjp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenameCityGMLZip(t *testing.T) {
	tests := []struct {
		ty            string
		prefix        string
		path          string
		expected      string
		expectedError bool
	}{
		{
			ty:            "bldg",
			path:          "udx/bldg/file.txt",
			expectedError: true,
		},
		{
			ty:            "bldg",
			prefix:        "udx/",
			path:          "udx/bldg/file.txt",
			expected:      "bldg/file.txt",
			expectedError: false,
		},
		{
			ty:            "bldg",
			prefix:        "udx/",
			path:          "udx/bldg/bldg/file.txt",
			expected:      "bldg/file.txt",
			expectedError: false,
		},
		{
			ty:            "bldg",
			prefix:        "udx/",
			path:          "udx/xxx_bldg/file.txt",
			expected:      "bldg/file.txt",
			expectedError: false,
		},
		{
			ty:            "bldg",
			path:          "file.gml",
			expectedError: true,
		},
		{
			ty:            "bldg",
			path:          "bldg/file.txt",
			expectedError: false,
		},
		{
			ty:            "bldg",
			path:          "xxx_bldg/file.txt",
			expected:      "bldg/file.txt",
			expectedError: false,
		},
		{
			ty:            "bldg",
			path:          "xxx_bldg/bldg/file.txt",
			expected:      "bldg/file.txt",
			expectedError: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.path, func(t *testing.T) {
			fn, err := renameCityGMLZip(test.ty, test.prefix)(test.path)
			assert.Equal(t, test.expected, fn)
			if !test.expectedError {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
