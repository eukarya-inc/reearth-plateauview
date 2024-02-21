package datacatalogv3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCitygmlItemURLFrom(t *testing.T) {
	assert.Equal(t, "https://example.com/citygml/udx/fld/hoge", citygmlItemURLFrom("https://example.com/citygml.zip", "hoge", "fld"))
}
