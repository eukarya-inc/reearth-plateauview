package cmsintegrationv3

import (
	"testing"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/stretchr/testify/assert"
)

func TestCityItemFrom(t *testing.T) {
	item := &cms.Item{
		ID: "id",
		Fields: []*cms.Field{
			{
				Key:   "bldg",
				Type:  "reference",
				Value: "BLDG",
			},
		},
		MetadataFields: []*cms.Field{
			{
				Key:   "city_public",
				Type:  "bool",
				Value: true,
			},
			{
				Key:   "bldg_public",
				Type:  "bool",
				Value: true,
			},
		},
	}

	expected := &CityItem{
		ID: "id",
		References: map[string]string{
			"bldg": "BLDG",
		},
		Public: map[string]bool{
			"bldg": true,
		},
		CityPublic: true,
	}

	cityItem := CityItemFrom(item)
	assert.Equal(t, expected, cityItem)
	item2 := cityItem.CMSItem()
	assert.Equal(t, item, item2)
}

func TestFeatureItemFrom(t *testing.T) {
	item := &cms.Item{
		ID: "id",
		MetadataFields: []*cms.Field{
			{
				Key:  "conv_status",
				Type: "tag",
				Value: map[string]any{
					"id":   "xxx",
					"name": string(ConvertionStatusError),
				},
			},
		},
	}

	expected := &FeatureItem{
		ID: "id",
		ConvertionStatus: &cms.Tag{
			ID:   "xxx",
			Name: string(ConvertionStatusError),
		},
	}

	expected2 := &cms.Item{
		ID: "id",
		MetadataFields: []*cms.Field{
			{
				Key:   "conv_status",
				Type:  "tag",
				Value: "xxx",
			},
		},
	}

	featureItem := FeatureItemFrom(item)
	assert.Equal(t, expected, featureItem)
	item2 := featureItem.CMSItem()
	assert.Equal(t, expected2, item2)
}

func TestGenericItemFrom(t *testing.T) {
	item := &cms.Item{
		ID: "id",
		MetadataFields: []*cms.Field{
			{
				Key:   "public",
				Type:  "bool",
				Value: true,
			},
		},
	}

	expected := &GenericItem{
		ID:     "id",
		Public: true,
	}

	expected2 := &cms.Item{
		ID: "id",
		MetadataFields: []*cms.Field{
			{
				Key:   "public",
				Type:  "bool",
				Value: true,
			},
		},
	}

	genericItem := GenericItemFrom(item)
	assert.Equal(t, expected, genericItem)
	item2 := genericItem.CMSItem()
	assert.Equal(t, expected2, item2)
}

func TestRelatedItemFrom(t *testing.T) {
	item := &cms.Item{
		ID: "id",
		Fields: []*cms.Field{
			{
				Key:   "park",
				Type:  "asset",
				Value: []string{"PARK"},
			},
			{
				Key:   "park_conv",
				Type:  "asset",
				Value: []string{"PARK_CONV"},
			},
			{
				Key:   "landmark",
				Type:  "asset",
				Value: []string{"LANDMARK"},
			},
		},
		MetadataFields: []*cms.Field{
			{
				Key:   "conv_status",
				Type:  "tag",
				Value: map[string]any{"id": "xxx", "name": string(ConvertionStatusSuccess)},
			},
			{
				Key:   "public",
				Type:  "bool",
				Value: true,
			},
		},
	}

	expected := &RelatedItem{
		ID: "id",
		Assets: map[string][]string{
			"park":     {"PARK"},
			"landmark": {"LANDMARK"},
		},
		ConvertedAssets: map[string][]string{
			"park": {"PARK_CONV"},
		},
		ConvertStatus: &cms.Tag{
			ID:   "xxx",
			Name: string(ConvertionStatusSuccess),
		},
		Public: true,
	}

	expected2 := &cms.Item{
		ID: "id",
		Fields: []*cms.Field{
			{
				Key:   "park",
				Type:  "asset",
				Value: []string{"PARK"},
			},
			{
				Key:   "park_conv",
				Type:  "asset",
				Value: []string{"PARK_CONV"},
			},
			{
				Key:   "landmark",
				Type:  "asset",
				Value: []string{"LANDMARK"},
			},
		},
		MetadataFields: []*cms.Field{
			{
				Key:   "conv_status",
				Type:  "tag",
				Value: "xxx",
			},
			{
				Key:   "public",
				Type:  "bool",
				Value: true,
			},
		},
	}

	relatedItem := RelatedItemFrom(item)
	assert.Equal(t, expected, relatedItem)
	item2 := relatedItem.CMSItem()
	assert.Equal(t, expected2, item2)
}
