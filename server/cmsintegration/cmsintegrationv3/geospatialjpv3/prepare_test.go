package geospatialjpv3

import (
	"context"
	"testing"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/stretchr/testify/assert"
)

const jobName = "projects/xxxxx/locations/xxxxx/jobs/plateauview-api-worker"
const itemID = ""
const projectID = ""

func TestPrepare_RequestZip(t *testing.T) {
	t.Skip()

	ctx := context.Background()
	w := &cmswebhook.Payload{
		ItemData: &cmswebhook.ItemData{
			Item:   &cms.Item{ID: itemID},
			Schema: &cms.Schema{ProjectID: projectID},
		},
	}

	err := Prepare(ctx, w, jobName)
	assert.NoError(t, err)
}
