package datacatalogv2adapter

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/datacatalog/plateauapi"
	"github.com/stretchr/testify/assert"
)

func TestCache_GetDatasetTypes(t *testing.T) {
	c := &cache{
		plateauDatasetTypes: []plateauapi.PlateauDatasetType{
			{
				Code: "p1",
			},
			{
				Code: "p2",
			},
		},
		relatedDatasetTypes: []plateauapi.RelatedDatasetType{
			{
				Code: "r",
			},
		},
		genericDatasetTypes: []plateauapi.GenericDatasetType{
			{
				Code: "g",
			},
		},
	}

	assert.Equal(t, []string{"p1", "p2", "r", "g"}, c.getDatasetTypes(nil, nil))
	assert.Equal(t, []string{"p1", "g"}, c.getDatasetTypes([]string{"p1", "g"}, nil))
	assert.Equal(t, []string{"p1", "p2"}, c.getDatasetTypes(nil, []plateauapi.DatasetTypeCategory{
		plateauapi.DatasetTypeCategoryPlateau,
	}))
	assert.Equal(t, []string{"p2"}, c.getDatasetTypes([]string{"p2"}, []plateauapi.DatasetTypeCategory{
		plateauapi.DatasetTypeCategoryPlateau,
	}))
}
