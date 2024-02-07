package geospatialjpv3

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const jobName = "projects/xxxxx/locations/xxxxx/jobs/plateauview-api-worker"
const itemID = ""
const projectID = ""

func TestPrepare_RequestZip(t *testing.T) {
	t.Skip()

	ctx := context.Background()
	err := Prepare(ctx, itemID, projectID, jobName)
	assert.NoError(t, err)
}
