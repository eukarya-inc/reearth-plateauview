package plateauv2

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestDescFromAsset(t *testing.T) {
	a := AssetName{
		CityCode:       "000000",
		CityEn:         "hoge-shi",
		Year:           "2020",
		Format:         "mvt",
		Op:             "op",
		Feature:        "urf",
		UrfFeatureType: "urf",
	}

	desc := descFromAsset(a, []string{
		"bbb.zip\n\nBBB",
		"000000_hoge-shi_2020_mvt_op_urf_urf.zip\n\nAAA",
		"CCC",
	}, false)
	assert.Equal(t, Description{
		Desc: "AAA",
	}, desc)

	desc = descFromAsset(a, []string{
		"bbb.zip\n\nBBB",
		"aaa.zip\n\nAAA",
		"CCC",
	}, false)
	assert.Equal(t, Description{}, desc)

	desc = descFromAsset(a, []string{
		"CCC",
	}, false)
	assert.Equal(t, Description{}, desc)

	desc = descFromAsset(a, []string{
		"000000_hoge-shi_2020_mvt_op_urf_urf.zip\n@name: CCC\n\naaaa\nbbbb",
	}, false)
	assert.Equal(t, Description{
		Override: Override{
			Name: "CCC",
		},
		Desc: "aaaa\nbbbb",
	}, desc)

	desc = descFromAsset(a, []string{
		"000000_hoge-shi_2020_mvt_op_urf_urf.zip\n\n@name: CCC\naaaa\nbbbb",
	}, false)
	assert.Equal(t, Description{
		Override: Override{
			Name: "CCC",
		},
		Desc: "aaaa\nbbbb",
	}, desc)

	desc = descFromAsset(a, []string{
		"000000_hoge-shi_2020_mvt_op_urf_urf.zip\n@name:CCC\n@group:aaaa",
	}, false)
	assert.Equal(t, Description{
		Override: Override{
			Name:  "CCC",
			Group: "aaaa",
		},
	}, desc)

	desc = descFromAsset(a, []string{
		"@name: bbb\naaa",
		"@name: aaa\nbbb",
	}, true)
	assert.Equal(t, Description{
		Desc: "aaa",
		Override: Override{
			Name: "bbb",
		},
	}, desc)
}

func TestDescriptionFrom(t *testing.T) {
	assert.Equal(t, Description{
		Desc: "aaa",
		Override: Override{
			Name:         "bbb",
			DatasetOrder: lo.ToPtr(12),
		},
	}, DescriptionFrom("@name: bbb\n@datasetOrder: 12\n\naaa"))
}

func TestExtractTags(t *testing.T) {
	tags, rest := extractTags("\n\n@name: CCC\n@aaa: bbb\n\n@type: DDD\n\n@layer: aaa,bbb,ccc\n@order: 1\n\naaaa\nbbbb")
	assert.Equal(t, map[string]string{
		"name":  "CCC",
		"aaa":   "bbb",
		"type":  "DDD",
		"layer": "aaa,bbb,ccc",
		"order": "1",
	}, tags)
	assert.Equal(t, "aaaa\nbbbb", rest)

	tags, rest = extractTags("aaaa\nbbbb")
	assert.Equal(t, map[string]string{}, tags)
	assert.Equal(t, "aaaa\nbbbb", rest)
}
