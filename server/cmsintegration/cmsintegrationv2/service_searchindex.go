package cmsintegrationv2

import (
	"context"

	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
)

func buildSearchIndex(ctx context.Context, s *Services, w *cmswebhook.Payload) error {
	log.Infofc(ctx, "cmsintegrationv2: buildSearchIndex")
	// TODO

	return nil
}
