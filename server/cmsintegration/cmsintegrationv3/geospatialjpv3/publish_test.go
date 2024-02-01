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
		indexItemID = ""
		dataItemID  = ""
		cmsURL      = os.Getenv("REEARTH_PLATEAUVIEW_CMS_BASEURL")
		cmsToken    = os.Getenv("REEARTH_PLATEAUVIEW_CMS_TOKEN")
		ckanOrg     = os.Getenv("REEARTH_PLATEAUVIEW_CKAN_ORG")
		ckanBaseURL = os.Getenv("REEARTH_PLATEAUVIEW_CKAN_BASEURL")
		ckanToken   = os.Getenv("REEARTH_PLATEAUVIEW_CKAN_TOKEN")
	)

	ctx := context.Background()

	w := &cmswebhook.Payload{
		ItemData: &cmswebhook.ItemData{
			Item: &cms.Item{
				// ID: itemID,
				Fields: []*cms.Field{
					{
						Key:   "city_code",
						Value: "99999",
					},
					{
						Key:   "city_name",
						Value: "テスト市",
					},
					{
						Key:   "city_name_en",
						Value: "test-shi",
					},
					{
						Key:   "year",
						Value: "2023",
					},
					{
						Key:   "specification",
						Value: "4",
					},
					{
						Key:   "geospatialjp-index",
						Value: indexItemID,
					},
					{
						Key:   "geospatialjp-data",
						Value: dataItemID,
					},
				},
			},
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

func TestShouldReorder(t *testing.T) {
	pkg := &ckan.Package{
		Resources: []ckan.Resource{
			{
				Name: "Resource 1 (v1)",
			},
			{
				Name: "Resource 2 (v2)",
			},
			{
				Name: "Resource 3 (v3)",
			},
		},
	}

	result := shouldReorder(pkg, 4)
	assert.True(t, result)

	result = shouldReorder(pkg, 2)
	assert.False(t, result)
}
