package geospatialjpv3

import (
	"net/http"

	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/ckan"
	cms "github.com/reearth/reearth-cms-api/go"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
)

type Config struct{}

type handler struct {
	cms  cms.Interface
	ckan ckan.Interface
}

func (h *handler) Webhook(conf Config) (cmswebhook.Handler, error) {
	return func(req *http.Request, w *cmswebhook.Payload) error {
		// TODO
		return nil
	}, nil
}
