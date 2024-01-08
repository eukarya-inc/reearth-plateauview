package geospatialjpv3

import (
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/ckan"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
	"github.com/reearth/reearthx/log"
)

const (
	modelKey = "plateau"
)

type handler struct {
	cms  cms.Interface
	ckan ckan.Interface
}

func (h *handler) Webhook(conf Config) (cmswebhook.Handler, error) {
	return func(req *http.Request, w *cmswebhook.Payload) error {
		if req == nil || w == nil {
			log.Debug("geospatialjp webhook: invalid payload")
			return nil
		}

		ctx := req.Context()

		if !w.Operator.IsUser() && w.Operator.IsIntegrationBy(conf.CMSIntegration) {
			log.Debugfc(ctx, "geospatialjp webhook: invalid event operator: %+v", w.Operator)
			return nil
		}

		if w.Type != cmswebhook.EventItemCreate && w.Type != cmswebhook.EventItemUpdate && w.Type != cmswebhook.EventItemPublish {
			log.Debugfc(ctx, "geospatialjp webhook: invalid event type: %s", w.Type)
			return nil
		}

		if w.ItemData == nil || w.ItemData.Item == nil || w.ItemData.Model == nil {
			log.Debugfc(ctx, "geospatialjp webhook: invalid event data: %+v", w.Data)
			return nil
		}

		if w.ItemData.Model.Key != modelKey {
			log.Debugfc(ctx, "geospatialjp webhook: invalid model id: %s, key: %s", w.ItemData.Item.ModelID, w.ItemData.Model.Key)
			return nil
		}

		log.Debugfc(ctx, "geospatialjp webhook")
		// TODO prepare data (start async job)

		// TODO: create resources to ckan and publish

		// TODO: make resources private
		return nil
	}, nil
}
