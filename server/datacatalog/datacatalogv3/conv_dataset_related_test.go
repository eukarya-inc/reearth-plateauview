package datacatalogv3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRelatedAssetName(t *testing.T) {
	assert.Equal(t, &RelatedAssetName{
		Code: "13101",
		Name: "chiyoda-ku",
		Type: "shelter",
		Ext:  "geojson",
	}, ParseRelatedAssetName("13101_chiyoda-ku_shelter.geojson"))
	assert.Equal(t, &RelatedAssetName{
		Code: "13101",
		Name: "chiyoda-ku",
		Type: "border",
		Ext:  "czml",
	}, ParseRelatedAssetName("13101_chiyoda-ku_border.czml"))
	assert.Nil(t, ParseRelatedAssetName("invalid"))
}
