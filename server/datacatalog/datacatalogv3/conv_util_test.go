package datacatalogv3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NamesWithoutExtFromUrls(t *testing.T) {
	urls := []string{"https://example.com/foo/bar.json", "https://example.com/foo/baz"}
	expected := []string{"bar", "baz"}
	actual := namesWithoutExtFromUrls(urls)
	assert.Equal(t, expected, actual)
}
