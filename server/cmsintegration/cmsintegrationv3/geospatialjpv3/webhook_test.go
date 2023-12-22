package geospatialjpv3

import (
	"testing"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/ckan"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsmock"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Webhook(t *testing.T) {
	packages := []ckan.Package{}
	resources := []ckan.Resource{}

	cmsmock := &cmsmock.CMSMock{}
	ckanmock := ckan.NewMock("test", packages, resources)

	h := &handler{
		cms:  cmsmock,
		ckan: ckanmock,
	}

	wh, err := h.Webhook(Config{
		// TODO
	})
	assert.NoError(t, err)

	payload := &cmswebhook.Payload{
		// TODO
	}

	err = wh(nil, payload)
	assert.NoError(t, err)

	// TODO: assert ckan packages and resources
}
