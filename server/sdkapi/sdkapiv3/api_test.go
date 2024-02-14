package sdkapiv3

import (
	"testing"

	"github.com/k0kubun/pp/v3"
	"github.com/stretchr/testify/assert"
)

func TestQueryDatasets(t *testing.T) {
	baseURL := ""

	if baseURL == "" {
		t.Skip("baseURL is not set")
	}

	client, err := NewClient(Config{BaseURL: baseURL})
	assert.NoError(t, err)

	q, err := client.QueryDatasets()
	assert.NoError(t, err)

	t.Log(ppp.Sprint(q))
}

var ppp = pp.New()

func init() {
	ppp.SetColoringEnabled(false)
}
