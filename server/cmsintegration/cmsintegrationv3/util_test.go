package cmsintegrationv3

import (
	"testing"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/stretchr/testify/assert"
)

func TestFindFieldChangeByKey(t *testing.T) {
	w := &cmswebhook.Payload{
		ItemData: &cmswebhook.ItemData{
			Item: &cms.Item{
				Fields: []*cms.Field{
					{
						ID:  "field1",
						Key: "key1",
					},
					{
						ID:  "field2",
						Key: "key2",
					},
				},
			},
			Changes: []cms.FieldChange{
				{
					ID:           "field1",
					CurrentValue: "Value 1",
				},
				{
					ID:           "field2",
					CurrentValue: "Value 2",
				},
			},
		},
	}

	value, found := FindFieldChangeByKey(w, "key1")
	assert.True(t, found)
	assert.Equal(t, "Value 1", value)

	value, found = FindFieldChangeByKey(w, "key2")
	assert.True(t, found)
	assert.Equal(t, "Value 2", value)

	value, found = FindFieldChangeByKey(w, "key3")
	assert.False(t, found)
	assert.Nil(t, value)

	value, found = FindFieldChangeByKey(nil, "key1")
	assert.False(t, found)
	assert.Nil(t, value)
}
