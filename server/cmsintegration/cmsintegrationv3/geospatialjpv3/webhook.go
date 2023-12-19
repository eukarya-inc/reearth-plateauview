package geospatialjpv3

import (
	"net/http"

	"github.com/reearth/reearth-cms-api/go/cmswebhook"
)

type Config struct{}

func WebhookHandler(conf Config) (cmswebhook.Handler, error) {
	return func(req *http.Request, w *cmswebhook.Payload) error {
		// TODO
		return nil
	}, nil
}
