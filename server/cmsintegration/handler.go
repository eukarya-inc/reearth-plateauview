package cmsintegration

import (
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationcommon"
	"github.com/eukarya-inc/reearth-plateauview/server/cmsintegration/cmsintegrationv2"
	"github.com/labstack/echo/v4"
	"github.com/reearth/reearth-cms-api/go/cmswebhook"
)

type Config = cmsintegrationcommon.Config

func NotifyHandler(conf Config) (echo.HandlerFunc, error) {
	return cmsintegrationv2.NotifyHandler(conf)
}

func WebhookHandler(conf Config) (cmswebhook.Handler, error) {
	return cmsintegrationv2.WebhookHandler(conf)
}
