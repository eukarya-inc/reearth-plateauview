package geospatialjpv3

import (
	"context"
	"os"
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/ckan"
	"github.com/joho/godotenv"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	t.Skip()
	lo.Must0(godotenv.Load("../../../.env"))

	var (
		itemID      = ""
		projectID   = ""
		cmsURL      = os.Getenv("REEARTH_PLATEAUVIEW_CMS_BASEURL")
		cmsToken    = os.Getenv("REEARTH_PLATEAUVIEW_CMS_TOKEN")
		ckanOrg     = ""
		ckanBaseURL = ""
		ckanToken   = ""
	)

	ctx := context.Background()

	w := &cmswebhook.Payload{
		ItemData: &cmswebhook.ItemData{
			Item:   &cms.Item{ID: itemID},
			Schema: &cms.Schema{ProjectID: projectID},
		},
	}

	ckan, err := ckan.New(ckanBaseURL, ckanToken)
	assert.NoError(t, err)

	cms, err := cms.New(cmsURL, cmsToken)
	assert.NoError(t, err)

	h := &handler{
		cms:     cms,
		ckan:    ckan,
		ckanOrg: ckanOrg,
	}

	err = h.Publish(ctx, w)
	assert.NoError(t, err)

}
