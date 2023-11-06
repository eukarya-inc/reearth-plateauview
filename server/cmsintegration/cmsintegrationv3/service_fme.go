package cmsintegrationv3

import (
	"context"

	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
)

func sendRequestToFME(ctx context.Context, s *Services, w *cmswebhook.Payload) error {
	log.Infofc(ctx, "cmsintegrationv2: sendRequestToFME")
	// TODO

	return nil
}

func receiveResultFromFME(ctx context.Context, s *Services, f fmeResult) error {
	log.Infofc(ctx, "cmsintegrationv2: receiveResultFromFME")
	// TODO

	return nil
}
