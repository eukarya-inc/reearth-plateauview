package preparegspatialjp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckCityGML(t *testing.T) {
	tests := []struct {
		path          string
		expectedError bool
	}{
		{
			path:          "udx/file.txt",
			expectedError: true,
		},
		{
			path:          "file.gml",
			expectedError: true,
		},
		{
			path:          "dir/file.txt",
			expectedError: false,
		},
		{
			path:          "dir/file.gml",
			expectedError: false,
		},
	}

	for _, test := range tests {
		result := checkCityGMLZip(test.path)
		if !test.expectedError {
			assert.NoError(t, result)
		} else {
			assert.Error(t, result)
		}
	}
}
