package geospatialjpv3

import (
	"context"
	"testing"

	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
)

func TestPrepare_RequestZip(t *testing.T) {
	// t.Skip()

	ctx := context.Background()

	c := PrepareConfig{
		gcpProjectID: "reearth-plateau-dev",
		gcpLocation:  "asia-northeast1",
	}

	w := &cmswebhook.Payload{
		ItemData: &cmswebhook.ItemData{
			Item:   &cms.Item{ID: ""},
			Schema: &cms.Schema{ProjectID: ""},
		},
	}

	log.Debugfc(ctx, "geospatialjp webhook: RequestPreparing")

	_ = c.RequestPreparing(ctx, w)

}
