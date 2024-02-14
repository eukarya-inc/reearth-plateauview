package sdkapiv3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryDatasets(t *testing.T) {
	t.Skip()

	client, err := NewClient(Config{BaseURL: ""})
	assert.NoError(t, err)

	q, err := client.QueryDatasets()
	assert.NoError(t, err)

	// Add your assertions here to validate the query response
	// For example:
	assert.NotEmpty(t, q.Areas[0].ID)
	assert.NotEmpty(t, q.Areas[0].Name)
	assert.NotEmpty(t, q.Areas[0].Prefecture.Cities[0].ID)
	assert.NotEmpty(t, q.Areas[0].Prefecture.Cities[0].Name)
	assert.NotEmpty(t, q.Areas[0].Prefecture.Cities[0].Datasets[0].ID)
	assert.NotEmpty(t, q.Areas[0].Prefecture.Cities[0].Datasets[0].Name)
	assert.NotEmpty(t, q.Areas[0].Prefecture.Cities[0].Datasets[0].TypeCode)
	assert.NotEmpty(t, q.PlateauSpecs[0].MajorVersion)
	assert.NotEmpty(t, q.PlateauSpecs[0].MinorVersions[0].Version)
}
