package cmsintegrationv2

import (
	"net/http"

	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
)

const (
	modelKey = "plateau" // TODO
)

func WebhookHandler(conf Config) (cmswebhook.Handler, error) {
	// s, err := NewServices(conf)
	// if err != nil {
	// 	return nil, err
	// }

	return func(req *http.Request, w *cmswebhook.Payload) error {
		ctx := req.Context()

		if !w.Operator.IsUser() && w.Operator.IsIntegrationBy(conf.CMSIntegration) {
			log.Debugfc(ctx, "cmsintegrationv2 webhook: invalid event operator: %+v", w.Operator)
			return nil
		}

		if w.Type != cmswebhook.EventItemCreate && w.Type != cmswebhook.EventItemUpdate {
			log.Debugfc(ctx, "cmsintegrationv2 webhook: invalid event type: %s", w.Type)
			return nil
		}

		if w.ItemData == nil || w.ItemData.Item == nil || w.ItemData.Model == nil {
			log.Debugfc(ctx, "cmsintegrationv2 webhook: invalid event data: %+v", w.Data)
			return nil
		}

		if w.ItemData.Model.Key != modelKey {
			log.Debugfc(ctx, "cmsintegrationv2 webhook: invalid model id: %s, key: %s", w.ItemData.Item.ModelID, w.ItemData.Model.Key)
			return nil
		}

		// TODO

		return nil
	}, nil
}
