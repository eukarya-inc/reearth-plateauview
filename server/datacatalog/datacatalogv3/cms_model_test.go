package datacatalogv3

import (
	"testing"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/stretchr/testify/assert"
)

func TestPlateauFeatureItemFrom(t *testing.T) {
	item := &cms.Item{
		ID: "id",
		Fields: []*cms.Field{
			{
				Key: "data",
				Value: []any{
					map[string]any{"url": "url1"},
					map[string]any{"url": "url2"},
				},
			},
			{
				Key:   "maxlod",
				Value: map[string]any{"url": "url_maxlod"},
			},
			{
				Key:   "citygml",
				Value: map[string]any{"url": "url_citygml"},
			},
			{
				Key:   "items",
				Type:  "group",
				Value: []any{"item1", "item2"},
			},
			{
				Key:   "data",
				Group: "item1",
				Value: []any{map[string]any{"url": "url3"}},
			},
			{
				Key:   "data",
				Group: "item2",
				Value: []any{"url4"}, // string is ignored
			},
		},
	}

	expected := &PlateauFeatureItem{
		ID:      "id",
		Data:    []string{"url1", "url2"},
		CityGML: "url_citygml",
		MaxLOD:  "url_maxlod",
		Items: []PlateauFeatureItemDatum{
			{
				ID:   "item1",
				Data: []string{"url3"},
			},
			{
				ID: "item2",
			},
		},
	}

	assert.Equal(t, expected, PlateauFeatureItemFrom(item))
}

func TestGenericItemFrom(t *testing.T) {
	item := &cms.Item{
		ID: "id",
		Fields: []*cms.Field{
			{
				Key:   "data",
				Type:  "group",
				Value: []any{"item1"},
			},
			{
				Key:   "data",
				Group: "item1",
				Value: map[string]any{"url": "url1"},
			},
			{
				Key:   "desc",
				Group: "item1",
				Value: "desc1",
			},
		},
	}

	expected := &GenericItem{
		ID: "id",
		Data: []GenericItemDataset{
			{
				ID:   "item1",
				Data: "url1",
				Desc: "desc1",
			},
		},
	}

	assert.Equal(t, expected, GenericItemFrom(item))
}

func TestRelatedItemFrom(t *testing.T) {
	item := &cms.Item{
		ID: "id",
		Fields: []*cms.Field{
			{
				Key:   "hoge",
				Value: map[string]any{"url": "url1"},
			},
			{
				Key:   "foo_conv",
				Value: map[string]any{"url": "url2"},
			},
		},
	}

	ft := []FeatureType{
		{
			Code: "hoge",
		},
		{
			Code: "foo",
		},
	}

	expected := &RelatedItem{
		ID: "id",
		Assets: map[string][]string{
			"hoge": {"url1"},
		},
		ConvertedAssets: map[string][]string{
			"foo": {"url2"},
		},
	}

	assert.Equal(t, expected, RelatedItemFrom(item, ft))
}

func TestValueToAssetURLs(t *testing.T) {
	assert.Nil(t, valueToAssetURLs(cms.NewValeu("string")))
	assert.Nil(t, valueToAssetURLs(cms.NewValeu(map[string]string{"aaa": "bbb"})))
	assert.Equal(t, []string{"url"}, valueToAssetURLs(cms.NewValeu(map[string]any{"url": "url"})))
	assert.Equal(t, []string{"url"}, valueToAssetURLs(cms.NewValeu(map[any]any{"url": "url"})))
	assert.Equal(t, []string{"url", "url2"}, valueToAssetURLs(cms.NewValeu([]any{
		map[string]any{"url": "url"}, map[any]any{"url": "url2"}, map[string]any{},
	})))
}
