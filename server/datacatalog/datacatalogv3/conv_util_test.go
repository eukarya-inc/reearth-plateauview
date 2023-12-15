package datacatalogv3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NameWithoutExt(t *testing.T) {
	assert.Equal(t, "foo", nameWithoutExt("foo"))
	assert.Equal(t, "bar", nameWithoutExt("bar.json"))
}

func Test_NameFromUrl(t *testing.T) {
	assert.Equal(t, "foo", nameFromUrl("https://example.com/foo"))
	assert.Equal(t, "bar.json", nameFromUrl("bar.json"))
}
