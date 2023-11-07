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
				Key:   "conv_status",
				Type:  "select",
				Value: string(ConvertionStatusError),
			},
		},
	}

	expected := &FeatureItem{
		ID:               "id",
		ConvertionStatus: ConvertionStatusError,
	}

	expected2 := &cms.Item{
		ID: "id",
		MetadataFields: []*cms.Field{
			{
				Key:   "conv_status",
				Type:  "select",
				Value: ConvertionStatusError,
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
				Key:   "status",
				Type:  "select",
				Value: string(ManagementStatusDone),
			},
		},
	}

	expected := &GenericItem{
		ID:     "id",
		Status: ManagementStatusDone,
	}

	expected2 := &cms.Item{
		ID: "id",
		MetadataFields: []*cms.Field{
			{
				Key:   "status",
				Type:  "select",
				Value: ManagementStatusDone,
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
				Value: "PARK",
			},
			{
				Key:   "park_conv",
				Type:  "asset",
				Value: "PARK_CONV",
			},
			{
				Key:   "landmark",
				Type:  "asset",
				Value: "LANDMARK",
			},
		},
		MetadataFields: []*cms.Field{
			{
				Key:   "park_status",
				Type:  "select",
				Value: string(ConvertionStatusSuccess),
			},
			{
				Key:   "park_public",
				Type:  "bool",
				Value: true,
			},
			{
				Key:   "border_public",
				Type:  "bool",
				Value: true,
			},
		},
	}

	expected := &RelatedItem{
		ID: "id",
		Assets: map[string]string{
			"park":     "PARK",
			"landmark": "LANDMARK",
		},
		ConvertedAssets: map[string]string{
			"park": "PARK_CONV",
		},
		Public: map[string]bool{
			"park":   true,
			"border": true,
		},
		ConvertStatus: map[string]ConvertionStatus{
			"park": ConvertionStatusSuccess,
		},
	}

	relatedItem := RelatedItemFrom(item)
	assert.Equal(t, expected, relatedItem)
	item2 := relatedItem.CMSItem()
	assert.Equal(t, item, item2)
}
