package preparegspatialjp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInvalidFile(t *testing.T) {
	tests := []struct {
		filePath     string
		expectedBool bool
	}{
		{
			filePath:     "udx/file.txt",
			expectedBool: true,
		},
		{
			filePath:     "file.gml",
			expectedBool: true,
		},
		{
			filePath:     "dir/file.txt",
			expectedBool: false,
		},
		{
			filePath:     "dir/file.gml",
			expectedBool: false,
		},
	}

	for _, test := range tests {
		result := isInvalid(test.filePath)
		assert.Equal(t, test.expectedBool, result)
	}
}
