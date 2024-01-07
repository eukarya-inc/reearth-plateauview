package datacatalogv3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NameWithoutExt(t *testing.T) {
	assert.Equal(t, "foo", nameWithoutExt("foo"))
	assert.Equal(t, "bar", nameWithoutExt("bar.json"))
}

func Test_NameFromURL(t *testing.T) {
	assert.Equal(t, "foo", nameFromURL("https://example.com/foo"))
	assert.Equal(t, "bar.json", nameFromURL("bar.json"))
}
