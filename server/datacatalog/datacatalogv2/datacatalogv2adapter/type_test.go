package datacatalogv2adapter

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/datacatalogv2"
	"github.com/stretchr/testify/assert"
)

func TestGetSubName(t *testing.T) {
	assert.Equal(t, "名前", getSubName(datacatalogv2.DataCatalogItem{
		Type2: "名前",
	}))
	assert.Empty(t, getSubName(datacatalogv2.DataCatalogItem{
		Name: "名前（○○市）",
		Type: "名前",
	}))
	assert.Equal(t, "AAAA BBBB", getSubName(datacatalogv2.DataCatalogItem{
		Name: "名前 AAAA BBBB（○○市）",
		Type: "名前",
	}))
	assert.Equal(t, "名前 AAAA BBBB", getSubName(datacatalogv2.DataCatalogItem{
		Name: "名前 AAAA BBBB（○○市）",
	}))
}

func TestGetSubCode(t *testing.T) {
	assert.Equal(t, "name", getSubCode(datacatalogv2.DataCatalogItem{
		Type2En: "name",
	}))
	assert.Empty(t, getSubCode(datacatalogv2.DataCatalogItem{
		ID: "11111_hoge-shi_bldg",
	}))
	assert.Equal(t, "aaaa_bbbb", getSubCode(datacatalogv2.DataCatalogItem{
		ID: "11111_hoge-shi_urf_aaaa_bbbb",
	}))
	assert.Equal(t, "aaaa_bbbb", getSubCode(datacatalogv2.DataCatalogItem{
		ID: "11111_hoge-shi_fld_aaaa_bbbb_l1",
	}))
	assert.Equal(t, "aaaa_bbbb", getSubCode(datacatalogv2.DataCatalogItem{
		ID: "11111_hoge-shi_fld_aaaa_bbbb_l2",
	}))
}
